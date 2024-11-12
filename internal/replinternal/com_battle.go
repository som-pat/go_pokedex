package replinternal

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"

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
	var rng = rand.New(rand.NewSource(time.Now().UnixNano()))
	pokeDetails, err := cfg_state.PokeapiClient.InvokePokeCatch(tobattle)
	if err!=nil{
		return "", nil, fmt.Errorf("%s not a valid pokemon", tobattle)
	}else{
		carrier = append(carrier, pokeDetails.Name)
		ranlev:=rng.Intn(12)
		cmdseq.WriteString(fmt.Sprintf("You have encounterd a WILD LV%d %s.....\n",ranlev,pokeDetails.Name))
		cmdseq.WriteString(fmt.Sprintf("Initiating Battle sequence with %s.....\n\n",pokeDetails.Name))
		ascii_img, err := imagegen.AsciiGen(pokeDetails.Sprites.FrontDefault,56)
		if err != nil {
			cmdseq.WriteString(" [Image Unavailable]\n")
		}
		cmdseq.WriteString(ascii_img + "\n")
		carrier = append(carrier, ascii_img)
		
		carrier = append(carrier, fmt.Sprintf("LV%s",strconv.Itoa(ranlev)))
		for _,stats := range pokeDetails.Stats{
			statmult := 0.63
			nstat := int(float64(stats.BaseStat) * math.Pow((1+statmult), float64(ranlev)))
			carrier = append(carrier, strconv.Itoa(nstat))
		}		
	}

	cmdseq.WriteString("Engaging")
	return cmdseq.String(), carrier, nil
}