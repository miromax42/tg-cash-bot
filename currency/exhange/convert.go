package exhange

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"gitlab.ozon.dev/miromaxxs/telegram-bot/currency"
	"gitlab.ozon.dev/miromaxxs/telegram-bot/util"
)

type Converter struct {
	baseCurrency currency.Token
	data         map[currency.Token]float64
}

func New(ctx context.Context, cfg util.ConfigExchange) (*Converter, error) {
	data, err := getValues(ctx, cfg)
	if err != nil {
		return nil, err
	}

	return &Converter{
		baseCurrency: currency.MustParse(cfg.BaseCurrency),
		data:         data,
	}, nil
}

func (c *Converter) Base() currency.Token {
	return c.baseCurrency
}

func (c *Converter) Convert(_ context.Context, req currency.ConvertReq) (amount float64, err error) {
	_, ok1 := c.data[req.To]
	_, ok2 := c.data[req.From]
	if !ok1 || !ok2 {
		return 0, util.ErrUnsupported
	}

	coefficient := c.getCoefficient(req.From, req.To)
	if err != nil {
		return 0, err
	}

	return coefficient * req.Amount, nil
}

func (c *Converter) getCoefficient(from currency.Token, to currency.Token) float64 {
	var coefficient float64 = 1

	switch c.baseCurrency {
	case from:
		coefficient *= c.data[to]
	case to:
		coefficient /= c.data[from]
	default:
		toBase := c.getCoefficient(from, c.baseCurrency)
		fromBase := c.getCoefficient(c.baseCurrency, to)

		coefficient = toBase * fromBase
	}

	return coefficient
}

func getValues(ctx context.Context, cfg util.ConfigExchange) (map[currency.Token]float64, error) {
	response := struct {
		Quotes    map[string]float64 `json:"quotes"`
		Source    string             `json:"source"`
		Success   bool               `json:"success"`
		Timestamp int                `json:"timestamp"`
	}{}

	url := fmt.Sprintf("https://api.apilayer.com/currency_data/live?source=%s&currencies=%s",
		cfg.BaseCurrency, strings.Join(currency.Supported[:], "%2C"))

	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	req.Header.Set("apikey", cfg.Token)
	if err != nil {
		return nil, err
	}

	res, err := client.Do(req)
	if res.Body != nil {
		defer res.Body.Close()
	}
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	if !response.Success {
		return nil, util.ErrBadResponse
	}

	return parseValues(response.Source, response.Quotes)
}

func parseValues(baseCurrency string, data map[string]float64) (map[currency.Token]float64, error) {
	m := make(map[currency.Token]float64)

	for i, v := range currency.Supported {
		key := baseCurrency + v
		if _, ok := data[key]; !ok {
			if key == baseCurrency+baseCurrency {
				m[currency.Token(i)] = 1

				continue
			}
			return nil, util.ErrUnsupported
		}

		m[currency.Token(i)] = data[key]
	}

	return m, nil
}
