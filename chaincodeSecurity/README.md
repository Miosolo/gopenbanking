# TEST METHODS
## unit test
1. If you want to have unit test, please type `go test` under the catalogue of original file and test file in command line.
2. `go test -cover -covermode count -coverprofile ./cover.out` can run unit test while get the cover rate of unit test.

## Benchmark test
1. Type `go test -bench=. -benchmem` will run all benchmark with unit test.
2. If you do not want to run unit test while running benchmark,  please add argument `-run=none`, because there usually do not have unit test method called `none`.
3. If you want to change testing time(default time is 1s), you can add `-benchtime=3s` to change testing time to 3s.
4. `go test -benchmem -run=^$ -bench ^(Function name)$` to run benchmark for each specific funtion. For example, `go test -benchmem -run=^$ -bench ^(BenchmarkCreateAddReduceDelete)$` run BenchmarkCreateAddReduceDelete funtion.