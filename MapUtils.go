package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

func GetKeys(data map[string]interface{}) []string {
	j := 0
	keys := make([]string, len(data))
	for  k := range data {
		keys[j] = k
		j++
	}
	return keys
}

//KeyExists - checks to see if map key exists
func KeyExists( data map[string]interface{}, key string) (bool, interface{}){
	fmt.Println("------ Single Map ------")
	PrintMap(data)
	fmt.Println("------ key ------")
	fmt.Println(key)
	if val, ok := data["key"]; ok {
		fmt.Printf("Key : %+s | Value : %+v\n" , key, val)
		return true, val
	}
	return false, nil
}

//KVExists - checks to see if specific key with specific value exists in map
func KVExists (data map[string]interface{}, srchKey, srchVal string) bool {
	fmt.Println("Map passed in")
	PrintMap(data)

	fmt.Printf("Key : %+s | Value :%+s passed in\n", srchKey,srchVal)
	exists, mapVal := KeyExists(data, srchKey)
	if exists && strings.Compare(srchVal, mapVal.(string)) == 0 {
		return true
	}
	return false
}

//KVExistsArray - checks to see if specific map exists in array of maps
func KVExistsArray (srcMap []map[string]interface{}, srchVal map[string]interface{}) bool {
	fmt.Println("Search Map")
	PrintMapArray(srcMap)

	fmt.Println("Search Val")
	PrintMap(srchVal)
	for _, singleMap := range srcMap {
		res := reflect.DeepEqual(singleMap, srchVal)
		if res {
			fmt.Println("Map exists in array of Maps")
			return true
		}
	}	
	return false
}



//PrintMap - prints outs all the keys and values for single map
func PrintMap(data map[string]interface{}) {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
    		fmt.Println("error:", err)
	}
	fmt.Println(string(b))
}

//PrintMapArray - prints out all the keys and values in a for an array of maps
func PrintMapArray (data []map[string]interface{}) {
	b := make([]byte, 0)
	var err error
	for range data {
		b, err = json.MarshalIndent(data, "", "\t")
		if err != nil {
    		fmt.Println("error:", err)
		}
	}
	fmt.Println(string(b))
}

//RemoveDupMaps - removes any duplicate maps from an array of maps
func RemoveDupMaps(data []map[string]interface{}) []map[string]interface{} {
	fmt.Println("array of Maps Passed in")
	PrintMapArray(data)
	var noDups = make([]map[string]interface{},0)
	var genKeyList = make([]string, 0)
	var occurred = make(map[string]bool)
	
	for _, singleMap := range data {
		//generates a unique key for the map
		var checkKey = ""
		for _, vals := range singleMap {
			checkKey = checkKey+"_"+vals.(string) 
		}
		//creates array of all keys 
		genKeyList = append(genKeyList, checkKey)
	}
	fmt.Printf("List of all keys: %+v\n" , genKeyList)

	for _, singleMap := range data {
		//generate a unique key
		var genKey = ""
		for _, vals := range singleMap {
			genKey = genKey+"_"+vals.(string)
		}
		fmt.Println("=========")
		fmt.Println(genKey)
		fmt.Println(occurred)
		fmt.Println("=========")
		if !occurred[genKey]  {
			occurred[genKey] = true
			noDups = append(noDups, singleMap)
			PrintMapArray(noDups)
		}


		//iterates over array of map keys that have appeared
		// for _, curKey := range genKeyList {
		// 	fmt.Printf("curKey : %+v | genKey : %+v\n" , curKey, genKey)
		// 	if strings.Compare(curKey, genKey) != 0 {
		// 		fmt.Printf("Have not seen key: %+v\n", curKey)
		// 		noDups = append(noDups, singleMap)		
		// 	}
		// }
	}

	fmt.Println("Cleaned array of maps")
	PrintMapArray(noDups) 
	return noDups
}