package re2

import (
	"io"
	"regexp"
)

type Regexp struct {
	// 元の変数をそのまま使う
	origRe *regexp.Regexp
}

func Match(pattern string, b []byte) (matched bool, err error) {
	return regexp.Match(pattern, b)
}

func MatchReader(pattern string, r io.RuneReader) (matched bool, err error) {
	return regexp.MatchReader(pattern, r)
}

func MatchString(pattern string, s string) (matched bool, err error) {
	return regexp.MatchString(pattern, s)
}

func QuoteMeta(s string) string {
	return regexp.QuoteMeta(s)
}

// test
func Compile(expr string) (*Regexp, error) {
	origRe, err := regexp.Compile(expr)
	re := &Regexp{
		origRe: origRe,
	}
	return re, err
}

// test
func MustCompile(str string) *Regexp {
	re := &Regexp{
		origRe: regexp.MustCompile(str),
	}
	return re
}

// test
func (re *Regexp) Expand(dst []byte, template []byte, src []byte, match []int) []byte {
	return re.origRe.Expand(dst, template, src, match)
}

// test
func (re *Regexp) ExpandString(dst []byte, template string, src string, match []int) []byte {
	return re.origRe.ExpandString(dst, template, src, match)
}

// test
func (re *Regexp) Find(b []byte) []byte {
	return re.origRe.Find(b)
}

// test
func (re *Regexp) FindAll(b []byte, n int) [][]byte {
	return re.origRe.FindAll(b, n)
}

func (re *Regexp) FindAllIndex(b []byte, n int) [][]int {
	return re.origRe.FindAllIndex(b, n)
}

func (re *Regexp) FindAllString(s string, n int) []string {
	return re.origRe.FindAllString(s, n)
}

func (re *Regexp) FindAllStringIndex(s string, n int) [][]int {
	return re.origRe.FindAllStringIndex(s, n)
}

func (re *Regexp) FindAllStringSubmatch(s string, n int) [][]string {
	return re.origRe.FindAllStringSubmatch(s, n)
}

func (re *Regexp) FindAllStringSubmatchIndex(s string, n int) [][]int {
	return re.origRe.FindAllStringSubmatchIndex(s, n)
}

func (re *Regexp) FindAllSubmatch(b []byte, n int) [][][]byte {
	return re.origRe.FindAllSubmatch(b, n)
}

func (re *Regexp) FindAllSubmatchIndex(b []byte, n int) [][]int {
	return re.origRe.FindAllSubmatchIndex(b, n)
}

func (re *Regexp) FindIndex(b []byte) (loc []int) {
	return re.origRe.FindIndex(b)
}

func (re *Regexp) FindReaderIndex(r io.RuneReader) (loc []int) {
	return re.origRe.FindReaderIndex(r)
}

func (re *Regexp) FindReaderSubmatchIndex(r io.RuneReader) []int {
	return re.origRe.FindReaderSubmatchIndex(r)
}

func (re *Regexp) FindString(s string) string {
	return re.origRe.FindString(s)
}

func (re *Regexp) FindStringIndex(s string) (loc []int) {
	return re.origRe.FindStringIndex(s)
}

func (re *Regexp) FindStringSubmatch(s string) []string {
	return re.origRe.FindStringSubmatch(s)
}

func (re *Regexp) FindStringSubmatchIndex(s string) []int {
	return re.origRe.FindStringSubmatchIndex(s)
}

func (re *Regexp) FindSubmatch(b []byte) [][]byte {
	return re.origRe.FindSubmatch(b)
}

func (re *Regexp) FindSubmatchIndex(b []byte) []int {
	return re.origRe.FindSubmatchIndex(b)
}

func (re *Regexp) LiteralPrefix() (prefix string, complete bool) {
	return re.origRe.LiteralPrefix()
}

func (re *Regexp) Longest() {
	re.origRe.Longest()
}

func (re *Regexp) Match(b []byte) bool {
	return re.origRe.Match(b)
}

func (re *Regexp) MatchReader(r io.RuneReader) bool {
	return re.origRe.MatchReader(r)
}

func (re *Regexp) MatchString(s string) bool {
	return re.origRe.MatchString(s)
}

func (re *Regexp) NumSubexp() int {
	return re.origRe.NumSubexp()
}

func (re *Regexp) ReplaceAll(src, repl []byte) []byte {
	return re.origRe.ReplaceAll(src, repl)
}

func (re *Regexp) ReplaceAllFunc(src []byte, repl func([]byte) []byte) []byte {
	return re.origRe.ReplaceAllFunc(src, repl)
}

func (re *Regexp) ReplaceAllLiteral(src, repl []byte) []byte {
	return re.origRe.ReplaceAllLiteral(src, repl)
}

func (re *Regexp) ReplaceAllLiteralString(src, repl string) string {
	return re.origRe.ReplaceAllLiteralString(src, repl)
}

func (re *Regexp) ReplaceAllString(src, repl string) string {
	return re.origRe.ReplaceAllString(src, repl)
}

func (re *Regexp) ReplaceAllStringFunc(src string, repl func(string) string) string {
	return re.origRe.ReplaceAllStringFunc(src, repl)
}

func (re *Regexp) Split(s string, n int) []string {
	return re.origRe.Split(s, n)
}

func (re *Regexp) String() string {
	return re.origRe.String()
}

func (re *Regexp) SubexpNames() []string {
	return re.origRe.SubexpNames()
}
