package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

/*
 * The handler file, will be used for implementing handlers for each endpoint
 * @author Martin Iversen
 * @version 0.1
 * @date 09.03.2021
 */
func handle() {
	r := mux.NewRouter()
	r.HandleFunc("/corona/v1/country/{country_name}/", getCountryInfo)
	r.HandleFunc("/corona/v1/policy/{:country_name}{?begin_date-end_date}", getPolicy)
	r.HandleFunc("/corona/v1/diag/", getDiagnostics)
	r.HandleFunc("/corona/v1/notifications/", getNotification)
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
