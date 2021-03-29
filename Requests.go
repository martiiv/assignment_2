package main

import (
	"io/ioutil"
	"net/http"
)

/*
 * File containing a method for handling a get request
 * Provided with more time this file would contain functions used frequently through the app
 * author Martin Iversen
 * date 29ø.03.2021
 * version 1.0
 */
func invokeGet(w http.ResponseWriter, r *http.Request, url string) []byte {
	w.Header().Set("Content-Type", "application/json")
	client := &http.Client{}

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

	return body
}
