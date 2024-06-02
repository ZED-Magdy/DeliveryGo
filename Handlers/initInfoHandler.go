package Handlers

import (
	"ZED-Magdy/Delivery-go/Dtos"
	"ZED-Magdy/Delivery-go/Models"
	services "ZED-Magdy/Delivery-go/Services"
	"net/http"

	"github.com/gofiber/fiber/v2"
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

func (h *InitInfoHandler) GetInitialInformation(c *fiber.Ctx) error {

	lat := c.Query("lat")
	lng := c.Query("lng")

	if lat == "" || lng == "" {
		return c.Status(http.StatusBadRequest).JSON(map[string]string{"message": "lat and lng are required"})
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

	if region.ID == 0 || err != nil {
		return c.Status(http.StatusNotFound).JSON(map[string]string{"message": "Region not found"})
	}

	user, err := services.NewAuthService(c, h.db).User()

	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(map[string]string{"message": "Unauthorized"})
	}

	h.db.Exec(`Update users set current_region_id = ? where id = ?`, region.ID, user.ID)

	geofence := toGeofenceDto(region.Geofence)
	regionDto := Dtos.RegionDto{
		Name: region.Name,
	}
	priceList := Models.PriceList{}
	h.db.First(&priceList, region.PriceListId)
	priceListDto := Dtos.PriceListDto{
		KmCost:           priceList.KmCost,
		CancellationCost: priceList.CancellationCost,
	}

	return c.JSON(Dtos.InitInfoResponseDto{
		Region:    regionDto,
		Geofence:  geofence,
		PriceList: priceListDto,
	})
}

func toGeofenceDto(geofence Models.GeoPolygon) Dtos.GeofenceDto {
	geofenceDto := Dtos.GeofenceDto{}
	for _, point := range geofence.Coordinates {
		geofenceDto.Coordinates = append(geofenceDto.Coordinates, Dtos.GeoPointDto{Lat: point.Lat, Lng: point.Lng})
	}
	return geofenceDto
}
