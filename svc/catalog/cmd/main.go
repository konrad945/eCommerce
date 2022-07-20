//go:generate oapi-codegen -package api ../openapi.yaml > ../api/api.gen.go
package main

import (
	"context"
	"github.com/konrad945/eCommerce/svc/catalog/internal/app"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	a := app.NewApp()

	go func(a *app.App) {
		if err := a.Run(); err != nil {
			log.Println(err)
		}
	}(a)

	<-ctx.Done()

	if err := a.Shutdown(); err != nil {
		log.Println(err)
	}

}
