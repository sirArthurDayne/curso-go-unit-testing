package util

import (
	"catching-pokemons/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

type args struct {
	apiPokemon models.PokeApiPokemonResponse
}

func TestParserPokemonSuccess(t *testing.T) {
	var pokeApiData map[string]models.PokeApiPokemonResponse
	err := ReadFileAndUnmarshall(&pokeApiData, "samples/pokeapi_response.json")
	assert.NoError(t, err, "[ERROR] cannot load pokeapi_response.json")

	var pokemon map[string]models.Pokemon
	err = ReadFileAndUnmarshall(&pokemon, "samples/api_response.json")
	assert.NoError(t, err, "[ERROR] cannot load api_response.json")
	tests := []struct {
		name        string
		args        args
		expected    models.Pokemon
		expectedErr error
	}{
		{
			name:        "Success pikachu",
			args:        args{apiPokemon: pokeApiData["pikachu"]},
			expected:    pokemon["pikachu"],
			expectedErr: nil,
		},
		{
			name:        "Pokemon type not found",
			args:        args{apiPokemon: pokeApiData["pikachuWithNotType"]},
			expected:    models.Pokemon{},
			expectedErr: ErrNotFoundPokemonType,
		},
		{
			name:        "Pokemon reftype's type is empty",
			args:        args{apiPokemon: pokeApiData["pikachuWithRefTypeEmpty"]},
			expected:    models.Pokemon{},
			expectedErr: ErrNotFoundPokemonTypeName,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParsePokemon(tt.args.apiPokemon)
			if err == ErrNotFoundPokemonType || err == ErrNotFoundPokemonTypeName {
				assert.NotNil(t, err, "[ERROR] ParsePokemon() error = %v, wantErr %v", err, tt.expectedErr)
				assert.EqualError(t, err, tt.expectedErr.Error(), "[ERROR] ParsePokemon() error = %v, wantErr %v", err, tt.expectedErr)
			} else {
				assert.NoError(t, err, "[ERROR] ParsePokemon() error = %v, wantErr %v", err, tt.expectedErr)
				assert.Equal(t, tt.expected, got, "[ERROR] ParsePokemon() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func BenchmarkParser(b *testing.B) {
	var pokemonResponse models.PokeApiPokemonResponse

	err := ReadFileAndUnmarshall(&pokemonResponse, "../controller/samples/pokeapi_response.json")
	assert.NoError(b, err)

	for n := 0; n < b.N; n++ {
		_, err = ParsePokemon(pokemonResponse)
		assert.NoError(b, err)
	}
}
