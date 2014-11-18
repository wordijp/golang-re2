// Package testutil.
// テスト時のUtility関数
package testutil

import (
	"fmt"
	"runtime"
	"testing"
)

// testing.T.Errorに、関数名や行番号等を追加する
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
func DetailErrorf(t *testing.T, format string, args ...interface{}) {
	detailErrorfImpl(t, 1, format, args...)
}

func DetailError(t *testing.T, args ...interface{}) {
	detailErrorImpl(t, 1, args...)
}

// この関数を呼び出した親の位置の、関数名や行番号等を追加で表示
func DetailErrorfParent(t *testing.T, format string, args ...interface{}) {
	detailErrorfImpl(t, 1, format, args...)
	detailErrorfImpl(t, 2, "(called from here)")
}

func DetailErrorParent(t *testing.T, args ...interface{}) {
	detailErrorImpl(t, 1, args...)
	detailErrorImpl(t, 2, "(called from here)")
}
