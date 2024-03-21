package db

import (
	"context"
	dbmodel "sy_spatio-temporal_big_data_platform/db_model"

	"github.com/cloudwego/hertz/cmd/hz/util/logs"
)

func GetAccountInfo(ctx context.Context, accountNumber string) (*dbmodel.Account, error) {
	account := &dbmodel.Account{}
	err := db.Debug().Model(&dbmodel.Account{}).Where("account_number = ?", accountNumber).Find(account).Error
	if err != nil {
		logs.Error("db GetAccountInfo err: %v", err)
		return nil, err
	}
	return account, nil
}

func GetAllAccountInfo(ctx context.Context) ([]*dbmodel.Account, error) {
	account := make([]*dbmodel.Account, 0)
	err := db.Debug().Model(&dbmodel.Account{}).Find(&account).Error
	if err != nil {
		logs.Error("db GetAccountInfo err: %v", err)
		return nil, err
	}
	return account, nil
}
