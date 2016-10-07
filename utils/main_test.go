package utils

import (
	"regexp"
	"testing"
)

func TestGenerateANString(t *testing.T) {
	sizes := []int{8, 16, 32, 64}
	for _, size := range sizes {
		key, err := GenerateRandomANString(size)
		if err != nil {
			t.Error(err)
		}

		if len(key) != size {
			t.Errorf("Expect key with %d length, but got %d", size, len(key))
		}

		re := regexp.MustCompile("^[a-zA-Z0-9_]*$")
		if !re.MatchString(key) {
			t.Errorf("Expecting alphanumeric string, but got %s", key)
		}
	}
}

func TestStringInSlice(t *testing.T) {
	tests := []struct {
		slice []string
		word  string
	}{
		{[]string{"Hello", "World"}, "World"},
		{[]string{"The", "Quick", "Brown", "Fox"}, "Fox"},
		{[]string{"Hi"}, "Hello"},
	}
	for i, test := range tests {
		exists := StringInSlice(test.word, test.slice)
		if i != 2 {
			if exists == false {
				t.Error("Expecting true but got false")
			}
		} else {
			if exists == true {
				t.Error("Expecting false but got true")
			}
		}
	}
}
