package main

import (
	"fmt"
	"github.com/influxdata/influxdb/client/v2"
	"strings"
)

// A URI representing an influx connection
type influxConnection struct {
	addr string
	user string
	pass string
}

func (this *influxConnection) Grep(searchable string) bool {
	return strings.Contains(this.addr, searchable) || strings.Contains(this.user, searchable)
}

func (this *influxConnection) Data() shellBuffer {
	return shellBuffer{}
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

func (this *influxSeries) Grep(searchable string) bool {
	return false
}

func (this *influxSeries) Data() shellBuffer {
	return shellBuffer{}
}

func (this *influxSeries) Present() string {
	return fmt.Sprintf("Influx series: %s (tags: %s) partial: %t\nColumns: %s", this.Name, this.Tags, this.Partial, this.Columns)
}

type influxRow struct {
	Values []interface{}
}

func (this *influxRow) Grep(searchable string) bool {
	return strings.Contains(this.Present(), searchable)
}

func (this *influxRow) Data() shellBuffer {
	return shellBuffer{}
}

func (this *influxRow) Present() string {
	return fmt.Sprintf("%s\n", this.Values)
}

// Create an influx connection
type InfluxConnect struct{}

func (this InfluxConnect) Call(inChan chan shellData, outChan chan shellData, arguments []string) {
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

func (this InfluxQuery) Call(inChan chan shellData, outChan chan shellData, arguments []string) {
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

				for _, val := range series.Values {
					outChan <- &influxRow{
						Values: val,
					}
				}
			}
		}
	}
	close(outChan)
}
