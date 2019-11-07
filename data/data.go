package data

import (
	"encoding/json"

	"github.com/codechrysalis/go.pokemon-api/models"
	"github.com/gobuffalo/packr/v2"
)

var (
	attacks models.AttackList
	pokemon models.PokemonList
	types   models.TypeList
	box     = packr.New("assets", "../assets")
)

// Pokemon contains all available Pokemon
func Pokemon() *models.PokemonList {
	return &pokemon
}

// Attacks contains all available attacks
func Attacks() *models.AttackList {
	return &attacks
}

// Types contains all available Types
func Types() *models.TypeList {
	return &types
}

func loadFile(name string) []byte {
	res, err := box.Find(name)

	if err != nil {
		panic(err)
	}

	return res
}

func init() {
	Reload()
}

// Reload the data from the json files
func Reload() {
	json.Unmarshal(loadFile("attacks.json"), &attacks)
	json.Unmarshal(loadFile("pokemon.json"), &pokemon)
	json.Unmarshal(loadFile("types.json"), &types)
}
