package diff

// OpKind represents the type of operation performed on a line of text
type OpKind int

const (
	Insert OpKind = iota
	Delete
	Equal
)

func (o OpKind) String() string {
	switch o {
	case Insert:
		return "insert"
	case Delete:
		return "delete"
	case Equal:
		return "equal"
	default:
		return "Unknown"
	}
}

// Line represents a line of text in the diff with an associated operation
type Line struct {
	Kind OpKind // Op is the operation performed on the line
	Text string // Text is the content of the line
}

// The Differ interface defines the contract for diffing two slices of strings
// Use the factory functions from each package to create a Differ implementation
type Differ interface {
	// Diff returns a string representation of the differences between two
	// strings
	Diff(a, b string) string
	// DiffStrings returns a string representation of the differences between
	// two slices of strings
	DiffStrings(a, b []string) string
}
