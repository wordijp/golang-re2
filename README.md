golang-re2
==========

## Introduction
Googleの正規表現ライブラリre2のラッパーライブラリ
regexpパッケージのパフォーマンスを改善したい時に置き換える事によって解決出来る事を目的にしています。

## License

LICENSEに記載(The MIT License(MIT))

## Install

C++製のRE2と、[RE2のCラッパーライブラリであるCRE2](https://github.com/marcomaggi/cre2)の事前インストールが必要です。
なお、CRE2の一部関数は、そのままではcgoによるinclude時にエラーになりますので、cre2.hに対して下記のpatchファイル、cre2.patchを適用します

  $ cp ./cre2.patch CRE2インストールディレクトリ/cre2/
  $ cd CRE2インストールディレクトリ/cre2
  $ patch -p1 < cre2.patch

## Usage

使い方やメソッド名や引数等は、regexpパッケージとほぼ同じです。
regexpパッケージと動作の違う一部のメソッド(ReplaceAll等)は、prefixにRE2を付けて名前の差別化を図っています。

また、元のC++クラスをCにラップして使っている関係上、Regexpの終了処理が必要になっていますが、
それは(Must)?Compile(POSIX)?の戻り値に終了処理用のCloserを返す事によって対応しています。
また、戻り値のインターフェースを変更する事によって、regexpパッケージからの単純置き換え時に終了処理を追加するまでコンパイルエラーになるように図っています。

具体的な使い方はテストメソッドが書かれているre2_bench_test.goやre2_test.goを読んでください。

## Benchmark
### regexpパッケージと置き換え後のベンチマーク比較
$ benchcmp benchbefore.txt benchafter.txt
benchmark                               old ns/op     new ns/op     delta
BenchmarkCompile                        11357         27268         +140.10%
BenchmarkExpand                         1304          23216         +1680.37%
BenchmarkExpandString                   975           24266         +2388.82%
BenchmarkFind                           2687          1918          -28.62%
BenchmarkFindAll                        8065          4196          -47.97%
BenchmarkFindAllIndex                   8075          4414          -45.34%
BenchmarkFindAllString                  7990          4781          -40.16%
BenchmarkFindAllStringIndex             8025          4533          -43.51%
BenchmarkFindAllStringSubmatch          8310          6646          -20.02%
BenchmarkFindAllStringSubmatchIndex     8055          5619          -30.24%
BenchmarkFindAllSubmatch                8535          5511          -35.43%
BenchmarkFindAllSubmatchIndex           8155          5459          -33.06%
BenchmarkFindIndex                      2684          1974          -26.45%
BenchmarkFindString                     2672          2095          -21.59%
BenchmarkFindStringIndex                2674          2144          -19.82%
BenchmarkFindStringSubmatch             3014          2923          -3.02%
BenchmarkFindStringSubmatchIndex        2900          2554          -11.93%
BenchmarkFindSubmatch                   3052          2406          -21.17%
BenchmarkFindSubmatchIndex              2893          2406          -16.83%
BenchmarkLongest                        81.9          11037         +13376.19%
BenchmarkMatchStatic                    13854         59339         +328.32%
BenchmarkMatchStringStatic              13899         59499         +328.08%
BenchmarkMatch                          2380          1653          -30.55%
BenchmarkMatchString                    2360          1861          -21.14%
BenchmarkNumSubexp                      0.31          78.1          +25093.55%
BenchmarkReplaceAll                     2357          2887          +22.49%
BenchmarkReplaceAllLiteral              2078          2894          +39.27%
BenchmarkReplaceAllLiteralString        2145          3206          +49.46%
BenchmarkReplaceAllString               2332          3084          +32.25%
BenchmarkSplit                          3091          2141          -30.73%
BenchmarkString                         1.23          97.2          +7802.44%
BenchmarkQuoteMeta                      282           991           +251.42%
BenchmarkReplaceRE2Sequences            106           114           +7.55%

数Byte、数十Byteのバイト列に対してのベンチーマークです、
Find系は一律高速化されているのに対して、
ReplaceAll(Literal | LiteralString | String)はGoからCを呼ぶオーバーヘッドにより遅くなっています、
また、CompileやExpand、Longest、NumSubexp、Stringは著しく遅くなっています。

しかし、対象のバイト列が十数KB以上になると、ReplaceAll系はこのラッパーライブラリの方が高速になります、
例えば、走れメロスの全文が記載されている30,895Byteのテキストファイル内の「メ.ス」を「ドラ○もん」に置換する場合、

benchmark          old ns/op     new ns/op     delta
BenchmarkMelos     400067        211841        -47.05%

となります。

数MByte、数十MByteのバイト列に対しては、[ベンチマーク比較サイトのregexp-dnaの項目](http://benchmarksgame.alioth.debian.org/u32/performance.php?test=regexdna#about)
にて、[Goのregexpパッケージのベンチマーク](http://benchmarksgame.alioth.debian.org/u32/program.php?test=regexdna&lang=go&id=1)にて、

```
$ time ./regex-dna < regexdna-input5000000.txt
ilen: 50833411
agggtaaa|tttaccct 356
[cgt]gggtaaa|tttaccc[acg] 1250
a[act]ggtaaa|tttacc[agt]t 4252
ag[act]gtaaa|tttac[agt]ct 2894
agg[act]taaa|ttta[agt]cct 5435
aggg[acg]aaa|ttt[cgt]ccct 1537
agggt[cgt]aa|tt[acg]accct 1431
agggta[cgt]a|t[acg]taccct 1608
agggtaa[cgt]|[acg]ttaccct 2178

50833411
50000000
66800214

real    1m20.403s
user    0m0.000s
sys     0m0.015s
```

```
$ time ./regex-dna_wrap_re2 < regexdna-input5000000.txt
ilen: 50833411
agggtaaa|tttaccct 356
[cgt]gggtaaa|tttaccc[acg] 1250
a[act]ggtaaa|tttacc[agt]t 4252
ag[act]gtaaa|tttac[agt]ct 2894
agg[act]taaa|ttta[agt]cct 5435
aggg[acg]aaa|ttt[cgt]ccct 1537
agggt[cgt]aa|tt[acg]accct 1431
agggta[cgt]a|t[acg]taccct 1608
agggtaa[cgt]|[acg]ttaccct 2178

50833411
50000000
66800214

real    0m5.947s
user    0m0.000s
sys     0m0.000s
```

という結果になりました。

### ベンチマーク取得時のコミット
regexpパッケージでのベンチマーク結果取得時のコミット・ファイル
  「ベンチマーク処理をregexp版の使用に戻した」 : b6649e86d2d804fb716b6094a896ba149dd26614
  benchbefore.txt

置き換え後のベンチマーク結果取得時のコミット・ファイル
  「Revert "ベンチマーク処理をregexp版の使用に戻した"」 :  85f43a6b0a0ecaaa249c13983e6f493dc4b434d4
  benchafter.txt


[![Build Status](https://drone.io/github.com/wordijp/golang-re2/status.png)](https://drone.io/github.com/wordijp/golang-re2/latest)
