package main

import (
	"flag"
	tea "github.com/charmbracelet/bubbletea"
	"log"
	"os"
)

var token = flag.String("token", "", "Lepton API token")

func main() {
	flag.Parse()
	if *token == "" {
		flag.Usage()
		os.Exit(1)
	}
	p := tea.NewProgram(initialModel(*token))

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
