This task is a practical part of a Coursera's course
"Golang web services"
(https://www.coursera.org/learn/golang-webservices-1/home/welcome)
This is an english translation for a task from Week-3.

There is a function, that searches something in a file. But it do search
not so fast and you need to optimize it.

This is a task for working with go profiler 'pprof'.

Main goal of this task - learn how to work with pprof, find code hotspots,
can build CPU and memory usage, optimize a code using this information.
There is no need to write fastest solution in this task.

For generating graph you need graphviz. Windows users - don't forget add it
to PATH environment variable, command 'dot' should be available.

There is up to one dozen places you may optimize.

You need to get one of parameters ( ns/op, B/op, allocs/op ) just faster
than *BenchmarkSolution* ( fast < solution ) and another parameter better 
then *BenchmarkSolution* + 20% ( fast < solution * 1.2). For example
( fast allocs/op < 10422*1.2=12506 )

You can focus on next results - memory (B/op) and allocations per operation
(allocs/op) from results *BenchmarkSolution*, but not on time results (ns/op),
because they are dependent on system.

Result should be put in FastSearch (fast.go). Now there is the same as in SlowSearch.

Here is one of examples, which will be used to compare with your results:

```
$ go test -bench . -benchmem

goos: windows

goarch: amd64

BenchmarkSlow-8 10 142703250 ns/op 336887900 B/op 284175 allocs/op

BenchmarkSolution-8 500 2782432 ns/op 559910 B/op 10422 allocs/op

PASS

ok coursera/hw3 3.897s
```

Execution:
* `go test -v` - to check that everything is OK and not broken
* `go test -bench . -benchmem` - for monitoring performance

Hints:
* See where we allocate memory
* See where we accumulate results, though we don't need all values in the same time
* See where we make type conversions, which may be avoided
* See not only graphic representation of results, but and pprof wit text 
representation (list FastSearch)
- You can see in the code where are the problems
* Task assumes using easyjson lib. Generated with it code you should include
to file with your function.
* Task can be done without easyjson.


Note:
* easyjson based on reflection and can't work with main package.
  You need to move your structure to another packet, generate code there, then
  move code to your main.
