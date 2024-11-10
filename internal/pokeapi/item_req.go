package pokeapi

import (
	"net/http"
	"fmt"
	"io"
	"encoding/json"
	"math/rand"
)

func (c *Client) ItemRandomizer(numItems int) ([]string, error){
	end_point :="/item"
	full_url := baseURL + end_point
	
	req, err := http.NewRequest("GET", full_url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	
	defer resp.Body.Close()

	if resp.StatusCode > 399{
		return nil, fmt.Errorf("bad Status Encounterd: %v", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	
	itemResp := Itemsdef{}
	err = json.Unmarshal(data, &itemResp)
	
	if err != nil {
		return nil, err
	}
	if numItems > itemResp.Count{
		numItems = itemResp.Count
	}
	var selItems []string
	itemIndices := rand.Perm(len(itemResp.Results))[:numItems]
	for _, i := range itemIndices{
		selItems = append(selItems,itemResp.Results[i].Name )
	}

	return selItems, nil
}

func (c *Client) ItemFetch(item_name string) (ItemDescription, error){
	end_point := "/item/" + item_name
	full_url := baseURL + end_point

	req, err := http.NewRequest("GET", full_url, nil)
	if err != nil {
		return ItemDescription{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return ItemDescription{}, err
	}
	
	defer resp.Body.Close()

	if resp.StatusCode > 399{
		return ItemDescription{}, fmt.Errorf("bad Status Encounterd: %v", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return ItemDescription{}, err
	}
	
	item_details := ItemDescription{}
	err = json.Unmarshal(data, &item_details)
	
	if err != nil {
		return ItemDescription{}, err
	}
	return item_details, nil



}