package main

import (
	"fmt"
	"os"
)

//TODO: make function more generic to allow for pasing in options denoting what cleanup processes to kick off
//cleanup process in the error process when creation fo file fails
func ErrCleanUp(rule string,file *os.File, options ...string) {
	fmt.Printf("Error occured in creating new rule: %+s, cleaning up data!!!\n", rule)
	if options == nil {
		fmt.Println("No options passed in. Will completely clean up file")
	} else {
		fmt.Printf("options passed in: %+v\n" , options)
	}

	//TODO: add a check to keep trying X number of times till sucessful

	//1. Close the file if its not already closed 
	fmt.Println("Attempting to close file")
	err := file.Close()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully closed file")

	//2. Delete the file
	fmt.Println("Attempting to delete file") 
	err = os.Remove(rule + ".csv") 
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully deleted file")
}