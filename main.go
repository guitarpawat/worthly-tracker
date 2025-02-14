package main

import (
	"github.com/guitarpawat/worthly-tracker/cmd"
)

//go:generate templ generate -path ./internal/view

func main() {
	err := cmd.Execute()
	if err != nil {
		panic(err)
	}
}
