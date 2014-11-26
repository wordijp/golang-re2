package re2

import (
	"testing"
)

func TestReplaceRE2Sequences(t *testing.T) {
	var check = func(expr, answer []byte) {
		ret := replaceRE2Sequences(expr)
		if string(ret) != string(answer) {
			detailErrorfParent(t, "wrong. %s != %s\n", string(ret), string(answer))
		}
	}
	check([]byte("$1"), []byte("\\1"))
	check([]byte("$$1"), []byte("$\\1"))
	check([]byte("$1$2"), []byte("\\1\\2"))
	check([]byte("\\$1"), []byte("\\\\1"))

	check([]byte("$"), []byte("$"))
	check([]byte("\\"), []byte("\\"))
	check([]byte("\\1"), []byte("\\1"))
	check([]byte("$\\1"), []byte("$\\1"))
	check([]byte("1"), []byte("1"))

	check([]byte("$$ 1"), []byte("$$ 1"))
	check([]byte("$ 1$2"), []byte("$ 1\\2"))

	check([]byte("$1$"), []byte("\\1$"))
	check([]byte("$1$2$"), []byte("\\1\\2$"))
}

func TestReplaceRE2InvalidSequences(t *testing.T) {
	var check = func(repl, answer []byte) {
		ret := replaceRE2InvalidSequences(repl)
		if string(ret) != string(answer) {
			detailErrorfParent(t, "wrong. %s != %s\n", string(ret), string(answer))
		}
	}
	check([]byte("\\1"), []byte("\\\\1"))
	check([]byte("\\\\1"), []byte("\\\\\\1"))
	check([]byte("\\1\\2"), []byte("\\\\1\\\\2"))

	check([]byte("\\"), []byte("\\"))
	check([]byte("\\\\"), []byte("\\\\"))
	check([]byte("$\\1"), []byte("$\\\\1"))
	check([]byte("1"), []byte("1"))

	check([]byte("\\ 1"), []byte("\\ 1"))
	check([]byte("\\ 1\\2"), []byte("\\ 1\\\\2"))

	check([]byte("\\1\\"), []byte("\\\\1\\"))
	check([]byte("\\1\\2\\"), []byte("\\\\1\\\\2\\"))
}

const (
	repl = "@$1@"
)

func BenchmarkReplaceRE2Sequences(b *testing.B) {
	bytes := []byte(repl)
	for i := 0; i < b.N; i++ {
		_ = replaceRE2Sequences(bytes)
	}
}
