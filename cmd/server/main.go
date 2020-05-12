package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/lucasmls/backend-cacautime/application/server"
	"github.com/lucasmls/backend-cacautime/domain/customers"
	"github.com/lucasmls/backend-cacautime/infra"
	"github.com/lucasmls/backend-cacautime/infra/errors"
	"github.com/lucasmls/backend-cacautime/infra/postgres"
)

type config struct {
	goEnv                string
	dbConnectionString   string
	dbMaxConnectionsOpen int
}

func env() (*config, *infra.Error) {
	const opName infra.OpName = "cmd/server.env"

	c := &config{
		goEnv:              os.Getenv("GO_ENV"),
		dbConnectionString: os.Getenv("DB_CONNECTION_STRING"),
	}

	dbMaxConnectionsOpen, err := strconv.Atoi(os.Getenv("DB_MAX_CONNECTIONS_OPEN"))
	if err != nil {
		log.Fatal(err)
		return nil, errors.New(err, opName, infra.KindBadRequest)
	}

	c.dbMaxConnectionsOpen = dbMaxConnectionsOpen

	return c, nil
}

func main() {
	env, err := env()
	if err != nil {
		// @TODO => Create the error log method
		fmt.Println("Error when getting the environment variables.", err.Error())
		return
	}

	postgres, err := postgres.NewClient(postgres.ClientInput{
		ConnectionString:   env.dbConnectionString,
		MaxConnectionsOpen: env.dbMaxConnectionsOpen,
	})

	if err != nil {
		log.Panic(err)
		return
	}

	customers, err := customers.NewService(customers.ServiceInput{
		Db: postgres,
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	s, err := server.NewService(server.ServiceInput{
		CustomersRepo: customers,
	})

	if err != nil {
		// @TODO => Create the error log method
		fmt.Println("Error when creating the server.", err.Error())
		return
	}

	s.Run()
}
