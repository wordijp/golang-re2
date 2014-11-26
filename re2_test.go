package re2

import (
	"fmt"

	//	"strings"
	//	"regexp"
	"testing"
)

func TestCompile(t *testing.T) {

	// test wrong text
	var checkWrong = func(str string) {
		re, closer, err := Compile(str)
		defer closer.Close(re)
		if err == nil {
			detailErrorf(t, "text:%s is wrong. but err == nil", str)
		}
	}
	checkWrong("[a-z")
	checkWrong("+*")
	checkWrong("?")
	checkWrong("(")
	checkWrong("(()")

	// test correct string
	var checkCorrect = func(str string) {
		re, closer, err := Compile(str)
		defer closer.Close(re)
		if err != nil {
			detailError(t, err)
		}
	}
	checkCorrect("[a-z]?")
	checkCorrect(".+.*")
	checkCorrect(".?")
	checkCorrect("()")
	checkCorrect("(())")
}

func TestCompilePOSIX(t *testing.T) {
	// POSIXではposix_syntax, longest_matchをtrue
	var check = func(expr, input, answer string) {
		re, closer := MustCompilePOSIX(expr)
		defer closer.Close(re)
		equals_s(t, re.FindString(input), answer)
	}

	check(
		"a{1,4}?",
		"aaaaa",
		"aaaa",
	)
	check(
		"-.+?-",
		"-abc-eef-foo-hoge-",
		"-abc-eef-foo-hoge-",
	)
	check(
		".+[A-Z]{1,2}?",
		"aaBBBcccc",
		"aaBBBcccc",
	)
	check(
		".+[A-Z]{1,2}",
		"aaBBBcccc",
		"aaBBB",
	)
	check(
		"[a-z]+[A-Z]{1,2}?",
		"aaBBBcccc",
		"aaBB",
	)
	check(
		"[a-z]+[A-Z]{1,2}",
		"aaBBBcccc",
		"aaBB",
	)

}

func TestExpand(t *testing.T) {

	//	re, closer := MustCompile(`.*name\s+is\s+(?P<name>.+)\.`) // re2ではタグを使えない
	re, closer := MustCompile(`.*name\s+is\s+(.+)\.`)
	defer closer.Close(re)

	src := []byte(`
		my name is tom.
		my favorite food is sushi.
		hello, my name is bob.
		he name is hiroshi.
	`)
	answers := []string{
		"prefix name = tom",
		"prefix name = bob",
		"prefix name = hiroshi",
	}

	// Expand
	var ret []string
	for _, s := range re.FindAllSubmatchIndex(src, -1) {
		//		r := re.RE2Expand([]byte("prefix "), []byte("name = $name"), src, s) // re2ではタグを使えない
		r := re.RE2Expand([]byte("prefix "), []byte("name = $1"), src, s)
		ret = append(ret, string(r))
	}
	equals_as(t, ret, answers)

	// ExpandString
	ret = ret[:0]
	for _, s := range re.FindAllSubmatchIndex(src, -1) {
		//		r := re.RE2ExpandString([]byte("prefix "), "name = $1", string(src), s) // re2ではタグを使えない
		r := re.RE2ExpandString([]byte("prefix "), "name = $1", string(src), s)
		ret = append(ret, string(r))
	}
	equals_as(t, ret, answers)
}

func TestFind(t *testing.T) {
	re, closer := MustCompile(":([^: ]*)\\s*tom:")
	defer closer.Close(re)
	bytes := []byte("abc :super tom: abc :hyper bob: :tom: abc")
	//                   |     |    |                |    |
	//                   4    10   15               32   37

	{
		// Find
		ret := re.Find(bytes)
		if ret == nil {
			detailError(t, "Find is failed")
		}

		answer := ":super tom:"
		if string(ret) != answer {
			detailErrorf(t, "string: %s != %s", string(ret), answer)
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
		fmt.Println("* Find ReaderIndex not implemented *")
		//			r := strings.NewReader(string(bytes))
		//			ret := re.FindReaderIndex(r)
		//			answers := []int{
		//				4, 15,
		//			}
		//			equals_ai(t, ret, answers)
	}
	{
		// FindReaderSubmatchIndex
		fmt.Println("* FindReaderSubmatchIndex not implemented *")
		//		r := strings.NewReader(string(bytes))
		//		ret := re.FindReaderSubmatchIndex(r)
		//		answers := []int{
		//			4, 15, 5, 10,
		//		}
		//		equals_ai(t, ret, answers)
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
	fmt.Println("* LiteralPrefix not implemented *")
	//	re := MustCompile("prefix[a-z]?")
	//	prefix, _ := re.LiteralPrefix()
	//	equals_s(t, prefix, "prefix")
}

// 最長一致への切り替えテスト
// NOTE : POSIXへの切り替えはしない
func TestLongest(t *testing.T) {
	var check = func(expr, input, answer, answerLongested string) {
		re, closer := MustCompile(expr)
		defer closer.Close(re)
		if !equals_s(t, re.FindString(input), answer) {
			detailErrorParent(t, "wrong")
		}

		re.Longest()

		if !equals_s(t, re.FindString(input), answerLongested) {
			detailErrorParent(t, "wrong")
		}
	}

	check(
		"a{1,4}?",
		"aaaaa",
		"a",
		"aaaa",
	)
	check(
		"-.+?-",
		"-abc-eef-foo-hoge-",
		"-abc-",
		"-abc-eef-foo-hoge-",
	)
	check(
		"[a-z]+[A-Z]{1,2}?",
		"aaBBBcccc",
		"aaB",
		"aaBB",
	)
	check(
		"[a-z]+[A-Z]{1,2}",
		"aaBBBcccc",
		"aaBB",
		"aaBB",
	)
	check(
		".+[A-Z]{1,2}?",
		"aaBBBcccc",
		"aaBBB",
		"aaBBB",
	)
	check(
		".+[A-Z]{1,2}",
		"aaBBBcccc",
		"aaBBB",
		"aaBBB",
	)
}

func TestMatch(t *testing.T) {
	expr := ":([^: ]*)\\s*tom:"
	re, closer := MustCompile(expr)
	defer closer.Close(re)
	bytes := []byte("abc :super tom: abc :hyper bob: :tom: abc")

	var check = func(match bool) {
		if !match {
			detailErrorParent(t, "not match")
		}
	}

	{
		// Match (static method)
		matched, err := Match(expr, bytes)
		if err != nil {
			detailErrorf(t, "wrong. err:%s", err)
		}
		check(matched)
	}
	{
		// Match
		check(re.Match(bytes))
	}
	{
		fmt.Println("* MatchReader not implemented *")
		//		// MatchReader
		//		r := strings.NewReader(string(bytes))
		//		check(re.MatchReader(r))
	}
	{
		// MatchString
		check(re.MatchString(string(bytes)))
	}
}

func TestNumSubexp(t *testing.T) {
	var check = func(expr string, n int) {
		re, closer := MustCompile(expr)
		defer closer.Close(re)
		if re.NumSubexp() != n {
			detailErrorfParent(t, "num: %d != %d", re.NumSubexp(), n)
		}
	}

	check("[a-z]{2,4}", 0)
	check("([a-z])", 1)
	check("()", 1)
	check("()()", 2)
	check("(())", 2)
}

func TestReplace(t *testing.T) {
	re, closer := MustCompile("a(x*)b")
	defer closer.Close(re)
	{
		// ReplaceAll
		var check = func(str, repl, answer []byte) {
			ret := re.RE2ReplaceAll(str, repl)
			if !equals_ab(t, ret, answer) {
				detailErrorParent(t, "wrong")
			}
		}
		check([]byte("-ab-axxb-"), []byte("T"), []byte("-T-T-"))
		check([]byte("-ab-axxb-"), []byte("$1"), []byte("--xx-"))
		check([]byte("-ab-axxb-"), []byte("$1W"), []byte("-W-xxW-")) // regexpと結果が違う(regexpでは"---"となる)
		//		check([]byte("-ab-axxb-"), []byte("${1}W"), []byte("-W-xxW-")) // re2では{}で囲えない
	}
	{
		// ReplaceAllFunc
		fmt.Println("* ReplaceAllFunc not implemented *")
		//		ret := re.ReplaceAllFunc(bytes, func(dst []byte) []byte {
		//			return []byte(strings.ToLower(string(dst)))
		//		})
		//		answers := []byte("{t_name}:Tom. {t_age}:18.")
		//		equals_ab(t, ret, answers)
	}
	{
		// ReplaceAllLiteral
		ret := re.ReplaceAllLiteral([]byte("-ab-axxb-"), []byte("@$1@"))
		answers := []byte("-@$1@-@$1@-")
		equals_ab(t, ret, answers)

		ret = re.ReplaceAllLiteral([]byte("-ab-axxb-"), []byte("@\\1@"))
		answers = []byte("-@\\1@-@\\1@-")
		equals_ab(t, ret, answers)
	}
	{
		// ReplaceAllLiteralString
		ret := re.ReplaceAllLiteralString("-ab-axxb-", "@$1@")
		answers := "-@$1@-@$1@-"
		equals_s(t, ret, answers)

		ret = re.ReplaceAllLiteralString("-ab-axxb-", "@\\1@")
		answers = "-@\\1@-@\\1@-"
		equals_s(t, ret, answers)
	}

	{
		// ReplaceAllString
		var check = func(str, repl, answer string) {
			ret := re.RE2ReplaceAllString(str, repl)
			if !equals_s(t, ret, answer) {
				detailErrorParent(t, "wrong")
			}
		}
		check("-ab-axxb-", "T", "-T-T-")
		check("-ab-axxb-", "$1", "--xx-")
		//		check("-ab-axxb-", "$1W", "-W-xxW-") // regexpでは"---"となる
		//		check("-ab-axxb-", "${1}W", "-W-xxW-") // re2では{}で囲えない
	}
	{
		// ReplaceAllStringFunc
		fmt.Println("* ReplaceAllStringFunc not implemented *")
		//		ret := re.ReplaceAllStringFunc(string(bytes), func(dst string) string {
		//			return strings.ToLower(dst)
		//		})
		//		answers := "{t_name}:Tom. {t_age}:18."
		//		equals_s(t, ret, answers)
	}
}

// Split. String. SubexpNames.
func TestOther(t *testing.T) {
	{
		// Split
		var check = func(expr, str string, n int, answers []string) {
			re, closer := MustCompile(expr)
			defer closer.Close(re)
			//			re := regexp.MustCompile(expr)

			ret := re.Split(str, n)
			if !equals_as(t, ret, answers) {
				detailErrorParent(t, "wrong.")
			}
		}

		check(
			";",
			"abc;123;ABC;45;",
			5,
			[]string{
				"abc", "123", "ABC", "45", "",
			})
		check(
			"[0-9](;)",
			"abc;123;ABC;45;",
			2,
			[]string{
				"abc;12", "ABC;45;",
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
		check(
			"a*",
			"baabaccadaaae",
			5,
			[]string{
				"b", "b", "c", "c", "daaae",
			})
	}
	{
		// String
		str := ":([^: ]*)\\s*tom:"
		re, closer := MustCompile(str)
		defer closer.Close(re)
		equals_s(t, re.String(), str)
	}
	{
		// SubexpNames
		fmt.Println("* ReplaceAllStringFunc not implemented *")
		//		re := regexp.MustCompile("(?P<first>[a-zA-Z]+) (?P<last>[a-zA-Z]+)")
		//		ret := re.SubexpNames()
		//		answers := []string{
		//			"", "first", "last", // [0] is always the empty string.
		//		}
		//		equals_as(t, ret, answers)
	}
}
