# msgpack
MsgPack Decoding and Encoding

```bash
go test -bench=. -run=X -cpuprofile cpu.out -memprofile mem.out
PASS
BenchmarkDecodeBool     30000000                38.2 ns/op
BenchmarkDecodeFixInt   30000000                53.7 ns/op
BenchmarkDecodeUInt     20000000                74.7 ns/op
BenchmarkDecodeUInt16   20000000                80.0 ns/op
BenchmarkDecodeUInt32   20000000                70.5 ns/op
BenchmarkDecodeUInt64   20000000                85.6 ns/op
BenchmarkDecodeInt      30000000                56.0 ns/op
BenchmarkDecodeInt16    20000000                71.4 ns/op
BenchmarkDecodeInt32    20000000                70.1 ns/op
BenchmarkDecodeInt64    20000000                75.6 ns/op
BenchmarkDecodeFixStr   10000000               174 ns/op
BenchmarkDecodeStr8      2000000               658 ns/op
BenchmarkDecodeStr16           0                 0 ns/op
BenchmarkDecodeStr32           0                 0 ns/op
BenchmarkDecodeFixArray  1000000              1624 ns/op
BenchmarkDecodeArray16    200000             11262 ns/op
BenchmarkDecodeFixMap     200000              6486 ns/op
BenchmarkDecodeMap16           0                 0 ns/op

BenchmarkEncodeBool     10000000               191 ns/op
BenchmarkEncodeFixInt   10000000               224 ns/op
BenchmarkEncodeInt16    10000000               236 ns/op
BenchmarkEncodeInt32    10000000               233 ns/op
BenchmarkEncodeInt64    10000000               238 ns/op
BenchmarkEncodeUInt16   10000000               181 ns/op
BenchmarkEncodeUInt32   10000000               184 ns/op
BenchmarkEncodeUInt64   10000000               231 ns/op
BenchmarkEncodeStr8      5000000               324 ns/op
BenchmarkEncodeArrayString       5000000               268 ns/op
BenchmarkEncodeFixMap    2000000               791 ns/op
BenchmarkEncodeMap16      500000              3250 ns/op
```

## Features

| Type      | Decode Impl | Decode Test | Encode Impl | Encode Test | Issues |
| ----      | ----------- | ----------- | ----------- | ----------- | ------ |
| fixint    | X           |             | X           | X           | None   |
| fixmap    | X           |             | X           | X           | None   |
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
| array\*   | X           |             | X           | X           | None   |
| map\*     | X           |             | X           | X           | None   |

**WORK IN PROGRESS -- do not use this for anything requiring correctness**
