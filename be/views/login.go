package views

import (
	"context"
	"sy_spatio-temporal_big_data_platform/dal/db"

	"github.com/cloudwego/hertz/cmd/hz/util/logs"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type LoginReq struct {
	AccountNumber string `json:"account_number" form:"account_number"`
	Password      string `json:"password" form:"password"`
}

type LoginResp struct {
	AccountId int32  `json:"account_id"`
	Token     string `json:"token"`
	Refresh   string `json:"refresh"`
}

func Login(c context.Context, ctx *app.RequestContext) {
	req := &LoginReq{}
	err := ctx.Bind(req)
	if err != nil {
		logs.Error("Login Bind err: %v", err)
		return
	}
	logs.Info("account: %+v", req)

	accountInfo, err := db.GetAccountInfo(c, req.AccountNumber)
	if err != nil {
		logs.Error("Login db GetAccountInfo err: %v", err)
		ctx.JSON(consts.StatusInternalServerError, utils.H{"data": ""})
		return
	}
	if accountInfo.Password != req.Password {
		ctx.JSON(consts.StatusUnauthorized, utils.H{"data": "账号或密码错误"})
		return
	}

	ctx.JSON(consts.StatusOK, utils.H{"data": &LoginResp{AccountId: int32(accountInfo.Id), Token: "123"}})
}

type InfoReq struct {
	AccountId int32 `json:"account_id" form:"account_id"`
}

type InfoResp struct {
	Id            int32    `json:"id"`
	Roles         []string `json:"roles"`
	AccountNumber string   `json:"accountNumber"`
	Permissions   []string `json:"permissions"`
}

type ListAllResp struct {
	Id            int32  `json:"id"`
	AccountNumber string `json:"accountNumber"`
}

func GetAccountInfo(c context.Context, ctx *app.RequestContext) {
	req := &InfoReq{}
	err := ctx.Bind(req)
	if err != nil {
		logs.Error("Login Bind err: %v", err)
		return
	}

	logs.Info("%+v", req)

	// permissionCodes, err := db.GetAllPermissions(c, 3)
	// if err != nil {
	// 	logs.Error("Login db GetAllPermissions err: %v", err)
	// 	ctx.JSON(consts.StatusInternalServerError, utils.H{"data": ""})
	// }
	ctx.JSON(consts.StatusOK, utils.H{"data": &InfoResp{
		Id:          3,
		Roles:       make([]string, 0),
		Permissions: []string{"admin"},
	}})
}

func ListAll(c context.Context, ctx *app.RequestContext) {
	allAccounts, err := db.GetAllAccountInfo(c)
	if err != nil {
		logs.Error("Login db GetAccountInfo err: %v", err)
		ctx.JSON(consts.StatusInternalServerError, utils.H{"data": ""})
	}

	ctx.JSON(consts.StatusOK, utils.H{"data": allAccounts})
}
