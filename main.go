package main

// Renders expanding spheres of random colors.
// Spheres always expand until all LEDs are within sphere.

import (
	"flag"
	"fmt"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/ottojo/blnk2"
	"github.com/ottojo/blnk2/vector"
	"math"
	"math/rand"
	"time"
)

var filename = flag.String("c", "/home/jonas/clients.json", "blnk System config file")
var framesPerSecond = flag.Float64("fps", 20, "Framerate")
var rate = flag.Float64("expansionRate", 1, "sphere radius expansion in m/s")

func main() {
	flag.Parse()
	system := blnk2.CreateFromFile(*filename)
	go system.Discovery()
	system.WaitForDiscovery()
	fmt.Println("Clients connected")

	var maxLedDistance float64 = 0
	ledCount := 0
	for pixel := system.Stage.First; pixel != nil; pixel = pixel.Next {
		ledCount++
		if pixel.Data.Position.Length() > maxLedDistance {
			maxLedDistance = pixel.Data.Position.Length()
		}
	}

	for {
		sphereOrigin := randVector(maxLedDistance)
		coloredLeds := 0
		radius := 0.0
		sphereColor := colorful.HappyColor()
		for coloredLeds < ledCount {
			coloredLeds = 0
			for pixel := system.Stage.First; pixel != nil; pixel = pixel.Next {
				if pixel.Data.Position.Minus(sphereOrigin).Length() < radius {
					pixel.Data.Color = sphereColor
					coloredLeds++
				}
			}
			system.Commit()
			radius += *rate / *framesPerSecond
			time.Sleep(time.Duration(1000.0 / *framesPerSecond) * time.Millisecond)
		}
	}
}

func randVector(maxLength float64) vector.Vec3 {
	phi := rand.Float64() * 2 * math.Pi
	theta := rand.Float64() * math.Pi
	r := rand.Float64() * maxLength
	sint := math.Sin(theta)
	return vector.Vec3{r * sint * math.Cos(phi), r * sint * math.Sin(phi), r * math.Cos(theta)}
}
