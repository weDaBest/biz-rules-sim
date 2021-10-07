package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

//TODO: may need to add put to Update and entire row
func main(){

	fmt.Println("Inside of the main of the biz-rules-simulator")
	r := mux.NewRouter()
	
	//creates a new rule csv and 
	r.HandleFunc("/rule/{name}", CreateRule).Methods("POST") 
	
	//gets data for specific rule
	r.HandleFunc("/rule/{name}", GetRule).Methods("GET") 
	
	//updates an existing rule csv, adds a new row to the csv file
	r.HandleFunc("/rule/{name}", UpdateRule).Methods("PUT")
	
	//partially updates an row in an existing rule csv
	r.HandleFunc("/rule/{name}", UpdateRule).Methods("PATCH")
	
	//deletes an entire rule csv
	r.HandleFunc("/rule/{name}", DeleteRule).Methods("DELETE") 

	http.ListenAndServe(":80", r)

}