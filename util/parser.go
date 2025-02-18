package util

import (
	"catching-pokemons/models"
	"encoding/json"
	"errors"
	"io/ioutil"
)

var (
	// ErrNotFoundPokemonType occurs when the type array in pokeapi response it's not found
	ErrNotFoundPokemonType = errors.New("pokemon type array not found")
	// ErrNotFoundPokemonTypeName occurs when we found type struct but no name
	ErrNotFoundPokemonTypeName = errors.New("pokemon type name not found")
)

func ReadFileAndUnmarshall(data any, path string) error {
	body, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return err
	}
	return nil
}

func ParsePokemon(apiPokemon models.PokeApiPokemonResponse) (models.Pokemon, error) {
	if len(apiPokemon.PokemonType) < 1 {
		return models.Pokemon{}, ErrNotFoundPokemonType
	}

	if apiPokemon.PokemonType[0].RefType.Name == "" {
		return models.Pokemon{}, ErrNotFoundPokemonTypeName
	}

	pokemonType := apiPokemon.PokemonType[0].RefType.Name

	abilitiesMap := map[string]int{}

	for _, stat := range apiPokemon.Stats {
		parsedAbilityName, ok := models.AllowedAbilities[stat.Stat.Name]
		if !ok {
			continue
		}

		abilitiesMap[parsedAbilityName] = stat.BaseStat
	}

	parsedPokemon := models.Pokemon{
		Id:        apiPokemon.Id,
		Name:      apiPokemon.Name,
		Power:     pokemonType,
		Abilities: abilitiesMap,
	}

	return parsedPokemon, nil
}
