//+build wireinject

package main

import (
	"week04/internal/biz"
	"week04/internal/data"
	"week04/internal/service"

	"github.com/google/wire"
)

func InitializeVideoInfoService() *service.VideoInfoService {
	wire.Build(service.NewVideoInfoService,
		biz.NewVideoInfoBiz,
		wire.Bind(new(biz.VideoRepository), new(*data.VideoRepo)),
		data.NewVideoRepo,
		data.NewVideoCountStorage,
		data.NewVideoMetaStorage,
	)
	return &service.VideoInfoService{}
}
