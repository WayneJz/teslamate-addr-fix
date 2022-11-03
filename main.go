package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var (
	timeout  int
	proxy    string
	host     string
	port     string
	user     string
	db       string
	password string

	logger *log.Logger
)

const logName = "teslamate-addr-fix.log"

func init() {
	flag.StringVar(&proxy, "proxy", "", "http proxy (default use system proxy)")
	flag.IntVar(&timeout, "timeout", 5, "timeout of openstreetmap request")

	flag.StringVar(&host, "host", "127.0.0.1", "teslamate psql host")
	flag.StringVar(&port, "port", "5432", "teslamate psql port")
	flag.StringVar(&user, "user", "teslamate", "teslamate psql user")
	flag.StringVar(&db, "db", "teslamate", "teslamate psql database")
	flag.StringVar(&password, "password", "", "teslamate psql password")
}

func main() {
	flag.Parse()

	if password == "" {
		fmt.Println("must specify teslamate database password")
		return
	}

	if err := initPSql(host, port, user, password, db); err != nil {
		panic(err)
	}
	if err := initProxyCli(proxy, timeout); err != nil {
		panic(err)
	}

	log.SetFlags(log.LstdFlags)
	f, err := os.Create(logName)
	if err == nil {
		log.SetOutput(f)
	} else {
		log.SetOutput(os.Stdout)
	}

	saveBrokenAddr()
	fixAddrBroken()
}
