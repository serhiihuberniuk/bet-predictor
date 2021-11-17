package models

import (
	"strings"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"github.com/pkg/errors"
)

var nameSpaceTeam = uuid.NewSHA1(uuid.NameSpaceURL, []byte("team"))

type Team struct {
	ID          string `bson:"_id"`
	Name        string `bson:"name"`
	Slug        string `bson:"slug"`
	Country     string `bson:"country"`
	CountrySlug string `bson:"country_slug"`
}

type CreateTeamPayload struct {
	Name    string
	Country string
}

func NewTeam(payload CreateTeamPayload) (*Team, error) {
	t := &Team{
		Name:    strings.Trim(payload.Name, " "),
		Country: strings.Trim(payload.Country, " "),
	}

	if t.Name == "" {
		return nil, errors.New("name id not specified")
	}

	if t.Country == "" {
		return nil, errors.New("country id not specified")
	}

	t.CountrySlug = slug.Make(t.Country)
	t.Slug = slug.Make(t.Name)
	t.ID = uuid.NewSHA1(nameSpaceTeam, []byte(t.CountrySlug+t.Slug)).String()

	return t, nil
}
