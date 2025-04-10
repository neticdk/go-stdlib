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
	// OpKind is the operation performed on the line
	Kind OpKind

	// Text is the content of the line
	Text string
}

// The Differ interface defines the contract for diffing two slices of strings
// Use the factory functions from each package to create a Differ implementation
type Differ interface {
	// Diff returns a string representation of the differences between two
	// strings. It returns an error if invalid options are provided or diffing
	// fails.
	Diff(a, b string) (string, error)

	// DiffStrings returns a string representation of the differences between
	// two slices of strings. It returns an error if invalid options are
	// provided or diffing fails.
	DiffStrings(a, b []string) (string, error)
}
