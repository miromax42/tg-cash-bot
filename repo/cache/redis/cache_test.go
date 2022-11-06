package redis

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"gitlab.ozon.dev/miromaxxs/telegram-bot/repo"
	"gitlab.ozon.dev/miromaxxs/telegram-bot/repo/cache"
	rediscontainer "gitlab.ozon.dev/miromaxxs/telegram-bot/repo/cache/redis_container_test"
	"gitlab.ozon.dev/miromaxxs/telegram-bot/util"
)

type CacheTestSuite struct {
	suite.Suite
	ctx context.Context

	testCache *Cache
}

func (s *CacheTestSuite) SetupSuite() {
	s.ctx = context.Background()

	var err error
	s.testCache, err = NewCache(s.ctx, util.ConfigCache{
		Redis: util.ConfigRedis{
			SocketAddr: rediscontainer.GetRedisConnectionString(s.T()),
		},
		TTL:          time.Minute,
		ObjectsCount: 100,
		ObjectTTL:    time.Hour,
	})
	require.NoError(s.T(), err)
}

func TestCacheTestSuite(t *testing.T) {
	suite.Run(t, new(CacheTestSuite))
}

func (s *CacheTestSuite) TestSetGet() {
	const (
		userID = 1312

		currency = 1
		limit    = 2
	)

	value := repo.PersonalSettingsResp{
		Currency: currency,
		Limit:    limit,
	}
	err := s.testCache.Set(s.ctx, cache.UserSettingsToken(userID), value)
	require.NoError(s.T(), err)

	var getValue repo.PersonalSettingsResp
	err = s.testCache.Get(s.ctx, cache.UserSettingsToken(userID), &getValue)
	require.NoError(s.T(), err)

	require.EqualValues(s.T(), value, getValue)
}

func (s *CacheTestSuite) TestGetWithMiss() {
	const (
		userIDnotInCache = 1337
	)

	var getValue repo.PersonalSettingsResp
	err := s.testCache.Get(s.ctx, cache.UserSettingsToken(userIDnotInCache), &getValue)
	require.ErrorIs(s.T(), err, cache.ErrMiss)
}
