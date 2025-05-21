package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
	"snippetbox.glebich/internal/models"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	snippets      *models.SnippetModel
	templateCache map[string]*template.Template
}

func main() {

	// для того, чтобы передавать флаги через CLI типа - go run ./cmd/web -addr=":8000"
	addr := flag.String("addr", ":8000", "HTTP network address")
	flag.Parse()

	infoLog := log.New(os.Stdout, "[INFO]\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "[ERROR]\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Часть с подключением пула соединений к БД
	dsn := "postgres://postgres:postgres@db:5432/snippetbox?sslmode=disable"
	db, err := openDB(dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	// создаю новый темплейт кэш
	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	app := &application{
		infoLog:       infoLog,
		errorLog:      errorLog,
		snippets:      &models.SnippetModel{DB: db},
		templateCache: templateCache,
	}

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	// мини настройка веб-сервера - адрес порта, поток записи ошибок и
	// глобальный обработчик запросов
	srv := &http.Server{
		Addr:         *addr,
		ErrorLog:     errorLog,
		Handler:      app.routes(),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// просто информационное сообщение о запуске сервера
	infoLog.Printf("Starting server on %s", *addr)
	// запуск прослушивания порта - на этом шаге программа останавливается (не завершается),
	// пока не упадет сервер
	// (начало работы сервера)
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	// в случае сбоя работы сервера программа завершается
	// выводится ошибка и os.Exit(1)
	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

/*
	dsn := "postgres://postgres:postgres@localhost:5432/snippetbox"
	ctx := context.Background()
	dbpool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		app.errorLog.Fatal(err)
	}
	defer dbpool.Close()
*/
