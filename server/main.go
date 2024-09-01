package server

import {
	"log"
}

func main() {
	// Wrapping our server with the *net/http.Server in NewHTTPServer saved us from writing a bunch of code 
	// here—and anywhere else we’d create an HTTP server.
	srv  := server.NewHTTPServer(":8080") //pass in address to listen on
	log.Fatal(srv.ListenAndServe())
}