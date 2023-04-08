package source

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func ClientMain() {
	// set up the proxy URL
	proxyUrl, err := url.Parse("http://localhost:18080")
	if err != nil {
		panic(err)
	}

	// set up a custom transport that uses the proxy
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyUrl),
	}

	// set up an http client that uses the custom transport
	client := &http.Client{
		Transport: transport,
	}

	// make a request using the http client
	resp, err := client.Get("https://httpbin.org/get")
	if err != nil {
		panic(err)
	}

	//read the resp body into a string
	body, _ := ioutil.ReadAll(resp.Body)

	// print the response status code
	fmt.Println(string(body))
}
