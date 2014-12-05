package re2

import (
	"fmt"

	//	"strings"
	"regexp"
	"testing"
)

// XXX : テストの共通処理を個別に書いてる為、抜けの発生が起こる
// TODO : テストの共通処理をまとめる

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

	bytes := []byte(`
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
		var check = func(expr string, bytes []byte, answers []string) {
			re, closer := MustCompile(expr)
			defer closer.Close(re)

			prevRe := regexp.MustCompile(expr)

			var ret []string
			for _, x := range re.FindAllSubmatchIndex(bytes, -1) {
				r := re.RE2Expand([]byte("prefix "), []byte("name = $1"), bytes, x)
				ret = append(ret, string(r))
			}

			var prevRet []string
			for _, x := range prevRe.FindAllSubmatchIndex(bytes, -1) {
				r := prevRe.Expand([]byte("prefix "), []byte("name = $1"), bytes, x)
				prevRet = append(prevRet, string(r))
			}

			if !equals_as(t, ret, answers) {
				detailErrorParent(t, "wrong.")
			}
			if !equals_as(t, prevRet, answers) {
				detailErrorParent(t, "wrong.")
			}
			if !equals_as(t, ret, prevRet) {
				detailErrorParent(t, "wrong.")
			}
		}
		check(expr, bytes, answers)
		check("(a*)", []byte("abc"),
			[]string{
				"prefix name = a",
				"prefix name = ",
				"prefix name = ",
			},
		)
		check("a*", []byte("bc"),
			[]string{
				"prefix name = ",
				"prefix name = ",
				"prefix name = ",
			},
		)
		check("a*", []byte(""),
			[]string{
				"prefix name = ",
			},
		)
		check("a", []byte(""),
			[]string{},
		)
	}
	{
		// ExpandString
		var check = func(expr string, s string, answers []string) {
			re, closer := MustCompile(expr)
			defer closer.Close(re)

			prevRe := regexp.MustCompile(expr)

			var ret []string
			for _, x := range re.FindAllSubmatchIndex([]byte(s), -1) {
				r := re.RE2ExpandString([]byte("prefix "), "name = $1", s, x)
				ret = append(ret, string(r))
			}

			var prevRet []string
			for _, x := range prevRe.FindAllSubmatchIndex([]byte(s), -1) {
				r := prevRe.ExpandString([]byte("prefix "), "name = $1", s, x)
				prevRet = append(prevRet, string(r))
			}

			equals_as(t, ret, answers)
			equals_as(t, prevRet, answers)
			equals_as(t, ret, prevRet)
		}
		check(expr, string(bytes), answers)
		check("(a*)", "abc",
			[]string{
				"prefix name = a",
				"prefix name = ",
				"prefix name = ",
			},
		)
		check("a*", "bc",
			[]string{
				"prefix name = ",
				"prefix name = ",
				"prefix name = ",
			},
		)
		check("a*", "",
			[]string{
				"prefix name = ",
			},
		)
		check("a", "",
			[]string{},
		)

	}
}

func TestFind(t *testing.T) {

	// XXX : checkメソッドの引数名と被っているため、typo時に不思議なバグになってしまう
	expr := ":([^: ]*)\\s*tom:"
	bytes := []byte("abc :super tom: abc :hyper bob: :tom: abc")
	//                   |     |    |                |    |
	//                   4    10   15               32   37

	{
		// Find
		var check = func(expr string, bytes []byte, answer string) {
			re, closer := MustCompile(expr)
			defer closer.Close(re)

			prevRe := regexp.MustCompile(expr)

			ret := re.Find(bytes)
			prevRet := prevRe.Find(bytes)

			if (ret != nil) != (prevRet != nil) {
				detailError(t, "wrong. ret?:", ret != nil, " prevRet?:", prevRet != nil)
			}

			equals_s(t, string(ret), answer)
			equals_s(t, string(prevRet), answer)
			equals_s(t, string(ret), string(prevRet))
		}
		check(expr, bytes, ":super tom:")
		check("a*", []byte("abc"), "a")
		check("a*", []byte("bc"), "")
		check("a*", []byte(""), "")
		check("a", []byte(""), "")
	}
	{
		// FindAll
		var check = func(expr string, bytes []byte, answers []string) {
			re, closer := MustCompile(expr)
			defer closer.Close(re)

			prevRe := regexp.MustCompile(expr)

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
		check(expr, bytes, answers)
		check("a*", []byte("abc"), []string{"a", "", ""})
		check("a*", []byte("bc"), []string{"", "", ""})
		check("a*", []byte(""), []string{""})
		check("a", []byte(""), []string{})
	}
	{
		// FindAllIndex
		var check = func(expr string, bytes []byte, answers [][]int) {
			re, closer := MustCompile(expr)
			defer closer.Close(re)

			prevRe := regexp.MustCompile(expr)

			ret := re.FindAllIndex(bytes, -1)
			prevRet := prevRe.FindAllIndex(bytes, -1)

			if !equals_aai(t, ret, answers) {
				detailErrorParent(t, "wrong.")
			}
			if !equals_aai(t, prevRet, answers) {
				detailErrorParent(t, "wrong.")
			}
			if !equals_aai(t, ret, prevRet) {
				detailErrorParent(t, "wrong.")
			}
		}
		answers := [][]int{
			{4, 15},  // :super tom:
			{32, 37}, // :tom:
		}
		check(expr, bytes, answers)
		check("a*", []byte("abc"),
			[][]int{
				{0, 1}, // "a"
				{2, 2}, // ""
				{3, 3}, // ""
			},
		)
		check("a*", []byte("bc"),
			[][]int{
				{0, 0}, // ""
				{1, 1}, // ""
				{2, 2}, // ""
			},
		)
		check("a*", []byte(""),
			[][]int{
				{0, 0},
			},
		)
		check("a", []byte(""),
			[][]int{},
		)
	}
	{
		// FindAllString
		var check = func(expr string, s string, answers []string) {
			re, closer := MustCompile(expr)
			defer closer.Close(re)

			prevRe := regexp.MustCompile(expr)

			ret := re.FindAllString(s, -1)
			prevRet := prevRe.FindAllString(s, -1)

			equals_as(t, ret, answers)
			equals_as(t, prevRet, answers)
			equals_as(t, ret, prevRet)
		}
		answers := []string{
			":super tom:",
			":tom:",
		}
		check(expr, string(bytes), answers)
		check("a*", "abc", []string{"a", "", ""})
		check("a*", "bc", []string{"", "", ""})
		check("a*", "", []string{""})
		check("a", "", []string{})
	}
	{
		// FindAllStringIndex
		var check = func(expr string, s string, answers [][]int) {
			re, closer := MustCompile(expr)
			defer closer.Close(re)

			prevRe := regexp.MustCompile(expr)

			ret := re.FindAllStringIndex(s, -1)
			prevRet := prevRe.FindAllStringIndex(s, -1)

			equals_aai(t, ret, answers)
			equals_aai(t, prevRet, answers)
			equals_aai(t, ret, prevRet)
		}
		answers := [][]int{
			{4, 15},  // :super tom:
			{32, 37}, // :tom:
		}
		check(expr, string(bytes), answers)
		check("a*", "abc",
			[][]int{
				{0, 1}, // "a"
				{2, 2}, // ""
				{3, 3}, // ""
			},
		)
		check("a*", "bc",
			[][]int{
				{0, 0}, // ""
				{1, 1}, // ""
				{2, 2}, // ""
			},
		)
		check("a*", "",
			[][]int{
				{0, 0},
			},
		)
		check("a", "",
			[][]int{},
		)
	}
	{
		// FindAllStringSubmatch
		var check = func(expr string, bytes []byte, answers [][]string) {
			re, closer := MustCompile(expr)
			defer closer.Close(re)

			prevRe := regexp.MustCompile(expr)

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
		check(expr, bytes, answers)
		check("a*", []byte("abc"),
			[][]string{
				{"a"},
				{""},
				{""},
			},
		)
		check("a*", []byte("bc"),
			[][]string{
				{""},
				{""},
				{""},
			},
		)
		check("a*", []byte(""),
			[][]string{
				{""},
			},
		)
		check("a", []byte(""),
			[][]string{},
		)
	}
	{
		// FindAllStringSubmatchIndex
		var check = func(expr string, s string, answers [][]int) {
			re, closer := MustCompile(expr)
			defer closer.Close(re)

			prevRe := regexp.MustCompile(expr)

			ret := re.FindAllStringSubmatchIndex(s, -1)
			prevRet := prevRe.FindAllStringSubmatchIndex(s, -1)

			equals_aai(t, ret, answers)
			equals_aai(t, prevRet, answers)
			equals_aai(t, ret, prevRet)
		}
		answers := [][]int{
			{4, 15, 5, 10},   // {":super tom:", "super"},
			{32, 37, 33, 33}, // {":tom:", ""},
		}
		check(expr, string(bytes), answers)
		check("a*", "abc",
			[][]int{
				{0, 1}, // "a"
				{2, 2}, // ""
				{3, 3}, // ""
			},
		)
		check("a*", "bc",
			[][]int{
				{0, 0}, // ""
				{1, 1}, // ""
				{2, 2}, // ""
			},
		)
		check("a*", "",
			[][]int{
				{0, 0},
			},
		)
		check("a", "",
			[][]int{},
		)
	}
	{
		// FindAllSubmatch
		var check = func(expr string, bytes []byte, answers [][][]byte) {
			re, closer := MustCompile(expr)
			defer closer.Close(re)

			prevRe := regexp.MustCompile(expr)

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
		check(expr, bytes, answers)
		check("a*", []byte("abc"),
			[][][]byte{
				{[]byte("a")},
				{[]byte("")},
				{[]byte("")},
			},
		)
		check("a*", []byte("bc"),
			[][][]byte{
				{[]byte("")},
				{[]byte("")},
				{[]byte("")},
			},
		)
		check("a*", []byte(""),
			[][][]byte{
				{[]byte("")},
			},
		)
		check("a", []byte(""),
			[][][]byte{},
		)
	}
	{
		// FindAllSubmatchIndex
		var check = func(expr string, bytes []byte, answers [][]int) {
			re, closer := MustCompile(expr)
			defer closer.Close(re)

			prevRe := regexp.MustCompile(expr)

			ret := re.FindAllSubmatchIndex(bytes, -1)
			prevRet := prevRe.FindAllSubmatchIndex(bytes, -1)

			if !equals_aai(t, ret, answers) {
				detailErrorParent(t, "wrong.")
			}
			if !equals_aai(t, prevRet, answers) {
				detailErrorParent(t, "wrong.")
			}
			if !equals_aai(t, ret, prevRet) {
				detailErrorParent(t, "wrong.")
			}
		}
		answers := [][]int{
			{4, 15, 5, 10},   // {":super tom:", "super"},
			{32, 37, 33, 33}, // {":tom:", ""},
		}
		check(expr, bytes, answers)
		check("a*", []byte("abc"),
			[][]int{
				{0, 1}, // "a"
				{2, 2}, // ""
				{3, 3}, // ""
			},
		)
		check("a*", []byte("bc"),
			[][]int{
				{0, 0}, // ""
				{1, 1}, // ""
				{2, 2}, // ""
			},
		)
		check("a*", []byte(""),
			[][]int{
				{0, 0},
			},
		)
		check("a", []byte(""),
			[][]int{},
		)
	}
	{
		// FindIndex
		var check = func(expr string, bytes []byte, answers []int) {
			re, closer := MustCompile(expr)
			defer closer.Close(re)

			prevRe := regexp.MustCompile(expr)

			ret := re.FindIndex(bytes)
			prevRet := prevRe.FindIndex(bytes)

			equals_ai(t, ret, answers)
			equals_ai(t, prevRet, answers)
			equals_ai(t, ret, prevRet)
		}
		answers := []int{
			4, 15,
		}
		check(expr, bytes, answers)
		check("a*", []byte("abc"),
			[]int{0, 1}, // "a"
		)
		check("a*", []byte("bc"),
			[]int{0, 0}, // ""
		)
		check("a*", []byte(""),
			[]int{0, 0},
		)
		check("a", []byte(""),
			[]int{},
		)
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
		var check = func(expr string, s string, answers string) {
			re, closer := MustCompile(expr)
			defer closer.Close(re)

			prevRe := regexp.MustCompile(expr)

			ret := re.FindString(s)
			prevRet := prevRe.FindString(s)

			equals_s(t, ret, answers)
			equals_s(t, prevRet, answers)
			equals_s(t, ret, prevRet)
		}
		answers := ":super tom:"
		check(expr, string(bytes), answers)
		check("a*", "abc", "a")
		check("a*", "bc", "")
		check("a*", "", "")
		check("a", "", "")
	}
	{
		// FindStringIndex
		var check = func(expr string, s string, answers []int) {
			re, closer := MustCompile(expr)
			defer closer.Close(re)

			prevRe := regexp.MustCompile(expr)

			ret := re.FindStringIndex(s)
			prevRet := prevRe.FindStringIndex(s)

			equals_ai(t, ret, answers)
			equals_ai(t, prevRet, answers)
			equals_ai(t, ret, prevRet)
		}
		answers := []int{
			4, 15,
		}
		check(expr, string(bytes), answers)
		check("a*", "abc",
			[]int{0, 1}, // "a"
		)
		check("a*", "bc",
			[]int{0, 0}, // ""
		)
		check("a*", "",
			[]int{0, 0},
		)
		check("a", "",
			[]int{},
		)
	}
	{
		// FindStringSubmatch
		var check = func(expr string, s string, answers []string) {
			re, closer := MustCompile(expr)
			defer closer.Close(re)

			prevRe := regexp.MustCompile(expr)

			ret := re.FindStringSubmatch(s)
			prevRet := prevRe.FindStringSubmatch(s)

			equals_as(t, ret, answers)
			equals_as(t, prevRet, answers)
			equals_as(t, ret, prevRet)
		}
		answers := []string{
			":super tom:", "super",
		}
		check(expr, string(bytes), answers)
		check("a*", "abc",
			[]string{
				"a",
			},
		)
		check("a*", "bc",
			[]string{
				"",
			},
		)
		check("a*", "",
			[]string{
				"",
			},
		)
		check("a", "",
			[]string{},
		)
	}
	{
		// FindStringSubmatchIndex
		var check = func(expr string, s string, answers []int) {
			re, closer := MustCompile(expr)
			defer closer.Close(re)

			prevRe := regexp.MustCompile(expr)

			ret := re.FindStringSubmatchIndex(s)
			prevRet := prevRe.FindStringSubmatchIndex(s)

			equals_ai(t, ret, answers)
			equals_ai(t, prevRet, answers)
			equals_ai(t, ret, prevRet)
		}
		answers := []int{
			4, 15, 5, 10,
		}
		check(expr, string(bytes), answers)
		check("a*", "abc",
			[]int{0, 1}, // "a"
		)
		check("a*", "bc",
			[]int{0, 0}, // ""
		)
		check("a*", "",
			[]int{0, 0},
		)
		check("a", "",
			[]int{},
		)
	}
	{
		// FindSubmatch
		var check = func(expr string, bytes []byte, answers [][]byte) {
			re, closer := MustCompile(expr)
			defer closer.Close(re)

			prevRe := regexp.MustCompile(expr)

			ret := re.FindSubmatch(bytes)
			prevRet := prevRe.FindSubmatch(bytes)

			equals_aab(t, ret, answers)
			equals_aab(t, prevRet, answers)
			equals_aab(t, ret, prevRet)
		}
		answers := [][]byte{
			[]byte(":super tom:"), []byte("super"),
		}
		check(expr, bytes, answers)
		check("a*", []byte("abc"),
			[][]byte{
				[]byte("a"),
			},
		)
		check("a*", []byte("bc"),
			[][]byte{
				[]byte(""),
			},
		)
		check("a*", []byte(""),
			[][]byte{
				[]byte(""),
			},
		)
		check("a", []byte(""),
			[][]byte{},
		)
	}
	{
		// FindSubmatchIndex
		var check = func(expr string, bytes []byte, answers []int) {
			re, closer := MustCompile(expr)
			defer closer.Close(re)

			prevRe := regexp.MustCompile(expr)

			ret := re.FindSubmatchIndex(bytes)
			prevRet := prevRe.FindSubmatchIndex(bytes)

			equals_ai(t, ret, answers)
			equals_ai(t, prevRet, answers)
			equals_ai(t, ret, prevRet)
		}
		answers := []int{
			4, 15, 5, 10,
		}
		check(expr, bytes, answers)
		check("a*", []byte("abc"),
			[]int{0, 1}, // "a"
		)
		check("a*", []byte("bc"),
			[]int{0, 0}, // ""
		)
		check("a*", []byte(""),
			[]int{0, 0},
		)
		check("a", []byte(""),
			[]int{},
		)
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
	check(
		"a*",
		"abc",
		"a",
		"a",
	)
	check(
		"a*",
		"bc",
		"",
		"",
	)
	check(
		"a*",
		"",
		"",
		"",
	)
	check(
		"a",
		"",
		"",
		"",
	)
}

func TestMatch(t *testing.T) {
	expr := ":([^: ]*)\\s*tom:"
	bytes := []byte("abc :super tom: abc :hyper bob: :tom: abc")

	{
		// Match (static method)
		var check = func(expr string, bytes []byte, answer bool) {
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

			if ret != answer {
				detailErrorParent(t, "wrong")
			}
			if prevRet != answer {
				detailErrorParent(t, "wrong")
			}
			if ret != prevRet {
				detailErrorParent(t, "wrong")
			}
		}
		check(expr, bytes, true)
		check("a*", []byte("abc"), true)
		check("a*", []byte("bc"), true)
		check("a*", []byte(""), true)
		check("a", []byte(""), false)
	}
	{
		// MatchString (static method)
		var check = func(expr string, s string, answer bool) {
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

			if ret != answer {
				detailErrorParent(t, "wrong")
			}
			if prevRet != answer {
				detailErrorParent(t, "wrong")
			}
			if ret != prevRet {
				detailErrorParent(t, "wrong")
			}
		}
		check(expr, string(bytes), true)
		check("a*", "abc", true)
		check("a*", "bc", true)
		check("a*", "", true)
		check("a", "", false)
	}

	{
		// Match
		var check = func(expr string, bytes []byte, answer bool) {
			re, closer := MustCompile(expr)
			defer closer.Close(re)

			prevRe := regexp.MustCompile(expr)

			ret := re.Match(bytes)
			prevRet := prevRe.Match(bytes)

			if ret != answer {
				detailErrorParent(t, "wrong")
			}
			if prevRet != answer {
				detailErrorParent(t, "wrong")
			}
			if ret != prevRet {
				detailErrorParent(t, "wrong")
			}
		}
		check(expr, bytes, true)
		check("a*", []byte("abc"), true)
		check("a*", []byte("bc"), true)
		check("a*", []byte(""), true)
		check("a", []byte(""), false)
	}
	{
		fmt.Println("* MatchReader not implemented *")
		//		// MatchReader
		//		r := strings.NewReader(string(bytes))
		//		check(re.MatchReader(r))
	}
	{
		// MatchString
		var check = func(expr string, s string, answer bool) {
			re, closer := MustCompile(expr)
			defer closer.Close(re)

			prevRe := regexp.MustCompile(expr)

			ret := re.MatchString(s)
			prevRet := prevRe.MatchString(s)

			if ret != answer {
				detailErrorParent(t, "wrong")
			}
			if prevRet != answer {
				detailErrorParent(t, "wrong")
			}
			if ret != prevRet {
				detailErrorParent(t, "wrong")
			}
		}
		check(expr, string(bytes), true)
		check("a*", "abc", true)
		check("a*", "bc", true)
		check("a", "", false)
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
		check([]byte("-ab-axxb-"), []byte(""), []byte("---"))
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
		check("-ab-axxb-", "", "---")
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
		check(
			"a",
			"abc",
			-1,
			[]string{
				"", "bc",
			})
		check(
			"a",
			"bc",
			-1,
			[]string{
				"bc",
			})
		check(
			"a*",
			"",
			-1,
			[]string{
				"",
			})
		check(
			"a",
			"",
			-1,
			[]string{
				"",
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
		check("a")
		check("")
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
		check("", ``)
		//		check("[[:lower:]]", `\[\[:lower:\]\]`) // re2では`\[\[\:lower\:\]\]`となる
	}
}
