package Application

import (
	"ZED-Magdy/Delivery-go/Application/Dtos"
	"ZED-Magdy/Delivery-go/Models"
	services "ZED-Magdy/Delivery-go/Services"
	"ZED-Magdy/Delivery-go/infrastructure/database"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type InitInfoHandler struct {
	dbService *database.Service
}

func NewInitInfoHandler(dbService *database.Service) *InitInfoHandler {
	return &InitInfoHandler{dbService}
}

func (h *InitInfoHandler) GetInitialInformation(c *fiber.Ctx) error {

	lat := c.Query("lat")
	lng := c.Query("lng")

	if lat == "" || lng == "" {
		return c.Status(http.StatusBadRequest).JSON(map[string]string{"message": "lat and lng are required"})
	}

	region := Models.Region{}
	err := h.dbService.Db.Table("regions").Select("regions.id, price_lists.id, regions.name, price_list_id, ST_AsText(geofence) as geofence").Where("ST_Contains(geofence, ST_GeomFromText(?))", "POINT("+lat+" "+lng+")").Joins("left join price_lists on price_lists.id = regions.price_list_id ").Scan(&region).Error

	if region.ID == 0 || err != nil {
		return c.Status(http.StatusNotFound).JSON(map[string]string{"message": "Region not found"})
	}

	user, err := services.NewAuthService(c, *h.dbService.Db).User()

	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(map[string]string{"message": "Unauthorized"})
	}

	h.dbService.Db.Exec(`Update users set current_region_id = ? where id = ?`, region.ID, user.ID)

	geofence := toGeofenceDto(region.Geofence)
	regionDto := Dtos.RegionDto{
		Name: region.Name,
	}
	priceList := Models.PriceList{}
	h.dbService.Db.First(&priceList, region.PriceListId)
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
