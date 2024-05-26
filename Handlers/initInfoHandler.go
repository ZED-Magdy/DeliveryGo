package Handlers

import (
	"ZED-Magdy/Delivery-go/Dtos"
	"ZED-Magdy/Delivery-go/Models"
	"encoding/json"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

type InitInfoHandler struct {
	db gorm.DB
}

func NewInitInfoHandler(db gorm.DB) *InitInfoHandler {
	return &InitInfoHandler{
		db: db,
	}
}

func (h *InitInfoHandler) GetInitialInformation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	lat := r.URL.Query().Get("lat")
	lng := r.URL.Query().Get("lng")

	if lat == "" || lng == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "lat and lng are required"}`))
		return
	}

	region := Models.Region{}
	err := h.db.Raw(`
		SELECT
		id,
		name, 
		price_list_id, 
		ST_AsText(geofence) as geofence
		FROM regions 
		WHERE ST_Contains(geofence, ST_GeomFromText(?))`, "POINT("+lat+" "+lng+")").Scan(&region).Error

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "not found."}`))
		return
	}
	if region.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(strconv.FormatUint(uint64(region.PriceList.ID), 10)))
		w.Write([]byte(`{"error": "You are out of our coverage area."}`))
		return
	}
	geofence := toGeofenceDto(region.Geofence)
	regionDto := Dtos.RegionDto{
		Name:     region.Name,
		Geofence: geofence,
	}
	priceList := Models.PriceList{}
	h.db.First(&priceList, region.PriceListId)
	priceListDto := Dtos.PriceListDto{
		KmCost:           priceList.KmCost,
		CancellationCost: priceList.CancellationCost,
	}

	priceList_marshalled, _ := json.Marshal(priceListDto)
	region_marshalled, _ := json.Marshal(regionDto)

	w.Write([]byte("{\"region\":" + string(region_marshalled) + ",\"priceList\":" + string(priceList_marshalled) + "}"))
}

func toGeofenceDto(geofence Models.GeoPolygon) Dtos.GeofenceDto {
	geofenceDto := Dtos.GeofenceDto{}
	for _, point := range geofence.Coordinates {
		geofenceDto.Coordinates = append(geofenceDto.Coordinates, Dtos.GeoPointDto{Lat: point.Lat, Lng: point.Lng})
	}
	return geofenceDto
}
