package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

type WebhookRegistation struct {
	Url     string
	Timeout int
	Field   string
	Country string
	Trigger string
}

var Key = "something"

var Secret []byte

var webHooks []WebhookRegistation

func WebHookHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		webHook := WebhookRegistation{}
		err := json.NewDecoder(r.Body).Decode(&webHook)
		if err != nil {
			http.Error(w, "Not able to decode http Request "+err.Error(), http.StatusBadRequest)
		}

		webHooks = append(webHooks, webHook)
		fmt.Println("Webhook " + webHook.Url + " has been registered")

	case http.MethodGet:
		err := json.NewEncoder(w).Encode(webHooks)
		if err != nil {
			http.Error(w, "Not able to encode http Request "+err.Error(), http.StatusBadRequest)
		}
	case http.MethodDelete:
	}
}

func ServiceHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		fmt.Println("Received POST request...")
		for _, v := range webHooks {
			go CallUrl(v.Url, "Trigger event: Call to service endpoint with method "+v.Field)
		}
	default:
		http.Error(w, "Invalid method "+r.Method, http.StatusBadRequest)
	}
}

func CallUrl(url string, content string) {
	fmt.Println("invoking url:" + url)

	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader([]byte(content)))
	if err != nil {
		fmt.Errorf("%v", "Error during request invocation")
		return
	}

	hash := hmac.New(sha256.New, Secret)
	_, err = hash.Write([]byte(content))
	if err != nil {
		fmt.Errorf("%v", "Error occurred during hashing")
		return
	}

	request.Header.Add(Key, hex.EncodeToString(hash.Sum(nil)))

	client := http.Client{}
	res, err := client.Do(request)
	if err != nil {
		fmt.Println("Error when doing request" + err.Error())
	}
	response, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error when reading from response" + err.Error())
	}

	fmt.Println("Webhook invoked, status code: " + strconv.Itoa(res.StatusCode) + "body: " + string(response))
}
