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
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time" //Used for getting the current date
)

//Struct for Json object which will get saved onto firebase
type JSONWebHook struct {
	Id         string  `json: "id"`
	Confirmed  int     `json:"confirmed"`
	Stringency float64 `json:"stringency"`
	WebhookRegistation
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
var jsonHooks []JSONWebHook

/*
 * Method WebHookHandler
 * Can register webhooks, view them or delete them
 * Uses function:
 *				GetWebhookResponseObject
 */
func WebHooksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodPost: //Case for registering webhook
		webHook := WebhookRegistation{}
		err := json.NewDecoder(r.Body).Decode(&webHook) //Gets information from post request
		if err != nil {
			http.Error(w, "Not able to decode http Request "+err.Error(), http.StatusBadRequest)
		}
		webHooks = append(webHooks, webHook)             //Appends it to the WebhookRegistration object
		entry := getWebhookResponseObject(w, r, webHook) //Creates a responseobject for DB saving(see struct for fields)
		jsonHooks = append(jsonHooks, entry)

		entry.Id, _ = AddWebhook(entry)
		update(entry.Id, entry)

		fmt.Fprintf(w, "Id of registered webhook: %v \n", entry.Id)
		fmt.Fprintf(w, "Status code: %v \n", http.StatusCreated)
		infinityRunner(w, r, entry)

	case http.MethodGet: //Case for listing webhooks
		list, err := GetAll()
		if err != nil {
			http.Error(w, "Error occurred when listing webhooks", http.StatusBadRequest)
		}
		for _, doc := range list {
			webhook := JSONWebHook{}
			if err := doc.DataTo(&webhook); err != nil {
				http.Error(w, "Tried to iterate through webhooks but failed", http.StatusBadRequest)
			}
			fmt.Fprintf(w, "%v\n", webhook)
		}
	}

}

func singleHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"] //ID from url

	switch r.Method {
	case http.MethodGet:
		for i := 0; i < cap(jsonHooks); i++ {
			if jsonHooks[i].Id == id {
				object, err := client.Collection(Collection).Doc(id).Get(ctx)
				if err != nil {
					http.Error(w, "error when trying to list object from database", http.StatusInternalServerError)
				}
				m := object.Data()
				fmt.Fprintf(w, "%v", m)
			}
		}

	case http.MethodDelete: //Case for deleting webhook
		for i := 0; i < cap(jsonHooks); i++ {
			if jsonHooks[i].Id == id {
				DeleteWebhook(jsonHooks[i].Id)
				jsonHooks = append(jsonHooks[:i], jsonHooks[i+1:]...)
				webHooks = append(webHooks[:i], webHooks[i+1:]...)
				fmt.Printf("Deleted Webhook with ID: %s", jsonHooks[i].Id)
			}
		}
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
	dataStringency := getDataStringency(w, r, hook.Country, today) //gets stringency data

	dataConfirmed := getDataConfirmed(w, r, hook.Country) //gets confirmed cases
	webHookResponse.WebhookRegistation = hook             //Links the WebhookRegistration object with JSONWebhook object

	if dataStringency.StringencyData.Stringency_actual == 0 {
		webHookResponse.Stringency = dataStringency.StringencyData.Stringency
	} else {
		webHookResponse.Stringency = dataStringency.StringencyData.Stringency_actual
	}
	webHookResponse.Confirmed = dataConfirmed.All.Population

	return webHookResponse
}

/*
 * Function for checking if information from the api has gotten updated
 * If the field trigger is:
 *							ON_CHANGE : The application will notify the user if there has been a change in the field
 *							ON_TIMEOUT: The application will notify the user when the timeout is 0 regardless of change
 * Uses functions:
 *				getDataStringency() for getting the latest stringency_actual value
 * 				getDataConfirmed() for getting the latest confirmed cases value
 */
func infinityRunner(w http.ResponseWriter, r *http.Request, hook JSONWebHook) {
	currentTime := time.Now()
	today := currentTime.Format("2006-01-02")

	StringencyActual := hook.Stringency
	Confirmed := hook.Confirmed
	newStringencyActual := getDataStringency(w, r, hook.Country, today).StringencyData.Stringency_actual //gets stringency data
	newConfirmed := getDataConfirmed(w, r, hook.Country).All.Population                                  //gets confirmed cases

	switch hook.Trigger {

	case "ON_CHANGE":
		if (StringencyActual != newStringencyActual) && (hook.Field == "stringency") { //If field is stringency it gets checked
			update(hook.Id, hook)
			callUrl(hook.Url, hook)
			fmt.Fprintf(w, "Change occurred in stringency! \n")
			fmt.Fprintf(w, "New stringency value: %v \n", hook.Stringency)

		} else if (Confirmed != newConfirmed) && (hook.Field == "confirmed") { //If field is confirmed it gets checked
			update(hook.Id, hook)
			callUrl(hook.Url, hook)
			fmt.Fprintf(w, "Change occurred in confirmed cases! \n")
			fmt.Fprintf(w, "New confirmed value: %v \n", hook.Confirmed)
		}

	case "ON_TIMEOUT":
		fmt.Fprintf(w, "Timeout reached! \n Checking for updated values in field %v  for webhook with ID: %q\n", hook.Field, hook.Id)
		if hook.Field == "stringency" { //If field is stringency it gets checked
			update(hook.Id, hook)
			callUrl(hook.Id, hook)

		} else if hook.Field == "confirmed" { //If field is confirmed it gets checked
			update(hook.Id, hook)
			callUrl(hook.Id, hook)
		}
	}

	time.Sleep(time.Duration(hook.Timeout) * time.Second)
	go infinityRunner(w, r, hook)
}

func callUrl(url string, data interface{}) {
	content, err := json.Marshal(data)
	if err != nil {
		return
	}

	// Prepare request
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(content))
	if err != nil {
		log.Printf("%v", "Error during request creation.")
		return
	}

	// Send request
	req.Header.Set("Content-Type", "application/json")
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Println("Error in HTTP request: " + err.Error())
		return
	}
	response, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("Something is wrong with invocation response: " + err.Error())
		return
	}

	fmt.Println("Webhook invoked. Received status code " + strconv.Itoa(res.StatusCode) +
		" and body: " + string(response))
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
func getDataConfirmed(w http.ResponseWriter, r *http.Request, countryName string) All {
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
