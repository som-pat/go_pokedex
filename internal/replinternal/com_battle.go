package replinternal

import (
	"errors"
	"fmt"
	"strings"

	"github.com/som-pat/poke_dex/imagegen"
	"github.com/som-pat/poke_dex/internal/config"
)

func call_battle(cfg_state *config.ConfigState, args ...string) (string,[]string,error){
	if len(args) != 1{
		return "",nil,errors.New("no pokemon name provided")
	}
	tobattle := args[0]
	var cmdseq strings.Builder
	var carrier []string
	pokeDetails, err := cfg_state.PokeapiClient.InvokePokeCatch(tobattle)
	if err!=nil{
		return "", nil, fmt.Errorf("%s not a valid pokemon", tobattle)
	}else{
		carrier = append(carrier, pokeDetails.Name)
		cmdseq.WriteString(fmt.Sprintf("Initiating Battle sequence with %s.....",pokeDetails.Name))
		ascii_img, err := imagegen.AsciiGen(pokeDetails.Sprites.FrontDefault)
		if err != nil {
			cmdseq.WriteString(" [Image Unavailable]\n")
		}
		cmdseq.WriteString(ascii_img + "\n")
		carrier = append(carrier, ascii_img)	
	}

	cmdseq.WriteString("Engaging")
	return cmdseq.String(), carrier, nil
}