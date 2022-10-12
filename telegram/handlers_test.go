package telegram

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"gitlab.ozon.dev/miromaxxs/telegram-bot/mocks"
)

type HandlersTestSuite struct {
	suite.Suite

	s *Server
}

func TestRedisTestSuite(t *testing.T) {
	suite.Run(t, new(HandlersTestSuite))
}

func (h *HandlersTestSuite) TestHelp() {
	teleContext := &mocks.TelegramContext{}
	teleContext.On("Send",
		mock.MatchedBy(func(s string) bool {
			return strings.Contains(s, "available commands")
		}),
	).Return(nil)

	err := h.s.StartHelp(teleContext)
	require.NoError(h.T(), err)
}
