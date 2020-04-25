package projection

import (
	"github.com/andrepxx/sydney/coordinates"
	"math"
)

/*
 * Mathematical constants.
 */
const (
	MATH_HALF_PI    = 0.5 * math.Pi
	MATH_TWO_PI     = 2.0 * math.Pi
	MATH_QUARTER_PI = 0.25 * math.Pi
)

/*
 * Interface type representing a projection from a geographic location to a point
 * in a plane (surface of a map) and the other way round.
 */
type Projection interface {
	Forward(location coordinates.Geographic) coordinates.Cartesian
	Inverse(location coordinates.Cartesian) coordinates.Geographic
}

/*
 * Data structure representing the Mercator projection.
 */
type mercatorProjectionStruct struct {
}

/*
 * Project geographic coordinates in longitude and latitude to a point on a map using
 * the Mercator projection.
 */
func (this *mercatorProjectionStruct) Forward(loc coordinates.Geographic) coordinates.Cartesian {
	longitude := loc.Longitude()
	latitude := loc.Latitude()
	x := longitude / MATH_TWO_PI
	latA := 0.5 * latitude
	latB := MATH_QUARTER_PI + latA
	latC := math.Tan(latB)
	latD := math.Log(latC)
	y := latD / MATH_TWO_PI
	proj := coordinates.CreateCartesian(x, y)
	return proj
}

/*
 * Project a point on a map to geographic coordinates in longitude and latitude using
 * the Mercator projection.
 */
func (this *mercatorProjectionStruct) Inverse(loc coordinates.Cartesian) coordinates.Geographic {
	x := loc.X()
	y := loc.Y()
	longitude := MATH_TWO_PI * x
	yA := MATH_TWO_PI * y
	yB := math.Exp(yA)
	yC := math.Atan(yB)
	yD := 2.0 * yC
	latitude := yD - MATH_HALF_PI
	proj := coordinates.CreateGeographic(longitude, latitude)
	return proj
}

/*
 * Create a Mercator projection.
 */
func Mercator() Projection {
	proj := mercatorProjectionStruct{}
	return &proj
}
