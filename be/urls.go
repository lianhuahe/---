package main

import (
	"context"
	"sy_spatio-temporal_big_data_platform/views"
	"time"

	"github.com/hertz-contrib/cors"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func ServerInit(ctx context.Context) {
	h := server.Default(
		server.WithHostPorts("127.0.0.1:9999"),
		server.WithMaxRequestBodySize(4*1024*1024*1024),
	)
	h.Use(cors.New(cors.Config{
		//准许跨域请求网站,多个使用,分开,限制使用*
		AllowOrigins: []string{"*"},
		//准许使用的请求方式
		AllowMethods: []string{"PUT", "PATCH", "POST", "GET", "DELETE"},
		//准许使用的请求表头
		AllowHeaders: []string{"Origin", "Authorization", "Content-Type"},
		//显示的请求表头
		ExposeHeaders: []string{"Content-Type"},
		//凭证共享,确定共享
		AllowCredentials: true,
		//容许跨域的原点网站,可以直接return true就万事大吉了
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		//超时时间设定
		MaxAge: 24 * time.Hour,
	}))

	h.GET("/ping", func(c context.Context, ctx *app.RequestContext) {
		ctx.JSON(consts.StatusOK, utils.H{"message": "pong"})
	})

	LoginRegister(h)
	BusinessRegister(h)

	h.Spin()
}

func LoginRegister(h *server.Hertz) {
	h.POST("/auth/token/", views.Login)
	h.POST("/auth/account/info/", views.GetAccountInfo)
	h.GET("/auth/account/list_all/", views.ListAll)
}

func BusinessRegister(h *server.Hertz) {
	h.GET("/business/file/", views.SearchFiles)
	h.GET("/business/file/download/", views.Download)
	h.GET("/business/file/get_all/", views.FileGetAll)
	h.GET("/business/file/:id/generate_gis_view/", views.GenerateGisView)
	h.GET("/business/file/:id/get_file_status/", views.GetFileStatus)
	h.GET("/business/file/:id/get_gis_view/", views.GetGisView)

	h.POST("/business/file/", views.UploadFile)

	h.GET("/business/task/get_task_model_dict/", views.GetTaskModelDict)
	h.POST("/business/task/exists/", views.CheckTaskExists)
	h.POST("/business/task/", views.CreateTask)
	h.GET("/business/task/", views.SearchTasks)
}
