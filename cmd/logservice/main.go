package main

import (
	"context"
	"distributed/log"
	"distributed/registry"
	"distributed/service"
	"fmt"
	stlog "log"
)

func main() {
	log.Run("./distributed.log")
	host, port := "localhost", "8080"
	serviceAddress := fmt.Sprintf("http://%s:%s", host, port)
	reg := registry.Registration{
		ServiceName:      registry.LogService,
		ServiceURL:       serviceAddress,
		RequiredServices: make([]registry.ServiceName, 0),
		ServiceUpdateURL: serviceAddress + "/services",
		HeartbeatURL:     serviceAddress + "/heartbeat",
	}
	ctx, err := service.Start(
		context.Background(),
		reg,
		host,
		port,
		log.RegisterHandlers,
	)
	if err != nil {
		stlog.Fatalln(err) // 自定义的log无法使用时，需要使用标准库的log记录错误日志
	}

	<-ctx.Done() // concel 时
	fmt.Println("Shutting down log service.")

}
