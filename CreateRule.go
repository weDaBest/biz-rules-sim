package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
)

//TEMP create rule request
//Structure is array of updates. The array contains a map.
//The map will be structured as follows: field that needs to be updated, and the new/update value
type createReq struct {
	NewRuleData []map[string]interface{} `json:"newRuleData,required"`
	/*	
		  the rule to create
		  the keys in the rule
		  the values	
		{
		  	Organizaton: DHS
			Unit: CBP
			Subunit: AIREntry ----> 1 map = 1 row
			Activity: entry
			Service: Add derogatory Information
		},
		{
		  	Organizaton: DHS
			Unit: CBP
			Subunit: AIREntry ----> 1 map = 1 row
			Activity: entry
			Service: Delete Identity Flag
		}
	*/	
}

//Creates an entirely new rule
func CreateRule(w http.ResponseWriter, r *http.Request) {
	query := mux.Vars(r)
	rule := query["name"]
	
	//checks if rule if in url
	if len(rule) == 0  {
		err := errors.New("rule must be specified")
		fmt.Println(err)
		w.WriteHeader(400)
		fmt.Fprint(w, err)
		return 
	}

	//Try to Open Rule, if the rule exists then throw an error
	rulename := rule + ".csv"
	csvFile, err := os.Open(rulename)
	if err == nil {
		errStr := fmt.Sprintf("%+s already exist cannot create a duplicate rule", rule)
		csvFile.Close()
		w.WriteHeader(400)
		fmt.Fprint(w, errStr)	
		return
	}	
	defer csvFile.Close() 

	//create the new file
	fmt.Println("Creating a new rule")
	newRule, err := os.Create(rulename)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprint(w, err)
		//err cleanup process
		go ErrCleanUp(rule, newRule)	
		return
	}
	defer newRule.Close()

	//if parameters are present unmarshall the request
	var inCreatetReq createReq
	inCreatetReq.NewRuleData = make([]map[string]interface{}, 0)	
	inBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		fmt.Fprint(w, err)
		//err cleanup process
		go ErrCleanUp(rule, newRule)
		return 
	}

	err = json.Unmarshal(inBytes, &inCreatetReq)
	if err != nil {
		fmt.Printf("Failed to unmarshall: %+s\n" , err.Error())
		w.WriteHeader(500)
		fmt.Fprint(w, err)
		//err cleanup process
		go ErrCleanUp(rule, newRule)
		return 
	}
	PrintMapArray(inCreatetReq.NewRuleData)
	
	fmt.Println(inCreatetReq.NewRuleData)
	for _, val := range inCreatetReq.NewRuleData {
		fmt.Println(val)
	}


	//determine the columns for the rule
	keys := make([]string, 0, len(inCreatetReq.NewRuleData[0]))
	for k := range inCreatetReq.NewRuleData[0] {
		keys = append(keys, k)
	}
	fmt.Printf("columns for new rule: %+v\n" , keys)

	//iterate over the data an populate the column data
	for _, singMap := range inCreatetReq.NewRuleData {
		for _, val := range singMap {
			fmt.Println(val) 
		}
	}
	
	//write the corresponding data to the csv file
	//creating a new wrtier object
	writer := csv.NewWriter(newRule)
	defer writer.Flush()
	
	//writes updated rules to new file
	values := [][]string{}
	//first append the columns now append the actual values
	values = append(values, keys)

	
	//iterate over the unmarshalled request
	for _, maps := range inCreatetReq.NewRuleData {
		//tmp structure to hold each row of the the csv file
		tmpString := []string{}
		//iterate over the individual maps in the final response array
			for _, col := range values[0] {
				fmt.Printf("looking for value for key: %+v", col)
				//iterate over the map to find the value for key
				for key, mapVal := range maps {
					if strings.Compare(key, col) == 0 {
						fmt.Printf("found %+v for key: %+v\n", mapVal, col)
						tmpString = append(tmpString, fmt.Sprint(mapVal))
						fmt.Printf("value kv appended: %+v\n",tmpString)
					}
				}
			}
			
		//after populating row add to array 
		values = append(values, tmpString)
		fmt.Println(values)
	}
	
	fmt.Println("writing new rule data")
	for _, stringData := range values {
		err = writer.Write(stringData)
		if err != nil {
			fmt.Printf("Failed to write new rule data: %+v\n",  err)
			w.WriteHeader(500)
			fmt.Fprint(w, err)
			//err cleanup process
			go ErrCleanUp(rule, newRule)
			return
		}
	}
	fmt.Printf("Successfuly created rule:%+v\n" , rule)
	w.WriteHeader(200)
	fmt.Fprintf(w, "Successfully created rule: %+s", rule)
	return
}
