package main

import (
	"fmt"
	"log"
	"github.com/som-pat/poke_dex/internal/pokeapi"
)
func call_map() error{
	pokeapi_client := pokeapi.NewClient()
	resp, err := pokeapi_client.InvokeLocs()
	if err!= nil{
		log.Fatal(err)
	}
	fmt.Println("Locations:")
	for _, area := range resp. Results{
		fmt.Printf("- %s\n", area.Name)
	}
	return nil
}