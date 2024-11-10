package main

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
)

func call_search(cfg_state *ConfigState, args ...string) (string,[]string,error){
	if len(args) != 1{
		return "",nil,errors.New("no pokemon or item name provided")
	}
	if cfg_state.CurrentEncounterList == nil {
        return "",nil, errors.New("region not encountered")
    }
	toSearch := args[0]
	toSearchok := false
	for _, name := range *cfg_state.CurrentEncounterList{
		if name == toSearch{ 
			toSearchok = true
			break
		}
	}
	var encounter strings.Builder
	var encolist []string
	if !toSearchok{
		return "", nil, fmt.Errorf("%s not in this region",toSearch)
	}

	const commonWeight = 10
	const isItem = 12.3
	const babyMultiplier = 2.0
	const legendaryMultiplier = 0.1
	const mythicalMultiplier = 0.05

	totalWeight := 0.0
	weightMap := make(map[string]float64)

	for _,name := range *cfg_state.CurrentEncounterList{
		pokeSpecies, err := cfg_state.pokeapiClient.EncounterPoke(name)
		baseWeight := float64(commonWeight)
		if err !=nil{
			refactorItem := baseWeight * float64(isItem)
			weightMap[name] =refactorItem
			totalWeight += refactorItem
			continue
		}

		if pokeSpecies.IsBaby{
			baseWeight *= babyMultiplier
		}else if pokeSpecies.IsLegendary{
			baseWeight *= legendaryMultiplier
		}else if pokeSpecies.IsMythical{
			baseWeight *= mythicalMultiplier
		}

		normalizedCaptureRate :=  float64(pokeSpecies.CaptureRate) /255.0
		adjustedWeight := baseWeight * (1 + normalizedCaptureRate)

		weightMap[pokeSpecies.Name] = adjustedWeight
		totalWeight += adjustedWeight
	}

	encounter.WriteString("Encounters while searching the region: \n")
	maxencounter := 6
	encountered := make(map[string]bool)
	for i:=0; i<maxencounter;i++{
		threshold := (rand.Float64() * totalWeight)
		cumulativeWeight := 0.0	
		keys := make([]string, 0, len(weightMap))
		for name := range weightMap {
			keys = append(keys, name)
		}
		rand.Shuffle(len(keys), func(i, j int) { keys[i], keys[j] = keys[j], keys[i] })
		
		for _,name := range keys{
			//cumulative probability thresh
			weight := weightMap[name]
			cumulativeWeight += weight*2
			if cumulativeWeight >= threshold && !encountered[name]{
				encounter.WriteString(fmt.Sprintf(" - %s\n", name))
				encolist = append(encolist, name)
				encountered[name] = true
				break
			}
		}
	}
	if len(encountered) == 0 {
		encounter.WriteString("Nothing found in this region.\n")
	}
	

	return encounter.String(), encolist,nil
}