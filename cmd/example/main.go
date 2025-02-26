package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path"

	harwriter "github.com/oliverroer/go-har/writer"
)

func main() {
	err := example(3)
	if err != nil {
		log.Fatal(err)
	}
}

func example(sessionCount int) error {
	entryFiles := make([]string, 0, sessionCount)

	for number := range sessionCount {
		name, err := recordSession(number)
		if err != nil {
			return err
		}
		entryFiles = append(entryFiles, name)
	}

	dir := "out"
	perm := fs.FileMode(0750)
	err := os.MkdirAll(dir, perm)
	if err != nil {
		return err
	}

	name := harwriter.DefaultName()

	harPath := path.Join(dir, name+".har")

	return harwriter.EntriesToHar(harPath, entryFiles...)
}

func recordSession(number int) (string, error) {
	dir := "out"
	perm := fs.FileMode(0750)
	err := os.MkdirAll(dir, perm)
	if err != nil {
		return "", err
	}

	name := path.Join(dir, harwriter.DefaultName()+".jsonl")

	writer, err := harwriter.Open(name)
	if err != nil {
		return name, err
	}

	transport := writer.RoundTripper(http.DefaultTransport)
	client := http.Client{
		Transport: transport,
	}

	ip, err := getIP(client)
	if err != nil {
		return "", err
	}
	fmt.Printf("got ip: %s\n", ip)

	err = postNumber(client, number)
	if err != nil {
		return "", err
	}

	fmt.Printf("posted number %d\n", number)

	err = writer.Close()
	if err != nil {
		return "", err
	}

	return name, nil
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
