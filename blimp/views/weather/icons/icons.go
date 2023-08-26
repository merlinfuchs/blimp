package icons

import (
	"bytes"
	_ "embed"
	"fmt"
	"image"
	"image/png"
)

//go:embed 01d.png
var dayClearSky []byte

//go:embed 01n.png
var nightClearSky []byte

//go:embed 02d.png
var dayFewClouds []byte

//go:embed 02n.png
var nightFewClouds []byte

//go:embed 03d.png
var dayScatteredClouds []byte

//go:embed 03n.png
var nightScatteredClouds []byte

//go:embed 04d.png
var dayBrokenClouds []byte

//go:embed 04n.png
var nightBrokenClouds []byte

//go:embed 09d.png
var dayShowerRain []byte

//go:embed 09n.png
var nightShowerRain []byte

//go:embed 10d.png
var dayRain []byte

//go:embed 10n.png
var nightRain []byte

//go:embed 11d.png
var dayThunderstorm []byte

//go:embed 11n.png
var nightThunderstorm []byte

//go:embed 13d.png
var daySnow []byte

//go:embed 13n.png
var nightSnow []byte

//go:embed 50d.png
var dayMist []byte

//go:embed 50n.png
var nightMist []byte

var Icons = map[string][]byte{
	"01d": dayClearSky,
	"01n": nightClearSky,
	"02d": dayFewClouds,
	"02n": nightFewClouds,
	"03d": dayScatteredClouds,
	"03n": nightScatteredClouds,
	"04d": dayBrokenClouds,
	"04n": nightBrokenClouds,
	"09d": dayShowerRain,
	"09n": nightShowerRain,
	"10d": dayRain,
	"10n": nightRain,
	"11d": dayThunderstorm,
	"11n": nightThunderstorm,
	"13d": daySnow,
	"13n": nightSnow,
	"50d": dayMist,
	"50n": nightMist,
}

func IconImage(icon string) (image.Image, error) {
	raw, exists := Icons[icon]
	if !exists {
		return nil, fmt.Errorf("Icon %s does not exist", icon)
	}

	image, err := png.Decode(bytes.NewReader(raw))
	if err != nil {
		return nil, fmt.Errorf("Failed to decode icon %s, %w", icon, err)
	}

	return image, nil
}
