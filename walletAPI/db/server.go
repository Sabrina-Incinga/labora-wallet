package db_server

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/labora-wallet/walletAPI/controllers"
	"github.com/labora-wallet/walletAPI/db/variablesHandler"
	"github.com/labora-wallet/walletAPI/services"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
)

func validateDatabaseExistenceOrCreate(dbConfig variablesHandler.DbConfig, connection *sql.DB, err error) (bool, error) {
	var rowsAffected int64
	var response = rowsAffected != 0
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=disable",
		dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password)

	connection, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		return response, err
	}

	defer connection.Close()

	exists, err := checkDatabaseExists(connection, dbConfig.DbName)
	if err != nil {
		return response, err
	}

	if !exists {
		err = createDatabase(&response, connection, &rowsAffected)
		if err != nil {
			return response, err
		}
	} else {
		response = true
	}
	return response, nil
}

func checkDatabaseExists(db *sql.DB, dbname string) (bool, error) {
	var exists bool
	err := db.QueryRow("SELECT EXISTS (SELECT 1 FROM pg_database WHERE datname = $1)", dbname).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func createDatabase(response *bool, connection *sql.DB, rowsAffected *int64) error {
	scriptContent, err := ioutil.ReadFile("sql\\wallet_script_database.sql")
	if err != nil {
		*response = false
		return err
	}
	result, err := connection.Exec(string(scriptContent))
	if err != nil {
		*response = false
		return err
	}

	*rowsAffected, err = result.RowsAffected()
	if err != nil {
		*response = false
		return err
	}
	*response = true
	return nil
}

func getConnection() (*sql.DB, *variablesHandler.DbConfig, error) {
	var err error
	var connection *sql.DB
	dbConfig, err := variablesHandler.LoadEnvVariables()
	if err != nil {
		return nil, nil, err
	}

	validationResult, err := validateDatabaseExistenceOrCreate(dbConfig, connection, err)
	if err != nil {
		return nil, nil, err
	}

	if validationResult {
		psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.DbName)

		connection, err = sql.Open("postgres", psqlInfo)
		if err != nil {
			return nil, nil, err
		}
	}

	return connection, &dbConfig, nil
}

func createTables(connection *sql.DB) error {
	scriptContent, err := ioutil.ReadFile("sql\\wallet_script_tables.sql")
	if err != nil {
		return err
	}
	_, err = connection.Exec(string(scriptContent))

	if err != nil {
		return err
	}

	return nil
}

func startup(connection *sql.DB, dbConfig *variablesHandler.DbConfig) (*controllers.WalletController, *controllers.WalletTransactionController) {
	var mutex *sync.Mutex
	customerService := &services.PostgresCustomerDBHandler{Db: connection}
	walletService := &services.PostgresWalletDBHandler{Db: connection, Config: *dbConfig}
	walletTrackerService := &services.PostgresWalletTrackerDBHandler{Db: connection}
	walletAdministratorService := &services.PostgresWalletAdministrator{Db: connection, CustomerServiceImpl: customerService, WalletServiceImpl: walletService, WalletTrackerServiceImpl: walletTrackerService}
	walletController := &controllers.WalletController{CustomerServiceImpl: customerService, WalletServiceImpl: walletService, WalletTrackerServiceImpl: walletTrackerService, WalletCreationServiceImpl: walletAdministratorService}

	transactionService := &services.PostgresWalletTransactionDBHandler{Db: connection, WalletServiceImpl: walletService, Mutex: mutex}
	transactionController := &controllers.WalletTransactionController{WalletTransactionServiceImpl: transactionService}

	return walletController, transactionController
}

func StartServer() {
	connection, dbConfig, err := getConnection()
	if err != nil {
		log.Fatal(err)
	}

	defer connection.Close()

	err = createTables(connection)
	if err != nil {
		log.Fatal(err)
	}

	walletController, transactionController := startup(connection, dbConfig)

	router := mux.NewRouter()

	router.HandleFunc("/wallet", walletController.CreateWallet).Methods("POST")
	router.HandleFunc("/wallet/getStatusById/{id}", walletController.GetWalletStatus).Methods("GET")
	router.HandleFunc("/wallet/delete/{id}", walletController.DeleteWallet).Methods("DELETE")
	router.HandleFunc("/wallet/transaction/withdraw", transactionController.Withdraw).Methods("POST")
	router.HandleFunc("/wallet/transaction/add", transactionController.AddToAccount).Methods("POST")
	router.HandleFunc("/wallet/transaction/transfer", transactionController.Transfer).Methods("POST")

	corsOptions := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
	})

	handler := corsOptions.Handler(router)

	server := &http.Server{
		Addr:         ":8000",
		Handler:      handler,
		ReadTimeout:  40 * time.Second,
		WriteTimeout: 300 * time.Second,
	}
	fmt.Printf("api server listening at port %v", server.Addr)
	server.ListenAndServe()

}
