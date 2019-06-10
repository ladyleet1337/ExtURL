package main

import (
	"./util"
	"flag"
	"fmt"
	"mvdan.cc/xurls"
	"strconv"
	"strings"
)
var(
	baseURL      = flag.String("u","","Base URL to find URLs (required)")
	strictHost   = flag.Bool("h",true,"Whether to target only the same domain")
	threads      = flag.Int("t",10,"Number of threads")
	outputPath   = flag.String("o","output.txt","Output file path")
	scanTargets []string
	checkedURLs  = []string{*baseURL}
	baseHost    string
)
func main(){
	flag.Parse()
	fmt.Println("[*] ExtURL v1.0 by RyotaK")
	if *baseURL == ""{
		flag.PrintDefaults()
		return
	}
	fmt.Println("[*] Checking connection to target url")
	var response = util.SendHTTPGet(*baseURL)
	if response == nil{
		fmt.Println("[x] Failed to connect target")
		return
	}
	baseHost = strings.Split(*baseURL,"/")[2]
	scanTargets = append(scanTargets,xurls.Strict().FindAllString(string(response),-1)...)
	fmt.Println("[*] Successfully connected to target url")
	fmt.Println("[*] Starting extract url with "+strconv.Itoa(*threads)+" threads")
	jobs := make(chan string, 512)
	results := make(chan int, 512)
	var executedJobs int
	for w := 1; w <= *threads; w++ {
		go worker(jobs, results)
	}
	for num := range scanTargets {
		checkURL(scanTargets[num])
		//executedJobs++
		//jobs <- scanTargets[num]
	}
	close(jobs)
	for a := 1; a <= executedJobs; a++ {
		<-results
	}
	util.WriteToFile(scanTargets,*outputPath)
}

func worker(URLs <-chan string, results chan<- int){
	for URL := range URLs {
		checkURL(URL)
		results <- 0
	}
}

func checkURL(URL string){
	fmt.Println(URL)
	if len(strings.Split(URL,"/")) >= 3 {
		var host = strings.Split(URL, "/")[2]
		if !util.ArrayContains(checkedURLs, URL) && !(*strictHost && host != baseHost) {
			checkedURLs = util.AppendWithCheck(checkedURLs, URL)
			var response = util.SendHTTPGet(URL)
			var newURLs = xurls.Strict().FindAllString(string(response), -1)
			scanTargets = util.AppendWithCheck(scanTargets, newURLs...)
			for num := range scanTargets {
				checkURL(scanTargets[num])
			}
		}
	}
}
