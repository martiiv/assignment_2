package main

/*
 * This file contains the country endpoint
 * This endpoint is used for getting the corona cases of a given country
 * In a certain specified timeframe
 * It uses 3 functions:
 *					getCountryConfirmed() for getting confirmed cases
 *					getCountryRecovered() for getting recovered cases
 *					formatResponse() for formatting output
 * @author Martin Iversen
 * @version 1.0
 * @date 29.03.2021
 */
import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"math"
	"net/http"
)

/*
 * Defined structs for json parsing
 */
type All struct {
	All Country
}
type Country struct {
	Country    string
	Population int
	Continent  string
	Dates      map[string]int
}

/*
 * This method is used for formatting output data
 * This method uses 2 functions:
 *							getCountryConfirmed() for getting confirmed cases
 *							getCountryRecovered() for getting recovered cases
 */
func formatResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//Defining variables
	vars := mux.Vars(r)
	confirmed := getCountryConfirmed(w, r) //Object for confirmed cases
	recovered := getCountryRecovered(w, r) //Object for recovered cases
	startDate := vars["begin_date"]        //Start date from url
	endDate := vars["end_date"]            //End date from url

	//Calculates the values to be displayed
	confirmedInterval := float64(confirmed.All.Dates[endDate] - confirmed.All.Dates[startDate])
	recoveredInterval := recovered.All.Dates[endDate] - recovered.All.Dates[startDate]
	TotalPopulation := float64(confirmed.All.Population)
	populationPercentage := float64(confirmedInterval/TotalPopulation) * 100

	//Formatting output as specified in assignment
	fmt.Fprintf(w, "country:"+confirmed.All.Country+"\n")
	fmt.Fprintf(w, "continent:"+confirmed.All.Continent+"\n")

	if startDate == "" || endDate == "" { //If scope isn't specified
		fmt.Fprintf(w, "scope: total") //Sets the value to "total"
	} else {
		fmt.Fprintf(w, "scope: "+startDate+"-"+endDate+"\n")
	}
	fmt.Fprintf(w, "confirmed:%v \n", confirmedInterval)
	fmt.Fprintf(w, "recovered:%v\n", recoveredInterval)
	fmt.Fprintf(w, "population_percentage:%v\n", math.Ceil(populationPercentage*100)/100)
}

/*
 * This method is used for taking a http request and unmarshalling it
 * To get data from a request from api: covid-api.mmediagroup
 * URL to be invoked: https://covid-api.mmediagroup.fr/v1/history?country={countryName}&status=Recovered"
 * This method uses 1 function:
 *								invokeGet() for taking a REST GET request
 *											and returning the response in []byte format
 */
func getCountryRecovered(w http.ResponseWriter, r *http.Request) All {
	//Defining variables
	vars := mux.Vars(r)
	country := vars["country_name"] //Country name from url
	url := "https://covid-api.mmediagroup.fr/v1/history?country=" + country + "&status=Recovered"
	body := invokeGet(w, r, url) //Invoking request

	var countryInfo = All{} //Defining structure of object for unmarshalling
	err := json.Unmarshal([]byte(string(body)), &countryInfo)
	if err != nil {
		fmt.Fprintf(w, "error occured when unmarshalling request: %v", http.StatusBadRequest)
	}
	return countryInfo //Returns an object of type All{}
}

/*
 * This method is used for taking a http request and unmarshalling it
 * To get data from a country object
 * This object contains confirmed covid-19 cases
 * Url to be invoked: https://covid-api.mmediagroup.fr/v1/history?country={countryName}&status=Confirmed
 * This method uses 1 function:
 *								invokeGet() for taking a REST GET request
 *											and returning the response in []byte format
 */

func getCountryConfirmed(w http.ResponseWriter, r *http.Request) All {
	//Defining variables
	vars := mux.Vars(r)
	country := vars["country_name"] //Country name from url
	url := "https://covid-api.mmediagroup.fr/v1/history?country=" + country + "&status=Confirmed"
	body := invokeGet(w, r, url) //Invoking request

	var countryInfo = All{} //Object for unmarshalling
	err := json.Unmarshal([]byte(string(body)), &countryInfo)
	if err != nil {
		fmt.Fprintf(w, "Erro occurred when unmarshalling request: %v", http.StatusBadRequest)
	}
	return countryInfo
}
