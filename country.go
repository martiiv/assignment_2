package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

type CountryInfo struct {
	Country             string
	Population          int
	Sq_km_area          int
	Life_expectancy     string
	Elevation_in_meters int
	Continent           string
	Abbreviation        string
	Location            string
	Iso                 int
	Capital_city        string
	Dates               []string
}

type Date struct {
	Dates string
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
	client := &http.Client{}
	vars := mux.Vars(r)
	country := vars["country_name"]
	url := "https://covid-api.mmediagroup.fr/v1/history?country=" + country + "&status=Confirmed"
	fmt.Print(url)
	r, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		http.Error(w, "Error occurred when reading request", http.StatusBadRequest)
	}

	//Invokes request using the client
	res, err := client.Do(r)
	if err != nil {
		http.Error(w, "Error occurred handling request:", http.StatusConflict)
	}

	//Fetches json from the request
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		http.Error(w, "Error when parsing json response: ", http.StatusGone)
	}

	fmt.Fprint(w, string(body))
	if err != nil {
		http.Error(w, "Error ocurred when displaying content", http.StatusInternalServerError)
	}
}

func getDatesInCountry(w http.ResponseWriter, r *http.Request) {

}
