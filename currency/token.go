package currency

import (
	"database/sql/driver"
	"fmt"

	"gitlab.ozon.dev/miromaxxs/telegram-bot/util"
)

type Token int

const (
	TokenRUB Token = iota
	TokenCNY
	TokenEUR
	TokenUSD
)

var Supported = [...]string{
	"RUB",
	"CNY",
	"EUR",
	"USD",
}

func Parse(cur string) (Token, error) {
	for i, v := range Supported {
		if v == cur {
			return Token(i), nil
		}
	}

	return 0, util.ErrUnsupported
}

func MustParse(cur string) Token {
	t, err := Parse(cur)
	if err != nil {
		panic(fmt.Errorf("corventing %q: %w", cur, err))
	}

	return t
}

func (t Token) String() string {
	return Supported[t]
}

func (Token) Values() (kinds []string) {
	return Supported[:]
}

func (t Token) Value() (driver.Value, error) {
	return t.String(), nil
}

func (t *Token) Scan(val any) error {
	var s string
	switch v := val.(type) {
	case nil:
		return nil
	case string:
		s = v
	}

	var err error
	if *t, err = Parse(s); err != nil {
		return err
	}

	return nil
}
