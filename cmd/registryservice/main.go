package main

import (
	"context"
	"distributed/registry"
	"fmt"
	"log"
	"net/http"
)

func main() {

	registry.SetupRegistryService()

	http.Handle("/services", &registry.RegistryService{})
	ctx, concel := context.WithCancel(context.Background())
	defer concel()
	var srv http.Server
	srv.Addr = registry.ServerPort
	go func() {
		log.Println(srv.ListenAndServe())
		concel()
	}()

	go func() {
		fmt.Println("Registry service started. Press any key to stop.")
		var s string
		fmt.Scan(&s)
		srv.Shutdown(ctx)
		concel()
	}()

	<-ctx.Done()
	fmt.Println("Shutting down registry service.")
}
