package models

import (
	"fmt"

	"github.com/google/uuid"
)

func getNameSpace() (uuid.UUID, error) {
	nameSpace, err := uuid.Parse("750557cd-565d-4314-8f3b-0219be0b5a36")
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("error while parsing name-space: %w", err)
	}

	return nameSpace, nil
}
