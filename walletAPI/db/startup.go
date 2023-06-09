package db_server

import (
	"database/sql"
	"sync"

	"github.com/labora-wallet/walletAPI/controllers"
	"github.com/labora-wallet/walletAPI/db/variablesHandler"
	"github.com/labora-wallet/walletAPI/services"
)

func startup(connection *sql.DB, dbConfig *variablesHandler.DbConfig) (*controllers.WalletController, *controllers.WalletTransactionController) {
	var mutex sync.Mutex
	customerService := &services.PostgresCustomerDBHandler{Db: connection}
	walletService := &services.PostgresWalletDBHandler{Db: connection, Config: *dbConfig}
	walletTrackerService := &services.PostgresWalletTrackerDBHandler{Db: connection}
	walletAdministratorService := &services.PostgresWalletAdministrator{Db: connection, CustomerServiceImpl: customerService, WalletServiceImpl: walletService, WalletTrackerServiceImpl: walletTrackerService}
	walletController := &controllers.WalletController{CustomerServiceImpl: customerService, WalletServiceImpl: walletService, WalletTrackerServiceImpl: walletTrackerService, WalletAdministratorServiceImpl: walletAdministratorService}

	walletMovementService := &services.PostgresWalletMovementDBHandler{Db: connection}
	transactionService := &services.PostgresWalletTransactionDBHandler{Db: connection, WalletServiceImpl: walletService, WalletMovementServiceImpl: walletMovementService, WalletTrackerServiceImpl: walletTrackerService, Mutex: &mutex}
	transactionController := &controllers.WalletTransactionController{WalletTransactionServiceImpl: transactionService}

	return walletController, transactionController
}