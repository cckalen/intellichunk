package api

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func Run_api() {

	var router = mux.NewRouter()
	const port string = ":8080"

	//router.HandleFunc("/version", Version).Methods("GET")
	router.HandleFunc("/conversation", ConversationHandler).Methods("GET")
	router.HandleFunc("/intellichunk/add", IntellichunkHandler).Methods("POST")

	log.Println("Server listening on port ", port)
	handler := cors.Default().Handler(router)

	//the program will exit if there is an error starting the server and print the error message
	log.Fatal(http.ListenAndServe(port, handler))

}
