package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGeoIPCity(t *testing.T) {
	geoip, err := NewGeoIP("54.239.26.128")

	assert.Nil(t, err)
	assert.Equal(t, geoip.City(), "Ashburn")

	geoip, err = NewGeoIP("186.31.57.58")
	assert.Nil(t, err)
	assert.Equal(t, geoip.City(), GeoIPUnknown)
}

func TestGeoIPCountry(t *testing.T) {
	geoip, err := NewGeoIP("54.239.26.128")

	assert.Nil(t, err)
	assert.Equal(t, geoip.Country(), "United States")

	geoip, err = NewGeoIP("186.31.57.58")
	assert.Nil(t, err)
	assert.Equal(t, geoip.Country(), "Colombia")
}

func TestFormatIPAndLocation(t *testing.T) {
	formatted := FormatIPAndLocation("54.239.26.128")

	assert.Equal(t, formatted, "54.239.26.128 (Ashburn, United States)")
}
