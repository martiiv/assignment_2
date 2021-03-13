package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

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
func getCountryInfo(w http.ResponseWriter, r *http.Request) {
	client := &http.Client{}
	vars := mux.Vars(r)
	country := vars["country_name"]
	//date := vars["?begin_date-end_date"]
	url := "https://covid-api.mmediagroup.fr/v1/cases?country=" + country + ""
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
	output, err := ioutil.ReadAll(res.Body)
	if err != nil {
		http.Error(w, "Error when parsing json response: ", http.StatusGone)
	}

	fmt.Fprint(w, string(output))
	if err != nil {
		http.Error(w, "Error ocurred when displaying content", http.StatusInternalServerError)
	}

}
