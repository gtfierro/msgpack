# msgpack
MsgPack Decoding and Encoding

```bash
gabe@gabenix:~/src/msgpack$ go test -bench=. -run=X -cpuprofile=cpu.out
PASS
BenchmarkEncodeBool      5000000               304 ns/op
BenchmarkEncodeFixInt    5000000               324 ns/op
BenchmarkEncodeInt16     5000000               335 ns/op
BenchmarkEncodeInt32     5000000               333 ns/op
BenchmarkEncodeInt64     5000000               335 ns/op
BenchmarkEncodeUInt16    5000000               303 ns/op
BenchmarkEncodeUInt32    5000000               304 ns/op
BenchmarkEncodeUInt64    5000000               327 ns/op
BenchmarkEncodeStr8      3000000               487 ns/op
```

## Features

| Type      | Decode Impl | Decode Test | Encode Impl | Encode Test | Issues |
| ----      | ----------- | ----------- | ----------- | ----------- | ------ |
| fixint    | X           |             | X           | X           | None   |
| fixmap    | X           |             |             |             | None   |
| fixarray  | X           |             |             |             | None   |
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
