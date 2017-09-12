package utils

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Jeffail/gabs"
)

var (
	// GeoIPUnknown is returned when the city or country name can't be resolved.
	GeoIPUnknown = "<unknown>"

	errInvalidGeoip = errors.New("invalid geoip data")
)

// GeoIP contains information about the geolocalization.
type GeoIP struct {
	data *gabs.Container
}

// NewGeoIP returns a new instance of the Geoip struct
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

	if data == nil {
		return nil, errInvalidGeoip
	}

	geoip := &GeoIP{
		data: data,
	}

	return geoip, nil
}

// City returns the name of the city
func (geoip *GeoIP) City() string {
	if geoip.data == nil {
		return GeoIPUnknown
	}

	city, ok := geoip.data.Search("city").Data().(string)
	if !ok {
		return GeoIPUnknown
	}

	return city
}

// Country returns the name of the country
func (geoip *GeoIP) Country() string {
	if geoip.data == nil {
		return GeoIPUnknown
	}

	country, ok := geoip.data.Search("country", "name").Data().(string)
	if !ok {
		return GeoIPUnknown
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

// FormatIPAndLocation returns the ip with geo location for the given ip.
func FormatIPAndLocation(ip string) string {
	geoip, err := NewGeoIP(ip)

	if err != nil {
		return ip
	}

	return fmt.Sprintf("%s (%s, %s)", ip, geoip.City(), geoip.Country())
}

func parseGeoIPResponse(response *http.Response) (*gabs.Container, error) {
	defer func() {
		_ = response.Body.Close()
	}()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return gabs.ParseJSON(body)
}
