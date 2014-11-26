package re2

import (
	"testing"
)

func BenchmarkCompile(b *testing.B) {
	expr := `.*name\s+is\s+(.+)\.`

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		re, closer := MustCompile(expr)
		defer closer.Close(re)
	}
}

func BenchmarkExpand(b *testing.B) {
	expr := `.*name\s+is\s+(.+)\.`
	re, closer := MustCompile(expr)
	defer closer.Close(re)

	src := []byte(`
		my name is tom.
		my favorite food is sushi.
		hello, my name is bob.
		he name is hiroshi.
	`)
	dst := []byte("prefix ")
	template := []byte("name = $1")
	indexes := re.FindAllSubmatchIndex(src, -1)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, x := range indexes {
			_ = re.RE2Expand(dst, template, src, x)
		}
	}
}

func BenchmarkExpandString(b *testing.B) {
	expr := `.*name\s+is\s+(.+)\.`
	re, closer := MustCompile(expr)
	defer closer.Close(re)

	src := `
		my name is tom.
		my favorite food is sushi.
		hello, my name is bob.
		he name is hiroshi.
	`
	dst := []byte("prefix ")
	template := "name = $1"
	indexes := re.FindAllSubmatchIndex([]byte(src), -1)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, x := range indexes {
			_ = re.RE2ExpandString(dst, template, src, x)
		}
	}
}

func BenchmarkFind(b *testing.B) {
	expr := ":([^: ]*)\\s*tom:"
	re, closer := MustCompile(expr)
	defer closer.Close(re)

	bytes := []byte("abc :super tom: abc :hyper bob: :tom: abc")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = re.Find(bytes)
	}
}

func BenchmarkFindAll(b *testing.B) {
	expr := ":([^: ]*)\\s*tom:"
	re, closer := MustCompile(expr)
	defer closer.Close(re)

	bytes := []byte("abc :super tom: abc :hyper bob: :tom: abc")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = re.FindAll(bytes, -1)
	}
}

func BenchmarkFindAllIndex(b *testing.B) {
	expr := ":([^: ]*)\\s*tom:"
	re, closer := MustCompile(expr)
	defer closer.Close(re)

	bytes := []byte("abc :super tom: abc :hyper bob: :tom: abc")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = re.FindAllIndex(bytes, -1)
	}
}

func BenchmarkFindAllString(b *testing.B) {
	expr := ":([^: ]*)\\s*tom:"
	re, closer := MustCompile(expr)
	defer closer.Close(re)

	s := "abc :super tom: abc :hyper bob: :tom: abc"

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = re.FindAllString(s, -1)
	}
}

func BenchmarkFindAllStringIndex(b *testing.B) {
	expr := ":([^: ]*)\\s*tom:"
	re, closer := MustCompile(expr)
	defer closer.Close(re)

	s := "abc :super tom: abc :hyper bob: :tom: abc"

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = re.FindAllStringIndex(s, -1)
	}
}

func BenchmarkFindAllStringSubmatch(b *testing.B) {
	expr := ":([^: ]*)\\s*tom:"
	re, closer := MustCompile(expr)
	defer closer.Close(re)

	s := "abc :super tom: abc :hyper bob: :tom: abc"

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = re.FindAllStringSubmatch(s, -1)
	}
}

func BenchmarkFindAllStringSubmatchIndex(b *testing.B) {
	expr := ":([^: ]*)\\s*tom:"
	re, closer := MustCompile(expr)
	defer closer.Close(re)

	s := "abc :super tom: abc :hyper bob: :tom: abc"

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = re.FindAllStringSubmatchIndex(s, -1)
	}
}

func BenchmarkFindAllSubmatch(b *testing.B) {
	expr := ":([^: ]*)\\s*tom:"
	re, closer := MustCompile(expr)
	defer closer.Close(re)

	bytes := []byte("abc :super tom: abc :hyper bob: :tom: abc")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = re.FindAllSubmatch(bytes, -1)
	}
}

func BenchmarkFindAllSubmatchIndex(b *testing.B) {
	expr := ":([^: ]*)\\s*tom:"
	re, closer := MustCompile(expr)
	defer closer.Close(re)

	bytes := []byte("abc :super tom: abc :hyper bob: :tom: abc")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = re.FindAllSubmatchIndex(bytes, -1)
	}
}

func BenchmarkFindIndex(b *testing.B) {
	expr := ":([^: ]*)\\s*tom:"
	re, closer := MustCompile(expr)
	defer closer.Close(re)

	bytes := []byte("abc :super tom: abc :hyper bob: :tom: abc")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = re.FindIndex(bytes)
	}
}

func BenchmarkFindString(b *testing.B) {
	expr := ":([^: ]*)\\s*tom:"
	re, closer := MustCompile(expr)
	defer closer.Close(re)

	s := "abc :super tom: abc :hyper bob: :tom: abc"

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = re.FindString(s)
	}
}

func BenchmarkFindStringIndex(b *testing.B) {
	expr := ":([^: ]*)\\s*tom:"
	re, closer := MustCompile(expr)
	defer closer.Close(re)

	s := "abc :super tom: abc :hyper bob: :tom: abc"

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = re.FindStringIndex(s)
	}
}

func BenchmarkFindStringSubmatch(b *testing.B) {
	expr := ":([^: ]*)\\s*tom:"
	re, closer := MustCompile(expr)
	defer closer.Close(re)

	s := "abc :super tom: abc :hyper bob: :tom: abc"

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = re.FindStringSubmatch(s)
	}
}

func BenchmarkFindStringSubmatchIndex(b *testing.B) {
	expr := ":([^: ]*)\\s*tom:"
	re, closer := MustCompile(expr)
	defer closer.Close(re)

	s := "abc :super tom: abc :hyper bob: :tom: abc"

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = re.FindStringSubmatchIndex(s)
	}
}

func BenchmarkFindSubmatch(b *testing.B) {
	expr := ":([^: ]*)\\s*tom:"
	re, closer := MustCompile(expr)
	defer closer.Close(re)

	bytes := []byte("abc :super tom: abc :hyper bob: :tom: abc")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = re.FindSubmatch(bytes)
	}
}

func BenchmarkFindSubmatchIndex(b *testing.B) {
	expr := ":([^: ]*)\\s*tom:"
	re, closer := MustCompile(expr)
	defer closer.Close(re)

	bytes := []byte("abc :super tom: abc :hyper bob: :tom: abc")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = re.FindSubmatchIndex(bytes)
	}
}

func BenchmarkLongest(b *testing.B) {

	for i := 0; i < b.N; i++ {
		b.StopTimer()

		re, closer := MustCompile("hoge")
		defer closer.Close(re)

		b.StartTimer()

		re.Longest()

		b.StopTimer()
	}
}

func BenchmarkMatchStatic(b *testing.B) {
	expr := ":([^: ]*)\\s*tom:"
	bytes := []byte("abc :super tom: abc :hyper bob: :tom: abc")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = Match(expr, bytes)
	}
}

func BenchmarkMatchStringStatic(b *testing.B) {
	expr := ":([^: ]*)\\s*tom:"
	s := "abc :super tom: abc :hyper bob: :tom: abc"

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = MatchString(expr, s)
	}
}

func BenchmarkMatch(b *testing.B) {
	expr := ":([^: ]*)\\s*tom:"
	re, closer := MustCompile(expr)
	defer closer.Close(re)

	bytes := []byte("abc :super tom: abc :hyper bob: :tom: abc")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = re.Match(bytes)
	}
}

func BenchmarkMatchString(b *testing.B) {
	expr := ":([^: ]*)\\s*tom:"
	re, closer := MustCompile(expr)
	defer closer.Close(re)

	s := "abc :super tom: abc :hyper bob: :tom: abc"

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = re.MatchString(s)
	}
}

func BenchmarkNumSubexp(b *testing.B) {
	expr := "([a-z])"
	re, closer := MustCompile(expr)
	defer closer.Close(re)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = re.NumSubexp()
	}
}

func BenchmarkReplaceAll(b *testing.B) {
	expr := "a(x*)b"
	re, closer := MustCompile(expr)
	defer closer.Close(re)

	bytes := []byte("-ab-axxb-")
	repl := []byte("$1")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = re.RE2ReplaceAll(bytes, repl)
	}
}

func BenchmarkReplaceAllLiteral(b *testing.B) {
	expr := "a(x*)b"
	re, closer := MustCompile(expr)
	defer closer.Close(re)

	bytes := []byte("-ab-axxb-")
	repl := []byte("@$1@")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = re.ReplaceAllLiteral(bytes, repl)
	}
}

func BenchmarkReplaceAllLiteralString(b *testing.B) {
	expr := "a(x*)b"
	re, closer := MustCompile(expr)
	defer closer.Close(re)

	s := "-ab-axxb-"
	repl := "@$1@"

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = re.ReplaceAllLiteralString(s, repl)
	}
}

func BenchmarkReplaceAllString(b *testing.B) {
	expr := "a(x*)b"
	re, closer := MustCompile(expr)
	defer closer.Close(re)

	s := "-ab-axxb-"
	repl := "$1"

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = re.RE2ReplaceAllString(s, repl)
	}
}

func BenchmarkSplit(b *testing.B) {
	expr := "[0-9](;)"
	re, closer := MustCompile(expr)
	defer closer.Close(re)

	s := "abc;123;ABC;45;"

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = re.Split(s, -1)
	}
}

func BenchmarkString(b *testing.B) {
	expr := ":([^: ]*)\\s*tom:"
	re, closer := MustCompile(expr)
	defer closer.Close(re)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = re.String()
	}
}

func BenchmarkQuoteMeta(b *testing.B) {
	expr := "[foo]"

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = RE2QuoteMeta(expr)
	}
}
