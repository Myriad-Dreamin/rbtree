# rbtree
rbtree impl in go

## Benchmark

```plain
goos: windows
goarch: amd64
pkg: github.com/Myriad-Dreamin/rbtree
BenchmarkRBNode_Insert1e6-12                           3         475342933 ns/op
BenchmarkRBNode_Insert1e7-12                           1        5974000900 ns/op
BenchmarkRBNode_InsertDelete1e6-12                     2         577013400 ns/op
BenchmarkRBNode_RandomBehavior1e6-12                   2         847494500 ns/op
BenchmarkRBNode_RandomBehaviorPure1e6-12             100          18063947 ns/op
PASS
ok      github.com/Myriad-Dreamin/rbtree        22.492s
```