package pokeapi

import (
	"net/http"
	"time"

	"github.com/som-pat/poke_dex/internal/pokeapi/pokecache"
)

const baseURL = "https://pokeapi.co/api/v2/"

type Client struct{
	cache 	   pokecache.Cache
	httpClient http.Client
}

func NewClient(inter_time time.Duration) Client{
	return Client{
		cache: pokecache.CreateCache(inter_time),
		httpClient: http.Client{
			Timeout: time.Minute,
		},
	}
}


