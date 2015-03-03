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

**WORK IN PROGRESS -- do not use this for anything requiring correctness**
