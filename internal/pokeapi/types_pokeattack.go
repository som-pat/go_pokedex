package pokeapi

type PokeAbility struct {
	EffectChanges []struct {
		EffectEntries []struct {
			Effect   string `json:"effect"`
			Language struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"language"`
		} `json:"effect_entries"`
		VersionGroup struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"version_group"`
	} `json:"effect_changes"`
	EffectEntries []struct {
		Effect   string `json:"effect"`
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		ShortEffect string `json:"short_effect"`
	} `json:"effect_entries"`
	FlavorTextEntries []struct {
		FlavorText string `json:"flavor_text"`
		Language   struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		VersionGroup struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"version_group"`
	} `json:"flavor_text_entries"`
	Generation struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"generation"`
	ID           int    `json:"id"`
	IsMainSeries bool   `json:"is_main_series"`
	Name         string `json:"name"`
	Names        []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	Pokemon []struct {
		IsHidden bool `json:"is_hidden"`
		Pokemon  struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		Slot int `json:"slot"`
	} `json:"pokemon"`
}

type PokeMove struct {
	Accuracy      int `json:"accuracy"`
	ContestCombos struct {
		Normal struct {
			UseAfter []struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"use_after"`
			UseBefore any `json:"use_before"`
		} `json:"normal"`
		Super struct {
			UseAfter  any `json:"use_after"`
			UseBefore any `json:"use_before"`
		} `json:"super"`
	} `json:"contest_combos"`
	ContestEffect struct {
		URL string `json:"url"`
	} `json:"contest_effect"`
	ContestType struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"contest_type"`
	DamageClass struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"damage_class"`
	EffectChance  any   `json:"effect_chance"`
	EffectChanges []any `json:"effect_changes"`
	EffectEntries []struct {
		Effect   string `json:"effect"`
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		ShortEffect string `json:"short_effect"`
	} `json:"effect_entries"`
	FlavorTextEntries []struct {
		FlavorText string `json:"flavor_text"`
		Language   struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		VersionGroup struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"version_group"`
	} `json:"flavor_text_entries"`
	Generation struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"generation"`
	ID               int `json:"id"`
	LearnedByPokemon []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"learned_by_pokemon"`
	Machines []struct {
		Machine struct {
			URL string `json:"url"`
		} `json:"machine"`
		VersionGroup struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"version_group"`
	} `json:"machines"`
	Meta struct {
		Ailment struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"ailment"`
		AilmentChance int `json:"ailment_chance"`
		Category      struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"category"`
		CritRate     int `json:"crit_rate"`
		Drain        int `json:"drain"`
		FlinchChance int `json:"flinch_chance"`
		Healing      int `json:"healing"`
		MaxHits      any `json:"max_hits"`
		MaxTurns     any `json:"max_turns"`
		MinHits      any `json:"min_hits"`
		MinTurns     any `json:"min_turns"`
		StatChance   int `json:"stat_chance"`
	} `json:"meta"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	PastValues         []any `json:"past_values"`
	Power              int   `json:"power"`
	Pp                 int   `json:"pp"`
	Priority           int   `json:"priority"`
	StatChanges        []any `json:"stat_changes"`
	SuperContestEffect struct {
		URL string `json:"url"`
	} `json:"super_contest_effect"`
	Target struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"target"`
	Type struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"type"`
}

type FilPokeMove struct {
	Accuracy      int `json:"accuracy"`
	Pp            int   `json:"pp"`
	Priority      int   `json:"priority"`
	DamageClass struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"damage_class"`
	Power         int   `json:"power"`
	Meta struct {
		Ailment struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"ailment"`
		AilmentChance int `json:"ailment_chance"`
		Category      struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"category"`
		CritRate     int `json:"crit_rate"`
		Drain        int `json:"drain"`
		FlinchChance int `json:"flinch_chance"`
		Healing      int `json:"healing"`
		MaxHits      any `json:"max_hits"`
		MaxTurns     any `json:"max_turns"`
		MinHits      any `json:"min_hits"`
		MinTurns     any `json:"min_turns"`
		StatChance   int `json:"stat_chance"`
	} `json:"meta"`
	Name  string `json:"name"`
}