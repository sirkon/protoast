package protoast

import (
	"slices"
	"strings"

	"github.com/sirkon/protoast/internal/errors"
)

// normalizeOptionName валидирует и возвращает нормализованное название опции.
func normalizeOptionName(opt string) (string, error) {
	lexems, err := getLexems(opt)
	if err != nil {
		return "", errors.Wrap(err, "decompose option name")
	}

	if len(lexems) == 0 {
		return "", errors.New("invalid syntax")
	}

	rest, err := validateOption(lexems)
	if err != nil {
		return "", err
	}

	if len(rest) == len(lexems) {
		return "", errors.New("invalid syntax")
	}

	if len(rest) > 0 {
		return "", errors.New("invalid syntax")
	}

	nobracks := slices.DeleteFunc(lexems, func(l lexemType) bool {
		return l.typ == lexemTypeCodeOpen || l.typ == lexemTypeCodeClose
	})
	var res strings.Builder
	for _, nobrack := range nobracks {
		switch nobrack.typ {
		case lexemTypeCodeIdent:
			res.WriteString(nobrack.value)
		case lexemTypeCodeDot:
			res.WriteByte('.')
		}
	}

	return res.String(), nil
}

func validateOption(opt []lexemType) (rest []lexemType, err error) {
	if len(opt) == 0 {
		return nil, nil
	}

	for len(opt) > 0 {
		s := opt[0]
		switch s.typ {
		case lexemTypeCodeIdent:
			if len(opt) == 1 {
				return nil, nil
			}

			opt = opt[1:]
		case lexemTypeCodeOpen:
			if len(opt) == 1 {
				return nil, errors.New("unexpected open bracket")
			}
			rest, err = validateOption(opt[1:])
			if err != nil {
				return nil, err
			}
			if len(rest) == 0 {
				return nil, errors.New("missing close bracket")
			}
			if rest[0].typ != lexemTypeCodeClose {
				return nil, errors.New("missing close bracket")
			}
			opt = rest[1:]
			if len(opt) == 0 {
				return nil, nil
			}
		case lexemTypeCodeClose:
			return opt, nil
		case lexemTypeCodeDot:
			return nil, errors.New("invalid syntax: unexpected dot character")
		}

		switch opt[0].typ {
		case lexemTypeCodeDot:
			opt = opt[1:]
			if len(opt) == 0 {
				return nil, errors.New("invalid syntax: missing correct option expression after a dot")
			}
			continue
		case lexemTypeCodeClose:
			return opt, nil
		}

		return nil, errors.New("invalid syntax")
	}

	return nil, nil
}

type lexemTypeCode int

func (c lexemTypeCode) String() string {
	switch c {
	case lexemTypeCodeNone:
		return "none"
	case lexemTypeCodeIdent:
		return "ident"
	case lexemTypeCodeDot:
		return "dot"
	case lexemTypeCodeOpen:
		return "open"
	case lexemTypeCodeClose:
		return "close"
	default:
		panic(errors.Newf("unknown lexem type code %d", c))
	}
}

const (
	lexemTypeCodeNone lexemTypeCode = iota
	lexemTypeCodeIdent
	lexemTypeCodeDot
	lexemTypeCodeOpen
	lexemTypeCodeClose
)

type lexemType struct {
	value string
	typ   lexemTypeCode
}

func (l lexemType) String() string {
	return l.typ.String() + ":" + l.value
}

func getLexems(s string) (res []lexemType, err error) {
	for s != "" {
		switch s[0] {
		case '(':
			res = append(res, lexemType{typ: lexemTypeCodeOpen})
			s = s[1:]
		case ')':
			res = append(res, lexemType{typ: lexemTypeCodeClose})
			s = s[1:]
		case '.':
			res = append(res, lexemType{typ: lexemTypeCodeDot})
			s = s[1:]
		default:
			ident, rest := getIdentifier(s)
			if ident == "" {
				return nil, errors.Newf("invalid syntax: %q", s)
			}
			res = append(res, lexemType{typ: lexemTypeCodeIdent, value: ident})
			s = rest
		}
	}

	return res, nil
}

func getIdentifier(s string) (ident string, rest string) {
	var i int
	if ('a' <= s[0] && s[0] <= 'z') || ('A' <= s[0] && s[0] <= 'Z') || s[0] == '_' {
		i++
	} else {
		return "", s
	}

	for i < len(s) {
		if ('a' <= s[i] && s[i] <= 'z') || ('A' <= s[i] && s[i] <= 'Z') || ('0' <= s[i] && s[i] <= '9') || s[i] == '_' {
			i++
			continue
		}

		return s[:i], s[i:]
	}

	return s[:i], s[i:]
}
