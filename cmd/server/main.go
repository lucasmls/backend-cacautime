package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/lucasmls/backend-cacautime/application/server"
	"github.com/lucasmls/backend-cacautime/domain/customers"
	"github.com/lucasmls/backend-cacautime/infra/postgres"
)

type config struct {
	goEnv                string
	dbDriver             string
	dbConnectionString   string
	dbMaxConnectionsOpen int
}

func env() (*config, error) {
	c := &config{
		goEnv:              os.Getenv("GO_ENV"),
		dbConnectionString: os.Getenv("DB_CONNECTION_STRING"),
		dbDriver:           os.Getenv("DB_DRIVER"),
	}

	dbMaxConnectionsOpen, err := strconv.Atoi(os.Getenv("DB_MAX_CONNECTIONS_OPEN"))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	c.dbMaxConnectionsOpen = dbMaxConnectionsOpen

	return c, nil
}

func main() {
	env, err := env()
	if err != nil {
		fmt.Print("Error when getting the environment variables.")
		return
	}

	postgres, err := postgres.NewClient(postgres.ClientInput{
		ConnectionString:   env.dbConnectionString,
		Driver:             env.dbDriver,
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

	s := server.NewService(server.ServiceInput{
		CustomersRepo: customers,
	})

	s.Run()
}
