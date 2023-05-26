package db_server

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type DbConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DbName   string
	ApiKey   string
}

func loadEnvVariables() (DbConfig, error) {
	err := godotenv.Load("./config/global.env")
	if err != nil {
		return DbConfig{}, err
	}
	return DbConfig{
		Host:     os.Getenv("host"),
		Port:     os.Getenv("port"),
		User:     os.Getenv("user"),
		Password: os.Getenv("password"),
		DbName:   os.Getenv("dbname"),
		ApiKey:   os.Getenv("apiKey"),
	}, nil
}

func validateDatabaseExistenceOrCreate(dbConfig DbConfig, connection *sql.DB, err error) (bool, error) {
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

func getConnection() (*sql.DB, error) {
	var err error
	var connection *sql.DB
	dbConfig, err := loadEnvVariables()
	if err != nil {
		return connection, err
	}

	validationResult, err := validateDatabaseExistenceOrCreate(dbConfig, connection, err)
	if err != nil {
		return connection, err
	}

	if validationResult {
		psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.DbName)
	
		connection, err = sql.Open("postgres", psqlInfo)
		if err != nil {
			return connection, err
		}
	}

	return connection, nil
}

func createTables(connection *sql.DB) error{
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

func StartServer() {
	connection, err := getConnection()
	if err != nil {
		log.Fatal(err)
	}

	defer connection.Close()

	err = createTables(connection)
	if err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter()

	server := &http.Server{
		Addr:         ":8000",
		Handler:      router,
		ReadTimeout:  40 * time.Second,
		WriteTimeout: 300 * time.Second,
	}
	fmt.Printf("api server listening at port %v", server.Addr)
	server.ListenAndServe() 

}
