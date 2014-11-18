package re2

import (
	"./testutil"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func equals_i(t *testing.T, a, b int) bool {
	if a != b {
		testutil.DetailErrorfParent(t, "int: %d != %d", a, b)
		return false
	}

	return true
}

func equals_ai(t *testing.T, a, b []int) bool {
	if len(a) != len(b) {
		testutil.DetailErrorfParent(t, "length: %d != %d", len(a), len(b))
		return false
	}

	n := len(a)
	for i := 0; i < n; i++ {
		if !equals_i(t, a[i], b[i]) {
			testutil.DetailErrorfParent(t, "wrong. []int at(%d)", i)
			return false
		}
	}

	return true
}

func equals_aai(t *testing.T, a, b [][]int) bool {
	if len(a) != len(b) {
		testutil.DetailErrorfParent(t, "length: %d != %d", len(a), len(b))
		return false
	}

	n := len(a)
	for i := 0; i < n; i++ {
		if !equals_ai(t, a[i], b[i]) {
			testutil.DetailErrorfParent(t, "wrong. [][]int at(%d)", i)
			return false
		}
	}

	return true
}

func equals_s(t *testing.T, a, b string) bool {
	if a != b {
		testutil.DetailErrorfParent(t, "string: %s != %s", a, b)
		return false
	}

	return true
}

func equals_as(t *testing.T, a, b []string) bool {
	if len(a) != len(b) {
		testutil.DetailErrorfParent(t, "length: %d != %d", len(a), len(b))
		return false
	}

	n := len(a)
	for i := 0; i < n; i++ {
		if !equals_s(t, a[i], b[i]) {
			testutil.DetailErrorfParent(t, "wrong. []string at(%d)", i)
			return false
		}
	}

	return true
}

func equals_aas(t *testing.T, a, b [][]string) bool {
	if len(a) != len(b) {
		testutil.DetailErrorfParent(t, "length: %d != %d", len(a), len(b))
		return false
	}

	n := len(a)
	for i := 0; i < n; i++ {
		if !equals_as(t, a[i], b[i]) {
			testutil.DetailErrorfParent(t, "wrong. [][]string at(%d)", i)
			return false
		}
	}

	return true
}

func equals_ab(t *testing.T, a, b []byte) bool {
	if !reflect.DeepEqual(a, b) {
		testutil.DetailErrorfParent(t, "[]byte: %s != %s", fmt.Sprint(a), fmt.Sprint(b))
		testutil.DetailErrorfParent(t, "[]byte(to string): %s != %s", string(a), string(b))
		return false
	}

	return true
}

func equals_aab(t *testing.T, a, b [][]byte) bool {
	if len(a) != len(b) {
		testutil.DetailErrorfParent(t, "length: %d != %d", len(a), len(b))
		return false
	}

	n := len(a)
	for i := 0; i < n; i++ {
		if !equals_ab(t, a[i], b[i]) {
			testutil.DetailErrorfParent(t, "wrong. [][]byte at(%d)", i)
			return false
		}
	}

	return true
}

func equals_aaab(t *testing.T, a, b [][][]byte) bool {
	if len(a) != len(b) {
		testutil.DetailErrorfParent(t, "length: %d != %d", len(a), len(b))
		return false
	}

	n := len(a)
	for i := 0; i < n; i++ {
		if !equals_aab(t, a[i], b[i]) {
			testutil.DetailErrorfParent(t, "wrong. [][][]byte at(%d)", i)
			return false
		}
	}

	return true
}

func TestCompile(t *testing.T) {

	// test wrong text
	var checkWrong = func(str string) {
		_, err := Compile(str)
		if err == nil {
			testutil.DetailErrorf(t, "text:%s is wrong. but err == nil", str)
		}
	}
	checkWrong("[a-z")
	checkWrong("+*")
	checkWrong("?")
	checkWrong("(")
	checkWrong("(()")

	// test correct string
	var checkCorrect = func(str string) {
		_, err := Compile(str)
		if err != nil {
			testutil.DetailError(t, err)
		}
	}
	checkCorrect("[a-z]")
	checkCorrect(".+.*")
	checkCorrect(".?")
	checkCorrect("()")
	checkCorrect("(())")
}

func TestExpand(t *testing.T) {

	re := MustCompile(`.*name\s+is\s+(?P<name>.+)\.`)

	src := []byte(`
		my name is tom.
		my favorite food is sushi.
		hello, my name is bob.
		he name is hiroshi.
	`)
	answers := []string{
		"name = tom",
		"name = bob",
		"name = hiroshi",
	}

	// Expand
	var ret []string
	for _, s := range re.FindAllSubmatchIndex(src, -1) {
		r := re.Expand([]byte(""), []byte("name = $name"), src, s)
		ret = append(ret, string(r))
	}
	equals_as(t, ret, answers)

	// ExpandString
	ret = ret[:0]
	for _, s := range re.FindAllSubmatchIndex(src, -1) {
		r := re.ExpandString([]byte(""), "name = $name", string(src), s)
		ret = append(ret, string(r))
	}
	equals_as(t, ret, answers)
}

func TestFind(t *testing.T) {
	re := MustCompile(":([^: ]*)\\s*tom:")
	bytes := []byte("abc :super tom: abc :hyper bob: :tom: abc")
	//                   |     |    |                |    |
	//                   4    10   15               32   37

	{
		// Find
		ret := re.Find(bytes)
		if ret == nil {
			testutil.DetailError(t, "Find is failed")
		}

		answer := ":super tom:"
		if string(ret) != answer {
			testutil.DetailErrorf(t, "string: %s != %s", string(ret), answer)
		}
	}
	{
		// FindAll
		var ret []string
		for _, x := range re.FindAll(bytes, -1) {
			ret = append(ret, string(x))
		}

		answers := []string{
			":super tom:",
			":tom:",
		}
		equals_as(t, ret, answers)
	}
	{
		// FindAllIndex
		ret := re.FindAllIndex(bytes, -1)
		answers := [][]int{
			{4, 15},  // :super tom:
			{32, 37}, // :tom:
		}

		equals_aai(t, ret, answers)
	}
	{
		// FindAllString
		ret := re.FindAllString(string(bytes), -1)
		answers := []string{
			":super tom:",
			":tom:",
		}

		equals_as(t, ret, answers)
	}
	{
		// FindAllStringIndex
		ret := re.FindAllStringIndex(string(bytes), -1)
		answers := [][]int{
			{4, 15},  // :super tom:
			{32, 37}, // :tom:
		}

		equals_aai(t, ret, answers)
	}
	{
		// FindAllStringSubmatch
		ret := re.FindAllStringSubmatch(string(bytes), -1)
		answers := [][]string{
			{":super tom:", "super"},
			{":tom:", ""},
		}
		equals_aas(t, ret, answers)
	}
	{
		// FindAllStringSubmatchIndex
		ret := re.FindAllStringSubmatchIndex(string(bytes), -1)
		answers := [][]int{
			{4, 15, 5, 10},   // {":super tom:", "super"},
			{32, 37, 33, 33}, // {":tom:", ""},
		}
		equals_aai(t, ret, answers)
	}
	{
		// FindAllSubmatch
		ret := re.FindAllSubmatch(bytes, -1)
		answers := [][][]byte{
			{[]byte(":super tom:"), []byte("super")},
			{[]byte(":tom:"), []byte("")},
		}
		equals_aaab(t, ret, answers)
	}
	{
		// FindAllSubmatchIndex
		ret := re.FindAllSubmatchIndex(bytes, -1)
		answers := [][]int{
			{4, 15, 5, 10},   // {":super tom:", "super"},
			{32, 37, 33, 33}, // {":tom:", ""},
		}
		equals_aai(t, ret, answers)
	}
	{
		// FindIndex
		ret := re.FindIndex(bytes)
		answers := []int{
			4, 15,
		}
		equals_ai(t, ret, answers)
	}
	{
		// FindReaderIndex
		r := strings.NewReader(string(bytes))
		ret := re.FindReaderIndex(r)
		answers := []int{
			4, 15,
		}
		equals_ai(t, ret, answers)
	}
	{
		// FindReaderSubmatchIndex
		r := strings.NewReader(string(bytes))
		ret := re.FindReaderSubmatchIndex(r)
		answers := []int{
			4, 15, 5, 10,
		}
		equals_ai(t, ret, answers)
	}
	{
		// FindString
		ret := re.FindString(string(bytes))
		answers := ":super tom:"
		equals_s(t, ret, answers)
	}
	{
		// FindStringIndex
		ret := re.FindStringIndex(string(bytes))
		answers := []int{
			4, 15,
		}
		equals_ai(t, ret, answers)
	}
	{
		// FindStringSubmatch
		ret := re.FindStringSubmatch(string(bytes))
		answers := []string{
			":super tom:", "super",
		}
		equals_as(t, ret, answers)
	}
	{
		// FindStringSubmatchIndex
		ret := re.FindStringSubmatchIndex(string(bytes))
		answers := []int{
			4, 15, 5, 10,
		}
		equals_ai(t, ret, answers)
	}
	{
		// FindSubmatch
		ret := re.FindSubmatch(bytes)
		answers := [][]byte{
			[]byte(":super tom:"), []byte("super"),
		}
		equals_aab(t, ret, answers)
	}
	{
		// FindSubmatchIndex
		ret := re.FindSubmatchIndex(bytes)
		answers := []int{
			4, 15, 5, 10,
		}
		equals_ai(t, ret, answers)
	}
}

func TestLiteralPrefix(t *testing.T) {
	re := MustCompile("prefix[a-z]?")
	prefix, _ := re.LiteralPrefix()
	equals_s(t, prefix, "prefix")
}

// POSIX(最長一致)への切り替えテスト
func TestLongest(t *testing.T) {
	{
		re := MustCompile("a{1,4}?")
		str := "aaaaa"

		equals_s(t, re.FindString(str), "a")
		re.Longest()
		equals_s(t, re.FindString(str), "aaaa")
	}
	{
		re := MustCompile("-.+?-")
		str := "-abc-eef-foo-hoge-"

		equals_s(t, re.FindString(str), "-abc-")
		re.Longest()
		equals_s(t, re.FindString(str), "-abc-eef-foo-hoge-")
	}
}

func TestMatch(t *testing.T) {
	re := MustCompile(":([^: ]*)\\s*tom:")
	bytes := []byte("abc :super tom: abc :hyper bob: :tom: abc")

	var check = func(match bool) {
		if !match {
			testutil.DetailErrorParent(t, "not match")
		}
	}

	{
		// Match
		check(re.Match(bytes))
	}
	{
		// MatchReader
		r := strings.NewReader(string(bytes))
		check(re.MatchReader(r))
	}
	{
		// MatchString
		check(re.MatchString(string(bytes)))
	}

}

func TestNumSubexp(t *testing.T) {

	var check = func(expr string, n int) {
		re := MustCompile(expr)
		if re.NumSubexp() != n {
			testutil.DetailErrorfParent(t, "num: %d != %d", re.NumSubexp(), n)
		}
	}

	check("[a-z]{2,4}", 0)
	check("([a-z])", 1)
	check("()", 1)
	check("()()", 2)
	check("(())", 2)
}

func TestReplace(t *testing.T) {
	re := MustCompile("{T_([^}]+)}")
	bytes := []byte("{T_NAME}:Tom. {T_AGE}:18.")
	{
		// ReplaceAll
		ret := re.ReplaceAll(bytes, []byte("@$1@"))
		answers := []byte("@NAME@:Tom. @AGE@:18.")
		equals_ab(t, ret, answers)
	}
	{
		// ReplaceAllFunc
		ret := re.ReplaceAllFunc(bytes, func(dst []byte) []byte {
			return []byte(strings.ToLower(string(dst)))
		})
		answers := []byte("{t_name}:Tom. {t_age}:18.")
		equals_ab(t, ret, answers)
	}
	{
		// ReplaceAllLiteral
		ret := re.ReplaceAllLiteral(bytes, []byte("@$1@"))
		answers := []byte("@$1@:Tom. @$1@:18.")
		equals_ab(t, ret, answers)

	}
	{
		// ReplaceAllLiteralString
		ret := re.ReplaceAllLiteralString(string(bytes), "@$1@")
		answers := "@$1@:Tom. @$1@:18."
		equals_s(t, ret, answers)
	}
	{
		// ReplaceAllString
		ret := re.ReplaceAllString(string(bytes), "@$1@")
		answers := "@NAME@:Tom. @AGE@:18."
		equals_s(t, ret, answers)
	}
	{
		// ReplaceAllStringFunc
		ret := re.ReplaceAllStringFunc(string(bytes), func(dst string) string {
			return strings.ToLower(dst)
		})
		answers := "{t_name}:Tom. {t_age}:18."
		equals_s(t, ret, answers)
	}
}

// Split. String. SubexpNames.
func TestOther(t *testing.T) {
	{
		// Split
		var check = func(expr, str string, n int, answers []string) {
			re := MustCompile(expr)
			ret := re.Split(str, n)
			equals_as(t, ret, answers)
		}

		check(
			";",
			"abc;123;ABC;45;",
			5,
			[]string{
				"abc", "123", "ABC", "45", "",
			})

		check(
			";",
			"abc;;123;;ABC;;45;;",
			5,
			[]string{
				"abc", "", "123", "", "ABC;;45;;",
			})

		check(
			"a*",
			"abaabaccadaaae",
			5,
			[]string{
				"", "b", "b", "c", "cadaaae",
			})

	}
	{
		// String
		str := ":([^: ]*)\\s*tom:"
		re := MustCompile(str)
		equals_s(t, re.String(), str)
	}
	{
		// SubexpNames
		re := MustCompile("(?P<first>[a-zA-Z]+) (?P<last>[a-zA-Z]+)")
		ret := re.SubexpNames()
		answers := []string{
			"", "first", "last", // [0] is always the empty string.
		}
		equals_as(t, ret, answers)
	}
}
