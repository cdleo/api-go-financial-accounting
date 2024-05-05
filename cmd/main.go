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
	"github.com/swaggest/openapi-go/openapi3"
	"github.com/swaggest/rest/gorillamux"

	"github.com/gorilla/mux"
)

func main() {

	/*  CONFIGURATION  */
	generateOASFlag := flag.String("g", "", "generate OAS definition in the provided path")
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

	/*  REPOSITORIES  */
	accountsRepository := repository.NewAccountRepository(dbClient)
	budgetsRepository := repository.NewBudgetRepository(dbClient, accountsRepository)

	/*  SERVICES  */
	budgetCreate := service.NewBudgetCreate(budgetsRepository)
	budgetRetrieve := service.NewBudgetRetrieve(budgetsRepository)
	budgetUpdate := service.NewBudgetUpdate(budgetsRepository)

	budgetSummaryRetrieve := service.NewBudgetSummaryRetrieve(budgetRetrieve)

	accountCreate := service.NewAccountCreate(accountsRepository)
	accountRetrieve := service.NewAccountRetrieve(accountsRepository, budgetsRepository)
	accountUpdate := service.NewAccountUpdate(accountsRepository)

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

	if len(*generateOASFlag) > 0 {
		if err := generateOASFile(*generateOASFlag, httpRouter); err != nil {
			fmt.Printf("Unable to generate OAS file. Error: %s\n", err)
		}
	}

	srv := &http.Server{
		Handler:      httpRouter,
		Addr:         "127.0.0.1:" + strconv.Itoa(config.Server.Port),
		WriteTimeout: time.Duration(config.Server.WriteTimeout) * time.Second,
		ReadTimeout:  time.Duration(config.Server.ReadTimeout) * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

func generateOASFile(oasPath string, router *mux.Router) error {

	// Setup OpenAPI schema.
	refl := openapi3.NewReflector()
	refl.SpecSchema().SetTitle("Financial Accounting API")
	refl.SpecSchema().SetVersion("v0.1.0")
	refl.SpecSchema().SetDescription("OAS for the financial accounting API")
	// Walk the router with OpenAPI collector.
	c := gorillamux.NewOpenAPICollector(refl)
	_ = router.Walk(c.Walker)

	// Get the resulting schema.
	yml, err := refl.Spec.MarshalYAML()
	if err != nil {
		return err
	}

	oasFilePath := oasPath + "/oas.yaml"
	f, err := os.Create(oasFilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(yml)
	return err
}
