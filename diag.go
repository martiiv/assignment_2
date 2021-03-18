package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

var startTime time.Time

/*
 * This file contains code for the diag endpoint
 * This endpoint will provide a user with information regarding application diagnostics
 * @author Martin Iversen
 * @version 0.1
 * @date 09.03.2021
 */
//TODO Implement endpoint
//TODO Handle errors
func getDiagnostics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var casesCode int
	getCasesAPI, err := http.Get("https://mmediagroup.fr/covid-19") //Calls API
	if err != nil {                                                 //Error handling
		log.Printf("No/Bad response from API, %v", err)
		casesCode = 500
	} else {
		casesCode = getCasesAPI.StatusCode
		defer getCasesAPI.Body.Close()
	}

	var policyCode int
	getPolicyAPI, err := http.Get("https://covidtracker.bsg.ox.ac.uk") //Calls API
	if err != nil {                                                    //Error handling
		log.Printf("Something went wrong with the countries api, %v", err)
		policyCode = 500
	} else {
		policyCode = getPolicyAPI.StatusCode
		defer getPolicyAPI.Body.Close()
	}

	var countryCode int
	getCountriesAPI, err := http.Get("https://restcountries.eu") //Calls API
	if err != nil {                                              //Error handling
		log.Printf("Something went wrong with the countries api, %v", err)
		countryCode = 500
	} else {
		countryCode = getCountriesAPI.StatusCode
		defer getCountriesAPI.Body.Close()
	}

	fmt.Fprintf(w, `"CovidCasesAPI": "%v+","CovidPoliciesAPI": "%v"
							,"CountriesAPI:" "%v"
                              ,"version": "v1",
                               "uptime: "%v"`,
		casesCode, policyCode, countryCode, int64(time.Since(startTime).Seconds()))
}
