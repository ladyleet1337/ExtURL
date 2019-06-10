package util

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func ArrayContains(arr []string, str string) bool{
	for _, value := range arr{
		if value == str{
			return true
		}
	}
	return false
}

func SendHTTPPostJson(url string,json string) []byte{
	req,err := http.NewRequest(
		"POST",
		url,
		bytes.NewBuffer([]byte(json)),
	)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("Content-Type","application/json")
	client := http.Client{}
	res,err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	response,err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}
	res.Body.Close()
	return response
}

func SendHTTPGet(url string) []byte{
	resp,err := http.Get(url)
	if err != nil{
		fmt.Println(err)
		return nil
	}
	response,err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	return response
}

func AppendWithCheck(slice []string, elems ...string) []string{
	for num := range elems {
		if !ArrayContains(slice,elems[num]){
			slice = append(slice,elems[num])
		}
	}
	return slice
}

func ArrayRemove(slice []string, target string) []string {
	var result []string
	for _, value := range slice {
		if value != target {
			result = append(result, value)
		}
	}
	return result
}


func WriteToFile(texts []string, path string) {
	var writer *bufio.Writer
	file, _ := os.Create(path)
	writer = bufio.NewWriter(file)
	for idx := range texts {
		writer.WriteString(texts[idx] + "\n")
	}
	writer.Flush()
}
