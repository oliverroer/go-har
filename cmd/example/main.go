package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	harwriter "github.com/oliverroer/go-har/writer"
)

func main() {
	err := example()
	if err != nil {
		log.Fatal(err)
	}
}

func example() error {
	writer, err := harwriter.OpenDefault("out")
	if err != nil {
		return err
	}

	transport := writer.RoundTripper(http.DefaultTransport)
	client := http.Client{
		Transport: transport,
	}

	for number := range 2 {
		ip, err := getIP(client)
		if err != nil {
			return err
		}
		fmt.Printf("got ip: %s\n", ip)

		err = postNumber(client, number)
		if err != nil {
			return err
		}

		fmt.Printf("posted number %d\n", number)
	}

	return writer.Close()
}

func getIP(client http.Client) (string, error) {
	url := "https://httpbin.org/ip"

	res, err := client.Get(url)
	if err != nil {
		return "", err
	}

	var response struct {
		Origin string
	}
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return "", err
	}

	ip := response.Origin

	return ip, nil
}

func postNumber(client http.Client, number int) error {
	type Request struct {
		Number int `json:"number"`
	}

	request := Request{
		Number: number,
	}

	bodyJson, err := json.Marshal(request)
	if err != nil {
		return err
	}

	url := "https://httpbin.org/post"

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyJson))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	_, err = io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	return nil
}
