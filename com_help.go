package main

import "fmt"

func call_help() error{
	fmt.Println("below are the current commands")
	avail_coms := get_command()
	for _,com := range avail_coms{
		fmt.Println(com.name)
	}

	return nil
}