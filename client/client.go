package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type Client struct {
	http.Client
	Host string
	BasePath string
}

type Job struct {
	Text string 	`json:"text"`
	Tags []string 	`json:"tags"`
	Due time.Time 	`json:"due"`
}


func (c *Client) CreateTask(j *Job) error {
	url := c.Host + c.BasePath + "/"

	js, err := json.Marshal(j)
	if err != nil {
		return fmt.Errorf("error %v while creating task", err)
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(js)))
	if err != nil {
		return fmt.Errorf("error %v while setting request", err)
	}
	req.Header.Set("Content-Type", "application/json")
	_ , err = c.Do(req)

	if err != nil {
		return fmt.Errorf("error %v while getting respons", err)
	}
	if err != nil {
		return fmt.Errorf("unable to get response body %v", err)
	}
	return nil
}

func (c *Client) GetAllTasks() (*[]Job, error) {
	// req, err := http.NewRequest("GET", )

	resp, err := http.Get(c.Host + c.BasePath)
	if err != nil {
		return nil, fmt.Errorf("error %v while getting respons", err)
	}
	
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("unable to get response body %v", err)
	}

	defer resp.Body.Close()
	var data []Job

	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal json %v", err)
	}
	
	return &data, nil
}

func main() {
	c := Client {
		Host: "http://localhost:" + os.Args[1],
		BasePath: "/task",
	}
	jobs := make([]Job, 1)
	jobs[0] = Job {
			Text: "Kill Bill",
			Tags: []string{"fast", "blood", "katana"},
			Due: time.Date(2007, time.Month(10), 5, 4, 3, 2, 1, &time.Location{}),
	}
	c.CreateTask(&jobs[0])
	data, _  := c.GetAllTasks()
	fmt.Println(data)
}