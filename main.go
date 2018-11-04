package main

import (
	"fmt"

	"github.com/Boompyz/pokecord-catcher/pokemon"
)

func main() {
	fmt.Println("Hello World!")
	info := pokemon.NewPokemonInfo()
	info.FillFromWeb()
}
