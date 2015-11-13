package main

import (
	"fmt"
	"image"
	"image/gif"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"
)

// 64 x 48

func processImage(i image.Image) string {
	screen := [384]uint8{}
	bounds := i.Bounds()
	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y
	xr := float64(64) / float64(width)
	yr := float64(48) / float64(height)
	for y := 0; y < 48; y++ {
		for x := 0; x < 64; x++ {
			c := i.At(bounds.Min.X+int(float64(x)/xr), bounds.Min.Y+int(float64(y)/yr))
			r, g, b, _ := c.RGBA()
			on := int(float64(r+g+b)/3) > 0xffff/3
			if on {
				screen[x+(y/8)*64] |= 1 << uint(y%8)
			}
		}
	}
	cArray := ""
	for i, b := range screen {
		if i != 0 {
			cArray += ", "
			if i%16 == 0 {
				cArray += "\n"
			}
		}
		cArray += fmt.Sprintf("0x%02x", b)
	}
	return cArray
}

func main() {
	f, err := os.Open("test2.gif")
	if err != nil {
		panic(err)
	}

	header := [512]byte{}
	if _, err := f.Read(header[:]); err != nil {
		panic(err)
	}

	t := http.DetectContentType(header[:])
	log.Info(t)

	f.Seek(0, 0)

	cMatrix := "uint8_t bender [][384] = {{\n\t"
	if t == "image/gif" {
		g, err := gif.DecodeAll(f)
		if err != nil {
			panic(err)
		}
		for n, i := range g.Image {
			if n != 0 {
				cMatrix += "},\n{"
			}
			cMatrix += processImage(i)
		}
	} else if t != "application/octet-stream" {
		i, _, err := image.Decode(f)
		if err != nil {
			panic(err)
		}
		cMatrix += processImage(i)
	}
	cMatrix += "}};"
	fmt.Println(cMatrix)
}
