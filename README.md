# msgpack
MsgPack Decoding and Encoding

```bash
gabe@gabenix:~/src/msgpack$ go test -bench=. -run=X -cpuprofile=cpu.out
BenchmarkDecodeBool     30000000                37.3 ns/op
BenchmarkDecodeFixInt   30000000                53.2 ns/op
BenchmarkDecodeUInt     20000000                66.6 ns/op
BenchmarkDecodeUInt16   20000000                70.8 ns/op
BenchmarkDecodeUInt32   20000000                76.6 ns/op
BenchmarkDecodeUInt64   20000000                74.9 ns/op
BenchmarkDecodeInt      30000000                59.0 ns/op
BenchmarkDecodeInt16    20000000                62.6 ns/op
BenchmarkDecodeInt32    20000000                76.4 ns/op
BenchmarkDecodeInt64    20000000                79.9 ns/op
BenchmarkDecodeFixStr   10000000               167 ns/op
BenchmarkDecodeStr8      3000000               599 ns/op
BenchmarkDecodeStr16           0                 0 ns/op
BenchmarkDecodeStr32           0                 0 ns/op
BenchmarkDecodeFixArray  1000000              1593 ns/op
BenchmarkDecodeArray16    200000             11062 ns/op
BenchmarkDecodeFixMap     200000              6183 ns/op
BenchmarkDecodeMap16           0                 0 ns/op

BenchmarkEncodeBool     10000000               194 ns/op
BenchmarkEncodeFixInt   10000000               223 ns/op
BenchmarkEncodeInt16    10000000               233 ns/op
BenchmarkEncodeInt32    10000000               231 ns/op
BenchmarkEncodeInt64    10000000               239 ns/op
BenchmarkEncodeUInt16   10000000               185 ns/op
BenchmarkEncodeUInt32   10000000               186 ns/op
BenchmarkEncodeUInt64   10000000               231 ns/op
BenchmarkEncodeStr8      5000000               320 ns/op
BenchmarkEncodeArrayString       1000000              1103 ns/op
```

## Features

| Type      | Decode Impl | Decode Test | Encode Impl | Encode Test | Issues |
| ----      | ----------- | ----------- | ----------- | ----------- | ------ |
| fixint    | X           |             | X           | X           | None   |
| fixmap    | X           |             |             |             | None   |
| fixarray  | X           |             | X           | X           | None   |
| fixstr    | X           |             | X           | X           | None   |
| nil       | X           |             | X           | X           | None   |
| false     | X           |             | X           | X           | None   |
| true      | X           |             | X           | X           | None   |
| bin\*     |             |             |             |             | None   |
| ext\*     |             |             |             |             | None   |
| fixext\*  |             |             |             |             | None   |
| float\*   | X           |             | X           | X           | negatives   |
| int\*     | X           |             | X           | X           | negatives   |
| uint\*    | X           |             | X           | X           | None   |
| str\*     | X           |             | X           | X           | None   |
| array\*   | X           |             | X           |             | not generic yet   |
| map\*     | X           |             |             |             | None   |

**WORK IN PROGRESS -- do not use this for anything requiring correctness**
