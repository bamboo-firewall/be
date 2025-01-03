package parser

import (
	"strings"
)

// Labels defines the interface of labels that can be used by selector
type Labels interface {
	// Get returns value and presence of the given labelName
	Get(labelName string) (string, bool)
}

type MapAsLabels map[string]string

func (m MapAsLabels) Get(labelName string) (value string, present bool) {
	value, present = m[labelName]
	return
}

type node interface {
	Evaluate(labels Labels) bool
	collectFragments(fragments []string) []string
}

type selectorRoot struct {
	root         node
	cachedString *string
}

func (r *selectorRoot) EvaluateLabels(labels Labels) bool {
	return r.root.Evaluate(labels)
}

func (r *selectorRoot) Evaluate(labels map[string]string) bool {
	return r.EvaluateLabels(MapAsLabels(labels))
}

func (r *selectorRoot) String() string {
	if r.cachedString == nil {
		fragments := r.root.collectFragments([]string{})
		joined := strings.Join(fragments, "")
		r.cachedString = &joined
	}
	return *r.cachedString
}

type LabelEqValueNode struct {
	LabelName string
	Value     string
}

func (node *LabelEqValueNode) Evaluate(labels Labels) bool {
	val, ok := labels.Get(node.LabelName)
	if !ok {
		return false
	}
	return val == node.Value
}

func (node *LabelEqValueNode) collectFragments(fragments []string) []string {
	return appendLabelOpAndQuotedString(fragments, node.LabelName, " == ", node.Value)
}

type LabelContainsValueNode struct {
	LabelName string
	Value     string
}

func (node *LabelContainsValueNode) Evaluate(labels Labels) bool {
	val, ok := labels.Get(node.LabelName)
	if !ok {
		return false
	}
	return strings.Contains(val, node.Value)
}

func (node *LabelContainsValueNode) collectFragments(fragments []string) []string {
	return appendLabelOpAndQuotedString(fragments, node.LabelName, " contains ", node.Value)
}

type LabelStartsWithValueNode struct {
	LabelName string
	Value     string
}

func (node *LabelStartsWithValueNode) Evaluate(labels Labels) bool {
	val, ok := labels.Get(node.LabelName)
	if ok {
		return strings.HasPrefix(val, node.Value)
	}
	return false
}

func (node *LabelStartsWithValueNode) collectFragments(fragments []string) []string {
	return appendLabelOpAndQuotedString(fragments, node.LabelName, " starts with ", node.Value)
}

type LabelEndsWithValueNode struct {
	LabelName string
	Value     string
}

func (node *LabelEndsWithValueNode) Evaluate(labels Labels) bool {
	val, ok := labels.Get(node.LabelName)
	if ok {
		return strings.HasSuffix(val, node.Value)
	}
	return false
}

func (node *LabelEndsWithValueNode) collectFragments(fragments []string) []string {
	return appendLabelOpAndQuotedString(fragments, node.LabelName, " ends with ", node.Value)
}

type LabelNeValueNode struct {
	LabelName string
	Value     string
}

func (node *LabelNeValueNode) Evaluate(labels Labels) bool {
	val, ok := labels.Get(node.LabelName)
	if ok {
		return val != node.Value
	}
	return true
}

func (node *LabelNeValueNode) collectFragments(fragments []string) []string {
	return appendLabelOpAndQuotedString(fragments, node.LabelName, " != ", node.Value)
}

func appendLabelOpAndQuotedString(fragments []string, label, op, s string) []string {
	var quote string
	if strings.Contains(s, `"`) {
		quote = `'`
	} else {
		quote = `"`
	}
	return append(fragments, label, op, quote, s, quote)
}

type LabelInSetNode struct {
	LabelName string
	Value     StringSet
}

func (node *LabelInSetNode) Evaluate(labels Labels) bool {
	val, ok := labels.Get(node.LabelName)
	if ok {
		return node.Value.Contains(val)
	}
	return false
}

func (node *LabelInSetNode) collectFragments(fragments []string) []string {
	return collectInOpFragments(fragments, node.LabelName, "in", node.Value)
}

type LabelNotInSetNode struct {
	LabelName string
	Value     StringSet
}

func (node *LabelNotInSetNode) Evaluate(labels Labels) bool {
	val, ok := labels.Get(node.LabelName)
	if ok {
		return !node.Value.Contains(val)
	}
	return true
}

func (node *LabelNotInSetNode) collectFragments(fragments []string) []string {
	return collectInOpFragments(fragments, node.LabelName, "not in", node.Value)
}

func collectInOpFragments(fragments []string, labelName, op string, values StringSet) []string {
	var quote string
	fragments = append(fragments, labelName, " ", op, " {")
	first := true
	for _, s := range values {
		if strings.Contains(s, `"`) {
			quote = `'`
		} else {
			quote = `"`
		}
		if !first {
			fragments = append(fragments, ", ")
		} else {
			first = false
		}
		fragments = append(fragments, quote, s, quote)
	}
	fragments = append(fragments, "}")
	return fragments
}

type HasNode struct {
	LabelName string
}

func (node *HasNode) Evaluate(labels Labels) bool {
	_, ok := labels.Get(node.LabelName)
	if ok {
		return true
	}
	return false
}

func (node *HasNode) collectFragments(fragments []string) []string {
	return append(fragments, "has(", node.LabelName, ")")
}

type NotNode struct {
	Operand node
}

func (node *NotNode) Evaluate(labels Labels) bool {
	return !node.Operand.Evaluate(labels)
}

func (node *NotNode) collectFragments(fragments []string) []string {
	fragments = append(fragments, "!")
	return node.Operand.collectFragments(fragments)
}

type AndNode struct {
	Operands []node
}

func (node *AndNode) Evaluate(labels Labels) bool {
	for _, operand := range node.Operands {
		if !operand.Evaluate(labels) {
			return false
		}
	}
	return true
}

func (node *AndNode) collectFragments(fragments []string) []string {
	fragments = append(fragments, "(")
	fragments = node.Operands[0].collectFragments(fragments)
	for _, op := range node.Operands[1:] {
		fragments = append(fragments, " && ")
		fragments = op.collectFragments(fragments)
	}
	fragments = append(fragments, ")")
	return fragments
}

type OrNode struct {
	Operands []node
}

func (node *OrNode) Evaluate(labels Labels) bool {
	for _, operand := range node.Operands {
		if operand.Evaluate(labels) {
			return true
		}
	}
	return false
}

func (node *OrNode) collectFragments(fragments []string) []string {
	fragments = append(fragments, "(")
	fragments = node.Operands[0].collectFragments(fragments)
	for _, op := range node.Operands[1:] {
		fragments = append(fragments, " || ")
		fragments = op.collectFragments(fragments)
	}
	fragments = append(fragments, ")")
	return fragments
}

type AllNode struct {
}

func (node *AllNode) Evaluate(Labels) bool {
	return true
}

func (node *AllNode) collectFragments(fragments []string) []string {
	return append(fragments, "all()")
}

type GlobalNode struct {
}

func (node *GlobalNode) Evaluate(labels Labels) bool {
	return true
}

func (node *GlobalNode) collectFragments(fragments []string) []string {
	return append(fragments, "global()")
}
