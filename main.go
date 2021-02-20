package main

import (
	"fmt"
	"net/http"

	"persons.com/api/infrastructure/server"
)

func main() {

	errs := make(chan error, 2)
	go func() {
		fmt.Println("Listening on port :5000")
		errs <- http.ListenAndServe(server.HttpPort(), server.StartRouter())

	}()

	<-errs
}

//app flow: Domain -> Service -> useCases -> Repository -> Serializers(json, messagePack, grpc, soap, etc) -> Handlers(controllers) -> Transporter(http, websockets, GraphQl etc.)
