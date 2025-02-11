package main

import (
	"fmt"
	transport "go-jira/Transport"
	"net/http"
)

func main() {

	transport.HttpHandler()
	fmt.Println("Start")
	http.ListenAndServe(":8082", nil)
	fmt.Println("End")

}
