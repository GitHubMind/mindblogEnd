package lib

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func Post(urlStr string, json []byte) (body []byte, status int) {
	log.Println(string(json))
	req, err := http.NewRequest("POST", urlStr, bytes.NewBuffer(json))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{
		//Transport: netTransport,
		Timeout: 1 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("client", resp)
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	status = resp.StatusCode
	return
}
func Get(urlStr string, json []byte) (body []byte, status int) {
	log.Println(string(json))
	req, err := http.NewRequest("GET", urlStr, bytes.NewBuffer(json))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{
		//Transport: netTransport,
		Timeout: 1 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("client", resp)
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	status = resp.StatusCode
	return
}
