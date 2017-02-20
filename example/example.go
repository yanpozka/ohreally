package main

import (
	"fmt"
	"strconv"
)

type captain struct {
	id string
}

func (c *captain) checkErr() error {

	if err := doStuff(c.id); err != nil {
		return err
	}

	return nil
}

func doStuff(input string) error {
	_, err := strconv.Atoi(input)

	if err != nil {
		return err
	}

	return nil
}

func isValid() bool {
	ok := strconv.IsPrint('a')

	if ok {
		return true
	}

	return false
}

func main() {
	fmt.Println(doStuff("5656"))
	fmt.Println(isValid())

	c := captain{id: "abc"}
	fmt.Println(c.checkErr())

	var err error
	fmt.Println(func() error {
		if err != nil {
			return err
		}

		return nil
	}())
}
