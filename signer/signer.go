package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

var wg sync.WaitGroup

func main() {
	start := time.Now()
	// inputData := []int{0, 1}
	inputData := []int{0, 1, 1, 2, 3, 5, 8}
	flowJobs := []job{
		job(func(in, out chan interface{}) {
			for _, v := range inputData {
				out <- v
			}
		}),
		job(SingleHash),
		job(MultiHash),
		job(CombineResults),
		job(func(in, out chan interface{}) {
			for data := range in {
				fmt.Println(data.(string))
			}
		}),
	}
	ExecutePipeline(flowJobs...)
	end := time.Since(start)
	fmt.Println(end)
}

func ExecutePipeline(jobs ...job) {
	wg := &sync.WaitGroup{}
	in := make(chan interface{})
	for _, funcJob := range jobs {
		wg.Add(1)
		out := make(chan interface{})
		go func(funcJob job, in, out chan interface{}, wg *sync.WaitGroup) {
			// if i == len(jobs)-1 {
			// 	job(chanel[i], chanel[i+1])
			// 	wg.Wait()
			// 	close(chanel[i+1])
			// }
			funcJob(in, out)
			wg.Done()
			close(out)
		}(funcJob, in, out, wg)
		in = out
	}
	wg.Wait()
}

func SingleHash(in, out chan interface{}) {
	var counter int
	for v := range in {
		counter++
		// var result string
		data := strconv.Itoa(v.(int))
		md5 := DataSignerMd5(data)
		wg.Add(2)
		go func(counter int, data string) {
			out <- strconv.Itoa(counter) + " " + "1" + " " + DataSignerCrc32(data) + "~"
			wg.Done()
		}(counter, data)
		go func(counter int, md5 string) {
			out <- strconv.Itoa(counter) + " " + "2" + " " + DataSignerCrc32(md5)
			wg.Done()
		}(counter, md5)
	}
	wg.Wait()
}

func MultiHash(in, out chan interface{}) {
	var (
		n                        string = "s"
		result                   string
		recieved, middleRes, res []string
	)
	for data := range in {
		recieved = append(recieved, data.(string))
	}
	sort.Strings(recieved)
	for i, val := range recieved {
		middleRes = strings.Split(val, " ")
		result += middleRes[2]
		if i%2 == 1 && i != 0 {
			res = append(res, result)
			result = ""
		}
	}
	for _, v := range res {
		for i := 0; i < 6; i++ {
			wg.Add(1)
			go func(i int, n, data string) {
				th := strconv.Itoa(i)
				out <- n + " " + th + " " + DataSignerCrc32(th+data)
				wg.Done()
			}(i, n, v)
		}
		n += "s"
	}
	wg.Wait()
}

func CombineResults(in, out chan interface{}) {
	var resultSlice, resultSup []string
	var result string
	count := 0
	for data := range in {
		resultSlice = append(resultSlice, data.(string))
	}
	sort.Strings(resultSlice)
	for _, v := range resultSlice {
		count++
		s := strings.Split(v, " ")
		result += s[2]
		if count%6 == 0 {
			resultSup = append(resultSup, result)
			result = ""
		}
	}
	sort.Strings(resultSup)
	out <- strings.Join(resultSup, "_")
}
