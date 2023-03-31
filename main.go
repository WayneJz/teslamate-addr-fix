package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

var (
	timeout  int
	proxy    string
	host     string
	port     string
	user     string
	db       string
	password string
	interval int
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
	flag.IntVar(&interval, "interval", 0, "interval (minutes) for running in daemon mode")
}

func loadEnvInDocker() {
	if os.Getenv("TESLAMATE_ADDR_FIX_ENV") != "docker" {
		return
	}
	getEnvStr("PROXY", &proxy)
	getEnvInt("OSM_TIMEOUT", &timeout)
	getEnvStr("DATABASE_HOST", &host)
	getEnvStr("DATABASE_PORT", &port)
	getEnvStr("DATABASE_USER", &user)
	getEnvStr("DATABASE_NAME", &db)
	getEnvStr("DATABASE_PASS", &password)
	getEnvInt("INTERVAL", &interval)

	// below only for docker
	if interval == 0 {
		interval = 5
	}
}

func main() {
	flag.Parse()
	loadEnvInDocker()

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

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	f, err := os.Create(logName)
	if err == nil {
		log.SetOutput(f)
	} else {
		log.SetOutput(os.Stdout)
	}

	saveBrokenAddr()
	fixAddrBroken()

	if interval > 0 {
		for range time.Tick(time.Minute * time.Duration(interval)) {
			saveBrokenAddr()
			fixAddrBroken()
		}
	}
}

func getEnvStr(env string, str *string) {
	if val := os.Getenv(env); val != "" {
		*str = val
	}
}

func getEnvInt(env string, num *int) {
	if val := os.Getenv(env); val != "" {
		if n, err := strconv.Atoi(val); err == nil {
			*num = n
		}
	}
}
