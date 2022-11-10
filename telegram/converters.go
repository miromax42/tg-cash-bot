package telegram

import (
	"time"

	"gitlab.ozon.dev/miromaxxs/telegram-bot/repo"
)

func (r *ListUserExpenseReq) ToDB() repo.ListUserExpenseReq {
	dbReq := repo.ListUserExpenseReq{
		UserID:   r.UserID,
		FromTime: r.FromTime,
	}

	if r.ToTime != nil {
		dbReq.ToTime = *r.ToTime
	} else {
		dbReq.ToTime = time.Now()
	}

	return dbReq
}
