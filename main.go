package main

import (
	"fmt"

	"github.com/astridalia/lynxbot/mediawiki"
)

func main() {
	wikiService := mediawiki.NewWikiService()
	response, err := wikiService.WikiText("Pet:Stormzilla")
	if err != nil {
		panic(err)
	}
	fmt.Print(response)
}
