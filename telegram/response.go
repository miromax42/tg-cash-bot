package telegram

import (
	"strconv"
	"strings"
	"time"

	"gitlab.ozon.dev/miromaxxs/telegram-bot/repo"
)

func ListExpensesAnswer(data repo.ListUserExpenseResp, multiplier float64) string {
	b := strings.Builder{}

	for i := range data {
		b.WriteString(data[i].Category)
		b.WriteString(": ")
		b.WriteString(strconv.FormatFloat(data[i].Sum*multiplier, 'f', 2, 64)) //nolint:gomnd

		if i != len(data)-1 {
			b.WriteString("\n")
		}
	}

	if b.Len() == 0 {
		return "no expenses"
	}

	return b.String()
}
func CreateExpenseAnswer(data *repo.CreateExpenseResp, amount float64) string {
	if data == nil {
		return "internal error"
	}

	b := strings.Builder{}

	b.WriteString(strconv.FormatFloat(amount, 'f', 2, 64)) //nolint:gomnd
	b.WriteString(" on ")
	b.WriteString(data.Category + "\n")
	b.WriteString("(" + data.CreatedAt.Format(time.Kitchen) + ")")

	return b.String()
}
