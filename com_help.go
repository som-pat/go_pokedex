package main

import (
	"fmt"
	"strings"
)

func call_help(cfg_state *ConfigState, args ...string) (string,[]string,error){
	var result strings.Builder
	result.WriteString("Available Commands:\n\n")
	avail_coms := get_command()
	for _,com := range avail_coms{
		result.WriteString(fmt.Sprintf("- %s : %s \n",com.name, com.description))
	}

	return result.String(),nil,nil
}