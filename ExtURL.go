package main

import (
	"./util"
	"flag"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/panjf2000/ants"
	"mvdan.cc/xurls"
	"strconv"
	"strings"
	"sync"
)

var (
	baseURL     = flag.String("u", "", "Base URL to find URLs (required)")
	strictHost  = flag.Bool("s", true, "Whether to target only the same domain")
	threads     = flag.Int("t", 10, "Number of threads")
	outputPath  = flag.String("o", "output.txt", "Output file path")
	engine  = flag.String("e", "", "Select engine (\"xurls\" or \"exturl\")")
	scanTargets []string
	checkedURLs = []string{*baseURL}
	baseHost    string
	wg          sync.WaitGroup
)
var p *ants.PoolWithFunc

func main() {
	flag.Parse()
	fmt.Println("[*] ExtURL v1.0 by RyotaK")
	if *baseURL == "" {
		flag.PrintDefaults()
		return
	}
	fmt.Println("[*] Checking connection to target url")
	var response = util.SendHTTPGet(*baseURL)
	if response == nil {
		fmt.Println("[x] Failed to connect target")
		return
	}
	baseHost = strings.Split(*baseURL, "/")[2]
	scanTargets = append(scanTargets, xurls.Strict().FindAllString(string(response), -1)...)
	fmt.Println("[*] Successfully connected to target url")
	fmt.Println("[*] Starting extract url with " + strconv.Itoa(*threads) + " threads")
	pool, _ := ants.NewPoolWithFunc(*threads, func(url interface{}) {
		checkURL(url)
		wg.Done()
	})
	defer pool.Release()
	p = pool
	for num := range scanTargets {
		wg.Add(1)
		pool.Invoke(string(scanTargets[num]))
	}
	wg.Wait()
	fmt.Println("[*] Saving result to " + *outputPath)
	util.WriteToFile(scanTargets, *outputPath)
}
func checkURL(URLInterface interface{}) {
	var URL = URLInterface.(string)
	if len(strings.Split(URL, "/")) >= 3 {
		var host = strings.Split(URL, "/")[2]
		if !util.ArrayContains(checkedURLs, URL) && !(*strictHost && host != baseHost) {
			checkedURLs = util.AppendWithCheck(checkedURLs, URL)
			var newURLs = detectURLs(URL)
			scanTargets = util.AppendWithCheck(scanTargets, newURLs...)
			for num := range scanTargets {
				if !util.ArrayContains(checkedURLs, scanTargets[num]) {
					wg.Add(1)
					go p.Invoke(string(scanTargets[num]))
				}
			}
		}
	}
}

func detectURLs(URL string) []string{
	if *engine != "exturl"{
		var response = util.SendHTTPGet(URL)
		return xurls.Strict().FindAllString(string(response), -1)
	}
	var detectedURLs []string
	doc,err := goquery.NewDocument(URL)
	if err != nil{
		fmt.Println(err)
		return nil
	}

	doc.Find("a").Each(func(i int, selection *goquery.Selection) {
		URL,_ := selection.Attr("href")
		detectedURLs = append(detectedURLs,URL)
	})
	doc.Find("script").Each(func(i int, selection *goquery.Selection) {
		URL,_ := selection.Attr("src")
		detectedURLs = append(detectedURLs,URL)
	})
	doc.Find("img").Each(func(i int, selection *goquery.Selection) {
		URL,_ := selection.Attr("src")
		detectedURLs = append(detectedURLs,URL)
	})
	return detectedURLs
}
