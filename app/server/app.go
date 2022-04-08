package server

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/DarkSoul94/simple-rest-app/app"
	apphttp "github.com/DarkSoul94/simple-rest-app/app/delivery/http"
	apprepo "github.com/DarkSoul94/simple-rest-app/app/repo/mock"
	appusecase "github.com/DarkSoul94/simple-rest-app/app/usecase"
	micrologger "github.com/alexvelfr/micro-logger"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // required
	"github.com/spf13/viper"
)

// App ...
type App struct {
	appUC      app.Usecase
	appRepo    app.Repository
	httpServer *http.Server
}

// NewApp ...
func NewApp() *App {
	db := initDB()
	fmt.Println(db)
	repo := apprepo.NewRepo()
	uc := appusecase.NewUsecase(repo)
	return &App{
		appUC:   uc,
		appRepo: repo,
	}
}

// Run run application
func (a *App) Run(port string) error {
	defer a.appRepo.Close()

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(
		gin.RecoveryWithWriter(micrologger.GetWriter()),
	)
	if viper.GetBool("app.client.use") {
		router.Use(static.Serve("/", static.LocalFile(viper.GetString("app.client.dir"), true)))
	}

	apiRouter := router.Group("/api")
	apphttp.RegisterHTTPEndpoints(apiRouter, a.appUC)

	a.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	var l net.Listener
	var err error
	l, err = net.Listen("tcp", a.httpServer.Addr)
	if err != nil {
		panic(err)
	}

	go func(l net.Listener) {
		if err := a.httpServer.Serve(l); err != nil {
			log.Fatalf("Failed to listen and serve: %+v", err)
		}
	}(l)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	return a.httpServer.Shutdown(ctx)
}

func initDB() *sql.DB {
	dbString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		viper.GetString("app.db.host"),
		viper.GetString("app.db.port"),
		viper.GetString("app.db.login"),
		viper.GetString("app.db.pass"),
		viper.GetString("app.db.name"),
	)
	db, err := sql.Open(
		"postgres",
		dbString,
	)
	if err != nil {
		panic(err)
	}
	runMigrations(db)
	return db
}

func runMigrations(db *sql.DB) {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		panic(err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		viper.GetString("app.db.name"),
		driver)
	if err != nil {
		panic(err)
	}
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange && err != migrate.ErrNilVersion {
		fmt.Println(err)
	}
}
