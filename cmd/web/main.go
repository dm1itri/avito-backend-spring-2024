package main

import (
	"backend_spring_2024/internal/models"
	"database/sql"
	"flag"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
)

type Application struct {
	infoLog      *log.Logger
	errorLog     *log.Logger
	banners      *models.BannerModel
	secretKeyJWT []byte
}

func openDB(dsn *string) (*sql.DB, error) {
	db, err := sql.Open("postgres", *dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func main() {
	app := &Application{
		infoLog:  log.New(os.Stdout, "INFO\t", log.LUTC|log.Ldate|log.Ltime),
		errorLog: log.New(os.Stderr, "ERROR\t", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile),
	}

	port := flag.String("port", ":4000", "HTTP port")
	secretKeyJWT := flag.String("secret-key", "avito-backend-spring-secret-key", "JWT secret key")
	dsn := flag.String("dsn", "postgres://postgres:mysecretpassword@localhost:5432/avito_tech_backend?sslmode=disable", "PostgreSQL data source name")
	flag.Parse()
	app.secretKeyJWT = []byte(*secretKeyJWT)

	if adminToken, err := GenerateToken("admin", app.secretKeyJWT); err == nil {
		app.infoLog.Println("Admin TOKEN:", adminToken)
	} else {
		app.errorLog.Fatal(err)
	}
	if userToken, err := GenerateToken("user", app.secretKeyJWT); err == nil {
		app.infoLog.Println("User TOKEN:", userToken)
	} else {
		app.errorLog.Fatal(err)
	}

	db, err := openDB(dsn)
	if err != nil {
		app.errorLog.Fatal(err)
	}
	defer db.Close()

	app.banners = &models.BannerModel{DB: db}
	app.infoLog.Println("Starting server on", *port)
	server := http.Server{
		Addr:     *port,
		ErrorLog: app.errorLog,
		Handler:  app.routes(),
	}

	err = server.ListenAndServe()
	app.errorLog.Fatal(err)
}
