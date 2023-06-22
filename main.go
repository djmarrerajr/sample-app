package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
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

type Greeting struct {
	Name string `json:"name"  xml:"name"`
	Age  int    `json:"age"   xml:"age"   validate:"numeric,gt=10"`
}

func main() {
	env, err := utils.LoadEnv(*cfgPath)
	if err != nil {
		log.Fatalf("unable to parse environment: %v", err)
	}

	app, err := app.NewWithApiFromEnv(
		env,
		app.WithRouteHandler("/time", timeHandler, http.MethodGet),
		app.WithRequestHandler("/hello", helloHandler, nil, http.MethodGet),
		app.WithRequestHandler("/greet", greetHandler, Greeting{}, http.MethodPost),
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
	time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
	return "Hello World!"
}

func greetHandler(ctx context.Context, appCtx *shared.ApplicationContext, req any) any {
	type HelloResponse struct {
		Message string
	}

	rqst := req.(*Greeting)

	return HelloResponse{
		Message: fmt.Sprintf("Hello %s!", rqst.Name),
	}
}
