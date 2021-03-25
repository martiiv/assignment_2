package main

/*
 * File containing all webhooks related functionality
 * Conatains the following functions:
 *									WebHookHandler 				For handling the request from user
 *									GetWebhookResponseObject 	For filling information into the JSONWebhook(see struct)
 *									GetDataStringency 			For getting stringency as of today for JSONWebhook
 *									GetDataConfirmed			For getting confirmed cases as of today for JSONWebhook
 * Martin Iversen
 * 25.03.2020
 * version 0.7
 */
import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	guuid "github.com/google/uuid" //Used for generating unique ID for webhook
	"io/ioutil"
	"net/http"
	"strconv"
	"time" //Used for getting the current date
)

//TODO implement timer to call function to check for change

//Struct for Json object which will get saved onto firebase
type JSONWebHook struct {
	Id guuid.UUID `json: "id"`
	WebhookRegistation
	Confirmed  int     `json:confirmed`
	Stringency float64 `json:stringency`
}

//Struct for the webhook info passed in by user
type WebhookRegistation struct {
	Url     string `json:"url"`     //URL to be invoked
	Timeout int64  `json:"timeout"` //How long between each change check
	Field   string `json:"field"`   //Which field to check for changes in (stringency or confirmed)
	Country string `json:"country"` //The country in question
	Trigger string `json:"trigger"` //ON CHANGE or ON TIMEOUT deciding when to look for updated information
}

var Key = "something"

var Secret []byte

var webHooks []WebhookRegistation

/*
 * Method WebHookHandler
 * Can register webhooks, view them or delete them
 * Uses function:
 *				GetWebhookResponseObject
 */
func WebHookHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodPost: //Case for registering webhook
		webHook := WebhookRegistation{}
		err := json.NewDecoder(r.Body).Decode(&webHook) //Gets information from post request
		if err != nil {
			http.Error(w, "Not able to decode http Request "+err.Error(), http.StatusBadRequest)
		}

		webHooks = append(webHooks, webHook) //Appends it to the WebhookRegistration object
		fmt.Println("Webhook " + webHook.Url + " has been registered")

		getWebhookResponseObject(w, r, webHook) //Creates a responseobject for DB saving(see struct for fields)

	case http.MethodGet: //Case for listing webhooks
		err := json.NewEncoder(w).Encode(webHooks)
		if err != nil {
			http.Error(w, "Not able to encode http Request "+err.Error(), http.StatusBadRequest)
		}
	case http.MethodDelete: //Case for deleting webhook
	}
}

/*
 * Method for creating a JSONWebHook object
 * Fills the object with information
 * Uses functions:
 *					getDataConfirmed() for confirmed cases today
 *					getDataStringency() for todays stringency trend
 */
func getWebhookResponseObject(w http.ResponseWriter, r *http.Request, hook WebhookRegistation) JSONWebHook {
	currentTime := time.Now()
	today := currentTime.Format("2006-01-02")

	webHookResponse := JSONWebHook{}
	webHookResponse.Id = guuid.New()                                          //Gives the object an unique ID
	webHookResponse.WebhookRegistation = hook                                 //Links the WebhookRegistration object with JSONWebhook object
	dataStringency := getDataStringency(w, r, webHookResponse.Country, today) //gets stringency data
	dataConfirmed := getWebhookDataConfirmed(w, r, webHookResponse.Country)   //gets confirmed cases
	webHookResponse.Stringency = dataStringency.StringencyData.Stringency_actual
	webHookResponse.Confirmed = dataConfirmed.All.Population

	fmt.Fprintf(w, "Id of webhook: %v \n", webHookResponse.Id)
	fmt.Fprintf(w, "Registered stringency:%v \n", webHookResponse.Stringency)
	fmt.Fprintf(w, "Registered confirmed cases: %v \n", webHookResponse.Confirmed)

	return webHookResponse
}

/*
 * Method for getting an object containing information about the stringency of a country
 * This function is almost indentical to one found in the policy file however this one takes in parameters
 */
func getDataStringency(w http.ResponseWriter, r *http.Request, countryName string, date string) Stringency {
	//Defining variables
	Code := getCountryCode(w, r, countryName)
	url := "https://covidtrackerapi.bsg.ox.ac.uk/api/v2/stringency/actions/" + Code + "/" + date + ""
	body := invokeGet(w, r, url) //Invoking request

	var webHookdata = Stringency{} //Object from policy file (see policy.go for struct)
	err := json.Unmarshal([]byte(string(body)), &webHookdata)
	if err != nil {
		fmt.Println("error:", err)
	}
	return webHookdata
}

/*
 * Function for getting the amount of confirmed cases for a country
 * This function is almost indentical to one found in the country file however this one takes in parameters
 */
func getWebhookDataConfirmed(w http.ResponseWriter, r *http.Request, countryName string) All {
	//Defining variables
	url := "https://covid-api.mmediagroup.fr/v1/history?country=" + countryName + "&status=Confirmed"
	body := invokeGet(w, r, url) //Invoking request

	var countryInfo = All{} //Object for unmarshalling
	err := json.Unmarshal([]byte(string(body)), &countryInfo)
	if err != nil {
		fmt.Println("error:", err)
	}

	return countryInfo
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
