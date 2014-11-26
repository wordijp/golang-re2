// テスト時のUtility関数
package re2

import (
	"fmt"
	"reflect"
	"runtime"
	"testing"
)

// testing.T.Error(x)に、関数名や行番号等を追加する
func detailErrorfImpl(t *testing.T, skip int, format string, args ...interface{}) {
	pc, _, line, _ := runtime.Caller(skip + 1) // +1 : skip self method
	name := runtime.FuncForPC(pc).Name()
	msg := fmt.Sprintf(format, args...)
	t.Error(fmt.Sprintf("%s():%d: %s", name, line, msg))
}

func detailErrorImpl(t *testing.T, skip int, args ...interface{}) {
	pc, _, line, _ := runtime.Caller(skip + 1) // +1 : skip self method
	name := runtime.FuncForPC(pc).Name()
	msg := fmt.Sprint(args...)
	t.Error(fmt.Sprintf("%s():%d: %s", name, line, msg))
}

// この関数を呼び出した位置の、関数名や行番号等を追加で表示
func detailErrorf(t *testing.T, format string, args ...interface{}) {
	detailErrorfImpl(t, 1, format, args...)
}

func detailError(t *testing.T, args ...interface{}) {
	detailErrorImpl(t, 1, args...)
}

// この関数を呼び出した親の位置の、関数名や行番号等を追加で表示
func detailErrorfParent(t *testing.T, format string, args ...interface{}) {
	detailErrorfImpl(t, 1, format, args...)
	detailErrorfImpl(t, 2, "(called from here)")
}

func detailErrorParent(t *testing.T, args ...interface{}) {
	detailErrorImpl(t, 1, args...)
	detailErrorImpl(t, 2, "(called from here)")
}

// aとbの等価チェック
func equals_i(t *testing.T, a, b int) bool {
	if a != b {
		detailErrorfParent(t, "int: %d != %d", a, b)
		return false
	}

	return true
}

func equals_ai(t *testing.T, a, b []int) bool {
	if len(a) != len(b) {
		detailErrorfParent(t, "length: %d != %d", len(a), len(b))
		return false
	}

	for i := 0; i < len(a); i++ {
		if !equals_i(t, a[i], b[i]) {
			detailErrorfParent(t, "wrong. []int at(%d)", i)
			return false
		}
	}

	return true
}

func equals_aai(t *testing.T, a, b [][]int) bool {
	if len(a) != len(b) {
		detailErrorfParent(t, "length: %d != %d", len(a), len(b))
		return false
	}

	for i := 0; i < len(a); i++ {
		if !equals_ai(t, a[i], b[i]) {
			detailErrorfParent(t, "wrong. [][]int at(%d)", i)
			return false
		}
	}

	return true
}

func equals_s(t *testing.T, a, b string) bool {
	if a != b {
		detailErrorfParent(t, "string: %s != %s", a, b)
		return false
	}

	return true
}

func equals_as(t *testing.T, a, b []string) bool {
	if len(a) != len(b) {
		detailErrorfParent(t, "length: %d != %d", len(a), len(b))
		return false
	}

	for i := 0; i < len(a); i++ {
		if !equals_s(t, a[i], b[i]) {
			detailErrorfParent(t, "wrong. []string at(%d)", i)
			return false
		}
	}

	return true
}

func equals_aas(t *testing.T, a, b [][]string) bool {
	if len(a) != len(b) {
		detailErrorfParent(t, "length: %d != %d", len(a), len(b))
		return false
	}

	for i := 0; i < len(a); i++ {
		if !equals_as(t, a[i], b[i]) {
			detailErrorfParent(t, "wrong. [][]string at(%d)", i)
			return false
		}
	}

	return true
}

func equals_ab(t *testing.T, a, b []byte) bool {
	if !reflect.DeepEqual(a, b) {
		detailErrorfParent(t, "[]byte: %s != %s", fmt.Sprint(a), fmt.Sprint(b))
		detailErrorfParent(t, "[]byte(to string): %s != %s", string(a), string(b))
		return false
	}

	return true
}

func equals_aab(t *testing.T, a, b [][]byte) bool {
	if len(a) != len(b) {
		detailErrorfParent(t, "length: %d != %d", len(a), len(b))
		return false
	}

	for i := 0; i < len(a); i++ {
		if !equals_ab(t, a[i], b[i]) {
			detailErrorfParent(t, "wrong. [][]byte at(%d)", i)
			return false
		}
	}

	return true
}

func equals_aaab(t *testing.T, a, b [][][]byte) bool {
	if len(a) != len(b) {
		detailErrorfParent(t, "length: %d != %d", len(a), len(b))
		return false
	}

	for i := 0; i < len(a); i++ {
		if !equals_aab(t, a[i], b[i]) {
			detailErrorfParent(t, "wrong. [][][]byte at(%d)", i)
			return false
		}
	}

	return true
}
