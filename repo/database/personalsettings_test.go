package database

import (
	"context"
	"reflect"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"gitlab.ozon.dev/miromaxxs/telegram-bot/ent"
	"gitlab.ozon.dev/miromaxxs/telegram-bot/repo"
)

type PersonalSettingsSuite struct {
	suite.Suite
	c repo.PersonalSettings
}

func (s *PersonalSettingsSuite) SetupSuite() {
	db, err := ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	require.NoError(s.T(), err)

	err = db.Schema.Create(context.Background())
	require.NoError(s.T(), err)

	s.c = NewPersonalSettings(db)
}

func TestPersonalSettingsSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(PersonalSettingsSuite))
}

func (s *PersonalSettingsSuite) TestGet() {
	const (
		idNotExists = 9999
	)

	type args struct {
		ctx context.Context
		id  int64
	}
	tests := []struct {
		name    string
		arg     int64
		want    *repo.PersonalSettingsResp
		wantErr bool
	}{
		{
			"default",
			idNotExists,
			repo.DefaultPersonalSettingsReq(),
			false,
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			got, err := s.c.Get(context.Background(), tt.arg)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}
