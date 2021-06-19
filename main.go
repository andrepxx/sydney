package main

import (
	"fmt"
	"github.com/andrepxx/sydney/color"
	"github.com/andrepxx/sydney/coordinates"
	"github.com/andrepxx/sydney/scene"
	"math/rand"
	"image"
	imagecolor "image/color"
	"image/draw"
	"image/png"
	"os"
)

/*
 * Sample program demonstrating sydney graphics library.
 */
func main() {
	scn := scene.Create(800, 800, -5.0, 5.0, -5.0, 5.0)
	data := make([]coordinates.Cartesian, 1000)
	
	/*
	 * Create a total of a hundred thousand data points.
	 */
	for j := 0; j < 100; j++ {
	
		/*
		 * Generate some data.
		 */
		for i := range data {
			x := rand.NormFloat64()
			y := rand.NormFloat64()
			data[i] = coordinates.CreateCartesian(x, y)
		}
		
		scn.Aggregate(data)
	}
	
	scn.Spread(1)
	mapping := color.DefaultMapping()
	img, err := scn.Render(mapping)
	
	/*
	 * Check if an error occured during rendering.
	 */
	if err != nil {
		msg := err.Error()
		fmt.Printf("Something went wrong: %s\n", msg)
	} else {
		dim := image.Rect(0, 0, 800, 800)
		target := image.NewNRGBA(dim)
		
		/*
		 * The background color.
		 */
		c := imagecolor.NRGBA{
			R: 0,
			G: 0,
			B: 0,
			A: 255,
		}
		
		uniform := image.NewUniform(c)
		draw.Draw(target, dim, uniform, image.ZP, draw.Over)
		draw.Draw(target, dim, img, image.ZP, draw.Over)
		
		/*
		 * The PNG encoder.
		 */
		enc := png.Encoder{
			CompressionLevel: png.BestCompression,
		}
		
		fd, err := os.Create("output.png")
		
		/*
		 * Check if there was an error creating the file.
		 */
		if err != nil {
			msg := err.Error()
			fmt.Printf("Error creating output file: %s", msg)
		} else {
			enc.Encode(fd, target)
			fd.Close()
		}
		
	}
	
}
