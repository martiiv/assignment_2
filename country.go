package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type All struct {
	All Country
}
type Country struct {
	Confirmed  int
	Recovered  int
	Country    string
	Continent  string
	Population int
}

type Response struct {
}

/*
 * This file contains the country endpoint
 * This endpoint is used for getting the corona cases of a given country
 * In a certain specified timeframe (this is optional)
 * @author Martin Iversen
 * @version 0.1
 * @date 09.03.2021
 */
//TODO Implement endpoint
//TODO Implement date
//TODO Handle errors
//TODO implement other endpoiint for date functionality
func getCountryInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	country := vars["country_name"]
	//startDate := vars["begin_date"]
	//endDate := vars["end_date"]
	url := "https://covid-api.mmediagroup.fr/v1/cases?country=" + country
	body := invokeGet(w, r, url)

	var countryInfo = All{}
	err := json.Unmarshal([]byte(string(body)), &countryInfo)
	if err != nil {
		fmt.Println("error:", err)
	}

	fmt.Fprintf(w, "country:"+countryInfo.All.Country+"\n")
	fmt.Fprintf(w, "continent:"+countryInfo.All.Continent+"\n")
	fmt.Fprintf(w, "scope: total")
	fmt.Fprintf(w, "confirmed:%v\n", countryInfo.All.Confirmed)
	fmt.Fprintf(w, "recovered:%v\n", +countryInfo.All.Recovered)

}
