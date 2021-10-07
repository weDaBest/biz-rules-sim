package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

//TODO: differentiate between PATCH/PUT

//update rule request
//Structure is array of updates. The array contains a map.
//The map will be structured as follows: field that needs to be updated, and the new/update value
type updateReq struct {
	Updates []map[string]interface{} `json:"updates,required"`
	/*	
		  the rule to update
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


var versionNum int
//Modifies an existing rule
func UpdateRule(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*** Updating Rule ***")
	query := mux.Vars(r)
	rule := query["name"]
	
	
	//determine the method of update request: PUT/PATCH
	//PUT - adds a new row/updates existing row to existing CSV file
	//PATCH - partially updates an existing row in CSV FILE
	uptType := r.Method
	fmt.Printf("Update Request method: %+s\n" , uptType)
	
	
	//checks if rule if in url
	if len(rule) == 0  {
		w.WriteHeader(400)
		fmt.Fprint(w, errors.New("rule must be specified"))
		return 
	}
	fmt.Printf("Rule sheet requested: %+s\n", rule)
	
	//Try to Open Rule
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
	
	//tries to reads request body 
	var inUpdtReq updateReq
	inUpdtReq.Updates = make([]map[string]interface{}, 0)
	inBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		fmt.Fprint(w, err)
		csvFile.Close()
		return 
	}

	//Unmarshalls request body into json Request structure
	err = json.Unmarshal(inBytes, &inUpdtReq)
	if err != nil {
		fmt.Printf("Failed to unmarshall: %+s\n" , err.Error())
		w.WriteHeader(500)
		fmt.Fprint(w, err)
		csvFile.Close()
		return 
	}
	PrintMapArray(inUpdtReq.Updates)
	
	//check to see if request is nil, cannot be an 
	if inUpdtReq.Updates == nil{
		w.WriteHeader(400)
		fmt.Fprint(w, "Invalid request")
		csvFile.Close()
		return
	}

	//start reading
	csvLines, err := csv.NewReader(csvFile).ReadAll()
    if err != nil {
        fmt.Println(err)
		w.WriteHeader(500)
		fmt.Fprint(w, err)
		csvFile.Close()
		return
    } 
	
	//reads the csv file and stores into object
	//iterate over the first line to get the keys/column(s)
	keys := []string{}
	iter := 0
	finalResp := []map[string]interface{}{}
    for _, line := range csvLines {
		if iter == 0 {
			//fmt.Println(line)
			keys = append(keys, line...)
			//fmt.Printf("keys in rule %+v\n" , keys)
			iter = 1
			continue
		}

		//sets value for given key for tmp array
		tmp := map[string]interface{}{}
		for index, value := range line {
			tmp[keys[index]] = value
		}

		//append succesfully assemble row/tmp object to the final response array
		finalResp = append(finalResp, tmp)
    }
	fmt.Printf("Before update # of rows: %+d\n" , len(finalResp))

	//iterate over the update request array of maps
	for _, val := range inUpdtReq.Updates {
		//create a new temp object to hold data
		tmp := map[string]interface{}{}
		//iterates over the actual single map in the update: the keys and values
		for key, mapVal := range val {
			fmt.Printf("update values -- Key: %+v | Value: %+v\n" , key,mapVal)
			//see if key in request is valid, found in key retrieve from csv file
			for _, columns := range keys {
				if strings.Compare(key, columns) == 0{
					//if found add to new temp object
					fmt.Printf("Key: %+s | value: %+s \n",key,mapVal )
					tmp[key] = mapVal
					break
				}
			}
		
		}
		finalResp = append(finalResp, tmp)
	}
	fmt.Printf("final rule sheet: %+v\n" , finalResp)
	fmt.Printf("Post update # of rows: %+d\n" , len(finalResp))
	PrintMapArray(finalResp)
	finalResp = RemoveDupMaps(finalResp)

	//first close the old file
	err = csvFile.Close()
	if err != nil {
		fmt.Printf("Failed to close old rule file %+v\n", err)
		w.WriteHeader(500)
		fmt.Fprint(w, err)
		return
	}

	//TODO: add cleanup up process to remove extermely old versions of the file
	//renames old file before writing new one 
	//Keeps older version of file
	versionNum++
	newName := rule + "_" + time.Now().Format(("01-02-2006"))+"_"+ strconv.Itoa(versionNum) + ".csv"
	fmt.Printf("Renaming file to: %+s\n", newName)
	err = os.Rename(rulename, newName)
	if err != nil {
		fmt.Printf("Failed to rename an old rule file %+v\n", err)
		w.WriteHeader(500)
		fmt.Fprint(w, err)
		return 
	}

	//creating a new file
	fmt.Printf("Creating new version of rule: %+v\n" , rulename)
	updateFile, err := os.Create(rulename)
	if err != nil {
		fmt.Printf("Failed to create a new file %+v\n", err)
		w.WriteHeader(500)
		fmt.Fprint(w, err)
		return 
	}
	defer updateFile.Close()

	//creating a new wrtier object
	writer := csv.NewWriter(updateFile)
	defer writer.Flush()

	//writes updated rules to new file
	values := [][]string{}
	//first append the columns now append the actual values
	values = append(values, keys)
	fmt.Println("==================")
	fmt.Printf("first row of the csv file/the keys: %+v\n", values[0])
	fmt.Printf("# of columns: %+v\n" , len(values[0]))
	fmt.Println("==================")

	//iterate over the final response
	for _, maps := range finalResp {
		//tmp structure to hold each row of the the csv file
		tmpString := []string{}
		//iterate over the individual maps in the final response array
		for _, col := range values[0] {
			fmt.Printf("looking for value for key: %+v\n", col)
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
	
	//writes data to file
	fmt.Println("writing new rule data")
	for _, stringData := range values {
		err = writer.Write(stringData)
		if err != nil {
			fmt.Printf("Failed to write to new file: %+v\n",  err)
			w.WriteHeader(500)
			fmt.Fprint(w, err)
			return
		}
	}
	

	w.WriteHeader(200)
	fmt.Fprintf(w, "Succesfully update rule %+s\n" ,rule)
}
