package scene

import (
	"fmt"
	"github.com/andrepxx/sydney/color"
	"github.com/andrepxx/sydney/coordinates"
	"image"
	"math"
)

/*
 * A scene is a plane onto which points are drawn.
 */
type Scene interface {
	Aggregate(data []coordinates.Cartesian)
	Clear()
	Render(mapping color.Mapping) (*image.NRGBA, error)
	Spread(amount uint8)
}

/*
 * Data structure representing a scene.
 */
type sceneStruct struct {
	bins   []uint64
	height uint32
	maxX   float64
	maxY   float64
	minX   float64
	minY   float64
	width  uint32
}

/*
 * Calculate a bin index based on a pair of (integer) coordinates.
 */
func (this *sceneStruct) index(x uint32, y uint32) (uint64, bool) {
	width := this.width
	height := this.height

	/*
	 * Check if coordinates are in valid range.
	 */
	if (x >= width) || (y >= height) {
		return 0, false
	} else {
		width64 := uint64(width)
		x64 := uint64(x)
		y64 := uint64(y)
		idx := (width64 * y64) + x64
		return idx, true
	}

}

/*
 * Aggregate data into the scene.
 */
func (this *sceneStruct) Aggregate(data []coordinates.Cartesian) {
	minX := this.minX
	maxX := this.maxX
	width := this.width
	widthFloat := float64(width)
	scaleX := widthFloat / (maxX - minX)
	minY := this.minY
	maxY := this.maxY
	height := this.height
	heightFloat := float64(height)
	scaleY := heightFloat / (maxY - minY)

	/*
	 * Iterate over all data points.
	 */
	for i := range data {
		point := &data[i]
		x := point.X()
		y := point.Y()

		/*
		 * Check if point lies within plot bounds.
		 */
		if ((x >= minX) && (x < maxX)) && ((y > minY) && (y <= maxY)) {
			plotX := uint32((x - minX) * scaleX)
			plotY := uint32((maxY - y) * scaleY)
			idx, ok := this.index(plotX, plotY)

			/*
			 * Check if point can be mapped to bin.
			 */
			if ok {
				val := this.bins[idx]

				/*
				 * Make sure we are not exceeding datatype bounds.
				 */
				if val < math.MaxUint32 {
					this.bins[idx] = val + 1
				}

			}

		}

	}

}

/*
 * Clear all data from the scene.
 */
func (this *sceneStruct) Clear() {
	bins := this.bins

	/*
	 * Reset the count in each bin to zero.
	 */
	for i := range bins {
		bins[i] = 0
	}

}

/*
 * Render a set of data points into an image using a color mapping.
 *
 * Generates an NRGBA-image of width times height pixels displaying
 * the data points with minX <= x < maxX and minY <= y < maxY.
 */
func (this *sceneStruct) Render(mapping color.Mapping) (*image.NRGBA, error) {

	/*
	 * Verify that color mapping is non-nil.
	 */
	if mapping == nil {
		return nil, fmt.Errorf("%s", "Color mapping must not be nil when rendering an image!")
	} else {
		data := this.bins
		colors := mapping.Map(data)

		/*
		 * Verify that color mapping returned non-nil slice.
		 */
		if colors == nil {
			return nil, fmt.Errorf("%s", "Color mapping must not map to nil slice when rendering an image!")
		} else {
			width := this.width
			widthInt := int(width)
			height := this.height
			heightInt := int(height)
			numColors := len(colors)
			expectedNumColors := widthInt * heightInt

			/*
			 * Verify that the color mapping returned a result of the
			 * expected length.
			 */
			if numColors != expectedNumColors {
				return nil, fmt.Errorf("%s", "Color mapping returned %d pixels, but expected %d for a (%d * %d) image.", numColors, expectedNumColors, width, height)
			} else {
				rect := image.Rect(0, 0, widthInt, heightInt)
				img := image.NewNRGBA(rect)

				/*
				 * Iterate over the rows of the image.
				 */
				for y := uint32(0); y < height; y++ {
					yy := int(y)

					/*
					 * Iterate over the columns of the image and set pixel data.
					 */
					for x := uint32(0); x < width; x++ {
						xx := int(x)
						idx, ok := this.index(x, y)

						/*
						 * Check if index is valid.
						 */
						if ok {
							c := colors[idx]
							img.SetNRGBA(xx, yy, c)
						}

					}

				}

				return img, nil
			}

		}

	}

}

/*
 * Spreads data over multiple cells.
 */
func (this *sceneStruct) Spread(amount uint8) {

	/*
	 * Only spread if needed.
	 */
	if amount > 0 {
		bins := this.bins
		numBins := len(bins)
		binsNew := make([]uint64, numBins)
		height := this.height
		width := this.width
		amount64 := int64(amount)

		/*
		 * Iterate over the target rows.
		 */
		for y := uint32(0); y < height; y++ {
			y64 := int64(y)

			/*
			 * Iterate over the target columns.
			 */
			for x := uint32(0); x < width; x++ {
				x64 := int64(x)
				sum := uint64(0)

				/*
				 * Spread across rows.
				 */
				for j := -amount64; j <= amount64; j++ {

					/*
					 * Spread across columns.
					 */
					for i := -amount64; i <= amount64; i++ {
						xx64 := x64 + i
						yy64 := y64 + j

						/*
						 * Check if values are in range.
						 */
						if xx64 >= 0 && xx64 <= math.MaxUint32 && yy64 >= 0 && yy64 <= math.MaxUint32 {
							xx := uint32(xx64)
							yy := uint32(yy64)
							idxSource, ok := this.index(xx, yy)
							sumOld := sum

							/*
							 * Check if index is in range.
							 */
							if ok {
								sum += bins[idxSource]

								/*
								 * Check for overflow.
								 */
								if sum < sumOld {
									sum = math.MaxUint64
								}

							}

						}

					}

				}

				idxTarget, ok := this.index(x, y)

				/*
				 * Check if index was calculated.
				 */
				if ok {
					binsNew[idxTarget] = sum
				}

			}

		}

		this.bins = binsNew
	}

}

/*
 * Create a new scene.
 */
func Create(width uint32, height uint32, minX float64, maxX float64, minY float64, maxY float64) Scene {
	width64 := uint64(width)
	height64 := uint64(height)
	numBins := width64 * height64
	bins := make([]uint64, numBins)

	/*
	 * Create scene data structure.
	 */
	scn := sceneStruct{
		bins:   bins,
		height: height,
		maxX:   maxX,
		maxY:   maxY,
		minX:   minX,
		minY:   minY,
		width:  width,
	}

	return &scn
}
