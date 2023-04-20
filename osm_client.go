package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"crypto/tls"
	"time"
)

const (
	osmReverseURL = "https://nominatim.openstreetmap.org/reverse?lat=%.6f&lon=%.6f&format=json"
)

var cli *http.Client

func initProxyCli(proxy string, timeoutSec int) error {

	timeout := time.Duration(timeoutSec) * time.Second
	proxyfunc := http.ProxyFromEnvironment
	if proxy != "" {
		fmt.Printf("using specific proxy: %v\n", proxy)
		u, err := url.Parse(proxy)
		if err != nil {
			return err
		}
		proxyfunc = http.ProxyURL(u)
	}
	cli = &http.Client{
		Transport: &http.Transport{
			Proxy: proxyfunc,
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: timeout,
	}
	return nil
}

type OsmRevAddress struct {
	PlaceID     int                    `json:"place_id"`
	Licence     string                 `json:"licence"`
	OsmType     string                 `json:"osm_type"`
	OsmID       int                    `json:"osm_id"`
	Lat         string                 `json:"lat"`
	Lon         string                 `json:"lon"`
	DisplayName string                 `json:"display_name"`
	Address     map[string]interface{} `json:"address"`
	Boundingbox []string               `json:"boundingbox"`
}

func getAddressByProxy(latitude, longitude float64) (*OsmRevAddress, error) {
	reqURL := fmt.Sprintf(osmReverseURL, latitude, longitude)
	rsp, err := cli.Get(reqURL)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	if rsp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("osm returns error status, status=%d, url=%s", rsp.StatusCode, reqURL)
	}
	b, err := io.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}

	address := &OsmRevAddress{}
	err = json.Unmarshal(b, address)
	if err != nil {
		return nil, err
	}
	return address, nil
}
