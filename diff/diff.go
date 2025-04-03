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
