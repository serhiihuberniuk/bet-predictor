package models

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
)

var nameSpaceLeague = uuid.NewSHA1(uuid.NameSpaceURL, []byte("league"))

type League struct {
	ID              string `bson:"_id,omitempty"`
	Name            string `bson:"name,omitempty"`
	Slug            string `bson:"slug,omitempty"`
	Country         string `bson:"country,omitempty"`
	CountrySlug     string `bson:"country_slug,omitempty"`
	CurrentSeasonID int    `bson:"current_season_id"`
}

type CreateLeaguePayload struct {
	Name            string
	Country         string
	CurrentSeasonID int
}

func NewLeague(payload CreateLeaguePayload) (*League, error) {
	league := &League{
		Name:            strings.Trim(payload.Name, " "),
		Country:         strings.Trim(payload.Country, " "),
		CurrentSeasonID: payload.CurrentSeasonID,
	}
	if league.Name == "" {
		return nil, errors.New("name is not specified")
	}

	league.Slug = slug.Make(league.Name)
	league.CountrySlug = slug.Make(league.Country)
	league.ID = uuid.NewSHA1(nameSpaceLeague, []byte(league.CountrySlug+league.Slug)).String()

	return league, nil
}
