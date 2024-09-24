package parser

import (
	"errors"
	"fmt"

	"github.com/bamboo-firewall/be/cmd/server/pkg/selector/tokenizer"
)

func Parse(selector string) (*selectorRoot, error) {
	tokens, err := tokenizer.Tokenize(selector)
	if err != nil {
		return nil, err
	}
	if tokens[0].Kind == tokenizer.TokenEOF {
		return &selectorRoot{root: &AllNode{}}, nil
	}

	// The "||" operator has the lowest precedence so we start with that.
	sel, remainingTokens, err := parseOrExpression(tokens)
	if err != nil {
		return nil, err
	}
	// EOF token
	if len(remainingTokens) != 1 {
		err = fmt.Errorf("unexpected content at the end of selector %+v", remainingTokens)
		return nil, err
	}
	return &selectorRoot{root: sel}, err
}

// parseOrExpression parses a one or more "&&" terms, separated by "||" operators.
func parseOrExpression(tokens []tokenizer.Token) (sel node, remainingTokens []tokenizer.Token, err error) {
	// look for the first expression
	andNodes := make([]node, 0)
	sel, remainingTokens, err = parseAndExpression(tokens)
	if err != nil {
		return
	}
	andNodes = append(andNodes, sel)

	// Then loop looking for "||" followed by an <expression>
	for {
		switch remainingTokens[0].Kind {
		case tokenizer.TokenOr:
			remainingTokens = remainingTokens[1:]
			sel, remainingTokens, err = parseAndExpression(remainingTokens)
			if err != nil {
				return
			}
			andNodes = append(andNodes, sel)
		default:
			if len(andNodes) == 1 {
				sel = andNodes[0]
			} else {
				sel = &OrNode{andNodes}
			}
			return
		}
	}
}

func parseAndExpression(tokens []tokenizer.Token) (sel node, remainingTokens []tokenizer.Token, err error) {
	// Look for first operator
	opNodes := make([]node, 0)
	sel, remainingTokens, err = parseOperator(tokens)
	if err != nil {
		return
	}
	opNodes = append(opNodes, sel)

	// Then loop looking for "&&" followed by another operator
	for {
		switch remainingTokens[0].Kind {
		case tokenizer.TokenAnd:
			remainingTokens = remainingTokens[1:]
			sel, remainingTokens, err = parseOperator(remainingTokens)
			if err != nil {
				return
			}
			opNodes = append(opNodes, sel)
		default:
			if len(opNodes) == 1 {
				sel = opNodes[0]
			} else {
				sel = &AndNode{opNodes}
			}
			return
		}
	}
}

func parseOperator(tokens []tokenizer.Token) (sel node, remainingTokens []tokenizer.Token, err error) {
	if len(tokens) == 0 {
		err = errors.New("unexpected and of a string looking for op")
		return
	}
	// First, collapse any leading "!" operators to a single boolean
	negated := false
	for {
		if tokens[0].Kind == tokenizer.TokenNot {
			negated = !negated
			tokens = tokens[1:]
		} else {
			break
		}
	}

	// Then, look for the various types of operator
	switch tokens[0].Kind {
	case tokenizer.TokenHas:
		sel = &HasNode{tokens[0].Value.(string)}
		remainingTokens = tokens[1:]
	case tokenizer.TokenAll:
		sel = &AllNode{}
		remainingTokens = tokens[1:]
	case tokenizer.TokenGlobal:
		sel = &GlobalNode{}
		remainingTokens = tokens[1:]
	case tokenizer.TokenLabel:
		// should have an operator and a literal.
		if len(tokens) < 3 {
			err = errors.New(fmt.Sprint("unexpected end of string in middle of op", tokens))
			return
		}
		switch tokens[1].Kind {
		case tokenizer.TokenEq:
			if tokens[2].Kind == tokenizer.TokenStringLiteral {
				sel = &LabelEqValueNode{tokens[0].Value.(string), tokens[2].Value.(string)}
				remainingTokens = tokens[3:]
			} else {
				err = errors.New("expected string")
			}
		case tokenizer.TokenNe:
			if tokens[2].Kind == tokenizer.TokenStringLiteral {
				sel = &LabelNeValueNode{tokens[0].Value.(string), tokens[2].Value.(string)}
				remainingTokens = tokens[3:]
			} else {
				err = errors.New("expected string")
			}
		case tokenizer.TokenContains:
			if tokens[2].Kind == tokenizer.TokenStringLiteral {
				sel = &LabelContainsValueNode{tokens[0].Value.(string), tokens[2].Value.(string)}
				remainingTokens = tokens[3:]
			} else {
				err = errors.New("expected string")
			}
		case tokenizer.TokenStartWith:
			if tokens[2].Kind == tokenizer.TokenStringLiteral {
				sel = &LabelStartsWithValueNode{tokens[0].Value.(string), tokens[2].Value.(string)}
				remainingTokens = tokens[3:]
			} else {
				err = errors.New("expected string")
			}
		case tokenizer.TokenEndsWith:
			if tokens[2].Kind == tokenizer.TokenStringLiteral {
				sel = &LabelEndsWithValueNode{tokens[0].Value.(string), tokens[2].Value.(string)}
				remainingTokens = tokens[3:]
			} else {
				err = errors.New("expected string")
			}
		case tokenizer.TokenIn, tokenizer.TokenNotIn:
			if tokens[2].Kind == tokenizer.TokenLBrace {
				remainingTokens = tokens[3:]
				var values []string
				for {
					if remainingTokens[0].Kind == tokenizer.TokenStringLiteral {
						values = append(values, remainingTokens[0].Value.(string))
						remainingTokens = remainingTokens[1:]
						if remainingTokens[0].Kind == tokenizer.TokenComma {
							remainingTokens = remainingTokens[1:]
						} else {
							break
						}
					} else {
						break
					}
				}
				if remainingTokens[0].Kind != tokenizer.TokenRBrace {
					err = errors.New("expected }")
				} else {
					// Skip over the }
					remainingTokens = remainingTokens[1:]

					labelName := tokens[0].Value.(string)
					set := ConvertToStringSetInPlace(values)
					if tokens[1].Kind == tokenizer.TokenIn {
						sel = &LabelInSetNode{LabelName: labelName, Value: set}
					} else {
						sel = &LabelNotInSetNode{LabelName: labelName, Value: set}
					}
				}
			} else {
				err = errors.New("expected set literal")
			}
		default:
			err = errors.New(fmt.Sprint("expected == or != not ", tokens[1]))
		}
	case tokenizer.TokenLParen:
		// We hit a paren, skip past it, then recurse
		sel, remainingTokens, err = parseOrExpression(tokens[1:])
		if err != nil {
			return
		}
		// After parsing the nested expression, there should be a matching paren
		if len(remainingTokens) < 1 || remainingTokens[0].Kind != tokenizer.TokenRParen {
			err = errors.New("expected )")
			return
		}
		remainingTokens = remainingTokens[1:]
	default:
		err = errors.New(fmt.Sprint("unexpected token: ", tokens[0]))
		return
	}
	if negated && err == nil {
		sel = &NotNode{sel}
	}
	return

}
