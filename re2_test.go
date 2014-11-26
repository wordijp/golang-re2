package re2

import (
	"fmt"

	//	"strings"
	"regexp"
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

		_, prevErr := regexp.Compile(str)
		if prevErr == nil {
			detailErrorf(t, "text:%s is wrong. but err == nil", str)
		}

		if (err != nil) != (prevErr != nil) {
			detailErrorf(t, "wrong. err != oldErr")
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

		_, prevErr := regexp.Compile(str)
		if prevErr != nil {
			detailError(t, prevErr)
		}

		if (err != nil) != (prevErr != nil) {
			detailError(t, "wrong. err != prevErr")
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

		prevRe := regexp.MustCompilePOSIX(expr)

		ret := re.FindString(input)
		prevRet := prevRe.FindString(input)

		equals_s(t, ret, answer)
		equals_s(t, prevRet, answer)
		equals_s(t, ret, prevRet)
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

	//	expr := `.*name\s+is\s+(?P<name>.+)\.` // re2ではタグを使えない
	expr := `.*name\s+is\s+(.+)\.`
	re, closer := MustCompile(expr)
	defer closer.Close(re)

	prevRe := regexp.MustCompile(expr)

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

	{
		// Expand
		var check = func(src []byte, answers []string) {
			var ret []string
			for _, x := range re.FindAllSubmatchIndex(src, -1) {
				r := re.RE2Expand([]byte("prefix "), []byte("name = $1"), src, x)
				ret = append(ret, string(r))
			}

			var prevRet []string
			for _, x := range prevRe.FindAllSubmatchIndex(src, -1) {
				r := prevRe.Expand([]byte("prefix "), []byte("name = $1"), src, x)
				prevRet = append(prevRet, string(r))
			}

			equals_as(t, ret, answers)
			equals_as(t, prevRet, answers)
			equals_as(t, ret, prevRet)
		}
		check(src, answers)
	}
	{
		// ExpandString
		var check = func(src []byte, answers []string) {
			var ret []string
			for _, x := range re.FindAllSubmatchIndex(src, -1) {
				r := re.RE2ExpandString([]byte("prefix "), "name = $1", string(src), x)
				ret = append(ret, string(r))
			}

			var prevRet []string
			for _, x := range prevRe.FindAllSubmatchIndex(src, -1) {
				r := prevRe.ExpandString([]byte("prefix "), "name = $1", string(src), x)
				prevRet = append(prevRet, string(r))
			}

			equals_as(t, ret, answers)
			equals_as(t, prevRet, answers)
			equals_as(t, ret, prevRet)
		}
		check(src, answers)
	}
}

func TestFind(t *testing.T) {
	expr := ":([^: ]*)\\s*tom:"
	bytes := []byte("abc :super tom: abc :hyper bob: :tom: abc")
	//                   |     |    |                |    |
	//                   4    10   15               32   37

	re, closer := MustCompile(expr)
	defer closer.Close(re)

	prevRe := regexp.MustCompile(expr)

	{
		// Find
		var check = func(bytes []byte, answer string) {
			ret := re.Find(bytes)
			if ret == nil {
				detailError(t, "Find is failed")
			}
			prevRet := prevRe.Find(bytes)
			if prevRet == nil {
				detailError(t, "Find is failed")
			}

			equals_s(t, string(ret), answer)
			equals_s(t, string(prevRet), answer)
			equals_s(t, string(ret), string(prevRet))
		}
		check(bytes, ":super tom:")
	}
	{
		// FindAll
		var check = func(bytes []byte, answers []string) {
			var ret []string
			for _, x := range re.FindAll(bytes, -1) {
				ret = append(ret, string(x))
			}

			var prevRet []string
			for _, x := range prevRe.FindAll(bytes, -1) {
				prevRet = append(prevRet, string(x))
			}

			equals_as(t, ret, answers)
			equals_as(t, prevRet, answers)
			equals_as(t, ret, prevRet)
		}
		answers := []string{
			":super tom:",
			":tom:",
		}
		check(bytes, answers)

	}
	{
		// FindAllIndex
		var check = func(bytes []byte, answers [][]int) {
			ret := re.FindAllIndex(bytes, -1)
			prevRet := prevRe.FindAllIndex(bytes, -1)

			equals_aai(t, ret, answers)
			equals_aai(t, prevRet, answers)
			equals_aai(t, ret, prevRet)
		}
		answers := [][]int{
			{4, 15},  // :super tom:
			{32, 37}, // :tom:
		}
		check(bytes, answers)
	}
	{
		// FindAllString
		var check = func(bytes []byte, answers []string) {
			ret := re.FindAllString(string(bytes), -1)
			prevRet := prevRe.FindAllString(string(bytes), -1)

			equals_as(t, ret, answers)
			equals_as(t, prevRet, answers)
			equals_as(t, ret, prevRet)
		}
		answers := []string{
			":super tom:",
			":tom:",
		}
		check(bytes, answers)
	}
	{
		// FindAllStringIndex
		var check = func(bytes []byte, answers [][]int) {
			ret := re.FindAllStringIndex(string(bytes), -1)
			prevRet := prevRe.FindAllStringIndex(string(bytes), -1)

			equals_aai(t, ret, answers)
			equals_aai(t, prevRet, answers)
			equals_aai(t, ret, prevRet)
		}
		answers := [][]int{
			{4, 15},  // :super tom:
			{32, 37}, // :tom:
		}
		check(bytes, answers)
	}
	{
		// FindAllStringSubmatch
		var check = func(bytes []byte, answers [][]string) {
			ret := re.FindAllStringSubmatch(string(bytes), -1)
			prevRet := prevRe.FindAllStringSubmatch(string(bytes), -1)

			equals_aas(t, ret, answers)
			equals_aas(t, prevRet, answers)
			equals_aas(t, ret, prevRet)
		}
		answers := [][]string{
			{":super tom:", "super"},
			{":tom:", ""},
		}
		check(bytes, answers)
	}
	{
		// FindAllStringSubmatchIndex
		var check = func(bytes []byte, answers [][]int) {
			ret := re.FindAllStringSubmatchIndex(string(bytes), -1)
			prevRet := prevRe.FindAllStringSubmatchIndex(string(bytes), -1)

			equals_aai(t, ret, answers)
			equals_aai(t, prevRet, answers)
			equals_aai(t, ret, prevRet)
		}
		answers := [][]int{
			{4, 15, 5, 10},   // {":super tom:", "super"},
			{32, 37, 33, 33}, // {":tom:", ""},
		}
		check(bytes, answers)
	}
	{
		// FindAllSubmatch
		var check = func(bytes []byte, answers [][][]byte) {
			ret := re.FindAllSubmatch(bytes, -1)
			prevRet := prevRe.FindAllSubmatch(bytes, -1)

			equals_aaab(t, ret, answers)
			equals_aaab(t, prevRet, answers)
			equals_aaab(t, ret, prevRet)
		}
		answers := [][][]byte{
			{[]byte(":super tom:"), []byte("super")},
			{[]byte(":tom:"), []byte("")},
		}
		check(bytes, answers)
	}
	{
		// FindAllSubmatchIndex
		var check = func(bytes []byte, answers [][]int) {
			ret := re.FindAllSubmatchIndex(bytes, -1)
			prevRet := prevRe.FindAllSubmatchIndex(bytes, -1)

			equals_aai(t, ret, answers)
			equals_aai(t, prevRet, answers)
			equals_aai(t, ret, prevRet)
		}
		answers := [][]int{
			{4, 15, 5, 10},   // {":super tom:", "super"},
			{32, 37, 33, 33}, // {":tom:", ""},
		}
		check(bytes, answers)
	}
	{
		// FindIndex
		var check = func(bytes []byte, answers []int) {
			ret := re.FindIndex(bytes)
			prevRet := prevRe.FindIndex(bytes)

			equals_ai(t, ret, answers)
			equals_ai(t, prevRet, answers)
			equals_ai(t, ret, prevRet)
		}
		answers := []int{
			4, 15,
		}
		check(bytes, answers)
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
		var check = func(bytes []byte, answers string) {
			ret := re.FindString(string(bytes))
			prevRet := prevRe.FindString(string(bytes))

			equals_s(t, ret, answers)
			equals_s(t, prevRet, answers)
			equals_s(t, ret, prevRet)
		}
		answers := ":super tom:"
		check(bytes, answers)
	}
	{
		// FindStringIndex
		var check = func(bytes []byte, answers []int) {
			ret := re.FindStringIndex(string(bytes))
			prevRet := prevRe.FindStringIndex(string(bytes))

			equals_ai(t, ret, answers)
			equals_ai(t, prevRet, answers)
			equals_ai(t, ret, prevRet)
		}
		answers := []int{
			4, 15,
		}
		check(bytes, answers)
	}
	{
		// FindStringSubmatch
		var check = func(bytes []byte, answers []string) {
			ret := re.FindStringSubmatch(string(bytes))
			prevRet := prevRe.FindStringSubmatch(string(bytes))

			equals_as(t, ret, answers)
			equals_as(t, prevRet, answers)
			equals_as(t, ret, prevRet)
		}
		answers := []string{
			":super tom:", "super",
		}
		check(bytes, answers)
	}
	{
		// FindStringSubmatchIndex
		var check = func(bytes []byte, answers []int) {
			ret := re.FindStringSubmatchIndex(string(bytes))
			prevRet := prevRe.FindStringSubmatchIndex(string(bytes))

			equals_ai(t, ret, answers)
			equals_ai(t, prevRet, answers)
			equals_ai(t, ret, prevRet)
		}
		answers := []int{
			4, 15, 5, 10,
		}
		check(bytes, answers)
	}
	{
		// FindSubmatch
		var check = func(bytes []byte, answers [][]byte) {
			ret := re.FindSubmatch(bytes)
			prevRet := prevRe.FindSubmatch(bytes)

			equals_aab(t, ret, answers)
			equals_aab(t, prevRet, answers)
			equals_aab(t, ret, prevRet)
		}
		answers := [][]byte{
			[]byte(":super tom:"), []byte("super"),
		}
		check(bytes, answers)
	}
	{
		// FindSubmatchIndex
		var check = func(bytes []byte, answers []int) {
			ret := re.FindSubmatchIndex(bytes)
			prevRet := prevRe.FindSubmatchIndex(bytes)

			equals_ai(t, ret, answers)
			equals_ai(t, prevRet, answers)
			equals_ai(t, ret, prevRet)
		}
		answers := []int{
			4, 15, 5, 10,
		}
		check(bytes, answers)
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

		prevRe := regexp.MustCompile(expr)

		{
			ret := re.FindString(input)
			prevRet := prevRe.FindString(input)

			if !equals_s(t, ret, answer) {
				detailErrorParent(t, "wrong")
			}
			if !equals_s(t, prevRet, answer) {
				detailErrorParent(t, "wrong")
			}
			if !equals_s(t, ret, prevRet) {
				detailErrorParent(t, "wrong")
			}
		}
		re.Longest()
		prevRe.Longest()
		{
			ret := re.FindString(input)
			prevRet := prevRe.FindString(input)

			if !equals_s(t, ret, answerLongested) {
				detailErrorParent(t, "wrong")
			}
			if !equals_s(t, prevRet, answerLongested) {
				detailErrorParent(t, "wrong")
			}
			if !equals_s(t, ret, prevRet) {
				detailErrorParent(t, "wrong")
			}
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
	bytes := []byte("abc :super tom: abc :hyper bob: :tom: abc")

	re, closer := MustCompile(expr)
	defer closer.Close(re)

	prevRe := regexp.MustCompile(expr)

	{
		// Match (static method)
		var check = func(expr string, bytes []byte) {
			ret, err := Match(expr, bytes)
			if err != nil {
				detailErrorfParent(t, "wrong. err:%s", err)
			}

			prevRet, prevErr := regexp.Match(expr, bytes)
			if prevErr != nil {
				detailErrorfParent(t, "wrong, err:%s", prevErr)
			}

			if (err != nil) != (prevErr != nil) {
				detailErrorParent(t, "wrong")
			}

			if !ret {
				detailErrorParent(t, "not match")
			}
			if !prevRet {
				detailErrorParent(t, "not match")
			}
			if ret != prevRet {
				detailErrorParent(t, "wrong")
			}
		}
		check(expr, bytes)
	}
	{
		// MatchString (static method)
		var check = func(expr string, s string) {
			ret, err := MatchString(expr, s)
			if err != nil {
				detailErrorfParent(t, "wrong. err:%s", err)
			}

			prevRet, prevErr := regexp.MatchString(expr, s)
			if prevErr != nil {
				detailErrorfParent(t, "wrong, err:%s", prevErr)
			}

			if (err != nil) != (prevErr != nil) {
				detailErrorParent(t, "wrong")
			}

			if !ret {
				detailErrorParent(t, "not match")
			}
			if !prevRet {
				detailErrorParent(t, "not match")
			}
			if ret != prevRet {
				detailErrorParent(t, "wrong")
			}
		}
		check(expr, string(bytes))
	}

	{
		// Match
		var check = func(bytes []byte) {
			ret := re.Match(bytes)
			prevRet := prevRe.Match(bytes)

			if !ret {
				detailErrorParent(t, "not match")
			}
			if !prevRet {
				detailErrorParent(t, "not match")
			}
			if ret != prevRet {
				detailErrorParent(t, "wrong")
			}
		}
		check(bytes)
	}
	{
		fmt.Println("* MatchReader not implemented *")
		//		// MatchReader
		//		r := strings.NewReader(string(bytes))
		//		check(re.MatchReader(r))
	}
	{
		// MatchString
		var check = func(s string) {
			ret := re.MatchString(s)
			prevRet := prevRe.MatchString(s)

			if !ret {
				detailErrorParent(t, "not match")
			}
			if !prevRet {
				detailErrorParent(t, "not match")
			}
			if ret != prevRet {
				detailErrorParent(t, "wrong")
			}
		}
		check(string(bytes))
	}
}

func TestNumSubexp(t *testing.T) {
	var check = func(expr string, n int) {
		re, closer := MustCompile(expr)
		defer closer.Close(re)

		prevRe := regexp.MustCompile(expr)

		ret := re.NumSubexp()
		prevRet := prevRe.NumSubexp()

		if ret != n {
			detailErrorfParent(t, "num: %d != %d", re.NumSubexp(), n)
		}
		if prevRet != n {
			detailErrorfParent(t, "num: %d != %d", re.NumSubexp(), n)
		}
		if ret != prevRet {
			detailErrorParent(t, "wrong")
		}
	}

	check("[a-z]{2,4}", 0)
	check("([a-z])", 1)
	check("()", 1)
	check("()()", 2)
	check("(())", 2)
}

func TestReplace(t *testing.T) {
	expr := "a(x*)b"
	re, closer := MustCompile(expr)
	defer closer.Close(re)

	prevRe := regexp.MustCompile(expr)
	{
		// ReplaceAll
		var check = func(bytes, repl, answer []byte) {
			ret := re.RE2ReplaceAll(bytes, repl)
			prevRet := prevRe.ReplaceAll(bytes, repl)

			if !equals_ab(t, ret, answer) {
				detailErrorParent(t, "wrong")
			}
			if !equals_ab(t, prevRet, answer) {
				detailErrorParent(t, "wrong")
			}
			if !equals_ab(t, ret, prevRet) {
				detailErrorParent(t, "wrong")
			}
		}
		check([]byte("-ab-axxb-"), []byte("T"), []byte("-T-T-"))
		check([]byte("-ab-axxb-"), []byte("$1"), []byte("--xx-"))
		//		check([]byte("-ab-axxb-"), []byte("$1W"), []byte("-W-xxW-")) // regexpでは"---"となる
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
		var check = func(bytes, repl, answer []byte) {
			ret := re.ReplaceAllLiteral(bytes, repl)
			prevRet := prevRe.ReplaceAllLiteral(bytes, repl)

			if !equals_ab(t, ret, answer) {
				detailErrorParent(t, "wrong")
			}
			if !equals_ab(t, prevRet, answer) {
				detailErrorParent(t, "wrong")
			}
			if !equals_ab(t, ret, prevRet) {
				detailErrorParent(t, "wrong")
			}
		}
		check([]byte("-ab-axxb-"), []byte("@$1@"), []byte("-@$1@-@$1@-"))
		check([]byte("-ab-axxb-"), []byte("@\\1@"), []byte("-@\\1@-@\\1@-"))
		check([]byte("-ab-axxb-"), []byte("$1W"), []byte("-$1W-$1W-"))
		check([]byte("-ab-axxb-"), []byte("${1}W"), []byte("-${1}W-${1}W-"))
	}
	{
		// ReplaceAllLiteralString
		var check = func(str, repl, answer string) {
			ret := re.ReplaceAllLiteralString(str, repl)
			prevRet := prevRe.ReplaceAllLiteralString(str, repl)

			if !equals_s(t, ret, answer) {
				detailErrorParent(t, "wrong")
			}
			if !equals_s(t, prevRet, answer) {
				detailErrorParent(t, "wrong")
			}
			if !equals_s(t, ret, prevRet) {
				detailErrorParent(t, "wrong")
			}
		}
		check("-ab-axxb-", "@$1@", "-@$1@-@$1@-")
		check("-ab-axxb-", "@\\1@", "-@\\1@-@\\1@-")
		check("-ab-axxb-", "$1W", "-$1W-$1W-")
		check("-ab-axxb-", "${1}W", "-${1}W-${1}W-")
	}
	{
		// ReplaceAllString
		var check = func(str, repl, answer string) {
			ret := re.RE2ReplaceAllString(str, repl)
			prevRet := prevRe.ReplaceAllString(str, repl)

			if !equals_s(t, ret, answer) {
				detailErrorParent(t, "wrong")
			}
			if !equals_s(t, prevRet, answer) {
				detailErrorParent(t, "wrong")
			}
			if !equals_s(t, ret, prevRet) {
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

// Split. String. SubexpNames. QuoteMeta
func TestOther(t *testing.T) {
	{
		// Split
		var check = func(expr, str string, n int, answers []string) {
			re, closer := MustCompile(expr)
			defer closer.Close(re)
			prevRe := regexp.MustCompile(expr)

			ret := re.Split(str, n)
			prevRet := prevRe.Split(str, n)

			if !equals_as(t, ret, answers) {
				detailErrorParent(t, "wrong")
			}
			if !equals_as(t, prevRet, answers) {
				detailErrorParent(t, "wrong")
			}
			if !equals_as(t, ret, prevRet) {
				detailErrorParent(t, "wrong")
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
			";",
			"abc;123;ABC;45;",
			-1,
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
			"[0-9](;)",
			"abc;123;ABC;45;",
			-1,
			[]string{
				"abc;12", "ABC;4", "",
			})
		check(
			";",
			"abc;;123;;ABC;;45;;",
			5,
			[]string{
				"abc", "", "123", "", "ABC;;45;;",
			})
		check(
			";",
			"abc;;123;;ABC;;45;;",
			-1,
			[]string{
				"abc", "", "123", "", "ABC", "", "45", "", "",
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
			"abaabaccadaaae",
			-1,
			[]string{
				"", "b", "b", "c", "c", "d", "e",
			})
		check(
			"a*",
			"baabaccadaaae",
			5,
			[]string{
				"b", "b", "c", "c", "daaae",
			})
		check(
			"a*",
			"baabaccadaaae",
			-1,
			[]string{
				"b", "b", "c", "c", "d", "e",
			})

	}
	{
		// String
		var check = func(expr string) {
			re, closer := MustCompile(expr)
			defer closer.Close(re)

			prevRe := regexp.MustCompile(expr)

			ret := re.String()
			prevRet := prevRe.String()

			equals_s(t, ret, expr)
			equals_s(t, prevRet, expr)
			equals_s(t, ret, prevRet)
		}
		check(":([^: ]*)\\s*tom:")
	}
	{
		// SubexpNames
		fmt.Println("* SubexpNames not implemented *")
		//		re := regexp.MustCompile("(?P<first>[a-zA-Z]+) (?P<last>[a-zA-Z]+)")
		//		ret := re.SubexpNames()
		//		answers := []string{
		//			"", "first", "last", // [0] is always the empty string.
		//		}
		//		equals_as(t, ret, answers)
	}
	{
		// QuoteMeta
		var check = func(expr, answer string) {
			ret := RE2QuoteMeta(expr)
			prevRet := regexp.QuoteMeta(expr)

			if !equals_s(t, ret, answer) {
				detailErrorParent(t, "wrong")
			}
			if !equals_s(t, prevRet, answer) {
				detailErrorParent(t, "wrong")
			}
			if !equals_s(t, ret, prevRet) {
				detailErrorParent(t, "wrong")
			}
		}
		check("[foo]", `\[foo\]`)
		check("a*", `a\*`)
		check("()", `\(\)`)
		check("?", `\?`)
		//		check("[[:lower:]]", `\[\[:lower:\]\]`) // re2では`\[\[\:lower\:\]\]`となる
	}
}
