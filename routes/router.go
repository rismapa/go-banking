package routes

import (
	"fmt"

	"net/http"

	hand "github.com/rismapa/go-banking/adapter/handler"
	repo "github.com/rismapa/go-banking/adapter/repository"
	conf "github.com/rismapa/go-banking/config"
	"github.com/rismapa/go-banking/domain"
	serv "github.com/rismapa/go-banking/service"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	repoauth "github.com/rismapa/go-banking-auth/adapter/repository"
	middleware "github.com/rismapa/go-banking-auth/middleware"
	servauth "github.com/rismapa/go-banking-auth/service"
	logger "github.com/rismapa/go-banking-lib/config"
)

func NewRouter(router *mux.Router, db *sqlx.DB) {
	// apply middleware to all routes
	router.Use(middleware.ApiKeyMiddleware)
	customerRepo := repo.NewCustomerRepositoryDB(db)
	CustomerService := serv.NewCustomerService(customerRepo)
	customerHandler := hand.NewCustomerHandlerDB(CustomerService)
	accountRepo := repo.NewAccountRepositoryDB(db)
	accountService := serv.NewAccountService(accountRepo, customerRepo)
	accountHandler := hand.NewAccountHandlerDB(accountService)
	transactionRepo := repo.NewTransactionRepositoryDB(db)
	transactionService := serv.NewTransactionService(transactionRepo, accountRepo)
	transactionHandler := hand.NewTransactionHandlerDB(transactionService)
	authService := servauth.NewAuthService(repoauth.NewAccountRepositoryDB(db))

	router.Handle("/customers", middleware.AuthMiddleware(authService, http.HandlerFunc(customerHandler.GetCustomers))).Methods("GET")
	router.Handle("/customers/add", http.HandlerFunc(customerHandler.CreateCustomer)).Methods("POST")
	router.Handle("/customers/{id}", middleware.AuthMiddleware(authService, http.HandlerFunc(customerHandler.GetCustomerByID))).Methods("GET")
	router.Handle("/customers/{id}/edit", middleware.AuthMiddleware(authService, http.HandlerFunc(customerHandler.UpdateCustomer))).Methods("PUT")
	router.Handle("/accounts", middleware.AuthMiddleware(authService, http.HandlerFunc(accountHandler.GetAccounts))).Methods("GET")
	router.Handle("/accounts/add", middleware.AuthMiddleware(authService, http.HandlerFunc(accountHandler.CreateAccount))).Methods("POST")
	router.Handle("/accounts/{id}", middleware.AuthMiddleware(authService, http.HandlerFunc(accountHandler.GetAccountByID))).Methods("GET")
	router.Handle("/accounts/customer/{id}", middleware.AuthMiddleware(authService, http.HandlerFunc(accountHandler.GetAccountByCustomerID))).Methods("GET")
	router.Handle("/accounts/{id}/edit", middleware.AuthMiddleware(authService, http.HandlerFunc(accountHandler.UpdateAccount))).Methods("PUT")
	router.Handle("/accounts/{id}/delete", middleware.AuthMiddleware(authService, http.HandlerFunc(accountHandler.SoftDeleteAccount))).Methods("PUT")
	router.Handle("/transactions/add", middleware.AuthMiddleware(authService, http.HandlerFunc(transactionHandler.CreateTransaction))).Methods("POST")
	router.Handle("/transactions", middleware.AuthMiddleware(authService, http.HandlerFunc(transactionHandler.GetAllTransaction))).Methods("GET")
	router.Handle("/transactions/account/{id}", middleware.AuthMiddleware(authService, http.HandlerFunc(transactionHandler.GetTransactionByAccountID))).Methods("GET")
	/*
	 * Datanya diambil dari mock data
	 */
	repoCustMock := repo.NewCustomerRepositoryMock()
	svcCustMock := serv.NewCustomerService(repoCustMock)
	handCustMock := hand.NewCustomerHandler(svcCustMock)

	router.Handle("/mock/customers", middleware.AuthMiddleware(authService, http.HandlerFunc(handCustMock.GetCustomers))).Methods("GET")
	router.Handle("/mock/customers/add", middleware.AuthMiddleware(authService, http.HandlerFunc(handCustMock.AddCustomer))).Methods("POST")
}

/*
 * implementasi routing dari Kang Ari
 */
func StartServer() {

	// Start of log setup
	logger.InitiateLog()
	defer logger.CloseLog() // Close log when application is stopped
	// End of log setup

	config, _ := domain.GetConfig()
	port := config.Server.Port

	db, _ := conf.NewDBConnectionENV()

	defer db.Close()

	router := mux.NewRouter()

	NewRouter(router, db)

	fmt.Println("starting server on port " + port)

	http.ListenAndServe(":"+port, router)
}
