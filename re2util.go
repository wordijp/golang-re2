package re2

import (
	"bytes"
)

// $nを\\nへ置換
// (Goのregexpのsequenceは$nだが、re2のsequenceは\\n)
func replaceRE2Sequences(b []byte) []byte {
	var buf bytes.Buffer

	escape := false
	digit := false

	for _, x := range b {

		n := int(x) - int('0')

		switch {
		case x == '$':
			digit = false
			if escape {
				buf.WriteByte('$')
			}
			escape = true
			break

		case n >= 0 && n <= 9:
			if !digit && escape {
				buf.WriteByte('\\')
				digit = true
				escape = false
			}
			buf.WriteByte(x)
			break

		default:
			digit = false
			if escape {
				buf.WriteByte('$')
				escape = false
			}
			buf.WriteByte(x)
			break
		}
	}

	if escape {
		buf.WriteByte('$')
	}

	ret := make([]byte, buf.Len(), buf.Len()+1)
	buf.Read(ret)
	return ret
}

func replaceRE2InvalidSequences(b []byte) []byte {
	var buf bytes.Buffer

	escape := false
	digit := false

	for _, x := range b {

		n := int(x) - int('0')

		switch {
		case x == '\\':
			digit = false
			if escape {
				buf.WriteRune('\\')
			}
			escape = true
			break

		case n >= 0 && n <= 9:
			if !digit && escape {
				buf.WriteRune('\\')
				buf.WriteRune('\\')
				digit = true
				escape = false
			}
			buf.WriteByte(x)
			break

		default:
			digit = false
			if escape {
				buf.WriteRune('\\')
				escape = false
			}
			buf.WriteByte(x)
			break
		}
	}

	if escape {
		buf.WriteRune('\\')
	}

	tmp := make([]byte, buf.Len(), buf.Len()+1)
	buf.Read(tmp)

	return tmp
}
