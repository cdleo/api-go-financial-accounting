package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"time"

	"github.com/cdleo/api-go-financial-accounting/cmd/config"
	"github.com/cdleo/api-go-financial-accounting/internal/handler"
	"github.com/cdleo/api-go-financial-accounting/internal/repository"
	"github.com/cdleo/api-go-financial-accounting/internal/service"
	"github.com/cdleo/api-go-financial-accounting/pkg/database"

	"github.com/gorilla/mux"
)

func main() {

	/*  CONFIGURATION  */

	configPathFlag := flag.String("p", "", "a string containing the full path of the configuration file")
	configFileNameFlag := flag.String("f", "", "a string containing the configuration filename")
	flag.Parse()

	configPath := *configPathFlag
	if configPath == "" {
		configPath = path.Join(filepath.Dir(os.Args[0]), "config")
	}
	configFileName := *configFileNameFlag
	if configFileName == "" {
		configFileName = "config.json"
	}
	configFilePath := path.Join(configPath, configFileName)

	config, err := config.GetAPIConfig(configFilePath)
	if err != nil {
		fmt.Printf("Unable to load configuration from %s. Error: %s\n", configPath, err)
		os.Exit(1)
	} else {
		fmt.Println("Configuration loaded")
	}

	dbClient := database.NewMongoDBClient(config.DB)
	if err := dbClient.Connect(context.TODO()); err != nil {
		fmt.Printf("Unable to connect to database. Error: %s\n", err)
		os.Exit(1)
	}
	defer dbClient.Disconnect()

	/*sqlProxy := sqlproxy.NewSQLProxyBuilder(connector.NewPostgreSqlConnector(config.DB.Host, config.DB.Port, config.DB.User, config.DB.Password, config.DB.DBName)).
		WithAdapter(adapter.NewNoopAdapter()).
		WithLogger(logger.NewNoLogLogger()).
		Build()

	sqlDB, err := sqlProxy.Open()
	if err != nil {
		fmt.Printf("Unable to connect to DB. Error: %v", err)
		os.Exit(1)
	}
	defer sqlProxy.Close()*/

	/*  REPOSITORIES  */
	accountsRepository := repository.NewAccountRepository(dbClient)
	budgetsRepository := repository.NewBudgetRepository(dbClient, accountsRepository)

	/*  SERVICES  */
	accountCreate := service.NewAccountCreate(accountsRepository)
	accountRetrieve := service.NewAccountRetrieve(accountsRepository)
	accountUpdate := service.NewAccountUpdate(accountsRepository)

	budgetCreate := service.NewBudgetCreate(budgetsRepository)
	budgetRetrieve := service.NewBudgetRetrieve(budgetsRepository, accountRetrieve)
	budgetUpdate := service.NewBudgetUpdate(budgetsRepository, accountUpdate)

	budgetSummaryRetrieve := service.NewBudgetSummaryRetrieve(budgetRetrieve)

	makeTransaction := service.NewMakeTransaction(accountsRepository)
	makeTransference := service.NewMakeTransference(accountsRepository)

	/*  HTTP HANDLERS  */
	var httpHandlers []handler.Handler
	httpHandlers = append(httpHandlers, handler.NewBudgetHandler(budgetCreate, budgetRetrieve, budgetUpdate))
	httpHandlers = append(httpHandlers, handler.NewBudgetSummaryHandler(budgetSummaryRetrieve))
	httpHandlers = append(httpHandlers, handler.NewAccountHandler(accountCreate, accountRetrieve, accountUpdate))
	httpHandlers = append(httpHandlers, handler.NewTransactionHandler(makeTransaction))
	httpHandlers = append(httpHandlers, handler.NewTransferenceHandler(makeTransference))

	httpRouter := mux.NewRouter()
	for _, handler := range httpHandlers {
		handler.RegisterEndpoints(httpRouter)
	}

	srv := &http.Server{
		Handler: httpRouter,
		Addr:    "127.0.0.1:" + strconv.Itoa(config.ServerPort),
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
