package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"
)

var username = flag.String("username", "root", "mysql username")
var password = flag.String("password", "", "password")
var mysql_host = flag.String("host", "", "mysql hostname")
var mysql_port = flag.Int("port", 3306, "mysql port")
var mysql_database = flag.String("database", "", "database name")
var statsd_host = flag.String("statsd_host", "127.0.0.1:8125", "statsd host")

func reportMetric(key string, value string) error {
	i, err := strconv.Atoi(value)
	if err != nil {
		return nil
	}
	conn, err := net.Dial("udp", *statsd_host)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Write([]byte(fmt.Sprintf("%s:%d|g", key, i)))

	return nil
}

func main() {
	flag.Parse()

	http.HandleFunc("/services/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("PONG"))
	})
	go func() {
		http.ListenAndServe(":3000", nil)
	}()

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@%s/%s", *username, *password, *mysql_host, *mysql_database))
	if err != nil {
		fmt.Printf("Error in opening mysql connection: %v\n", err)
		os.Exit(1)
	}

	if err = db.Ping(); err != nil {
		fmt.Printf("Error in database conncetion: %s!\n", err)
		os.Exit(1)
	}

	defer db.Close()

	for {
		rows, err := db.Query("show status")
		if err != nil {
			fmt.Printf("Error in query: %v\n", err)
		} else {
			var key string
			var value string

			for rows.Next() {
				err = rows.Scan(&key, &value)
				if err != nil {
					fmt.Printf("error in scan: %v\n", err)
					continue
				}
				err = reportMetric(key, value)
				if err != nil {
					fmt.Printf("error in reporting: %v\n", err)
				}
			}

			rows.Close()
		}
		time.Sleep(30)
	}
}
