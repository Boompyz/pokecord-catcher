package pokemon

import (
	"fmt"
	"strconv"

	"github.com/antchfx/xquery/html"
)

// PokemonInfo represents information about all pokemon
type PokemonInfo struct {
	pokemons []Pokemon
}

// NewPokemonInfo creates a new PokemonInfo
func NewPokemonInfo() *PokemonInfo {
	return &PokemonInfo{nil}
}

// FillFromWeb gets the pokemons from bulbapedia
func (p *PokemonInfo) FillFromWeb() error {
	indexPage := "https://bulbapedia.bulbagarden.net/wiki/List_of_Pok%C3%A9mon_by_National_Pok%C3%A9dex_number"
	pokemons, err := getLinksFromIndexPage(indexPage)
	if err != nil {
		return err
	}

	for _, pokemon := range pokemons {
		pokemon.resolve() // Some pokemon appear twice!
	}

	p.pokemons = pokemons

	return nil
}

func getLinksFromIndexPage(indexPage string) ([]Pokemon, error) {
	doc, err := htmlquery.LoadURL(indexPage)
	if err != nil {
		return nil, err
	}

	links := htmlquery.Find(doc, "//div[@id=\"mw-content-text\"]/table[@align=\"center\"]/tbody/tr/td/a[@href]/img")
	fmt.Println("Found " + strconv.Itoa(len(links)) + " pokemon!")
	ret := make([]Pokemon, 0, len(links))
	for _, link := range links {
		link = link.Parent
		ret = append(ret, Pokemon{0, htmlquery.SelectAttr(link, "href"), htmlquery.InnerText(link)})
	}
	return ret, nil
}
