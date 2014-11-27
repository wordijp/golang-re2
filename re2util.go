package re2

import ()

// $nを\\nへ置換
// (Goのregexpのsequenceは$nだが、re2のsequenceは\\n)
func replaceRE2Sequences(b []byte) []byte {

	var buf []byte

	escape := false
	digit := false

	for _, x := range b {

		n := int(x) - int('0')

		switch {
		case x == '$':
			digit = false
			if escape {
				buf = append(buf, '$')
			}
			escape = true
			break

		case n >= 0 && n <= 9:

			if !digit && escape {
				buf = append(buf, '\\')
				digit = true
				escape = false
			}
			buf = append(buf, x)
			break

		default:
			digit = false
			if escape {
				buf = append(buf, '$')
				escape = false
			}
			buf = append(buf, x)
			break
		}
	}

	if escape {
		buf = append(buf, '$')
	}

	return buf
}

func replaceRE2InvalidSequences(b []byte) []byte {
	var buf []byte

	escape := false
	digit := false

	for _, x := range b {

		n := int(x) - int('0')

		switch {
		case x == '\\':
			digit = false
			if escape {
				buf = append(buf, '\\')
			}
			escape = true
			break

		case n >= 0 && n <= 9:
			if !digit && escape {
				buf = append(buf, '\\')
				buf = append(buf, '\\')
				digit = true
				escape = false
			}
			buf = append(buf, x)
			break

		default:
			digit = false
			if escape {
				buf = append(buf, '\\')
				escape = false
			}
			buf = append(buf, x)
			break
		}
	}

	if escape {
		buf = append(buf, '\\')
	}

	return buf
}
