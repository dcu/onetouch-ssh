package utils

import (
	"fmt"
	"github.com/Jeffail/gabs"
	"io/ioutil"
	"net/http"
)

// GeoIP contains information about the geolocalization.
type GeoIP struct {
	data *gabs.Container
}

func NewGeoIP(ip string) (*GeoIP, error) {
	var data *gabs.Container
	var err error

	if ip == "::1" || ip == "127.0.0.1" {
		data, err = lookupMyIP()
	} else {
		data, err = lookupIP(ip)
	}

	if err != nil {
		return nil, err
	}

	geoip := &GeoIP{
		data: data,
	}

	return geoip, nil
}

func (geoip *GeoIP) City() string {
	city, ok := geoip.data.Search("city").Data().(string)
	if !ok {
		return "<unknown>"
	}

	return city
}

func (geoip *GeoIP) Country() string {
	country, ok := geoip.data.Search("country", "name").Data().(string)
	if !ok {
		return "<unknown>"
	}

	return country
}

func lookupMyIP() (*gabs.Container, error) {
	response, err := http.Get("https://geoip.nekudo.com/api")
	if err != nil {
		return nil, err
	}
	return parseGeoIPResponse(response)
}

func lookupIP(ip string) (*gabs.Container, error) {
	response, err := http.Get(fmt.Sprintf("https://geoip.nekudo.com/api/%s", ip))
	if err != nil {
		return nil, err
	}

	return parseGeoIPResponse(response)
}

func FormatIPAndLocation(ip string) string {
	geoip, err := NewGeoIP(ip)

	if err != nil {
		return fmt.Sprintf("%s (%s, %s)", ip, geoip.City(), geoip.Country())
	}

	return ip
}

func parseGeoIPResponse(response *http.Response) (*gabs.Container, error) {
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return gabs.ParseJSON(body)
}
