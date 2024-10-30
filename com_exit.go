package main

import "fmt"
import "os"

func call_exit(cfg_state *config_state, args ...string) error{
	fmt.Println("Exiting")
	os.Exit(0)
	return nil
}