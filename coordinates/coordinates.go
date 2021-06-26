package coordinates

/*
 * Data structure representing geographic coordinates as longitude and latitude.
 *
 * By convention, values are stored in radians.
 *
 * Geographic locations are immutable.
 */
type Geographic struct {
	latitude  float64
	longitude float64
}

/*
 * Data structure representing a 2-dimensional vector in Cartesian coordinates.
 *
 * Vectors are immutable.
 */
type Cartesian struct {
	x float64
	y float64
}

/*
 * Returns the latitude value of this geographic location.
 * By convention, this value is in radians.
 */
func (this *Geographic) Latitude() float64 {
	return this.latitude
}

/*
 * Returns the longitude value of this geographic location.
 * By convention, this value is in radians.
 */
func (this *Geographic) Longitude() float64 {
	return this.longitude
}

/*
 * Returns the abscissa (x-coordinate) of this two-dimensional Cartesian vector.
 */
func (this *Cartesian) X() float64 {
	return this.x
}

/*
 * Returns the ordinate (y-coordinate) of this two-dimensional Cartesian vector.
 */
func (this *Cartesian) Y() float64 {
	return this.y
}

/*
 * Creates an immutable data structure storing geographic coordinates as longitude
 * and latitude.
 */
func CreateGeographic(longitude float64, latitude float64) Geographic {

	/*
	 * Create a new geographic location with longitude and latitude.
	 */
	geo := Geographic{
		latitude:  latitude,
		longitude: longitude,
	}

	return geo
}

/*
 * Creates an immutable data structure representing a two-dimensional vector in
 * Cartesian coordinates.
 */
func CreateCartesian(x float64, y float64) Cartesian {

	/*
	 * Create a new two-dimensional vector in Cartesian coordinates.
	 */
	vec := Cartesian{
		x: x,
		y: y,
	}

	return vec
}
