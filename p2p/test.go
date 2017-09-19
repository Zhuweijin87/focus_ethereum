package main

import (
	"flag"
	"fmt"
	"os"
	"log"
	"net/http"
	"io/ioutil"
	"github.com/bitly/go-simplejson"
	"encoding/csv"
	"io"
	"time"
	"sync"
	"net/url"
	"strconv"
	"errors"
)

var (
	concurrency int
	timeout int
	infile string
	outfile string
)
var usage = `Usage:%s [options]
 	Options are:
 		-c concurrency Number of request to preform
 		-t timeout Request timeout
 		-i infile Input file
 		-o outfile Output file
`

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, usage, os.Args[0])
	}
	flag.IntVar(&concurrency, "c", 10, "")
	flag.IntVar(&timeout, "t", 60, "")
	flag.StringVar(&infile, "i", "", "")
	flag.StringVar(&outfile, "o", "", "")
	flag.Parse()
	_, err := os.Stat(infile)
	if err != nil {
		log.Fatalln("error:", err)
	}
	f, err := os.Create(outfile)
	if err != nil {
		log.Fatalln("error:", err)
	}
	defer f.Close()
	var lock sync.Mutex
	w := &Worker{
		concurrency:concurrency,
		timeout:timeout,
		infile:infile,
		lw:&LockWriter{
			m:lock,
			writer:f,
		},
	}
	w.Run()
}

type LockWriter struct {
	m      sync.Mutex
	writer io.Writer
}

func (lw LockWriter) write(b []byte) (n int, err error) {
	lw.m.Lock()
	defer lw.m.Unlock()
	return lw.writer.Write(b)
}

type CarInfo struct {
	Carno string
	Ecode string
	Vcode string
}

type Worker struct {
	concurrency int
	timeout     int
	infile      string
	jobs        chan *CarInfo
	lw          *LockWriter
}

func (w *Worker) Run() {
	var wg sync.WaitGroup
	wg.Add(w.concurrency + 1) 
	w.jobs = make(chan *CarInfo, w.concurrency)
	go func() {
		w.loadJobs()
		wg.Done()
	}()
	//并发数
	for i := 1; i <= w.concurrency; i++ {
		go func(n int) {
			w.doWork(n)
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func (w *Worker) loadJobs() {
	fin, err := os.Open(w.infile)
	if err != nil {
		log.Fatalln("file open failed,error:", err.Error())
	}
	defer fin.Close()
	reader := csv.NewReader(fin)
	for {
		row, err := reader.Read()
		if err != nil {
			if err != io.EOF {
				log.Println("file read error:", err)
			}
			break
		}
		time.Sleep(time.Duration(1) * time.Second)
		w.jobs <- &CarInfo{Carno:row[3], Vcode:row[4], Ecode:row[5]}
	}
	close(w.jobs)
}

func (w *Worker) doWork(num int) {
	total := 0
	log.Println("worker:", num)
	uriParams := url.Values{}
	uriParams.Add("openUDID", "3122dcf3-3a2a-34e9-8da5-e9dde29579a4")
	uriParams.Add("appid", "1")
	uriParams.Add("cartype", "02")
	uriParams.Add("os", "android")
	uriParams.Add("appVersion", "6.6.6")
	uriParams.Add("prefetch", "1")
	uriParams.Add("reqfrom", "1")
	uriParams.Add("secret", "UBjUFDL9kZSDUqivn4wb063QI4Es3mZhvWvT")
	httpClient := &http.Client{
		Timeout:time.Duration(w.timeout) * time.Second,
	}
	for carInfo := range w.jobs {
		log.Println(carInfo.Carno, "|", carInfo.Vcode, "|", carInfo.Ecode)
		uriParams.Set("carno", carInfo.Carno)
		uriParams.Set("vcode", carInfo.Vcode)
		uriParams.Set("ecode", carInfo.Ecode)
		uri := fmt.Sprintf("http://xxx.cn/common_prefix?%s", uriParams.Encode())
		fmt.Println(uri)
		code, err := w.doRequest(httpClient, uri)
		if err != nil {
			log.Println("url request error:", err)
		}
		switch code {
		case 203, 9999:
			w.writeResult(carInfo)
		}
		total++
	}
	log.Println("worker[", num, "] total:", total)
}

func (w *Worker) doRequest(client *http.Client, uri string) (int, error) {
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		log.Println("get new request failed!")
		return -1, err
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("error:", err)
		return -1, err
	}
	defer resp.Body.Close()
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("read respose body failed")
		return -1, err
	}
	log.Println("resonse:", string(result))
	if resp.StatusCode != 200 {
		log.Println("status code:", resp.StatusCode)
		return -1, errors.New("stats code error:" + strconv.Itoa(resp.StatusCode))
	}
	js, err := simplejson.NewJson(result)
	if err != nil {
		return -1, err
	}
	code, err := js.Get("code").Int()
	return code, nil
}

func (w *Worker) writeResult(carInfo *CarInfo) {
	line := fmt.Sprintf("%s,%s,%s\n", carInfo.Carno, carInfo.Vcode, carInfo.Ecode)
	w.lw.write([]byte(line))
}