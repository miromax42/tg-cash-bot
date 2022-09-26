package repo

import (
	"context"
	"strconv"
	"strings"
	"time"
)

type Expense interface {
	CreateExpense(ctx context.Context, req CreateExpenseReq) (*CreateExpenseResp, error)
	ListUserExpense(ctx context.Context, req ListUserExpenseReq) (ListUserExpenseResp, error)
}

type CreateExpenseReq struct {
	UserID   int64
	Amount   int
	Category string
}

type CreateExpenseResp struct {
	Amount    int
	Category  string
	CreatedAt time.Time
}

func (r *CreateExpenseResp) String() string {
	b := strings.Builder{}

	b.WriteString(strconv.Itoa(r.Amount))
	b.WriteString(" on ")
	b.WriteString(r.Category + "\n")
	b.WriteString("(" + r.CreatedAt.Format(time.Kitchen) + ")")

	return b.String()
}

type ListUserExpenseReq struct {
	UserID   int64
	FromTime time.Time
}

type ListUserExpenseResp []struct {
	Category string `json:"category"`
	Sum      int    `json:"sum"`
}

func (r ListUserExpenseResp) String() string {
	b := strings.Builder{}

	for i := range r {
		b.WriteString(r[i].Category)
		b.WriteString(": ")
		b.WriteString(strconv.Itoa(r[i].Sum))

		if i != len(r)-1 {
			b.WriteString("\n")
		}
	}

	if b.Len() == 0 {
		return "no expenses"
	}

	return b.String()
}
