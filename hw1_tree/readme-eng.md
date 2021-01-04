This task is a practical part of a Coursera's course
"Golang web services"
(https://www.coursera.org/learn/golang-webservices-1/home/welcome)
This is an english translation for a task from Week-1.

All tests of all tasks in this course are integral part of a task.
Don't try to understand a task only by reading its description.
You should also watch tests.


A tool "Tree".

Shows a tree of directories and files (if option -f specified)

You need to implement function `dirTree` inside `main.go` file.
You can start from https://golang.org/pkg/os/#Open and look further
at the result methods.

Code should be written in main.go file.

In order to run tests through `go test -v`, your current dir must
be a task's dir. After execution you should see these results:


-------8<-------
$ go test -v
=== RUN   TestTreeFull
--- PASS: TestTreeFull (0.00s)
=== RUN   TestTreeDir
--- PASS: TestTreeDir (0.00s)
PASS
ok      coursera/homework/tree     0.127s
------->8-------

-------8<-------
go run main.go . -f
├───main.go (1881b)
├───main_test.go (1318b)
└───testdata
	├───project
	│	├───file.txt (19b)
	│	└───gopher.png (70372b)
	├───static
	│	├───css
	│	│	└───body.css (28b)
	│	├───html
	│	│	└───index.html (57b)
	│	└───js
	│		└───site.js (10b)
	├───zline
	│	└───empty.txt (empty)
	└───zzfile.txt (empty)

------->8-------

-------8<-------
go run main.go .
└───testdata
	├───project
	├───static
	│	├───css
	│	├───html
	│	└───js
	└───zline
------->8-------

Notices:


- Line breaks - unix-style ( \n )

- Indents - graphics symbol + tab ( \t )

- For graphic symbol calculation in indents try to think about
last element and prefix of previous level. There is very simple
condition. It's a good to speak aloud what you're seeing at the screen.

- When using Windows - remember about directories separator. They are
different. It's better to use `string(os.PathSeparator)`.

- Recursive algorithm is the simplest. But task can be done also without
recursion.

- You may implement any functions you need, you don't limited with a
single dirTree function. If you need more function parameters -
create another function and call it recursively. dirTree must act here
only as an entry point.

- Graphics symbols should be copied from sources (main_test.go), not from
this text. They are really different, believe it.

- (!!!) Results (list of directories-files) must be alphabetically sorted.
I.e. you should have a code, which sorts specific level. See at package
sort. This is a most obvious cause of failed tests. All tests run in linux
environment at remote server. Dockerfile for tests is included. It runs
tests in linux environment and shows all problems immediately.

- You way want to use global variables, but recursive version is simpler
without them, and non-recursive version doesn't need them at all.

- You shouldn't change a signature of dirTree function (number and order
of parameters). In this case tests on coursera server will fail.

- At MacOS you may encounter a problem with system file `.DS_Store` -
it can be ignored by your program itself or you can set up dockerignore:

-------8<-------
**/*.DS_Store
**/.git
------->8-------

Additional materials for reading:

- https://habrahabr.ru/post/306914/ - io package
- https://golang.org/pkg/sort/
- https://golang.org/pkg/io/
- https://golang.org/pkg/io/ioutil/