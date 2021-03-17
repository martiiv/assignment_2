package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type PolicyActions struct {
	PolicyActions StringencyData
}

type StringencyData struct {
	Date_value        string
	Country_code      string
	Confirmed         int
	Deaths            int
	Stringency_actual float64
	Stringency        float64
}

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

func formatOutput(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	countryName := vars["country_name"]
	startDate := vars["begin_date"] //Start date from url
	endDate := vars["end_date"]     //End date from url
	policy := getPolicy(w, r)

	//Formatting output as specified in assignment
	fmt.Fprintf(w, "country:"+countryName+"\n")
	fmt.Fprintf(w, "scope: "+startDate+"-"+endDate+"\n")
	fmt.Fprintf(w, "stringency:%v \n", policy.PolicyActions.Stringency)
	fmt.Fprintf(w, "trend:%v\n")
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
func getPolicy(w http.ResponseWriter, r *http.Request) PolicyActions {
	//Defining variables
	vars := mux.Vars(r)
	country := vars["country_name"] //Country name from url
	startDate := vars["begin_date"] //Start date from url
	//endDate := vars["end_date"]            //End date from url
	Code := getCountryCode(w, r, country)
	url := "https://covidtrackerapi.bsg.ox.ac.uk/api/v2/stringency/actions/" + Code + "/" + startDate + ""
	body := invokeGet(w, r, url) //Invoking request

	var policyInfo = PolicyActions{}
	err := json.Unmarshal([]byte(string(body)), &policyInfo)
	if err != nil {
		fmt.Println("error:", err)
	}
	return policyInfo
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
