package main

import (
	"github.com/guitarpawat/worthly-tracker/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		panic(err)
	}
}
