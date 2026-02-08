package validators

import (
	"errors"
	"strings"
)

func ValidateSportInput(name string, minPlayers, maxPlayers int) error {

	if strings.TrimSpace(name) == "" {
		return errors.New("name is required")
	}

	if minPlayers < 1 {
		return errors.New("minPlayers must be at least 1")
	}

	if maxPlayers < 1 {
		return errors.New("maxPlayers must be at least 1")
	}

	if minPlayers > maxPlayers {
		return errors.New("minPlayers cannot be greater than maxPlayers")
	}

	if maxPlayers > 100 { // reasonable upper limit
		return errors.New("maxPlayers seems too high")
	}

	return nil
}
