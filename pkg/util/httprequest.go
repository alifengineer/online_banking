package util

import (
	"crypto/tls"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func HttpRequest(req *http.Request, header map[string]string) (resp []byte, status_code int, err error) {

	log.Println(":::HttpRequest", req.URL, req.Method)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr, Timeout: time.Duration(10) * time.Second}

	for k, v := range header {
		req.Header.Add(k, v)
	}

	req.Header.Add("Content-Type", "application/json")
	clientDoRes, err := client.Do(req)
	if err != nil {
		return
	}
	defer clientDoRes.Body.Close()

	resp, err = ioutil.ReadAll(clientDoRes.Body)
	return resp, clientDoRes.StatusCode, err
}
