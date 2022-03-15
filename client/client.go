package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	// "os"
)

type Client struct {
	Host string
	BasePath string
}

func (c *Client) MakeGetRequest() {
	resp, err := http.Get(c.Host + c.BasePath)
	if err != nil {
		log.Printf("error %v while getting respons", err)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("unable to get response body %v", err)
		return
	}

	fmt.Println(string(body))
}

func main() {
	c := Client {
		Host: "http://localhost:" + os.Args[1],
		BasePath: "/task",
	}
	c.MakeGetRequest()
}