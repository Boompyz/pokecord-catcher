package pokemon

import (
	"fmt"

	"github.com/antchfx/xquery/html"
)

// Pokemon represents a pokemon
type Pokemon struct {
	hash uint64
	page string
	name string
}

func (p *Pokemon) resolve() error {
	pokemonPage := "https://bulbapedia.bulbagarden.net" + p.page
	doc, err := htmlquery.LoadURL(pokemonPage)
	if err != nil {
		return err
	}

	picElement := htmlquery.FindOne(doc, "//table[@class=\"roundy\"]/tbody/tr/td[@colspan=\"4\"]/a[@class=\"image\"]/img")
	picSource := htmlquery.SelectAttr(picElement, "src")

	fmt.Println(picSource)
	return nil
}
