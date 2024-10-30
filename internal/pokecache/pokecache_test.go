package pokecache

import (
	"fmt"
	"testing"
	"time"
)

func TestCache(t *testing.T){
	cache := CreateCache(time.Millisecond)
	if cache.cache == nil{
		t.Error("cache is null")
	}
}

func TestAddgetcache(t *testing.T){
	cache := CreateCache(time.Millisecond)
	new_case := []struct{
		inputkey string
		inputval []byte
	}{
		{
			inputkey: "Wizard",
			inputval: []byte("Virgin"),
		},
	}
	for _, cas := range new_case{
		cache.Add(cas.inputkey, cas.inputval)
		actual, ok := cache.Get(cas.inputkey)

		if !ok {
			t.Errorf("%s not found", cas.inputkey)
		}
		if string(actual) !=string(cas.inputval){
			t.Errorf("%s dont match %s",
				string(actual),
				string(cas.inputval),
			)
			continue
		}

	}
}

func TestPurge(t *testing.T){
	interval := time.Millisecond *10
	cache := CreateCache(interval)
	new_case := []struct{
		inputkey string
		inputval []byte
	}{
		{
			inputkey: "Wizard",
			inputval: []byte("Virgin"),
		},
	}
	for _, cas := range(new_case){
		cache.Add(cas.inputkey, cas.inputval)
		time.Sleep(interval+ time.Millisecond)
		
		_,ok := cache.Get(cas.inputkey)
		if ok{
			t.Errorf("%s should be purged",
				cas.inputkey)
		}
	}

}

func TestPurgeFail(t *testing.T){
	interval := time.Millisecond *10
	cache := CreateCache(interval)
	new_case := []struct{
		inputkey string
		inputval []byte
	}{
		{
			inputkey: "Wizard",
			inputval: []byte("Virgin"),
		},
	}
	for _, cas := range(new_case){
		cache.Add(cas.inputkey, cas.inputval)
		time.Sleep(interval /2)
		
		cachem,ok := cache.Get(cas.inputkey)
		if !ok{
			t.Errorf("%s should not be purged",
				cas.inputkey)
			fmt.Println(cachem)
		}
	}

}