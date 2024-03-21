package main

import (
	"context"
	"sy_spatio-temporal_big_data_platform/dal/db"
)

func main() {
	ctx := context.Background()
	db.Init()
	ServerInit(ctx)
}
