package color

import (
	"image/color"
	"math"
)

/*
 * Maps a distribution to a series of colors.
 */
type Mapping interface {
	Map(counts []uint64) []color.NRGBA
}

/*
 * Restricts a value to an interval, so that min <= value <= max.
 */
func clamp(value float64, min float64, max float64) float64 {

	/*
	 * Decide on the value.
	 */
	if value < min {
		return min
	} else if value > max {
		return max
	} else {
		return value
	}

}

/*
 * Data structure representing a simple color mapping.
 */
type simpleMappingStruct struct {
	foreground color.NRGBA
}

/*
 * Data structure representing the default color mapping.
 */
type defaultMappingStruct struct {
}

/*
 * Map each count to a color value.
 */
func (this *simpleMappingStruct) Map(counts []uint64) []color.NRGBA {
	n := len(counts)
	colors := make([]color.NRGBA, n)
	fg := this.foreground

	/*
	 * Map each count in the distribution to a color value.
	 */
	for i, count := range counts {

		/*
		 * Check if there are dots in this cell.
		 */
		if count > 0 {
			colors[i] = fg
		}

	}

	return colors
}

/*
 * Map each count to a color value.
 */
func (this *defaultMappingStruct) Map(counts []uint64) []color.NRGBA {
	max := uint64(0)

	/*
	 * Iterate over the distribution.
	 */
	for _, count := range counts {

		/*
		 * If we found a larger value, make this the new maximum.
		 */
		if count > max {
			max = count
		}

	}

	maxFloat := float64(max)
	maxLog := math.Log(maxFloat)
	n := len(counts)
	colors := make([]color.NRGBA, n)

	/*
	 * Map each count in the distribution to a color value.
	 */
	for i, count := range counts {
		countFloat := float64(count)
		countLog := math.Log(countFloat)

		/*
		 * If the logarithm is finite, map to color scale.
		 */
		if !math.IsInf(countLog, 0) {
			frac := countLog / maxLog
			redFloat := float64(0.0)
			greenFloat := float64(0.0)
			blueFloat := float64(0.0)

			/*
			 * Map to a color.
			 */
			if frac <= 0.25 {
				diff := frac - 0.0
				greenFloat = 4.0 * diff
				blueFloat = 1.0
			} else if frac <= 0.5 {
				diff := frac - 0.25
				greenFloat = 1.0
				blueFloat = 1.0 - (4.0 * diff)
			} else if frac <= 0.75 {
				diff := frac - 0.5
				redFloat = 4.0 * diff
				greenFloat = 1.0
			} else if frac <= 1.0 {
				diff := frac - 0.75
				redFloat = 1.0
				greenFloat = 1.0
				blueFloat = 4.0 * diff
			} else {
				redFloat = 1.0
				greenFloat = 1.0
				blueFloat = 1.0
			}

			redFloat = math.Round(255.0 * redFloat)
			greenFloat = math.Round(255.0 * greenFloat)
			blueFloat = math.Round(255.0 * blueFloat)
			redFloat = clamp(redFloat, 0.0, 255.0)
			greenFloat = clamp(greenFloat, 0.0, 255.0)
			blueFloat = clamp(blueFloat, 0.0, 255.0)
			red := uint8(redFloat)
			green := uint8(greenFloat)
			blue := uint8(blueFloat)

			/*
			 * The resulting color.
			 */
			colors[i] = color.NRGBA{
				R: red,
				G: green,
				B: blue,
				A: 255,
			}

		}

	}

	return colors
}

/*
 * Create a new simple color mapping, which maps all cells with hits to a
 * predefined color.
 */
func SimpleMapping(red uint8, green uint8, blue uint8) Mapping {

	/*
	 * Create forground color.
	 */
	c := color.NRGBA {
		R: red,
		G: green,
		B: blue,
		A: 255,
	}

	/*
	 * Create simple color mapping.
	 */
	m := simpleMappingStruct{
		foreground: c,
	}

	return &m
}

/*
 * Create a new default color mapping.
 */
func DefaultMapping() Mapping {
	m := defaultMappingStruct{}
	return &m
}
