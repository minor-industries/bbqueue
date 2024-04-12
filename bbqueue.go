package main

import "github.com/minor-industries/bbqueue/app"

func main() {
	err := app.Run()
	if err != nil {
		panic(err)
	}
}
