package filmpg_test

import (
	"context"
	"testing"
	"time"

	"github.com/andrew-hayworth22/critiquefi-service/internal/models"
	"github.com/andrew-hayworth22/critiquefi-service/internal/store/postgres/filmpg"
	"github.com/andrew-hayworth22/critiquefi-service/internal/testutil"
)

func TestFilmPg_CreateFilm(t *testing.T) {
	t.Run("create film", func(t *testing.T) {
		testutil.PrepareDB(t, testDB)
		s := filmpg.New(testDB)

		description := "An insomniac meets himself"
		runtimeMinutes := 95
		externalReferences := []models.ExternalReference{
			{
				Name: "IMDB",
				URL:  "https://www.imdb.com/title/tt0137523/",
			},
		}

		cases := []struct {
			name        string
			film        models.NewFilm
			expectedErr error
		}{
			{
				name: "success: creates film",
				film: models.NewFilm{
					FilmType:           models.FeatureFilm,
					Title:              "Fight Club",
					Description:        &description,
					ReleaseDate:        time.Time{},
					RuntimeMinutes:     &runtimeMinutes,
					ExternalReferences: externalReferences,
				},
				expectedErr: nil,
			},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				id, err := s.CreateFilm(context.Background(), tc.film)
				testutil.CheckErr(err, tc.expectedErr, t)

				if tc.expectedErr == nil && id == 0 {
					t.Fatal("expected non-zero film ID")
				}
			})
		}
	})
}

func TestFilmPg_GetFilmByID(t *testing.T) {
	t.Run("get film by id", func(t *testing.T) {
		testutil.PrepareDB(t, testDB)
		s := filmpg.New(testDB)

		description := "A movie about a fight between two characters"
		runtimeMinutes := 95
		releaseDate, err := time.Parse(time.RFC3339, "2022-01-01T00:00:00Z")
		if err != nil {
			t.Fatal(err)
		}

		cases := []struct {
			name         string
			filmID       int64
			expectedFilm models.Film
			expectedErr  error
		}{
			{
				name:   "success: gets film",
				filmID: 1,
				expectedFilm: models.Film{
					ID:             1,
					FilmType:       models.FeatureFilm,
					Title:          "Fight Club",
					Description:    &description,
					ReleaseDate:    releaseDate,
					RuntimeMinutes: &runtimeMinutes,
					ExternalReferences: []models.ExternalReference{
						{
							Name: "IMDB",
							URL:  "https://www.imdb.com/title/tt0137523/",
						},
					},
				},
				expectedErr: nil,
			},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				film, err := s.GetFilmByID(context.Background(), tc.filmID)
				testutil.CheckErr(err, tc.expectedErr, t)

				if film.ID != tc.expectedFilm.ID {
					t.Errorf("expected film ID: %v, got: %v", tc.expectedFilm.ID, film.ID)
				}
				if film.FilmType.String() != tc.expectedFilm.FilmType.String() {
					t.Errorf("expected film type: %v, got: %v", tc.expectedFilm.FilmType.String(), film.FilmType.String())
				}
			})
		}
	})
}
