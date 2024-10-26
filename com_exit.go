package main

import "fmt"
import "os"

func call_exit() error{
	fmt.Println("Exiting")
	os.Exit(0)
	return nil
}