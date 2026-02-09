package main

import "time"

var startup_time = time.Now()

func getRuntime() string {
	return time.Since(startup_time).String()
}
