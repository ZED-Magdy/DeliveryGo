package Dtos

type RegionDto struct {
	Name string `json:"name"`
	Geofence GeofenceDto `json:"geofence"`
}

type GeofenceDto struct {
	Coordinates []GeoPointDto `json:"coordinates"`
}

type GeoPointDto struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}