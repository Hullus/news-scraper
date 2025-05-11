package main

import (
	"fmt"
	"strings"
)

func getSource(input string) string {
	strings.CutSuffix(input, "https://")
	fmt.Println(input)
	return input
}
