package pokemon

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/Boompyz/pokecord-catcher/imagecomparer"

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
func (p *PokemonInfo) FillFromWeb(threadCount int) error {
	indexPage := "https://bulbapedia.bulbagarden.net/wiki/List_of_Pok%C3%A9mon_by_National_Pok%C3%A9dex_number"
	pokemons, err := getLinksFromIndexPage(indexPage)
	if err != nil {
		return err
	}
	p.pokemons = pokemons

	var wg sync.WaitGroup
	wg.Add(threadCount)

	resolver := func(pi *PokemonInfo, getWork chan int) {
		for idx := range getWork {
			pi.pokemons[idx].resolve() // Some pokemon appear twice!
		}
		wg.Done()
	}

	giveWork := make(chan int)
	for i := 0; i < threadCount; i++ {
		go resolver(p, giveWork)
	}
	for idx := range pokemons {
		giveWork <- idx
		fmt.Println("[", idx, "/", len(pokemons), "]")
	}
	close(giveWork)
	wg.Wait()

	return nil
}

// FindPokemon attempts to guess what pokemon is in the image on the url
func (p *PokemonInfo) FindPokemon(imageURL string) string {
	resp, err := http.Get(imageURL)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	image, err := imagecomparer.NewComparedImage(resp.Body)
	var bestPokemon Pokemon
	var distance float64 = 2147483647

	for _, pokemon := range p.pokemons {
		d := pokemon.GetDistance(image)
		if d < distance {
			distance = d
			bestPokemon = pokemon
		}
	}

	fmt.Println(distance)
	return bestPokemon.name
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
		ret = append(ret, Pokemon{nil, htmlquery.SelectAttr(link, "href"), htmlquery.SelectAttr(link, "title")})
	}
	return ret, nil
}
