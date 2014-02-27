package augustine

import (
	"unicode"
)

type Normalizer interface {
	Normalize(tok []byte) (terms []string)
}

type simpleNormalizer struct {
	badChar [256]bool
}

const simpleBadChars = `!@#$%^&*()_+-={}[]:";'<>,.` + "`"

func NewSimpleNormalizer() Normalizer {
	n := new(simpleNormalizer)
	for _, b := range []byte(simpleBadChars) {
		n.badChar[int(b)] = true
	}
	return n
}

func (n *simpleNormalizer) Normalize(tok []byte) (terms []string) {
	term := make([]rune, 0, len(tok))
	for _, r := range string(tok) {
		if r < 256 && n.badChar[int(r)] {
			continue
		}
		term = append(term, unicode.ToLower(r))
	}
	if len(term) > 0 {
		return []string{string(term)}
	} else {
		return nil
	}
}
