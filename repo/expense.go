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
	Amount   float64
	Category string
}

type CreateExpenseResp struct {
	Amount    float64
	Category  string
	CreatedAt time.Time
}

func (r *CreateExpenseResp) String() string {
	b := strings.Builder{}

	b.WriteString(strconv.FormatFloat(r.Amount, 'f', 2, 64))
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
	Category string  `json:"category"`
	Sum      float64 `json:"sum"`
}
