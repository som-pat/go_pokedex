package main

import (
	"fmt"
	"log"
	"github.com/som-pat/poke_dex/internal/pokeapi"
)

func call_mapb() error{
	pokeapi_client := pokeapi.NewClient()
	resp, err := pokeapi_client.InvokeLocs()
	if err!= nil{
		log.Fatal(err)
	}
	fmt.Println(resp)
	return nil
}