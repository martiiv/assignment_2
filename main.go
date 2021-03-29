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
func main() {
	err := Init()
	if err != nil {
		fmt.Println("Error occurred when initializing the database!", err.Error())
	}
	startTime = time.Now()
	handle()
}
