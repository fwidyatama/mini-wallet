package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"mini-wallet/domain/constant"
	"mini-wallet/domain/middleware"
	"mini-wallet/internal/features/wallet"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func generatePath(path string) string {
	return fmt.Sprintf("/api/v1/%s", path)
}

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})

	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()

	port := viper.Get(constant.DBPort).(string)

	dbPort, _ := strconv.Atoi(port)
	dbUser := viper.Get(constant.DBUser)
	dbName := viper.Get(constant.DBName)
	dbPassword := viper.Get(constant.DBPassword)
	dbHost := viper.Get(constant.DBHost)

	dsn := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		logrus.Error(err)
		panic(err)
	}

	logrus.Println("Success connect to db")

	walletRepository := wallet.NewRepository(db)

	walletService := wallet.NewService(walletRepository)

	walletHandler := wallet.NewHandler(walletService)

	r := httprouter.New()

	r.POST(generatePath("init"), walletHandler.InitAccount)
	r.POST(generatePath("wallet"), middleware.AuthMiddleware(walletHandler.EnableWallet))
	r.PATCH(generatePath("wallet"), middleware.AuthMiddleware(walletHandler.DisableWallet))
	r.GET(generatePath("wallet"), middleware.AuthMiddleware(walletHandler.GetWallet))
	r.POST(generatePath("wallet/deposits"), middleware.AuthMiddleware(walletHandler.DepositWallet))
	r.POST(generatePath("wallet/withdrawals"), middleware.AuthMiddleware(walletHandler.WithdrawWallet))
	r.GET(generatePath("wallet/transactions"), middleware.AuthMiddleware(walletHandler.GetTransaction))

	go func() {
		if err := http.ListenAndServe("127.0.0.1:3000", r); err != nil {
			logrus.Error(err)
		}
	}()

	gracefulStop := make(chan os.Signal, 1)
	signal.Notify(gracefulStop, os.Interrupt)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)

	<-gracefulStop

}
