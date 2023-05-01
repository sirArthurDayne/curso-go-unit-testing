package util

import (
	"catching-pokemons/models"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParserPokemonSuccess(t *testing.T) {
	body, err := ioutil.ReadFile("samples/pokeapi_response.json")
	assert.NoError(t, err)

	var response models.PokeApiPokemonResponse

	err = json.Unmarshal(body, &response)
	assert.NoError(t, err)

	parsedPokemon, err := ParsePokemon(response)
	assert.NoError(t, err)

	body, err = ioutil.ReadFile("samples/api_response.json")
	assert.NoError(t, err)
	var expectedPokemon models.Pokemon

	err = json.Unmarshal(body, &expectedPokemon)
	assert.NoError(t, err)

	assert.Equal(t, expectedPokemon, parsedPokemon)
}
