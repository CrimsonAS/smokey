package influx

import (
	"fmt"
	"github.com/CrimsonAS/smokey/lib"
	"github.com/influxdata/influxdb/client/v2"
)

// A URI representing an influx connection
type influxConnection struct {
	addr string
	user string
	pass string
}

func (this *influxConnection) Data() lib.ShellBuffer {
	return lib.ShellBuffer{}
}

func (this *influxConnection) Present() string {
	return fmt.Sprintf("influx:%s@%s\n", this.user, this.addr)
}

type influxSeries struct {
	Name    string
	Tags    map[string]string
	Columns []string
	Partial bool
}

func (this *influxSeries) Data() lib.ShellBuffer {
	return lib.ShellBuffer{}
}

func (this *influxSeries) Present() string {
	return fmt.Sprintf("Influx series: %s (tags: %s) partial: %t\nColumns: %s", this.Name, this.Tags, this.Partial, this.Columns)
}

type influxRow struct {
	Values []interface{}
	colMap map[string]int
}

func (this *influxRow) SelectProperty(prop string) lib.ShellData {
	if col, ok := this.colMap[prop]; ok {
		return lib.ShellString(fmt.Sprintf("%s", this.Values[col]))
	}
	return nil
}

func (this *influxRow) SelectColumn(col int) lib.ShellData {
	if col >= 0 && col < len(this.Values) {
		return lib.ShellString(fmt.Sprintf("%s", this.Values[col]))
	}
	return nil
}

func (this *influxRow) Data() lib.ShellBuffer {
	return lib.ShellBuffer{}
}

func (this *influxRow) Present() string {
	return fmt.Sprintf("%s\n", this.Values)
}

// Create an influx connection
type InfluxConnect struct{}

func (this InfluxConnect) Call(inChan chan lib.ShellData, outChan chan lib.ShellData, arguments []string) {
	if len(arguments) < 3 {
		panic("What do you want to connect to?")
	}

	outChan <- &influxConnection{
		addr: arguments[0],
		user: arguments[1],
		pass: arguments[2],
	}
	close(outChan)
}

// Execute an Influx query on every provided connection
type InfluxQuery struct{}

func (this InfluxQuery) Call(inChan chan lib.ShellData, outChan chan lib.ShellData, arguments []string) {
	if len(arguments) < 3 {
		panic("What do you want to query?")
	}

	for conn := range inChan {
		influxConn := conn.(*influxConnection)
		c, err := client.NewHTTPClient(client.HTTPConfig{
			Addr:     influxConn.addr,
			Username: influxConn.user,
			Password: influxConn.pass,
		})
		if err != nil {
			panic(fmt.Sprintf("Can't create HTTP client %s: %s", influxConn.addr, err))
		}

		q := client.NewQuery(arguments[0], arguments[1], arguments[2])
		resp, err := c.Query(q)
		if err != nil {
			panic(fmt.Sprintf("Error querying %s: %s", influxConn.addr, err))
		}
		if resp.Error() != nil { // because apparently returning an error isn't enough.
			panic(fmt.Sprintf("Error querying %s: %s", influxConn.addr, resp.Error()))
		}

		for _, result := range resp.Results {
			for _, series := range result.Series {
				outChan <- &influxSeries{
					Name:    series.Name,
					Tags:    series.Tags,
					Columns: series.Columns,
					Partial: series.Partial,

					// is Messages of interest?
				}

				colMap := make(map[string]int)
				for idx, col := range series.Columns {
					colMap[col] = idx
				}

				for _, val := range series.Values {
					outChan <- &influxRow{
						Values: val,
						colMap: colMap,
					}
				}
			}
		}
	}
	close(outChan)
}
