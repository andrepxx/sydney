package projection

import (
	"fmt"
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
 * Interface type representing a projection from geographic locations to points
 * in a plane (surface of a map) and the other way round.
 */
type Projection interface {
	Forward(dst []coordinates.Cartesian, src []coordinates.Geographic) error
	ForwardSingle(dst *coordinates.Cartesian, src *coordinates.Geographic) error
	Inverse(dst []coordinates.Geographic, src []coordinates.Cartesian) error
	InverseSingle(dst *coordinates.Geographic, src *coordinates.Cartesian) error
}

/*
 * Data structure representing the Mercator projection.
 */
type mercatorProjectionStruct struct {
}

/*
 * Project geographic coordinates in longitude and latitude to points on a map
 * using the Mercator projection.
 */
func (this *mercatorProjectionStruct) Forward(dst []coordinates.Cartesian, src []coordinates.Geographic) error {
	numSrc := len(src)
	numDst := len(dst)

	/*
	 * Check if source and destination have same length.
	 */
	if numSrc != numDst {
		return fmt.Errorf("%s", "Source and destination must have same length")
	} else {

		/*
		 * Project all data points.
		 */
		for i := range src {
			srcPtr := &src[i]
			dstPtr := &dst[i]
			this.ForwardSingle(dstPtr, srcPtr)
		}

		return nil
	}

}

/*
 * Project geographic coordinates in longitude and latitude to a point on a map
 * using the Mercator projection.
 *
 * If src == nil or dst == nil, this is a no-op.
 */
func (this *mercatorProjectionStruct) ForwardSingle(dst *coordinates.Cartesian, src *coordinates.Geographic) error {

	/*
	 * Make sure source and destination are valid.
	 */
	if src == nil || dst == nil {
		return fmt.Errorf("%s", "Src and dst must be non-nil")
	} else {
		longitude := src.Longitude()
		latitude := src.Latitude()
		x := longitude / MATH_TWO_PI
		latA := 0.5 * latitude
		latB := MATH_QUARTER_PI + latA
		latC := math.Tan(latB)
		latD := math.Log(latC)
		y := latD / MATH_TWO_PI
		*dst = coordinates.CreateCartesian(x, y)
		return nil
	}

}

/*
 * Project points on a map to geographic coordinates in longitude and latitude
 * using the Mercator projection.
 */
func (this *mercatorProjectionStruct) Inverse(dst []coordinates.Geographic, src []coordinates.Cartesian) error {
	numSrc := len(src)
	numDst := len(dst)

	/*
	 * Check if source and destination have same length.
	 */
	if numSrc != numDst {
		return fmt.Errorf("%s", "Source and destination must have same length")
	} else {

		/*
		 * Project all data points.
		 */
		for i := range src {
			srcPtr := &src[i]
			dstPtr := &dst[i]
			this.InverseSingle(dstPtr, srcPtr)
		}

		return nil
	}

}

/*
 * Project a point on a map to geographic coordinates in longitude and latitude
 * using the Mercator projection.
 *
 * If src == nil or dst == nil, this is a no-op.
 */
func (this *mercatorProjectionStruct) InverseSingle(dst *coordinates.Geographic, src *coordinates.Cartesian) error {

	/*
	 * Make sure source and destination are valid.
	 */
	if src == nil || dst == nil {
		return fmt.Errorf("%s", "Src and dst must be non-nil")
	} else {
		x := src.X()
		y := src.Y()
		longitude := MATH_TWO_PI * x
		yA := MATH_TWO_PI * y
		yB := math.Exp(yA)
		yC := math.Atan(yB)
		yD := 2.0 * yC
		latitude := yD - MATH_HALF_PI
		*dst = coordinates.CreateGeographic(longitude, latitude)
		return nil
	}

}

/*
 * Create a Mercator projection.
 */
func Mercator() Projection {
	proj := mercatorProjectionStruct{}
	return &proj
}
