package models

import (
	"errors"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
)

type League struct {
	ID          string `bson:"_id,omitempty"`
	Name        string `bson:"name,omitempty"`
	Slug        string `bson:"slug,omitempty"`
	Country     string `bson:"country,omitempty"`
	CountrySlug string `bson:"country_slug,omitempty"`
}

type CreateLeaguePayload struct {
	Name    string
	Country string
}

func (l *League) SetSlug() error {
	if l.Name == "" {
		return errors.New("field name is empty")
	}

	l.Slug = slug.Make(l.Name)

	return nil
}

func (l *League) SetCountrySlug() {
	l.CountrySlug = slug.Make(l.Country)

	return
}

func (l *League) SetID() error {
	if l.Slug == "" {
		return errors.New("cannot create ID: field slug is empty")
	}

	l.ID = uuid.NewSHA1(uuid.NameSpaceURL, []byte(l.CountrySlug+l.Slug)).String()

	return nil
}
