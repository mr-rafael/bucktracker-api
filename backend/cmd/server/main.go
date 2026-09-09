package main

import (
	"github.com/Mr-Rafael/bucktracker-api/internal"
)

func main() {
	app := internal.New()
	app.Run()
}
