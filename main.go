package main

import (
	"fmt"
	"time"
)

/*
 * The main file used for running the application
 *
 * @author Martin Iversen
 * @version 1.0
 * @date 29.03.2021
 */
//TODO Implement endpoint
//TODO Handle errors
func main() {
	err := Init()
	if err != nil {
		fmt.Println("Error occurred when initializing the database!")
	}
	startTime = time.Now()
	handle()
}
