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
	fmt.Println("Starting")
	Init()

	startTime = time.Now()
	fmt.Println("initialized handler")
	handle()
}
