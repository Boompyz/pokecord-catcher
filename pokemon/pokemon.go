package pokemon

import (
	"fmt"
	"net/http"

	"github.com/Boompyz/pokecord-catcher/imagecomparer"
	"github.com/antchfx/xquery/html"
)

// Pokemon represents a pokemon
type Pokemon struct {
	Image *imagecomparer.ComparedImage `json:"image"`
	page  string
	Name  string `json:"Name"`
}

func (p *Pokemon) resolve() error {
	pokemonPage := "https://bulbapedia.bulbagarden.net" + p.page
	doc, err := htmlquery.LoadURL(pokemonPage)
	if err != nil {
		return err
	}

	picElement := htmlquery.FindOne(doc, "//table[@class=\"roundy\"]/tbody/tr/td[@colspan=\"4\"]/a[@class=\"image\"]/img")
	picSource := "https:" + htmlquery.SelectAttr(picElement, "src")

	resp, err := http.Get(picSource)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	p.Image, err = imagecomparer.NewComparedImage(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
	}

	return nil
}

func (p *Pokemon) GetDistance(image *imagecomparer.ComparedImage) float64 {
	return image.GetDistance(p.Image)
}
