package main

/*
 * This file contains code for the policy endpoint
 * This endpoint will be used to get information about what kind of policies
 * A given country has regarding the coronavirus pandemic
 * You can if desired get policies in a specified timeframe as well
 * It uses 3 functions:
 *  				getCountryCode()	for getting a given country's 3 digit code
 *					getPolicy()    		for getting stringency information
 * 					formatOutput() 		for formatting output
 * @author Martin Iversen
 * @version 0.8
 * @date 18.03.2021
 */
//TODO Implement trend value in formatOutput function
//TODO Handle errors
import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

/*
 * Structs for parsing json object from policy api:
 * Structure derived from https://covidtrackerapi.bsg.ox.ac.uk
 *
 */
type Stringency struct {
	StringencyData Data
}

type Data struct {
	Date_value        string
	Country_code      string
	Confirmed         int
	Deaths            int
	Stringency_actual float64
	Stringency        float64
}

/**
 * Structs for parsing Json object from restcountries api
 * Structure derived from: https://restcountries.eu/
 */
type CountryCode struct {
	Alpha3Code string
}

func formatOutput(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	countryName := vars["country_name"]
	startDate := vars["begin_date"] //Start date from url
	endDate := vars["end_date"]     //End date from url
	policy := getPolicy(w, r)
	fmt.Println(policy)

	//Formatting output as specified in assignment
	fmt.Fprintf(w, "country:"+countryName+"\n")
	fmt.Fprintf(w, "scope: "+startDate+"-"+endDate+"\n")
	fmt.Fprintf(w, "stringency:%v \n", policy.StringencyData.Stringency)
	fmt.Fprintf(w, "trend:%v\n")
}

/*
 * Method for getting a json object containing information about stringency trends
 * This method uses the https://covidtrackerapi.bsg.ox.ac.uk api
 * It returns a Stringency object with information
 */
func getPolicy(w http.ResponseWriter, r *http.Request) Stringency {
	//Defining variables
	vars := mux.Vars(r)
	startDate := vars["begin_date"] //Start date from url
	Code := getCountryCode(w, r)
	url := "https://covidtrackerapi.bsg.ox.ac.uk/api/v2/stringency/actions/" + Code + "/" + startDate + ""
	body := invokeGet(w, r, url) //Invoking request

	var policyInfo = Stringency{}
	err := json.Unmarshal([]byte(string(body)), &policyInfo)
	if err != nil {
		fmt.Println("error:", err)
	}

	return policyInfo
}

/*
 * Method for getting a country code from a country name
 * This method uses the restcountries api from assignment 1
 * https://exchangeratesapi.io/
 */
func getCountryCode(w http.ResponseWriter, r *http.Request) string {
	//Defining variables
	vars := mux.Vars(r)
	countryName := vars["country_name"] //Country name from url
	url := "https://restcountries.eu/rest/v2/name/" + countryName + "?fields=alpha3Code"

	body := invokeGet(w, r, url)

	var country []CountryCode
	//Defines an instance of the Country struct
	if err := json.Unmarshal([]byte(string(body)), &country); err != nil {
		log.Printf("error decoding sakura response: %v", err)
		if e, ok := err.(*json.SyntaxError); ok {
			log.Printf("syntax error at byte offset %d", e.Offset)
		}
		log.Printf("sakura response: %q", body)
	}

	return string(country[0].Alpha3Code)
}
