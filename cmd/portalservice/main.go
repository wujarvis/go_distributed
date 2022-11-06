package main

import (
	"context"
	"distributed/log"
	"distributed/portal"
	"distributed/registry"
	"distributed/service"
	"fmt"
	stlog "log"

	"net/http"
	_ "net/http/pprof"
)

func main() {
	err := portal.ImportTemplates()
	if err != nil {
		stlog.Fatal(err)
	}
	host, port := "localhost", "8090"
	serviceAddress := fmt.Sprintf("http://%s:%s", host, port)

	r := registry.Registration{
		ServiceName: registry.PortalService,
		ServiceURL:  serviceAddress,
		RequiredServices: []registry.ServiceName{
			registry.LogService,
			registry.GradingService,
		},
		ServiceUpdateURL: serviceAddress + "/services",
		HeartbeatURL:     serviceAddress + "/heartbeat",
	}

	ctx, err := service.Start(
		context.Background(),
		r,
		host,
		port,
		portal.RegistryHandlers,
	)
	if err != nil {
		stlog.Fatal(err)
	}

	if logProvider, err := registry.GetProvider(registry.LogService); err != nil {
		log.SetClientLogger(logProvider, r.ServiceName)
	}

	go http.ListenAndServe("localhost:9000", nil)

	<-ctx.Done()
	fmt.Println("Shutting down portal.")
}
