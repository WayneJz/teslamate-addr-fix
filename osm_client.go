package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
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
		u, err := url.Parse(proxy)
		if err != nil {
			return err
		}
		proxyfunc = http.ProxyURL(u)
	}
	cli = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			Proxy:           proxyfunc,
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

var lastQuery time.Time = time.Now()

func getAddressByProxy(latitude, longitude float64) (*OsmRevAddress, error) {

	if time.Now().Sub(lastQuery) < time.Second { // since osm banned frequent use HTTP 429
		time.Sleep(time.Second)
		lastQuery = time.Now()
	}

	reqURL := fmt.Sprintf(osmReverseURL, latitude, longitude)

	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", ua)
	req.Header.Set("Accept", "*/*")

	rsp, err := cli.Do(req)
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
