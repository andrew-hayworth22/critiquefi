package models

import (
	"fmt"
	"strings"
	"time"
)

// FilmType represents the type of film
type FilmType struct {
	value string
}

var (
	FeatureFilm = FilmType{"Feature Film"}
	ShortFilm   = FilmType{"Short Film"}
)

var validFilmTypes = map[string]FilmType{
	"FEATURE FILM": FeatureFilm,
	"SHORT FILM":   ShortFilm,
}

// ParseFilmType parses a string into a FilmType
func ParseFilmType(s string) (FilmType, error) {
	s = strings.ToUpper(s)
	ft, ok := validFilmTypes[s]
	if !ok {
		return FilmType{}, fmt.Errorf("invalid film type: %s", s)
	}
	return ft, nil
}

// String returns the string representation of the FilmType
func (ft *FilmType) String() string {
	return ft.value
}

// MarshalJSON implements the json.Marshaler interface
func (ft *FilmType) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, ft.value)), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (ft *FilmType) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	parsed, err := ParseFilmType(s)
	if err != nil {
		return err
	}
	*ft = parsed
	return nil
}

// Scan implements sql.Scanner so DB values can be scanned in directly
func (ft *FilmType) Scan(src any) error {
	s, ok := src.(string)
	if !ok {
		return fmt.Errorf("unsupported type for FilmType: %T", src)
	}
	parsed, err := ParseFilmType(s)
	if err != nil {
		return err
	}
	*ft = parsed
	return nil
}

type Film struct {
	ID                 int64
	FilmType           FilmType
	Title              string
	Description        *string
	ReleaseDate        time.Time
	RuntimeMinutes     *int
	ExternalReferences []ExternalReference
}

// NewFilm represents the data needed to create a new film
type NewFilm struct {
	FilmType           FilmType
	Title              string
	Description        *string
	ReleaseDate        time.Time
	RuntimeMinutes     *int
	ExternalReferences []ExternalReference
}

func (nf *NewFilm) Validate() error {
	ve := ValidationErrors{}

	if len(nf.Title) == 0 || len(nf.Title) > 250 {
		ve.Add("title", "title is required")
	}
	if nf.Description != nil && len(*nf.Description) > 500 {
		ve.Add("description", "description must be less than 500 characters")
	}

	if ve.Any() {
		return ve
	}
	return nil
}
