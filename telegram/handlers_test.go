package telegram

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gopkg.in/telebot.v3"

	"gitlab.ozon.dev/miromaxxs/telegram-bot/currency"
	"gitlab.ozon.dev/miromaxxs/telegram-bot/mocks"
	"gitlab.ozon.dev/miromaxxs/telegram-bot/repo"
	"gitlab.ozon.dev/miromaxxs/telegram-bot/telegram/tools"
	"gitlab.ozon.dev/miromaxxs/telegram-bot/util"
)

type HandlersTestSuite struct {
	suite.Suite
}

func TestRedisTestSuite(t *testing.T) {
	suite.Run(t, new(HandlersTestSuite))
}

func (s *HandlersTestSuite) getServerWithMocks() *Server {
	return &Server{
		logger:       mocks.NewLogger(s.T()),
		expense:      mocks.NewExpense(s.T()),
		userSettings: mocks.NewPersonalSettings(s.T()),
		exchange:     mocks.NewExchange(s.T()),
	}
}

func (s *HandlersTestSuite) TestHelp() {
	teleContext := mocks.NewTelegramContext(s.T())
	teleContext.On("Send",
		mock.MatchedBy(func(s string) bool {
			return strings.Contains(s, "available commands")
		}),
	).Return(nil)

	err := StartHelp(teleContext)
	require.NoError(s.T(), err)
}

func (s *HandlersTestSuite) TestCreateExpense() {
	const (
		userID      = 1337
		amount      = "12.881"
		amountFloat = 12.881
		category    = "test"
	)

	testCases := []struct {
		name       string
		setup      func() *mocks.TelegramContext
		buildStubs func(*Server)
		wantErr    bool
		err        error
	}{
		{
			"happy path",
			func() *mocks.TelegramContext {
				c := mocks.NewTelegramContext(s.T())
				c.On("Get", SettingsKey.String()).Return(&repo.PersonalSettingsResp{
					Currency: currency.TokenRUB,
				})
				c.On("Sender").Return(&telebot.User{ID: userID})
				c.On("Args").Return([]string{amount, category})
				c.On("Send", mock.AnythingOfType("string")).Return(nil)

				return c
			},
			func(srv *Server) {
				exchange := mocks.NewExchange(s.T())
				exchange.On("Base").Return(currency.TokenRUB)
				exchange.On("Convert",
					mock.Anything,
					mock.MatchedBy(func(req currency.ConvertReq) bool {
						return reflect.DeepEqual(req, currency.ConvertReq{
							Amount: amountFloat,
							From:   currency.TokenRUB,
							To:     currency.TokenRUB,
						})
					})).Return(amountFloat, nil)
				srv.exchange = exchange

				expense := mocks.NewExpense(s.T())
				expense.On("CreateExpense",
					mock.Anything,
					mock.MatchedBy(func(req repo.CreateExpenseReq) bool {
						return reflect.DeepEqual(req, repo.CreateExpenseReq{
							Amount:   amountFloat,
							UserID:   userID,
							Category: category,
						})
					})).Return(&repo.CreateExpenseResp{}, nil)
				srv.expense = expense
			},
			false,
			nil,
		},
		{
			"limit exceed",
			func() *mocks.TelegramContext {
				c := mocks.NewTelegramContext(s.T())
				c.On("Get", SettingsKey.String()).Return(&repo.PersonalSettingsResp{
					Currency: currency.TokenRUB,
				})
				c.On("Sender").Return(&telebot.User{ID: userID})
				c.On("Args").Return([]string{amount, category})
				c.On("Send",
					mock.MatchedBy(func(s string) bool {
						return strings.Contains(s, "limit")
					}),
				).Return(nil)

				return c
			},
			func(srv *Server) {
				exchange := mocks.NewExchange(s.T())
				exchange.On("Base").Return(currency.TokenRUB)
				exchange.On("Convert",
					mock.Anything,
					mock.MatchedBy(func(req currency.ConvertReq) bool {
						return reflect.DeepEqual(req, currency.ConvertReq{
							Amount: amountFloat,
							From:   currency.TokenRUB,
							To:     currency.TokenRUB,
						})
					})).Return(amountFloat, nil)
				srv.exchange = exchange

				expense := mocks.NewExpense(s.T())
				expense.On("CreateExpense",
					mock.Anything,
					mock.MatchedBy(func(req repo.CreateExpenseReq) bool {
						return reflect.DeepEqual(req, repo.CreateExpenseReq{
							Amount:   amountFloat,
							UserID:   userID,
							Category: category,
						})
					})).Return(nil, util.ErrLimitExceed)
				srv.expense = expense
			},
			true,
			tools.ErrLimitBlockExpense,
		},
	}

	for i := range testCases {
		tc := testCases[i]

		s.T().Run(tc.name, func(t *testing.T) {
			srv := s.getServerWithMocks()
			tc.buildStubs(srv)

			err := srv.CreateExpense(tc.setup())
			if tc.wantErr {
				assert.ErrorIs(t, err, tc.err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func (s *HandlersTestSuite) TestSetLimit() {
	const (
		userID = 1337
		limit  = "12.881"
	)

	var limitFloat = 12.881

	testCases := []struct {
		name       string
		setup      func() *mocks.TelegramContext
		buildStubs func(*Server)
		wantErr    bool
		err        error
	}{
		{
			"happy path",
			func() *mocks.TelegramContext {
				c := mocks.NewTelegramContext(s.T())
				c.On("Get", SettingsKey.String()).Return(&repo.PersonalSettingsResp{
					Currency: currency.TokenRUB,
				})
				c.On("Sender").Return(&telebot.User{ID: userID})
				c.On("Data").Return(limit)
				c.On("Send",
					mock.MatchedBy(func(s string) bool {
						return strings.Contains(s, fmt.Sprintf("%.2f", limitFloat))
					})).Return(nil)

				return c
			},
			func(srv *Server) {
				exchange := mocks.NewExchange(s.T())
				exchange.On("Base").Return(currency.TokenRUB)
				exchange.On("Convert",
					mock.Anything,
					mock.MatchedBy(func(req currency.ConvertReq) bool {
						return reflect.DeepEqual(req, currency.ConvertReq{
							Amount: limitFloat,
							From:   currency.TokenRUB,
							To:     currency.TokenRUB,
						})
					})).Return(limitFloat, nil)
				srv.exchange = exchange

				settings := mocks.NewPersonalSettings(s.T())
				settings.On("Set",
					mock.Anything,
					mock.MatchedBy(func(req repo.PersonalSettingsReq) bool {
						return reflect.DeepEqual(req, repo.PersonalSettingsReq{
							UserID: userID,
							Limit:  &limitFloat,
						})
					})).Return(nil)
				srv.userSettings = settings
			},
			false,
			nil,
		},
	}

	for i := range testCases {
		tc := testCases[i]

		s.T().Run(tc.name, func(t *testing.T) {
			srv := s.getServerWithMocks()
			tc.buildStubs(srv)

			err := srv.SetLimit(tc.setup())
			if tc.wantErr {
				assert.ErrorIs(t, err, tc.err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
