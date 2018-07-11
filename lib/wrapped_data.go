package lib

// Data pretending to be something else for presentation purposes.
// Essentially this acts as a proxy to RealData while displaying FakeData.
type WrappedData struct {
	RealData ShellData
	FakeData ShellData
}

func (this *WrappedData) Data() ShellBuffer {
	return this.RealData.Data()
}

func (this *WrappedData) Present() string {
	return this.FakeData.Present()
}

func (this *WrappedData) SelectColumn(col int) ShellData {
	if ld, ok := this.FakeData.(ListyShellData); ok {
		return ld.SelectColumn(col)
	} else {
		return nil
	}
}

func (this *WrappedData) SelectProperty(prop string) ShellData {
	if ld, ok := this.FakeData.(AssociativeShellData); ok {
		return ld.SelectProperty(prop)
	} else {
		return nil
	}
}

func (this *WrappedData) Explode() []ShellData {
	if ld, ok := this.FakeData.(ExplodableData); ok {
		return ld.Explode()
	} else {
		return nil
	}
}
