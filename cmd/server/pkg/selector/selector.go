package selector

import "github.com/bamboo-firewall/be/cmd/server/pkg/selector/parser"

type Selector interface {
	// Evaluate evaluates the selector against the given labels expressed as a concrete map
	Evaluate(labels map[string]string) bool

	// String returns a string that represents this selector
	String() string
}

// Parse a string representation of a selector expression into a Selector.
func Parse(selector string) (Selector, error) {
	return parser.Parse(selector)
}
