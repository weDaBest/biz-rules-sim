package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

//TODO: add get for specific OrgCode

/*
	Gets all data for specific rule
	specific rule will be in the URL params
	response will be a json/string
*/
func GetRule(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*** Getting Rule ***")
	query := mux.Vars(r)
	rule := query["name"]

	
	//checks if rule if in url
	if len(rule) == 0 {
		err := errors.New("rule must be specified must be specified")
		w.WriteHeader(400)
		fmt.Fprint(w, err)
		return 
	}


	fmt.Printf("Rule requested: %+s\n", rule)

	//Retrieve Data from the CSV file
	rulename := rule + ".csv"
	csvFile, err := os.Open(rulename)
	if err != nil {
    	fmt.Println(err)
		w.WriteHeader(404)
		fmt.Fprint(w, err)
		return	
	}

	fmt.Println("Successfully Opened CSV file")
	defer csvFile.Close()

	csvLines, err := csv.NewReader(csvFile).ReadAll()
    if err != nil {
        fmt.Println(err)
		w.WriteHeader(500)
		fmt.Fprint(w, err)
		return
    } 

	keys := []string{}
	iter := 0
	finalResp := []map[string]interface{}{}
    for _, line := range csvLines {
		if iter == 0 {
			keys = append(keys, line...)
			fmt.Printf("keys in rule %+v" , keys)
			iter = 1
			continue
		}

		tmp := map[string]interface{}{}
		for index, value := range line {	
			tmp[keys[index]] = value
		}

		finalResp = append(finalResp, tmp)
		
    }
	fmt.Println(finalResp)	

	w.WriteHeader(200)
	b, _ := json.MarshalIndent(finalResp, "", " ")
	fmt.Fprint(w, string(b))
}