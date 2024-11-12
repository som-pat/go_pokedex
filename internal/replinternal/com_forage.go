package replinternal

import (
	"errors"
	"fmt"
	"strings"

	"github.com/som-pat/poke_dex/internal/config"
)


func call_forage(cfg_state *config.ConfigState, args ...string) (string,[]string,error){
	if len(args) != 1{
		return "",nil,errors.New("no item name provided")
	}
	toForage := args[0]
	_,exists := cfg_state.ItemsHeld[toForage]
	if exists{
		return "",nil, fmt.Errorf("%s already exists in your inventory",toForage)
	}
	itemdes, itemErr := cfg_state.PokeapiClient.ItemFetch(toForage)
	if itemErr != nil{
		return "", nil, errors.New("no such item found")
	}
	var itemScoured strings.Builder
	cfg_state.ItemsHeld[itemdes.Name] = itemdes
	itemScoured.WriteString(fmt.Sprintf("Foraged %s \n",itemdes.Name))
	itemScoured.WriteString(fmt.Sprintf("len of %d \n",len(cfg_state.ItemsHeld)))
	return itemScoured.String(),nil, nil	

}