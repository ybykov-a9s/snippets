This task is a practical part of a Coursera's course
"Golang web services"
(https://www.coursera.org/learn/golang-webservices-1/home/welcome)
This is an english translation for a task from Week-2

In this task we implement an unix pipeline equivalent, something like:
```
grep 127.0.0.1 | awk '{print $2}' | sort | uniq -c | sort -nr
```
When STDOUT of one program passes as STDIN of another program.

But in our case these roles act channels, which we pass from one function
to another.

The task consist basically of two parts
* Implementing function ExecutePipeline, which grant us conveyor processing
of Workers, that do something.
* Implementing several functions, which calculate us some hash-sum of
input data.

The hash-sum calculation is implemented with this sequence:
* SingleHash calculates value of crc32(data)+"~"+crc32(md5(data))
(concatenation of two lines through ~), where data us our input
(basically - numbers from first function).
* MultiHash calculates value crc32(th+data)) (concatenation of number,
converted to a string and a string), where th=0..5 (i.e. 6 hashes per every
input value). Than it takes concatenation of ordered results (0..5), where
data - is input data for a function (and SingleHash output).
* CombineResults gets all results, sorts them (https://golang.org/pkg/sort/),
combines sorted result with _ (underline symbol) to one line.
* crc32 calculates with DataSignerCrc32 function (included in common.go).
* md5 calculates with DataSignerMd5 function (included in common.go).

What's the catch:
* DataSignerMd5 can run simultaneuosly only 1 time, calculates 10 ms.
If it runs in parallel - you get 1 second Overheat penalty.
* DataSignerCrc32 - calculates 1 second.
* There are maximum 3 seconds for all calculations.
* If all calculations run straight forward in a chain, it takes around
57 seconds, so they should be some way parallelised.

There are results, you should get on a screen, when run a program:

```
0 SingleHash data 0
0 SingleHash md5(data) cfcd208495d565ef66e7dff9f98764da
0 SingleHash crc32(md5(data)) 502633748
0 SingleHash crc32(data) 4108050209
0 SingleHash result 4108050209~502633748
4108050209~502633748 MultiHash: crc32(th+step1)) 0 2956866606
4108050209~502633748 MultiHash: crc32(th+step1)) 1 803518384
4108050209~502633748 MultiHash: crc32(th+step1)) 2 1425683795
4108050209~502633748 MultiHash: crc32(th+step1)) 3 3407918797
4108050209~502633748 MultiHash: crc32(th+step1)) 4 2730963093
4108050209~502633748 MultiHash: crc32(th+step1)) 5 1025356555
4108050209~502633748 MultiHash result: 29568666068035183841425683795340791879727309630931025356555

1 SingleHash data 1
1 SingleHash md5(data) c4ca4238a0b923820dcc509a6f75849b
1 SingleHash crc32(md5(data)) 709660146
1 SingleHash crc32(data) 2212294583
1 SingleHash result 2212294583~709660146
2212294583~709660146 MultiHash: crc32(th+step1)) 0 495804419
2212294583~709660146 MultiHash: crc32(th+step1)) 1 2186797981
2212294583~709660146 MultiHash: crc32(th+step1)) 2 4182335870
2212294583~709660146 MultiHash: crc32(th+step1)) 3 1720967904
2212294583~709660146 MultiHash: crc32(th+step1)) 4 259286200
2212294583~709660146 MultiHash: crc32(th+step1)) 5 2427381542
2212294583~709660146 MultiHash result: 4958044192186797981418233587017209679042592862002427381542

CombineResults 29568666068035183841425683795340791879727309630931025356555_4958044192186797981418233587017209679042592862002427381542
```

Code should be written in signer.go. There nothing needs to be add from common.go.

Run it as `go test -v -race`

Hints:

* The task is built such way, you can use only sources from a course.
Watch samples and apply them practically. golang_mailru_coursera_2.zip from:
https://www.coursera.org/learn/golang-webservices-1/supplement/Eu14h/kod-i-konspiekty-ko-2-mu-uroku
There is no need to search something in google or stackoverflow.

* you shouldn't accumulate any data - just immediately pass it forward
(for example awk from code up do it such way. There is separate test,
that checks such behaviour. As exception may be a function, which accumulates
results by itself, like CombineResults in our case, or sort from code above).

* Think about finishing a function, if data is finite. What need to be done
to catch such situation?

* If you've got a race condition (-race option), examine its output - when we
read, when we write, in which codelines. Usually there is enough information
for detection source of problems.

* Before writing a parallelised code, in order to meet a deadline, write
first a linear code, which gives a right result. Try to begin with lesser
input values, to ease checks at first stage. Results should match program
output above.

* You can expect always get no more than 100 elements in input values.
* Answer to a question "When closes a loop through a channel?" helps with
implementation ExecutePipeline function.
* Answer to a question "Do I need results of previous calculations?" helps
with parallelising SingleHash and MultiHash.

* It is a good idea to visualize calculations with diagram.
* Of course you are not allowed to implement your own functions for calculating
hashes. It will be checked on a coursera server.
* Reference solution takes 130 lines of code, including debug you may see above.
