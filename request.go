package main

import (
	"io/ioutil"
	"net/http"
	"time"
)

var client = http.Client{
	Timeout: 5 * time.Second,
}

// this is dumb and i just stole it from my other projects so i didnt have to write the 3 lines again

// HTTPGet gets a url and error checks, loops request on error
func HTTPGet(url string) (string, int) {
	r, httperr := client.Get(url)
	if httperr != nil {
		time.Sleep(3 * time.Second)
		return HTTPGet(url)
	}

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return HTTPGet(url)
	}

	return string(body), r.StatusCode
}
