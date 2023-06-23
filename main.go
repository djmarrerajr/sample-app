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
	"github.com/djmarrerajr/common-lib/observability/tracing"
	"github.com/djmarrerajr/common-lib/shared"
	"github.com/djmarrerajr/common-lib/utils"
)

var (
	cfgPath = flag.String("cfg", "./config", "path in which to find the .env files")
)

var (
	dbServer    *Database
	emailServer *EmailServer
	emailVendor *EmailVendor
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

	dbServer = &Database{appCtx: *app.AppContext}
	emailServer = &EmailServer{appCtx: *app.AppContext}
	emailVendor = &EmailVendor{appCtx: *app.AppContext}

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

func helloHandler(ctx context.Context, appCtx *shared.ApplicationContext, req any) (any, int) {
	time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
	return "Hello World!", http.StatusOK
}

func greetHandler(ctx context.Context, appCtx *shared.ApplicationContext, req any) (any, int) {
	type HelloResponse struct {
		Message string
	}

	rqst := req.(*Greeting)

	dbServer.PerformQuery(ctx, randomString(12))
	time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)

	emailServer.SendEmail(ctx, rqst.Name)
	time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)

	return HelloResponse{
		Message: fmt.Sprintf("Hello %s!", rqst.Name),
	}, http.StatusOK
}

type Database struct {
	appCtx shared.ApplicationContext
}

func (d Database) PerformQuery(ctx context.Context, queryName string) {
	span, _ := tracing.StartChildSpan(ctx, "PerformQuery")
	defer tracing.FinishChildSpan(span)

	span.SetTag("queryName", queryName)

	time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
}

type EmailServer struct {
	appCtx shared.ApplicationContext
}

func (e EmailServer) SendEmail(ctx context.Context, recipient string) {
	span, parentCtx := tracing.StartChildSpan(ctx, "SendEmail")
	defer tracing.FinishChildSpan(span)

	emailVendor.SendEmail(parentCtx, recipient)

	span.SetTag("recipient", recipient)

	time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
}

type EmailVendor struct {
	appCtx shared.ApplicationContext
}

func (e EmailVendor) SendEmail(ctx context.Context, recipient string) {
	span, _ := tracing.StartChildSpan(ctx, "TransmitEmail")
	defer tracing.FinishChildSpan(span)

	span.SetTag("recipient", recipient)

	time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
}

func randomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}

	return string(s)
}
