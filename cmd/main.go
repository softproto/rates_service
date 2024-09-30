package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"rates_service/configs"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type Pair struct {
	ID        int
	Timestamp time.Time
	Ask       float64
	Bid       float64
}

var (
	db    *sql.DB
	table string
)

func statusHandler(w http.ResponseWriter, r *http.Request) {

	json.NewEncoder(w).Encode("status")
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {

	json.NewEncoder(w).Encode("default")

}

func main() {
	var err error

	configs.Load()

	ctx, cancel := InterceptOsSignals()
	defer cancel()

	table = os.Getenv("DB_TABLE")
	log.Println(table)

	db, err = sql.Open("postgres", fmt.Sprintf(
		"postgres://%s:%s@db:%s/%s?sslmode=disable", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME")))
	if err != nil {
		log.Fatal("could not connect to datbase")
	}
	defer db.Close()

	log.Println(readPairsfromDb(fmt.Sprintf("SELECT * FROM %s", table)))
	storePairIntoDb(table, Pair{0, time.Now(), 4.22, 4.33})
	log.Println(readPairsfromDb(fmt.Sprintf("SELECT * FROM %s", table)))

	run(ctx, cancel)

}

func run(context.Context, context.CancelFunc) {

	log.Print("Listening 8000")
	r := mux.NewRouter()
	r.HandleFunc("/", defaultHandler)
	r.HandleFunc("/status", statusHandler)
	log.Fatal(http.ListenAndServe(":8000", handlers.LoggingHandler(os.Stdout, r)))
}

func storePairIntoDb(table string, pair Pair) error {

	_, err := db.Exec(fmt.Sprintf("INSERT INTO %s (timestamp, ask, bid) VALUES ($1, $2, $3);", table), pair.Timestamp, pair.Ask, pair.Bid)

	if err != nil {
		return err
	}

	return nil
}

func readPairsfromDb(query string) (pairs []Pair) {

	rows, err := db.Query(query)
	if err != nil {
		log.Println("Error reading from db: ", err)
		return nil
	}

	for rows.Next() {
		var pair Pair
		err = rows.Scan(&pair.ID, &pair.Timestamp, &pair.Ask, &pair.Bid)
		if err != nil {
			log.Println("Error scaning rows from db: ", err)
			return nil
		}
		pairs = append(pairs, pair)
	}
	return pairs
}

func InterceptOsSignals() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		alarm := make(chan os.Signal, 1)
		signal.Notify(alarm, syscall.SIGINT, syscall.SIGTERM)
		<-alarm
		log.Println("catch SIGINT/SIGTERM signal...")
		cancel()
	}()
	return ctx, cancel
}
