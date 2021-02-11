package main

import (
	"fmt"
	"net/http"

	"persons.com/api/app"
)

func main() {

	errs := make(chan error, 2)
	go func() {
		fmt.Println("Listening on port :5000")
		errs <- http.ListenAndServe(app.HttpPort(), app.StartRouter())

	}()

	<-errs
}
