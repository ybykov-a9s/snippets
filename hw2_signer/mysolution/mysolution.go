package main

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"
)

var finalRes []string

//ExecutePipeline func
func ExecutePipeline(jobs ...job) {
	wg := &sync.WaitGroup{}
	x, y := make(chan interface{}), make(chan interface{})

	for _, f := range jobs {
		wg.Add(1)
		go JobWorker(f, x, y, wg)
		x, y = y, make(chan interface{})
	}
	wg.Wait()
}

// JobWorker func
func JobWorker(job job, in, out chan interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(out)
	job(in, out)
}

// SingleHash function
func SingleHash(in, out chan interface{}) {
	var wg sync.WaitGroup
	quotaCh := make(chan struct{}, 1) // Queue channel for maximum 1 concurent worker

	for s := range in {
		s := fmt.Sprintf("%v", s)
		wg.Add(1)
		go SingleHashWorker(s, out, quotaCh, &wg)
	}
	wg.Wait()
}

// SingleHashWorker func
func SingleHashWorker(s string, out chan interface{}, q chan struct{}, wg *sync.WaitGroup) {

	var wgcrc32 sync.WaitGroup
	var c32data, c32md5, res string
	defer wg.Done()

	n := "SingleHash"
	fmt.Printf("%s %s data %s\n", s, n, s)

	q <- struct{}{} // block until next free slot in Queue channel
	md5hash := DataSignerMd5(s)
	<-q // freeing slot in Queue channel
	fmt.Printf("%s %s md5(data) %s\n", s, n, md5hash)

	wgcrc32.Add(1)
	go func(wgcrc32 *sync.WaitGroup) {
		defer wgcrc32.Done()
		c32md5 = DataSignerCrc32(md5hash)
	}(&wgcrc32)

	c32data = DataSignerCrc32(s)
	fmt.Printf("%s %s, crc32(data) %s\n", s, n, c32data)

	wgcrc32.Wait()
	res = c32data + "~" + c32md5
	fmt.Printf("%s %s crc32(data) %s\n", s, n, res)
	out <- res
}

// MultiHash func
func MultiHash(in, out chan interface{}) {
	var wg sync.WaitGroup

	for s := range in {
		s := fmt.Sprintf("%v", s)
		wg.Add(1)
		go MultiHashWorker(s, out, &wg)
	}
	wg.Wait()
}

// MultiHashWorker - make N*goroutines, where N=len(InitialDataSlice)
func MultiHashWorker(s string, out chan interface{}, wg *sync.WaitGroup) {
	var cmb []string
	var res string
	var m = map[int]string{}
	mu := &sync.Mutex{}
	var wg2 sync.WaitGroup
	defer wg.Done()

	for i := 0; i < 6; i++ {
		wg2.Add(1)
		go func(i int, wg *sync.WaitGroup) {
			defer wg.Done()
			a := fmt.Sprintf("%d", i) + s
			mhash := DataSignerCrc32(fmt.Sprintf("%d", i) + s)
			mu.Lock()
			m[i] = mhash
			cmb = append(cmb, mhash)
			mu.Unlock()
			fmt.Printf("%s MiltiHash: crc32(%s)) %d %s\n", s, a, i, mhash)
			return
		}(i, &wg2)
	}
	wg2.Wait()
	for i := 0; i < 6; i++ {
		res = res + m[i]
	}
	out <- res
	fmt.Printf("%s MultiHash result:\n%s\n\n", s, res)
}

// CombineResults func
func CombineResults(in, out chan interface{}) {
	cmb := &finalRes

	for s := range in {
		s := fmt.Sprintf("%v", s)
		*cmb = append(*cmb, s)
	}
	sort.Strings(finalRes)
	out <- strings.Join(finalRes[:], "_")
}

func main() {
	start := time.Now()

	jobsList := []job{
		job(func(in, out chan interface{}) {
			inputData := []int{0, 1, 1, 2, 3, 5, 8}
			for _, f := range inputData {
				out <- f
			}
		}),
		job(SingleHash),
		job(MultiHash),
		job(CombineResults),
		job(func(in, out chan interface{}) {
			fmt.Println("Conveyer out:   ", <-in)
		}),
	}

	ExecutePipeline(jobsList...)

	elapsed := time.Since(start)
	fmt.Printf("Program execution took %s", elapsed)
}
