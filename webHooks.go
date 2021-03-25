package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	guuid "github.com/google/uuid"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

//TODO Call policy request to get data (line 52), store it in data field in JSONWebhook struct, implement timer to call function to check for change

type JSONWebHook struct {
	Id guuid.UUID `json: "id"`
	WebhookRegistation
	Confirmed  int     `json:confirmed`
	Stringency float64 `json:stringency`
}

type WebhookRegistation struct {
	Url     string `json:"url"`
	Timeout int64  `json:"timeout"`
	Field   string `json:"field"`
	Country string `json:"country"`
	Trigger string `json:"trigger"`
}

var Key = "something"

var Secret []byte

var webHooks []WebhookRegistation

func WebHookHandler(w http.ResponseWriter, r *http.Request) {
	currentTime := time.Now()
	today := currentTime.Format("2006-01-02")
	fmt.Println(today)
	switch r.Method {
	case http.MethodPost:
		webHook := WebhookRegistation{}
		err := json.NewDecoder(r.Body).Decode(&webHook)
		if err != nil {
			http.Error(w, "Not able to decode http Request "+err.Error(), http.StatusBadRequest)
		}

		webHooks = append(webHooks, webHook)
		fmt.Println("Webhook " + webHook.Url + " has been registered")

		webHookResponse := JSONWebHook{}
		webHookResponse.Id = guuid.New()
		webHookResponse.WebhookRegistation = webHook
		dataStringency := getWebhookDataStringency(w, r, webHookResponse.Country, today)

		webHookResponse.Stringency = dataStringency.StringencyData.Stringency
		webHookResponse.Confirmed = data.StringencyData.Confirmed

		fmt.Println(data)
		fmt.Fprintf(w, "Id of webhook: %v", webHookResponse.Id)
		fmt.Fprintf(w, "Registered stringency:%v", webHookResponse.Stringency)
		fmt.Fprintf(w, "Registered confirmed cases: %v", webHookResponse.Confirmed)

	case http.MethodGet:
		err := json.NewEncoder(w).Encode(webHooks)
		if err != nil {
			http.Error(w, "Not able to encode http Request "+err.Error(), http.StatusBadRequest)
		}
	case http.MethodDelete:
	}
}

func getWebhookDataStringency(w http.ResponseWriter, r *http.Request, countryName string, date string) Stringency {
	//Defining variables
	Code := getCountryCode(w, r, countryName)
	url := "https://covidtrackerapi.bsg.ox.ac.uk/api/v2/stringency/actions/" + Code + "/" + date + ""
	body := invokeGet(w, r, url) //Invoking request

	var webHookdata = Stringency{}
	err := json.Unmarshal([]byte(string(body)), &webHookdata)
	if err != nil {
		fmt.Println("error:", err)
	}
	return webHookdata
}

func getWebhookDataConfirmed(w http.ResponseWriter, r *http.Request) {

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
