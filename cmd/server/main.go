package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/lucasmls/backend-cacautime/application/server"
	"github.com/lucasmls/backend-cacautime/domain/auth"
	"github.com/lucasmls/backend-cacautime/domain/candies"
	"github.com/lucasmls/backend-cacautime/domain/customers"
	"github.com/lucasmls/backend-cacautime/domain/sales"
	"github.com/lucasmls/backend-cacautime/domain/users"
	"github.com/lucasmls/backend-cacautime/infra"
	"github.com/lucasmls/backend-cacautime/infra/bcrypt"
	"github.com/lucasmls/backend-cacautime/infra/errors"
	"github.com/lucasmls/backend-cacautime/infra/jwt"
	"github.com/lucasmls/backend-cacautime/infra/log"
	"github.com/lucasmls/backend-cacautime/infra/postgres"
)

type config struct {
	goEnv                infra.Environment
	logLevel             string
	dbConnectionString   string
	dbMaxConnectionsOpen int
	jwtSecret            string
	jwtExpirationInHours int
}

func env() (*config, *infra.Error) {
	const opName infra.OpName = "cmd/server.env"

	c := &config{
		goEnv:              infra.Environment(os.Getenv("GO_ENV")),
		dbConnectionString: os.Getenv("DB_CONNECTION_STRING"),
		jwtSecret:          os.Getenv("JWT_SECRET"),
		logLevel:           os.Getenv("LOG_LEVEL"),
	}

	dbMaxConnectionsOpen, err := strconv.Atoi(os.Getenv("DB_MAX_CONNECTIONS_OPEN"))
	if err != nil {
		return nil, errors.New(err, opName, infra.KindBadRequest)
	}

	c.dbMaxConnectionsOpen = dbMaxConnectionsOpen

	jwtExpirationInHours, err := strconv.Atoi(os.Getenv("JWT_EXPIRATION_IN_HOURS"))
	if err != nil {
		return nil, errors.New(err, opName, infra.KindBadRequest)
	}

	c.jwtExpirationInHours = jwtExpirationInHours

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
		Log:                log,
		ConnectionString:   env.dbConnectionString,
		MaxConnectionsOpen: env.dbMaxConnectionsOpen,
	})

	if err != nil {
		errors.Log(log, err)
		return
	}

	bcrypt, err := bcrypt.NewClient(bcrypt.ClientInput{
		Log: log,
	})

	if err != nil {
		errors.Log(log, err)
		return
	}

	jwt, err := jwt.NewClient(jwt.ClientInput{
		Log:    log,
		Secret: env.jwtSecret,
		TTL:    env.jwtExpirationInHours,
	})

	if err != nil {
		errors.Log(log, err)
		return
	}

	customers, err := customers.NewService(customers.ServiceInput{
		Db:  postgres,
		Log: log,
	})

	if err != nil {
		errors.Log(log, err)
		return
	}

	candiesR, err := candies.NewService(candies.ServiceInput{
		Db:  postgres,
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

	usersR, err := users.NewService(users.ServiceInput{
		Log: log,
		Db:  postgres,
	})

	if err != nil {
		errors.Log(log, err)
		return
	}

	authR, err := auth.NewService(auth.ServiceInput{
		Log:    log,
		Users:  usersR,
		Crypto: bcrypt,
		JWT:    jwt,
	})

	if err != nil {
		errors.Log(log, err)
		return
	}

	s, err := server.NewService(server.ServiceInput{
		Log:           log,
		CustomersRepo: customers,
		CandiesRepo:   candiesR,
		SalesRepo:     salesR,
		UsersRepo:     usersR,
		AuthRepo:      authR,
		Validator:     validator.New(),
		JwtSecret:     env.jwtSecret,
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
