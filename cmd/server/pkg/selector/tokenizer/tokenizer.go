package tokenizer

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type tokenKind uint8

const (
	TokenNone tokenKind = iota
	TokenLabel
	TokenStringLiteral
	TokenLBrace
	TokenRBrace
	TokenComma
	TokenEq
	TokenNe
	TokenIn
	TokenNot
	TokenNotIn
	TokenContains
	TokenStartWith
	TokenEndsWith
	TokenAll
	TokenHas
	TokenLParen
	TokenRParen
	TokenAnd
	TokenOr
	TokenGlobal
	TokenEOF
)

const whitespace = " \t"

type Token struct {
	Kind  tokenKind
	Value interface{}
}

const (
	labelKeyMatcher = `[a-zA-Z0-9_./-]{1,512}`
)

var (
	identifierRegex = regexp.MustCompile("^" + labelKeyMatcher)
	containsRegex   = regexp.MustCompile(`^contains`)
	startsWithRegex = regexp.MustCompile(`^starts\s*with`)
	endsWithRegex   = regexp.MustCompile(`^ends\s*with`)
	hasRegex        = regexp.MustCompile(`^has\(\s*(` + labelKeyMatcher + `)\s*\)`)
	allRegex        = regexp.MustCompile(`^all\(\s*\)`)
	notInRegex      = regexp.MustCompile(`^not\s*in\b`)
	inRegex         = regexp.MustCompile(`^in\b`)
	globalRegex     = regexp.MustCompile(`^global\(\s*\)`)
)

// Tokenize transform string to token slice
func Tokenize(s string) ([]Token, error) {
	var tokens []Token
	for {
		startLen := len(s)
		s = strings.TrimLeft(s, whitespace)
		if len(s) == 0 {
			tokens = append(tokens, Token{Kind: TokenEOF, Value: nil})
			break
		}
		var lastTokenKind = TokenNone
		if len(tokens) > 0 {
			lastTokenKind = tokens[len(tokens)-1].Kind
		}
		switch s[0] {
		case '(':
			tokens = append(tokens, Token{Kind: TokenLParen, Value: nil})
			s = s[1:]
		case ')':
			tokens = append(tokens, Token{Kind: TokenRParen, Value: nil})
			s = s[1:]
		case '"':
			s = s[1:]
			index := strings.Index(s, `"`)
			if index == -1 {
				return nil, errors.New("unterminated double quote")
			}
			value := s[0:index]
			tokens = append(tokens, Token{Kind: TokenStringLiteral, Value: value})
			s = s[index+1:]
		case '\'':
			s = s[1:]
			index := strings.Index(s, `'`)
			if index == -1 {
				return nil, errors.New("unterminated single quote")
			}
			value := s[0:index]
			tokens = append(tokens, Token{Kind: TokenStringLiteral, Value: value})
			s = s[index+1:]
		case '{':
			tokens = append(tokens, Token{Kind: TokenLBrace, Value: nil})
			s = s[1:]
		case '}':
			tokens = append(tokens, Token{Kind: TokenRBrace, Value: nil})
			s = s[1:]
		case ',':
			tokens = append(tokens, Token{Kind: TokenComma, Value: nil})
			s = s[1:]
		case '=':
			if len(s) > 1 && s[1] == '=' {
				tokens = append(tokens, Token{Kind: TokenEq, Value: nil})
				s = s[2:]
			} else {
				return nil, errors.New("expect ==")
			}
		case '!':
			if len(s) > 1 && s[1] == '=' {
				tokens = append(tokens, Token{Kind: TokenNe, Value: nil})
				s = s[2:]
			} else {
				tokens = append(tokens, Token{Kind: TokenNot, Value: nil})
				s = s[1:]
			}
		case '&':
			if len(s) > 1 && s[1] == '&' {
				tokens = append(tokens, Token{Kind: TokenAnd, Value: nil})
				s = s[2:]
			} else {
				return nil, errors.New("expect &&")
			}
		case '|':
			if len(s) > 1 && s[1] == '|' {
				tokens = append(tokens, Token{Kind: TokenOr, Value: nil})
				s = s[2:]
			} else {
				return nil, errors.New("expect ||")
			}
		default:
			// Handle less-simple cases with regex matches. We're already stripped any whitespace
			if lastTokenKind == TokenLabel {
				// IF we just saw a label, look for a contains/starts with/ends with operator instead if another label
				if idxs := containsRegex.FindStringIndex(s); idxs != nil {
					// "contains"
					tokens = append(tokens, Token{Kind: TokenContains, Value: nil})
					s = s[idxs[1]:]
				} else if idxs = startsWithRegex.FindStringIndex(s); idxs != nil {
					// "starts with"
					tokens = append(tokens, Token{Kind: TokenStartWith, Value: nil})
					s = s[idxs[1]:]
				} else if idxs = endsWithRegex.FindStringIndex(s); idxs != nil {
					// "ends with"
					tokens = append(tokens, Token{Kind: TokenEndsWith, Value: nil})
					s = s[idxs[1]:]
				} else if idxs = notInRegex.FindStringIndex(s); idxs != nil {
					// "not in"
					tokens = append(tokens, Token{Kind: TokenNotIn, Value: nil})
					s = s[idxs[1]:]
				} else if idxs = inRegex.FindStringIndex(s); idxs != nil {
					// "in"
					tokens = append(tokens, Token{Kind: TokenIn, Value: nil})
					s = s[idxs[1]:]
				} else {
					return nil, fmt.Errorf("unexpected characters after label '%v', was expecting an operator", tokens[len(tokens)-1].Value)
				}
			} else if idxs := hasRegex.FindStringSubmatchIndex(s); idxs != nil {
				// "has(label)"
				wholeMatchEnd := idxs[1]
				labelNameMatchStart := idxs[2]
				labelNameMatchEnd := idxs[3]
				labelName := s[labelNameMatchStart:labelNameMatchEnd]
				tokens = append(tokens, Token{Kind: TokenHas, Value: labelName})
				s = s[wholeMatchEnd:]
			} else if idxs = allRegex.FindStringIndex(s); idxs != nil {
				// "all"
				tokens = append(tokens, Token{Kind: TokenAll, Value: nil})
				s = s[idxs[1]:]
			} else if idxs = globalRegex.FindStringIndex(s); idxs != nil {
				// "global"
				tokens = append(tokens, Token{Kind: TokenGlobal, Value: nil})
				s = s[idxs[1]:]
			} else if idxs = identifierRegex.FindStringIndex(s); idxs != nil {
				// "label"
				wholeMatchEnd := idxs[1]
				identifier := s[0:wholeMatchEnd]
				tokens = append(tokens, Token{Kind: TokenLabel, Value: identifier})
				s = s[wholeMatchEnd:]
			} else {
				return nil, errors.New("unexpected characters")
			}
		}
		if len(s) >= startLen {
			return nil, errors.New("infinite loop detected in tokenizer")
		}
	}
	return tokens, nil
}
