package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"rates_service/common"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func statusHandler(w http.ResponseWriter, r *http.Request) {

	json.NewEncoder(w).Encode("status")
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {

	json.NewEncoder(w).Encode("default")

}

func main() {

	err := common.Prepare()
	if err != nil {
		log.Fatal("can not start: ", err)
	}
	defer common.DB.Close()

	ctx, cancel := InterceptOsSignals()
	defer cancel()

	log.Println(readPairsfromDb(fmt.Sprintf("SELECT timestamp, ask, bid FROM %s", common.Market)))
	storePairIntoDb(common.Market, common.Pair{Timestamp: time.Now(), Ask: 4.22, Bid: 4.33})
	log.Println(readPairsfromDb(fmt.Sprintf("SELECT timestamp, ask, bid FROM %s", common.Market)))

	run(ctx, cancel)

}

func run(context.Context, context.CancelFunc) {

	log.Print("Listening 8000")
	r := mux.NewRouter()
	r.HandleFunc("/", defaultHandler)
	r.HandleFunc("/status", statusHandler)
	log.Fatal(http.ListenAndServe(":8000", handlers.LoggingHandler(os.Stdout, r)))
}

func storePairIntoDb(table string, pair common.Pair) error {

	_, err := common.DB.Exec(fmt.Sprintf("INSERT INTO %s (timestamp, ask, bid) VALUES ($1, $2, $3);", table), pair.Timestamp, pair.Ask, pair.Bid)

	if err != nil {
		return err
	}

	return nil
}

func readPairsfromDb(query string) (pairs []common.Pair) {

	rows, err := common.DB.Query(query)
	if err != nil {
		log.Println("Error reading from db: ", err)
		return nil
	}

	for rows.Next() {
		var pair common.Pair
		err = rows.Scan(&pair.Timestamp, &pair.Ask, &pair.Bid)
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
