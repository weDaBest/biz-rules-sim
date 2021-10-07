package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

//TODO:add option to delete specific row from csv file

func DeleteRule(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*** Deleting Rule***")
	query := mux.Vars(r)
	rule := query["name"]
	
	//checks if rule if in url
	if len(rule) == 0  {
		err := errors.New("rule must be specified")
		w.WriteHeader(400)
		fmt.Fprint(w, err)
		return 
	}
	fmt.Printf("Rule requested to be deleted: %+s\n", rule)

	//delete rule
	ruleToDelete := rule + ".csv"
	err := os.Remove(ruleToDelete)
	if err != nil {
		fmt.Printf("Failed to delete rule: %+s\n", ruleToDelete)
		w.WriteHeader(500)
		fmt.Fprint(w, err.Error())
		return
	}

	fmt.Printf("Successfully deleted rule: %+s\n" , rule)
	w.WriteHeader(200)
	fmt.Fprintf(w,"Successfully deleted rule: %+s\n" , rule)
}