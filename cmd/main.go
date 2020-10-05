/*=============================================================================
  * @FileName: main.go
  * @Desc:
  * @Version: 1.0.0
  * @Author: guocheng.zeng
  * @Email: guocheng.zeng@ucloud.cn
  * @Date: 2020/9/26 09:22
=============================================================================*/
package main

import (
	"context"
	"log"

	"github.com/sirupsen/logrus"

	"github.com/feixiaobo/go-xxl-job-client/v2"
	"github.com/feixiaobo/go-xxl-job-client/v2/logger"
	"github.com/feixiaobo/go-xxl-job-client/v2/option"
)

func main() {
	client := xxl.NewXxlClient(
		option.WithClientPort(8083),
		option.WithAdminAddress("http://localhost:8080/xxl-job-admin"),
		option.WithAppName("xxl-job-executor-sample"),
	)
	client.SetLogger(&logrus.Entry{
		Logger: logrus.New(),
		Level:  logrus.InfoLevel,
	})
	client.RegisterJob("testJob", JobTest)
	client.Run()
}

func JobTest(ctx context.Context) error {
	val, _ := xxl.GetParam(ctx, "test")
	log.Print(">>>>>>>>>>>>>>>>", val)
	logger.Info(ctx, "test job!!!!!")
	param, _ := xxl.GetParam(ctx, "name") //获取输入参数
	logger.Info(ctx, "the input param:", param)
	shardingIdx, shardingTotal := xxl.GetSharding(ctx) //获取分片参数
	logger.Info(ctx, "the sharding param: idx:", shardingIdx, ", total:", shardingTotal)
	return nil
}