package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

/**
 * Struct called Currency used for parsing JSON
 * Structure derived from: https://exchangeratesapi.io/
 */
type CountryCode struct {
	Currencies Currency
}

type Currency struct {
	Code   string
	Name   string
	Symbol string
}

/*
 * This file contains code for the policy endpoint
 * This endpoint will be used to get information about what kind of policies
 * A given country has regarding the coronavirus pandemic
 * You can if desired get policies in a specified timeframe as well
 * @author Martin Iversen
 * @version 0.1
 * @date 09.03.2021
 */
//TODO Implement endpoint
//TODO Handle errors
func getPolicy(w http.ResponseWriter, r *http.Request) {
	//Defining variables
	vars := mux.Vars(r)
	country := vars["country_name"] //Country name from url
	startDate := vars["begin_date"] //Start date from url
	//endDate := vars["end_date"]            //End date from url
	Code := getCountryCode(w, r, country)
	url := "https://covidtrackerapi.bsg.ox.ac.uk/api/v2/stringency/actions/" + Code + "/" + startDate + ""
	body := invokeGet(w, r, url) //Invoking request
	fmt.Println(body)

}

/*
 * Method for getting a country code from a country name
 * This method uses the exchangerate api from assignment 1
 * https://exchangeratesapi.io/
 */
func getCountryCode(w http.ResponseWriter, r *http.Request, countryName string) string {
	//Defining variables
	url := "https://restcountries.eu/rest/v2/name/" + countryName + "?fields=currencies"

	body := invokeGet(w, r, url)

	var country = CountryCode{}
	//Defines an instance of the Country struct
	if err := json.Unmarshal([]byte(string(body)), &country); err != nil {
		http.Error(w, "Unmarshalling error", http.StatusBadRequest)
	}

	//Getting our desired variable from the currency struct
	return country.Currencies.Code
}
