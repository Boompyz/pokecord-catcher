package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Boompyz/pokecord-catcher/pokemon"

	"github.com/bwmarrin/discordgo"
)

// Variables used for command line parameters
var (
	Token string
	info  *pokemon.PokemonInfo
)

func init() {

	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()

	info = pokemon.NewPokemonInfo()
	info.FillFromWeb(16)
}

func main() {

	//fmt.Println(info.FindPokemon("https://cdn.discordapp.com/attachments/508677201042997248/508682357948022794/PokecordSpawn.jpg"))

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()

}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	if len(m.Embeds) != 0 {
		embed := m.Embeds[0]
		if embed.Image != nil {
			imageURL := embed.Image.URL
			fmt.Println(imageURL)
			//proxyImage := embed.Image.ProxyURL
			//fmt.Println(proxyImage)
			pokemonName := info.FindPokemon(imageURL)
			s.ChannelMessageSend(m.ChannelID, "It's "+pokemonName+"!")
		}
	}

}
