package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

/*
 * The handler file, will be used for implementing handlers for each endpoint
 * Contains endpoint:
 *					country: For getting information about corona cases in a given country between dates
 *					policy: For getting information about corona policies in a given country between dates
 *					diagnostics: For displaying application diagnostics
 *					notifications: For invoking, registering, viewing and deleting webhooks
 * @author Martin Iversen
 * @version 1.0
 * @date 29.03.2021
 */
func handle() {

	r := mux.NewRouter()
	r.HandleFunc("/corona/v1/country/{country_name}/{begin_date}/{end_date}", formatResponse) //country endpoint
	r.HandleFunc("/corona/v1/policy/{country_name}/{begin_date}/{end_date}", formatOutput)    //Policy endpoint
	r.HandleFunc("/corona/v1/diag", getDiagnostics)                                           //Diagnostics enpoint
	r.HandleFunc("/corona/v1/notifications/", WebHooksHandler)                                //Webhook endpoint 1
	r.HandleFunc("/corona/v1/notifications/{id}/", singleHandler)                             //Webhook endpoint 2
	log.Fatal(http.ListenAndServe(getPort(), r))
}

/*
 * Function for setting the port
 * Takes no parameters and returns the port the application is listening to
 * Will use 8080 localhost for testing
 */
func getPort() string {
	var port = os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return ":" + port
}
