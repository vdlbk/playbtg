package utils

import (
	"io/ioutil"
	"regexp"
	"strings"
)

func ReadFromFile(path string) ([]string, error) {
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	text := string(dat)
	rgx := regexp.MustCompile("\n|\r|\t")
	text = rgx.ReplaceAllString(text, " ")
	words := strings.Split(text, " ")

	return words, nil
}
