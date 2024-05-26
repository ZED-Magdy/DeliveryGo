package Models

import (
	"database/sql/driver"
	"errors"
	"fmt"

	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/wkt"
	"gorm.io/gorm"
)

type Region struct {
	gorm.Model
	Name        string
	Geofence    GeoPolygon `gorm:"type:geometry"`
	PriceListId uint       `gorm:"foreignKey:Id,constraint:OnDelete:SET NULL"`
	PriceList   PriceList  `gorm:"foreignKey:PriceListId"`
}
type GeoPolygon struct {
	Coordinates []GeoPoint
}

type GeoPoint struct {
	Lat float64
	Lng float64
}

func (p GeoPolygon) Value() (driver.Value, error) {
	if len(p.Coordinates) == 0 {
		return nil, nil
	}

	coords := make([][]float64, len(p.Coordinates))
	for i, point := range p.Coordinates {
		coords[i] = []float64{point.Lng, point.Lat}
	}

	geomPolygon := geom.NewPolygonFlat(geom.XY, flatten(coords), []int{len(coords)})

	wktString, err := wkt.Marshal(geomPolygon)
	if err != nil {
		return nil, err
	}

	return wktString, nil
}

func flatten(coords [][]float64) []float64 {
	var flat []float64
	for _, coord := range coords {
		flat = append(flat, coord...)
	}
	return flat
}

func (p *GeoPolygon) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	wktString := string(bytes)

	geometry, err := wkt.Unmarshal(wktString)
	if err != nil {
		return fmt.Errorf("failed to unmarshal WKT: %w", err)
	}

	polygon, ok := geometry.(*geom.Polygon)
	if !ok {
		return errors.New("geometry is not a Polygon")
	}

	p.Coordinates = make([]GeoPoint, 0)
	for _, ring := range polygon.Coords() {
		for _, coord := range ring {
			p.Coordinates = append(p.Coordinates, GeoPoint{
				Lat: coord.Y(),
				Lng: coord.X(),
			})
		}
	}

	return nil
}
