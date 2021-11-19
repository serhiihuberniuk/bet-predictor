package models

import (
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
)

type League struct {
	ID                string `bson:"_id,omitempty"`
	Name              string `bson:"name,omitempty"`
	Slug              string `bson:"slug,omitempty"`
	Country           string `bson:"country,omitempty"`
	CountrySlug       string `bson:"country_slug,omitempty"`
	ESCurrentSeasonID string `bson:"es_current_season_id"`
}

type CreateLeaguePayload struct {
	Name            string
	Country         string
	CurrentSeasonID string
}

func NewLeague(payload CreateLeaguePayload) (*League, error) {
	league := &League{
		Name:              strings.Trim(payload.Name, " "),
		Country:           strings.Trim(payload.Country, " "),
		ESCurrentSeasonID: payload.CurrentSeasonID,
	}
	if league.Name == "" {
		return nil, errors.New("name is not specified")
	}

	if league.Country != "" {
		league.Slug = slug.Make(league.Country + " " + league.Name)
		league.CountrySlug = slug.Make(league.Country)
	} else {
		league.Slug = slug.Make(league.Name)
	}

	nameSpace, err := getNameSpace()
	if err != nil {
		return nil, fmt.Errorf("error while getting nameSpace: %w", err)
	}

	league.ID = uuid.NewSHA1(nameSpace, []byte(league.CountrySlug+league.Slug)).String()

	return league, nil
}
