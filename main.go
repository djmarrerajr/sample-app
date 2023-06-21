package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/djmarrerajr/common-lib/app"
	"github.com/djmarrerajr/common-lib/shared"
	"github.com/djmarrerajr/common-lib/utils"
)

var (
	cfgPath = flag.String("cfg", "./config", "path in which to find the .env files")
)

func main() {
	env, err := utils.LoadEnv(*cfgPath)
	if err != nil {
		log.Fatalf("unable to parse environment: %v", err)
	}

	app, err := app.NewWithApiFromEnv(
		env,
		app.WithRouteHandler("/time", timeHandler, http.MethodGet),
		app.WithRequestHandler("/hello", helloHandler, nil, http.MethodGet),
	)
	if err != nil {
		log.Fatalf("unable to instantiate application: %v", err)
	}

	if err = app.Run(); err != nil {
		log.Fatalf("application terminated in error: %v", err)
	}

	log.Print("application terminated successfully")
	os.Exit(0)
}

func timeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	// nolint: errcheck
	w.Write([]byte(time.Now().UTC().Format(time.RFC3339)))
}

func helloHandler(ctx context.Context, appCtx *shared.ApplicationContext, req any) any {
	type HelloResponse struct {
		Message string
	}

	return HelloResponse{
		Message: "Hello World!",
	}
}
