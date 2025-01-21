package main

import (
	"github.com/guitarpawat/worthly-tracker/cmd"
)

//go:generate templ generate -path ./internal/templ

func main() {
	err := cmd.Execute()
	if err != nil {
		panic(err)
	}
}
