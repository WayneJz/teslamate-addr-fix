package main

import "flag"

var (
	timeout  int
	proxy    string
	host     string
	port     string
	user     string
	db       string
	password string
)

func init() {
	flag.StringVar(&proxy, "proxy", "", "http proxy (default use system proxy)")
	flag.IntVar(&timeout, "timeout", 5, "timeout of openstreetmap request (default 5s")

	flag.StringVar(&host, "host", "127.0.0.1", "teslamate psql host (default 127.0.0.1)")
	flag.StringVar(&port, "port", "5432", "teslamate psql port (default 5432)")
	flag.StringVar(&user, "user", "teslamate", "teslamate psql user (default teslamate)")
	flag.StringVar(&db, "db", "teslamate", "teslamate psql database (default teslamate)")
	flag.StringVar(&password, "password", "", "teslamate psql password")
}

func main() {
	flag.Parse()
}