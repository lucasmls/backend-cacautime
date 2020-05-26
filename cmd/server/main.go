package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/lucasmls/backend-cacautime/application/server"
	"github.com/lucasmls/backend-cacautime/domain/candies"
	"github.com/lucasmls/backend-cacautime/domain/customers"
	"github.com/lucasmls/backend-cacautime/domain/duties"
	"github.com/lucasmls/backend-cacautime/domain/sales"
	"github.com/lucasmls/backend-cacautime/infra"
	"github.com/lucasmls/backend-cacautime/infra/errors"
	"github.com/lucasmls/backend-cacautime/infra/log"
	"github.com/lucasmls/backend-cacautime/infra/postgres"
	"github.com/lucasmls/backend-cacautime/infra/postgrex"
)

type config struct {
	goEnv                infra.Environment
	logLevel             string
	dbConnectionString   string
	dbMaxConnectionsOpen int
}

func env() (*config, *infra.Error) {
	const opName infra.OpName = "cmd/server.env"

	c := &config{
		goEnv:              infra.Environment(os.Getenv("GO_ENV")),
		dbConnectionString: os.Getenv("DB_CONNECTION_STRING"),
	}

	dbMaxConnectionsOpen, err := strconv.Atoi(os.Getenv("DB_MAX_CONNECTIONS_OPEN"))
	if err != nil {
		return nil, errors.New(err, opName, infra.KindBadRequest)
	}

	c.dbMaxConnectionsOpen = dbMaxConnectionsOpen

	return c, nil
}

func main() {
	ctx := context.Background()

	env, err := env()
	if err != nil {
		fmt.Println("Error when getting the environment variables.", err.Error())
		return
	}

	log, err := log.NewClient(log.ClientInput{
		GoEnv: infra.Environment(env.goEnv),
		Level: infra.Severity(env.logLevel),
	})

	if err != nil {
		errors.Log(log, err)
		return
	}

	postgres, err := postgres.NewClient(postgres.ClientInput{
		ConnectionString:   env.dbConnectionString,
		MaxConnectionsOpen: env.dbMaxConnectionsOpen,
	})

	if err != nil {
		errors.Log(log, err)
		return
	}

	postgrex, err := postgrex.NewClient(postgrex.ClientInput{
		Log:                log,
		ConnectionString:   env.dbConnectionString,
		MaxConnectionsOpen: env.dbMaxConnectionsOpen,
	})

	if err != nil {
		errors.Log(log, err)
		return
	}

	customers, err := customers.NewService(customers.ServiceInput{
		Db:  postgrex,
		Log: log,
	})

	if err != nil {
		errors.Log(log, err)
		return
	}

	duties, err := duties.NewService(duties.ServiceInput{
		Log: log,
		Db:  postgrex,
	})

	if err != nil {
		errors.Log(log, err)
		return
	}

	candiesR, err := candies.NewService(candies.ServiceInput{
		Db:  postgrex,
		Log: log,
	})

	if err != nil {
		errors.Log(log, err)
		return
	}

	salesR, err := sales.NewService(sales.ServiceInput{
		Db:  postgres,
		Log: log,
	})

	if err != nil {
		errors.Log(log, err)
		return
	}

	s, err := server.NewService(server.ServiceInput{
		Log:           log,
		CustomersRepo: customers,
		DutiesRepo:    duties,
		CandiesRepo:   candiesR,
		SalesRepo:     salesR,
	})

	if err != nil {
		errors.Log(log, err)
		return
	}

	ch := s.Run(ctx)
	for err := range ch {
		errors.Log(log, err)
		// @TODO -> Add some monitoring metrics here and in the whole application.
	}
}
