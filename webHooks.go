package main

/*
 * File containing all webhooks related functionality
 * Conatains the following functions:
 *									WebHookHandler 				For handling the request from user like:
 *																adding, deleting and displaying webhooks
 *
 *									SingleHandler				For handling specific webhooks
 *									GetWebhookResponseObject 	For filling information into the JSONWebhook(see struct)
 *									InfinityRunner				For checking an object for updated info using a timeot
 *									CallUrl						For invoking the webhook url and displaying updated info
 *																Based on trigger field from WebhookRegistration struct
 *
 *									GetDataStringency 			For getting stringency as of today for JSONWebhook
 *									GetDataConfirmed			For getting confirmed cases as of today for JSONWebhook
 *									Response					Formatting callUrl response based on field
 * Martin Iversen
 * 29.03.2021
 * version 1.0
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
	Id         string  //Same as database id
	Confirmed  int     //Confirmed cases
	Stringency float64 //Stringency value
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

var webHooks []WebhookRegistation //Initializing list for webhook
var jsonHooks []JSONWebHook       //Initializing list for information that will get stored in DB

/*
 * Method WebHookHandler
 * Can register webhooks, view them or delete them
 * Uses function:
 *				GetWebhookResponseObject()	To set object information
 *				infinityRunner()			To check for changes
 * 				AddWebhook()				For adding webhook to DB
 * 				update()					For setting ID
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
		jsonHooks = append(jsonHooks, entry)             //Appens the DB entry to the local list

		entry.Id, _ = AddWebhook(entry) //Sets the jsonHooks id to the DB id
		err = update(entry.Id, entry)   //Updates ID
		if err != nil {
			http.Error(w, "Error occurred when trying to update webhook!", http.StatusInternalServerError)
		}

		//Prints response
		fmt.Fprintf(w, "Id of registered webhook: %v \n", entry.Id)
		fmt.Fprintf(w, "Status code: %v \n", http.StatusCreated)
		infinityRunner(w, r, entry) //Runs check on timeout

	case http.MethodGet: //Case for listing webhooks
		list, err := GetAll()
		if err != nil {
			http.Error(w, "Error occurred when listing webhooks from databae", http.StatusInternalServerError)
		}

		for _, doc := range list {
			webhook := JSONWebHook{}
			if err := doc.DataTo(&webhook); err != nil {
				http.Error(w, "Tried to iterate through webhooks but failed", http.StatusBadRequest)
			}

			//Prints out webhook info from database
			fmt.Fprintf(w, "\n"+webhook.Id+"\n")
			fmt.Fprintf(w, webhook.Url+"\n")
			fmt.Fprintf(w, "%v\n", webhook.Timeout)
			fmt.Fprintf(w, webhook.Field+"\n")
			fmt.Fprintf(w, webhook.Country+"\n")
			fmt.Fprintf(w, webhook.Trigger+"\n")
		}
	}
}

/*
 * Function for getting a specific webhook or deleting a spesific webhook based on id
 * Uses function:
 *				DeleteWebhook() For deleting a specific webhook
 */
func singleHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"] //ID from url

	switch r.Method {
	case http.MethodGet:
		hook, _ := client.Collection(Collection).Doc(id).Get(ctx)
		fmt.Fprintf(w, "%v", hook.Data())

	case http.MethodDelete: //Case for deleting webhook
		err := DeleteWebhook(id)
		if err != nil {
			http.Error(w, "Error occurred when deleting webhook with id:"+id, http.StatusInternalServerError)
		}

		for i := range jsonHooks { //Should delete local webhooks although this isn't working properly
			if jsonHooks[i].Id == id {
				jsonHooks = append(jsonHooks[:i+1], jsonHooks[i+1:]...)
				fmt.Fprintf(w, "Deleted webhook with id:"+id)
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
	today := currentTime.Format("2006-01-02") //Sets date today

	webHookResponse := JSONWebHook{}
	dataStringency := getDataStringency(w, r, hook.Country, today) //gets stringency data

	dataConfirmed := getDataConfirmed(w, r, hook.Country) //gets confirmed cases
	webHookResponse.WebhookRegistation = hook             //Links the WebhookRegistration object with JSONWebhook object

	if dataStringency.StringencyData.Stringency_actual == 0 { //If the stringency actual data is 0 use stringency
		webHookResponse.Stringency = dataStringency.StringencyData.Stringency
	} else {
		webHookResponse.Stringency = dataStringency.StringencyData.Stringency_actual
	}
	webHookResponse.Confirmed = dataConfirmed.All.Population

	return webHookResponse //Returns a JSONWebHook object see struct
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
	time.Sleep(time.Duration(hook.Timeout) * time.Second)
	currentTime := time.Now()
	today := currentTime.Format("2006-01-02")

	StringencyActual := hook.Stringency
	Confirmed := hook.Confirmed
	newStringencyActual := getDataStringency(w, r, hook.Country, today).StringencyData.Stringency_actual //gets stringency data
	newConfirmed := getDataConfirmed(w, r, hook.Country).All.Population                                  //gets confirmed cases

	switch hook.Trigger {

	case "ON_CHANGE":
		if (StringencyActual != newStringencyActual) && (hook.Field == "stringency") { //If field is stringency it gets checked
			err := update(hook.Url, hook)
			if err != nil {
				http.Error(w, "Error occurred when trying to update webhook!", http.StatusInternalServerError)
			}
			callUrl(hook.Url, Response(hook))
			fmt.Fprintf(w, "Change occurred in stringency! \n")
			fmt.Fprintf(w, "New stringency value: %v \n", hook.Stringency)

		} else if (Confirmed != newConfirmed) && (hook.Field == "confirmed") { //If field is confirmed it gets checked
			err := update(hook.Url, hook)
			if err != nil {
				http.Error(w, "Error occurred when trying to update webhook!", http.StatusInternalServerError)
			}
			callUrl(hook.Url, Response(hook))
			fmt.Fprintf(w, "Change occurred in confirmed cases! \n")
			fmt.Fprintf(w, "New confirmed value: %v \n", hook.Confirmed)
		}

	case "ON_TIMEOUT":
		if hook.Field == "stringency" { //If field is stringency it gets checked
			err := update(hook.Url, hook)
			if err != nil {
				http.Error(w, "Error occurred when trying to update webhook!", http.StatusInternalServerError)
			}
			callUrl(hook.Url, Response(hook))

		} else if hook.Field == "confirmed" { //If field is confirmed it gets checked
			err := update(hook.Url, hook)
			if err != nil {
				http.Error(w, "Error occurred when trying to update webhook!", http.StatusInternalServerError)
			}
			callUrl(hook.Url, Response(hook))
		}
	}

	go infinityRunner(w, r, hook) //Calls the function again using go routines
}

/*
 * Method for invoking webhook url and notifying user
 */
func callUrl(url string, data interface{}) {
	content, err := json.Marshal(data)
	if err != nil {
		return
	}

	//Defines request
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(content))
	if err != nil {
		log.Printf("%v", "Error during request creation.")
		return
	}

	// Invokes request
	req.Header.Set("Content-Type", "application/json")
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Println("Error in HTTP request: " + err.Error())
		return
	}

	//Reads the request
	response, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("Something is wrong with invocation response: " + err.Error())
		return
	}

	//Prints response to console
	fmt.Println("Webhook invoked. Received status code " + strconv.Itoa(res.StatusCode) +
		" and body: " + string(response))
}

/*
 * Function for formatting response based on field
 * For stringency prints stringency and contry
 * For confirmed prints confirmed cases and country
 * Returns a JSONWebHook with the information
 */
func Response(hook JSONWebHook) JSONWebHook {
	switch hook.Field {
	case "stringency":
		var stringencyResponse JSONWebHook
		stringencyResponse.Country = hook.Country
		stringencyResponse.Stringency = hook.Stringency
		return stringencyResponse
	case "confirmed":
		var confirmedResponse JSONWebHook
		confirmedResponse.Country = hook.Country
		confirmedResponse.Confirmed = hook.Confirmed
		return confirmedResponse
	}
	return hook
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
		fmt.Fprintf(w, "error occurred when unmarshalling:%v", http.StatusBadRequest)
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
		fmt.Fprintf(w, "error occurred when unmarshalling:%v", http.StatusBadRequest)
	}
	return countryInfo
}
