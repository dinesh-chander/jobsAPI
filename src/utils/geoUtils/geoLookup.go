package geoUtils

import (
	"config"
	"log"
	"logger"

	"golang.org/x/net/context"
	"googlemaps.github.io/maps"
)

var loggerInstance *log.Logger
var googleMapsClient *maps.Client

func GetLocationFromCoordinates(lat float64, lon float64, locationMap map[string]string) {

	if lat != 0.0 && lon != 0.0 {

		reverseGeoCodeRequest := &maps.GeocodingRequest{
			LatLng: &maps.LatLng{
				Lat: lat,
				Lng: lon,
			},
		}

		reverseGeoCodeResponse, requestErr := googleMapsClient.ReverseGeocode(context.Background(), reverseGeoCodeRequest)

		if requestErr != nil {
			loggerInstance.Println(requestErr.Error())
		} else {

			addressComponentsList := reverseGeoCodeResponse[0].AddressComponents

			for _, addressComponent := range addressComponentsList {
				for _, addressType := range addressComponent.Types {
					locationMap[addressType] = addressComponent.LongName
				}
			}
		}
	}
}

func GetLocationFromPlaceName(placeName string, locationMap map[string]string) {

	if placeName != "" {

		geoCodeRequest := &maps.GeocodingRequest{
			Address: placeName,
		}

		geoCodeResponse, requestErr := googleMapsClient.Geocode(context.Background(), geoCodeRequest)

		if requestErr != nil {
			loggerInstance.Println(requestErr.Error())
		} else {

			addressComponentsList := geoCodeResponse[0].AddressComponents

			for _, addressComponent := range addressComponentsList {
				for _, addressType := range addressComponent.Types {
					locationMap[addressType] = addressComponent.LongName
				}
			}
		}
	}
}

func init() {
	var googleMapsClientErr error

	loggerInstance = logger.Logger

	googleMapsClient, googleMapsClientErr = maps.NewClient(maps.WithAPIKey(config.GetConfig("googleGeoAPIKey")))

	if googleMapsClientErr != nil {
		loggerInstance.Panicln(googleMapsClientErr.Error())
	}
}
