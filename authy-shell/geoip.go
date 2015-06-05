package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Geoip contains information about the geolocalization.
type Geoip struct {
	IP      string `json:"ip"`
	Country string `json:"country"`
	City    string `json:"city"`
}

func lookupMyIP() *Geoip {
	response, err := http.Get("http://www.telize.com/geoip")
	if err != nil {
		return nil
	}
	return parseGeoipResponse(response)
}

func lookupIP(ip string) *Geoip {
	response, err := http.Get(fmt.Sprintf("http://www.telize.com/geoip/%s", ip))
	if err != nil {
		return nil
	}

	return parseGeoipResponse(response)
}

func formatIPAndLocation(ip string) string {
	var geoip *Geoip
	if ip == "::1" || ip == "127.0.0.1" {
		geoip = lookupMyIP()
		if geoip != nil {
			ip = geoip.IP
		}
	} else {
		geoip = lookupIP(ip)
	}

	if geoip != nil {
		return fmt.Sprintf("%s (%s, %s)", ip, geoip.City, geoip.Country)
	}

	return ip
}

func parseGeoipResponse(response *http.Response) *Geoip {
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil
	}

	geoip := &Geoip{}
	err = json.Unmarshal(body, geoip)
	if err != nil {
		return nil
	}

	return geoip

}
