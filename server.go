package main

import (
	"fmt"
	"net/http"
	"io"
)

type Response struct {
	URL string
	Bytes []byte
	Err error
}

func fetchFile(url string, channel chan string){
	resp, err := http.Get(url)

	if err != nil{
		fmt.Printf("Error fetching file from url %v: %v", url, err)
		response := Response{URL: url, Bytes: nil, Err: err}
		channel <- response
		return
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil{
		fmt.Printf("Error reading response body from url %v: %v", url, err)
		response := Response{URL: url, Bytes: nil, Err: err}
		channel <- response
		return
	}

	response := Response{URL: url, Bytes: body, Err: nil}
	channel <- response
}

func main(){

	urls := []string{"https://httpbin.org/image","https://httpbin.org/image/jpeg","https://httpbin.org/image/png",
	"https://httpbin.org/image/svg","https://httpbin.org/image/webp"}

	channel := make(chan Response, len(urls))

	sem := make(chan int, 10)

	
	for _, url := range urls{
		sem <- 1
		go func(u string,c chan Response){
			fetchFile(url ,channel)
			<- sem
		}(u, channel)
	}

	for i:=0;i<len(urls);i++{
		result := <- channel
		if result.Err != nil{
			fmt.Printf("Error fetching file from url %v: %v", result.URL, result.Err)
		}else{
			fmt.Printf("Successfully fetched file from url %v, bytes read: %v\n", result.URL, len(result.Bytes))
		}
	}
}