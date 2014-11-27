package re2

// #cgo LDFLAGS: -lcre2
// #include "cre2.h"
// #include <stdlib.h>
// int length_ptr(char *begin, char *end) {
//     return (int)(end - begin);
// }
// char *next_ptr(char *pos) {
//   return ++pos;
// }
import "C"

import (
	"errors"
	"fmt"
	"math"
	"unsafe"
	//	"runtime"
)

// XXX : Reader系は未実装
//       (ファイルハンドラをStringPieceを引数で受け取るre2ラッパーに渡す上手い方法が思いつかない)

// XXX : String系は、内部で[]byteへの変換をしているため、文字列の値コピー分のオーバーヘッドが余計にかかっている

// XXX : (Must)?Compile(POSIX)?や、Longestのパフォーマンスが、regexpパッケージに比べて著しく遅い

// NOTE : regexpパッケージのRegexpと同じ動作になるようにしているが、
//        動作が違う一部のメソッドはRE2というprefixを付けて差別化を図っている(ex:ReplaceAll)

// regexpインスタンスを解放する処理を担う
// NOTE : メソッドはregexpと統一しているが、このライブラリだと解放処理が必要なので、
//        返り値に差異を持たせ、単純置き換え時のClose呼び出し忘れを防ぐ
type Closer struct{}

func (c *Closer) Close(re *Regexp) {
	if re != nil {
		C.cre2_delete(unsafe.Pointer(re.cre2_re))
		C.cre2_delete(unsafe.Pointer(re.cre2_re_with_bracket))
		C.cre2_opt_delete(unsafe.Pointer(re.cre2_opt))
		C.free(unsafe.Pointer(re.c_expr))
		C.free(unsafe.Pointer(re.c_expr_with_bracket))

		// GCが追い付かずにメモリアロケートが失敗する時に有効化
		// 普段使いだとまず気にする必要はないはず
		// runtime.GC()
	}
}

// re2では、MatchやFindAndConsume時、ヒット箇所の抜出には外枠の()が必要で、Replace時は必要ない
// その為、その差異を吸収する為に()の有り・無しの2種類を持つ
type Regexp struct {
	cre2_re                    *C.cre2_regexp_t
	cre2_re_with_bracket       *C.cre2_regexp_t
	cre2_opt                   *C.cre2_options_t
	c_expr                     *C.char
	c_expr_length              C.int
	c_expr_with_bracket        *C.char
	c_expr_with_bracket_length C.int
}

func MatchString(expr string, s string) (matched bool, err error) {
	re, closer, err := Compile(expr)
	defer closer.Close(re)
	if err != nil {
		return false, err
	}
	return re.MatchString(s), nil
}

func Match(expr string, b []byte) (matched bool, err error) {
	re, closer, err := Compile(expr)
	defer closer.Close(re)
	if err != nil {
		return false, err
	}
	return re.Match(b), nil
}

// [ ]
//func MatchReader(pattern string, r io.RuneReader) (matched bool, err error) {
//    // not implemented
//}

func RE2QuoteMeta(s string) string {
	b := []byte(s)
	c_input := (*C.char)(unsafe.Pointer(&b[0]))
	cre2_input := C.cre2_string_t{
		data:   c_input,
		length: C.int(len(b)),
	}

	var cre2_quoted C.cre2_string_t
	result := C.cre2_quote_meta(&cre2_quoted, &cre2_input)
	if !bool(result == 0) {
		panic("bad_alloc: malloc failed.")
	}

	ret := C.GoBytes(unsafe.Pointer(cre2_quoted.data), cre2_quoted.length)
	C.free(unsafe.Pointer(cre2_quoted.data))

	return string(ret)
}

func makeOptions(posix, longest, literal bool) *C.cre2_options_t {
	// encodingはUTF-8のみ
	opt := C.cre2_opt_new()
	//	C.cre2_opt_set_encoding(opt, C.CRE2_UTF8) // デフォルトでutf
	if posix {
		C.cre2_opt_set_posix_syntax(opt, 1)
	}
	if longest {
		C.cre2_opt_set_longest_match(opt, 1)
	}
	if literal {
		C.cre2_opt_set_literal(opt, 1)
	}
	C.cre2_opt_set_log_errors(opt, 0)

	return (*C.cre2_options_t)(opt)
}

func compile(expr string, posix, longest bool) (*Regexp, *Closer, error) {
	c_expr := C.CString(expr)
	c_expr_length := C.int(len(expr))
	c_expr_with_bracket := C.CString("(" + expr + ")") // RE2のFindAndConsume用に()を付与
	c_expr_with_bracket_length := C.int(len(expr) + 2)
	cre2_opt := makeOptions(posix, longest, false)

	closer := &Closer{}

	// XXX : cre2_newのコストが数倍になっている

	// ()無し
	cre2_re := C.cre2_new(c_expr, c_expr_length, unsafe.Pointer(cre2_opt))
	errCode := C.cre2_error_code(cre2_re)
	if errCode != C.CRE2_NO_ERROR {

		C.cre2_delete(unsafe.Pointer(cre2_re))
		C.cre2_opt_delete(unsafe.Pointer(cre2_opt))
		C.free(unsafe.Pointer(c_expr))
		C.free(unsafe.Pointer(c_expr_with_bracket))

		return nil, closer, errors.New(fmt.Sprintf("Compile error code:%d", errCode))
	}

	// ()有り
	cre2_re_with_bracket := C.cre2_new(c_expr_with_bracket, c_expr_with_bracket_length, unsafe.Pointer(cre2_opt))
	errCode = C.cre2_error_code(cre2_re)
	if errCode != C.CRE2_NO_ERROR {

		C.cre2_delete(unsafe.Pointer(cre2_re))
		C.cre2_delete(unsafe.Pointer(cre2_re_with_bracket))
		C.cre2_opt_delete(unsafe.Pointer(cre2_opt))
		C.free(unsafe.Pointer(c_expr))
		C.free(unsafe.Pointer(c_expr_with_bracket))

		return nil, closer, errors.New(fmt.Sprintf("Compile error code:%d", errCode))
	}

	re := &Regexp{
		cre2_re:                    (*C.cre2_regexp_t)(cre2_re),
		cre2_re_with_bracket:       (*C.cre2_regexp_t)(cre2_re_with_bracket),
		cre2_opt:                   cre2_opt,
		c_expr:                     c_expr,
		c_expr_length:              c_expr_length,
		c_expr_with_bracket:        c_expr_with_bracket,
		c_expr_with_bracket_length: c_expr_with_bracket_length,
	}
	return re, closer, nil // regexp.Compileと違い、Closerも返す
}

func Compile(expr string) (*Regexp, *Closer, error) {
	return compile(expr, false, false)
}

func CompilePOSIX(expr string) (*Regexp, *Closer, error) {
	return compile(expr, true, true)
}

func MustCompile(expr string) (*Regexp, *Closer) {
	re, closer, err := Compile(expr)
	if err != nil {
		panic(err)
	}
	return re, closer // regexp.MustCompileと違い、Closerも返す
}

func MustCompilePOSIX(expr string) (*Regexp, *Closer) {
	re, closer, err := CompilePOSIX(expr)
	if err != nil {
		panic(err)
	}
	return re, closer // regexp.MustCompilePOSIXと違い、Closerも返す
}

func (re *Regexp) RE2Expand(dst []byte, template []byte, src []byte, match []int) []byte {

	text := make([]byte, match[1]-match[0])
	copy(text, src[match[0]:match[1]])

	c_text := (*C.char)(unsafe.Pointer(&text[0]))
	cre2_text := C.cre2_string_t{
		data:   c_text,
		length: C.int(len(text)),
	}

	re2_rewrite := replaceRE2Sequences(template) // re2用のsequenceへ($n -> \\n)
	c_re2_rewrite := (*C.char)(unsafe.Pointer(&re2_rewrite[0]))
	cre2_re2_rewrite := C.cre2_string_t{
		data:   c_re2_rewrite,
		length: C.int(len(re2_rewrite)),
	}

	cre2_out := C.cre2_string_t{}
	C.cre2_extract_re(unsafe.Pointer(re.cre2_re), &cre2_text, &cre2_re2_rewrite, &cre2_out)

	ret := C.GoBytes(unsafe.Pointer(cre2_out.data), cre2_out.length)

	C.free(unsafe.Pointer(cre2_out.data)) // cre2_outはメモリが再確保されている
	// 古いバッファ(text)は[]byteなので解放はGoのGCが担う

	return append(dst, ret...)
}

func (re *Regexp) RE2ExpandString(dst []byte, template string, src string, match []int) []byte {
	return re.RE2Expand(dst, []byte(template), []byte(src), match)
}

func (re *Regexp) Find(b []byte) []byte {
	ret := re.FindAll(b, 1)
	if ret == nil {
		return nil
	}
	if len(ret) != 1 {
		panic(fmt.Sprintf("len(ret):%d != 1", len(ret)))
	}
	return ret[0]
}

func (re *Regexp) FindAll(b []byte, n int) [][]byte {

	c_input := (*C.char)(unsafe.Pointer(&b[0]))
	cre2_input := C.cre2_string_t{
		data:   c_input,
		length: C.int(len(b)),
	}

	var _n int
	if n >= 0 {
		_n = n
	} else {
		// XXX : 無限ループにしなくて良いのか?
		_n = math.MaxInt32
	}

	var ret [][]byte
	for i := 0; i < _n; i++ {
		var cre2_match [1]C.cre2_string_t
		result := C.cre2_find_and_consume_re(unsafe.Pointer(re.cre2_re_with_bracket), &cre2_input, &cre2_match[0], 1)
		if !bool(result != 0) {
			break
		}

		bytes := C.GoBytes(unsafe.Pointer(cre2_match[0].data), cre2_match[0].length)
		ret = append(ret, bytes)
	}

	if len(ret) == 0 {
		return nil
	}
	return ret
}

func (re *Regexp) FindAllIndex(b []byte, n int) [][]int {

	c_input := (*C.char)(unsafe.Pointer(&b[0]))
	cre2_input := C.cre2_string_t{
		data:   c_input,
		length: C.int(len(b)),
	}

	var _n int
	if n >= 0 {
		_n = n
	} else {
		// XXX : 無限ループにしなくて良いのか?
		_n = math.MaxInt32
	}
	find := false

	var ret [][]int
	for i := 0; i < _n; i++ {
		var cre2_match [1]C.cre2_string_t
		result := C.cre2_find_and_consume_re(unsafe.Pointer(re.cre2_re_with_bracket), &cre2_input, &cre2_match[0], 1)
		if !bool(result != 0) {
			break
		}

		offset_begin := int(C.length_ptr(c_input, cre2_match[0].data))
		offset_end := offset_begin + int(cre2_match[0].length)

		// 見つからなかった場合は一つずらす
		if offset_end == offset_begin {
			cre2_input.data = C.next_ptr(cre2_input.data)
			cre2_input.length -= 1

			if find {
				i--
				find = false
				continue
			}
		} else {
			find = true
		}

		indexes := []int{
			offset_begin,
			offset_end,
		}
		ret = append(ret, indexes)
	}

	if len(ret) == 0 {
		return nil
	}
	return ret
}

func (re *Regexp) FindAllString(s string, n int) []string {
	bytes := re.FindAll([]byte(s), n)
	if bytes == nil {
		return nil
	}

	var ret []string
	for _, x := range bytes {
		ret = append(ret, string(x))
	}
	return ret
}

func (re *Regexp) FindAllStringIndex(s string, n int) [][]int {
	return re.FindAllIndex([]byte(s), n)
}

func (re *Regexp) FindAllStringSubmatch(s string, n int) [][]string {
	bytes := re.FindAllSubmatch([]byte(s), n)
	if bytes == nil {
		return nil
	}

	var ret [][]string
	for _, x := range bytes {
		var a []string
		for _, y := range x {
			a = append(a, string(y))
		}
		ret = append(ret, a)
	}
	return ret
}

func (re *Regexp) FindAllStringSubmatchIndex(s string, n int) [][]int {
	return re.FindAllSubmatchIndex([]byte(s), n)
}

func (re *Regexp) FindAllSubmatch(b []byte, n int) [][][]byte {

	c_input := (*C.char)(unsafe.Pointer(&b[0]))
	cre2_input := C.cre2_string_t{
		data:   c_input,
		length: C.int(len(b)),
	}

	var _n int
	if n >= 0 {
		_n = n
	} else {
		// XXX : 無限ループにしなくて良いのか?
		_n = math.MaxInt32
	}

	var ret [][][]byte
	for i := 0; i < _n; i++ {

		groups := C.cre2_num_capturing_groups(unsafe.Pointer(re.cre2_re_with_bracket))

		cre2_match := make([]C.cre2_string_t, int(groups))
		result := C.cre2_find_and_consume_re(unsafe.Pointer(re.cre2_re_with_bracket), &cre2_input, &cre2_match[0], groups)
		if !bool(result != 0) {
			break
		}

		var bytes [][]byte
		for i := 0; i < int(groups); i++ {
			b := C.GoBytes(unsafe.Pointer(cre2_match[i].data), cre2_match[i].length)
			bytes = append(bytes, b)
		}
		ret = append(ret, bytes)
	}

	if len(ret) == 0 {
		return nil
	}
	return ret
}

func (re *Regexp) FindAllSubmatchIndex(b []byte, n int) [][]int {
	c_input := (*C.char)(unsafe.Pointer(&b[0]))
	cre2_input := C.cre2_string_t{
		data:   c_input,
		length: C.int(len(b)),
	}

	var _n int
	if n >= 0 {
		_n = n
	} else {
		// XXX : 無限ループにしなくて良いのか?
		_n = math.MaxInt32
	}

	var ret [][]int
	for i := 0; i < _n; i++ {

		groups := C.cre2_num_capturing_groups(unsafe.Pointer(re.cre2_re_with_bracket))

		cre2_match := make([]C.cre2_string_t, int(groups))
		result := C.cre2_find_and_consume_re(unsafe.Pointer(re.cre2_re_with_bracket), &cre2_input, &cre2_match[0], groups)
		if !bool(result != 0) {
			break
		}

		var indexes []int
		for j := 0; j < int(groups); j++ {
			offset_begin := int(C.length_ptr(c_input, cre2_match[j].data))
			offset_end := offset_begin + int(cre2_match[j].length)

			indexes = append(indexes, offset_begin)
			indexes = append(indexes, offset_end)
		}
		ret = append(ret, indexes)
	}

	if len(ret) == 0 {
		return nil
	}
	return ret
}

func (re *Regexp) FindIndex(b []byte) (loc []int) {
	ret := re.FindAllIndex(b, 1)
	if ret == nil {
		return nil
	}
	if len(ret) != 1 {
		panic(fmt.Sprintf("len(ret):%d != 1", len(ret)))
	}
	return ret[0]
}

// [ ]
//func (re *Regexp) FindReaderIndex(r io.RuneReader) (loc []int) {
//    // not implemented
//}

// [ ]
//func (re *Regexp) FindReaderSubmatchIndex(r io.RuneReader) []int {
//    // not implemented
//}

func (re *Regexp) FindString(s string) string {
	return string(re.Find([]byte(s)))
}

func (re *Regexp) FindStringIndex(s string) (loc []int) {
	return re.FindIndex([]byte(s))
}

func (re *Regexp) FindStringSubmatch(s string) []string {
	bytes := re.FindSubmatch([]byte(s))

	var ret []string
	for _, x := range bytes {
		ret = append(ret, string(x))
	}
	return ret
}

func (re *Regexp) FindStringSubmatchIndex(s string) []int {
	return re.FindSubmatchIndex([]byte(s))
}

func (re *Regexp) FindSubmatch(b []byte) [][]byte {

	ret := re.FindAllSubmatch(b, 1)
	if ret == nil {
		return nil
	}
	if len(ret) != 1 {
		panic(fmt.Sprintf("len(ret):%d != 1", len(ret)))
	}
	return ret[0]
}

func (re *Regexp) FindSubmatchIndex(b []byte) []int {
	ret := re.FindAllSubmatchIndex(b, 1)
	if ret == nil {
		return nil
	}
	if len(ret) != 1 {
		panic(fmt.Sprintf("len(ret):%d !- 1", len(ret)))
	}
	return ret[0]
}

// [ ]
//func (re *Regexp) LiteralPrefix() (prefix string, complete bool) {
//    // not implemented
//}

// NOTE : re2のoptionsは変更出来ないので、インスタンスを作り変える
func (re *Regexp) Longest() {
	longest := C.cre2_opt_longest_match(unsafe.Pointer(re.cre2_opt))
	if bool(longest != 0) {
		return // do not need
	}
	posix := C.cre2_opt_posix_syntax(unsafe.Pointer(re.cre2_opt))

	// インスタンスの作り変え
	newOpt := makeOptions(bool(posix != 0), true, false)
	newOrigRe := C.cre2_new(re.c_expr, re.c_expr_length, unsafe.Pointer(newOpt))
	newOrigRe_with_bracket := C.cre2_new(re.c_expr_with_bracket, re.c_expr_with_bracket_length, unsafe.Pointer(newOpt))

	C.cre2_delete(unsafe.Pointer(re.cre2_re))
	C.cre2_delete(unsafe.Pointer(re.cre2_re_with_bracket))
	C.cre2_opt_delete(unsafe.Pointer(re.cre2_opt))

	re.cre2_re = (*C.cre2_regexp_t)(newOrigRe)
	re.cre2_re_with_bracket = (*C.cre2_regexp_t)(newOrigRe_with_bracket)
	re.cre2_opt = newOpt
}

func (re *Regexp) Match(b []byte) bool {
	c_input := (*C.char)(unsafe.Pointer(&b[0]))
	cre2_input := C.cre2_string_t{
		data:   c_input,
		length: C.int(len(b)),
	}

	var cre2_match [1]C.cre2_string_t
	result := C.cre2_partial_match_re(unsafe.Pointer(re.cre2_re_with_bracket), &cre2_input, &cre2_match[0], 1)
	return bool(result != 0)
}

// []
//func (re *Regexp) MatchReader(r io.RuneReader) bool {
//    // not implemented
//}

func (re *Regexp) MatchString(s string) bool {
	return re.Match([]byte(s))
}

func (re *Regexp) NumSubexp() int {
	groups := C.cre2_num_capturing_groups(unsafe.Pointer(re.cre2_re))
	return int(groups)
}

func (re *Regexp) RE2ReplaceAll(src, repl []byte) []byte {
	text_and_target := make([]byte, len(src))
	copy(text_and_target, src)

	c_text_and_target := (*C.char)(unsafe.Pointer(&text_and_target[0]))
	cre2_text_and_target := C.cre2_string_t{
		data:   c_text_and_target,
		length: C.int(len(text_and_target)),
	}

	re2_repl := replaceRE2Sequences(repl) // re2用のsequenceへ($n -> \\n)
	var cre2_re2_repl C.cre2_string_t
	if len(re2_repl) > 0 {
		cre2_re2_repl = C.cre2_string_t{
			data:   (*C.char)(unsafe.Pointer(&re2_repl[0])),
			length: C.int(len(re2_repl)),
		}
	} else {
		cre2_re2_repl = C.cre2_string_t{
			data:   nil,
			length: 0,
		}
	}

	C.cre2_global_replace_re(unsafe.Pointer(re.cre2_re), &cre2_text_and_target, &cre2_re2_repl)
	ret := C.GoBytes(unsafe.Pointer(cre2_text_and_target.data), cre2_text_and_target.length)

	C.free(unsafe.Pointer(cre2_text_and_target.data)) // cre2_text_and_targetはメモリが再確保されている
	// 古いバッファ(text_and_target)は[]byteなので解放はGoのGCが担う

	return ret
}

// [ ]
//func (re *Regexp) ReplaceAllFunc(src []byte, repl func([]byte) []byte) []byte {
//	  // not implemented
//}

func (re *Regexp) ReplaceAllLiteral(src, repl []byte) []byte {
	text_and_target := make([]byte, len(src))
	copy(text_and_target, src)

	c_text_and_target := (*C.char)(unsafe.Pointer(&text_and_target[0]))
	cre2_text_and_target := C.cre2_string_t{
		data:   c_text_and_target,
		length: C.int(len(text_and_target)),
	}

	re2_repl := replaceRE2InvalidSequences(repl) // sequencesの解析をしないように、\\を非sequencesへ置き換える
	var cre2_re2_repl C.cre2_string_t
	if len(re2_repl) > 0 {
		cre2_re2_repl = C.cre2_string_t{
			data:   (*C.char)(unsafe.Pointer(&re2_repl[0])),
			length: C.int(len(re2_repl)),
		}
	} else {
		cre2_re2_repl = C.cre2_string_t{
			data:   nil,
			length: 0,
		}
	}

	C.cre2_global_replace_re(unsafe.Pointer(re.cre2_re), &cre2_text_and_target, &cre2_re2_repl)
	ret := C.GoBytes(unsafe.Pointer(cre2_text_and_target.data), cre2_text_and_target.length)

	C.free(unsafe.Pointer(cre2_text_and_target.data)) // cre2_text_and_targetはメモリが再確保されている
	// 古いバッファ(text_and_target)は[]byteなので解放はGoのGCが担う

	return ret
}

func (re *Regexp) ReplaceAllLiteralString(src, repl string) string {
	return string(re.ReplaceAllLiteral([]byte(src), []byte(repl)))
}

// re2の置換方法に従う
// ({}で囲った記法が使えない)
func (re *Regexp) RE2ReplaceAllString(src, repl string) string {
	return string(re.RE2ReplaceAll([]byte(src), []byte(repl)))
}

// [ ]
//func (re *Regexp) ReplaceAllStringFunc(src string, repl func(string) string) string {
//    // not implemented
//}

func (re *Regexp) Split(s string, n int) []string {
	if n == 0 {
		return nil
	}

	if re.c_expr_length > 0 && len(s) == 0 {
		return []string{""}
	}

	matches := re.FindAllStringIndex(s, n)
	strings := make([]string, 0, len(matches))

	beg := 0
	end := 0
	for _, match := range matches {
		if n > 0 && len(strings) >= n-1 {
			break
		}

		end = match[0]
		if match[1] != 0 {
			strings = append(strings, s[beg:end])
		}
		beg = match[1]
	}

	if end != len(s) {
		strings = append(strings, s[beg:])
	}

	return strings
}

func (re *Regexp) String() string {
	return C.GoString(re.c_expr)
}

// [ ]
//func (re *Regexp) SubexpNames() []string {
//    // not implemented
//}
