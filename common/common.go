package common

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Pair struct {
	Timestamp time.Time
	Ask       float64
	Bid       float64
}

type Rate struct {
	Price  string `json:"price"`
	Volume string `json:"volume"`
	Amount string `json:"amount"`
	Factor string `json:"factor"`
	Type   string `json:"type"`
}

type Depth struct {
	Timestamp time.Time `json:"timestamp"`
	Asks      []Rate    `json:"asks"`
	Bids      []Rate    `json:"bids"`
}

var (
	DB        *sql.DB
	Market    string
	SourceAPI string
)

func Prepare() error {

	err := godotenv.Load("../.env")
	if err != nil {
		return err
	}

	Market = os.Getenv("MARKET")
	if Market == "" {
		return errors.New("could not get Market ID")
	}

	SourceAPI = os.Getenv("SOURCE_API")
	if SourceAPI == "" {
		return errors.New("could not get source URL")
	}

	DB, err = sql.Open("postgres", fmt.Sprintf(
		"postgres://%s:%s@db:%s/%s?sslmode=disable", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME")))
	if err != nil {
		return err
	}

	return nil
}
