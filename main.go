package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/brunoksato/golang-boilerplate/server"
	_ "github.com/heroku/x/hmetrics/onload"
)

func main() {
	server := server.Start()
	addr := ":" + os.Getenv("PORT")

	go func() {
		if err := server.Start(addr); err != nil {
			server.Logger.Info("shutting down the server")
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		server.Logger.Fatal(err)
	}
}
