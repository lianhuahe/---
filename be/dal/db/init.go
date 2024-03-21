package db

import (
	"fmt"

	"github.com/cloudwego/hertz/cmd/hz/util/logs"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	//变量db 通过init直接初始化
	db     *gorm.DB
	err    error
	dbinfo string
)

func Init() {
	var (
		dbuser     = "root"
		dbpassword = "lianhua123"
		dbip       = "127.0.0.1"
		dbport     = "3306"
		database   = "sy_spatio_temporal_big_data_platform"
	)
	dbinfo = fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", dbuser, dbpassword, dbip, dbport, database)

	//初始化db
	db, err = gorm.Open("mysql", dbinfo)
	if err != nil {
		logs.Error("mysql打开失败, err: %v", err)
		return
	}
	logs.Info("=====数据库初始化成功=====")
}
