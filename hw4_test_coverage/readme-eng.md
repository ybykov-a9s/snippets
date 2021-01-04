This task is a practical part of a Coursera's course
"Golang web services"
(https://www.coursera.org/learn/golang-webservices-1/home/welcome)
This is an english translation for a task from Week-4.

This is combined task: how to send requests, receive answers, work
with parameters, headers and write tests.

Task is not very complex. Main amount of work - writing different
conditions and tests, in order to satisfy these conditions.

We have some search engine:

* SearchClient - a structure with FindUsers method, which sends requests
to external system and gets returns result after some modification.
* SearchServer - a kind of external system. Do data search in
file `dataset.xml`. In production would be implemented as a separate
web-service.

Requirements:
* Write SearchServer function in file `client_test.go`, which should be run
through test server (`httptest.NewServer`).
* Cover with tests FindUsers method. Coverage has to be 100%. Tests should
be written in `client_test.go`.
* You also need to generate html-report with coverage (don't forget that
code should be located inside GOPATH environment variable)

Additional info:
* Work data is in `dataset.xml`
* `query` parameter searches by fields `Name` and `About`
* `order_field` parameter works with fields `Id`, `Age`, `Name`. If it
is empty, then `Name` should be used. If anything else - SearchServer
returns an error. `Name` - is first_name + last_name from xml.
* If `query` is empty, only sort should be done. I.e. all records should
be returned.
* Code should be written in client_test.go. There should be also tests and
SearchServer function.
* tests should be run as `go test -cover`
* Coverage building: `go test -coverprofile=cover.out && go tool
cover -html=cover.out -o cover.html`. Your code should be inside GOPATH.

Hints:
* Documentation https://golang.org/pkg/net/http/ may help
* Don't do one big test. Do many small tests instead.
* You don't have to limit yourself with SearchServer function while testing,
if you need to check some very tricky cases, like an error for example.
But there shouldn't be many such cases.
* For covering with test one of the errors, it is better to watch function sources,
which returns this error and see under what working conditions or input data
this error occurs.
* Performance, goroutines and other async stuff is not needed in this task
* Don't try to connect to unknown IP in order to get timeout error. In test
sandbox network connectivity is absent, so this method will return error
immediately 

