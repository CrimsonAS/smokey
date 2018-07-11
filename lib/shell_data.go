package lib

// Data passed in a shell pipeline.
type ShellData interface {
	// Return the underlying data for this shell item
	// For example, for a file, return the contents
	Data() ShellBuffer

	// Present the shell object in a nice way for display
	Present() string
}

// A type of data that has numeric columns
type ListyShellData interface {
	SelectColumn(col int) ShellData
}

// A type of data that has associative key:value data
type AssociativeShellData interface {
	SelectProperty(property string) ShellData
}

// A type of data that can be turned into multiple pieces of information
// For instance, a string can be exploded into line strings.
type ExplodableData interface {
	Explode() []ShellData
}
