package syspg_test

import (
	"context"
	"testing"

	"github.com/andrew-hayworth22/critiquefi/services/api/internal/store/postgres/syspg"
	"github.com/andrew-hayworth22/critiquefi/services/api/internal/testutil"
)

func TestSysPg_Ping(t *testing.T) {
	t.Run("ping", func(t *testing.T) {
		s := syspg.New(testDB)

		err := s.Ping(context.Background())
		testutil.CheckErr(err, nil, t)
	})
}
