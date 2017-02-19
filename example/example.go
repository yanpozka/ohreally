package main

import (
	"fmt"
	"strconv"
)

type captain struct {
	id string
}

func doStuff(input string) (string, error) {
	_, err := strconv.Atoi(input)

	if err != nil {
		return "abc", err
	}

	return "abc", nil
}

func isValid() (captain, bool) {
	ok := strconv.IsPrint('a')

	if ok {
		return captain{id: "t"}, true
	}

	return captain{id: "f"}, false
}

func main() {
	fmt.Println(doStuff("5656"))
	fmt.Println(isValid())
}
