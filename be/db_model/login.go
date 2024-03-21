package dbmodel

import "time"

type Account struct {
	Id            int64  `json:"id"`
	AccountNumber string `json:"account_number"`
	Password      string
	CreateTime    time.Time `json:"create_time"`
	UpdateTime    time.Time `json:"update_time"`
	LastLogin     time.Time `json:"last_login"`
	Mail          string    `json:"mail"`
}

func (t Account) TableName() string {
	return "tb_account"
}
