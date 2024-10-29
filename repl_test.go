package main

import (
	"fmt"
	"testing"
)

func TestInputClean(t *testing.T){
	tcs := []struct{
		input string
		expected []string
		
	}{
		{
			input: "hello World",
			expected: []string{
				"hello",
				"world",
			},
			
		},
		{
			input: "",
			expected: nil,
			
		},
	}

	for _, tc := range tcs {
		actual:= input_clean(tc.input)
		
		if len(actual) != len(tc.expected){
			t.Errorf("Length are not equal:%v NOT EQUAL TO %v,",
			len(actual),
			len(tc.expected),
		)
		continue
		}
		for i:= range actual{
			ac_word := actual[i]
			ex_word := tc.expected[i]
			if ac_word != ex_word{
				t.Errorf("%v not equal to %v", 
						ac_word, 
						ex_word, 
					)
			}
			fmt.Println(ac_word,ex_word)
		}
	}
}