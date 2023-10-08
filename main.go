package main

import (
	"image"
	"image/color"
	"image/draw"
	"os"
	"fmt"
	"time"

	mx "github.com/mcuadros/go-rpi-rgb-led-matrix"
	"github.com/pbnjay/pixfont"
	//"k60.in/go/tab/japfon"
)

func check(b []byte, e error, default_value []byte) {
	if e != nil {
		b = default_value
	}
}

func render(c *mx.Canvas) {
	now := time.Now()
	_, week := now.ISOWeek()
	date_str := fmt.Sprintf("W%02d %s", week, now.Format("2006-01-02"))
	time_str := now.Format("15:04:05")
	weather, err := os.ReadFile("/mnt/tmp/weather")
	check(weather, err, []byte{})
	playing, err := os.ReadFile("/mnt/tmp/track")
	check(playing, err, []byte("(silence)"))

	draw.Draw(c, c.Bounds(), &image.Uniform{color.RGBA{0, 0, 0, 255}}, image.ZP, draw.Src)
	pixfont.DrawString(c, 0, 0, time_str, color.RGBA{255, 255, 255, 255})
	pixfont.DrawString(c, 0, 9, date_str, color.RGBA{255, 255, 255, 255})
	pixfont.DrawString(c, 0, 18, string(weather), color.RGBA{20, 40, 255, 255})

	for pos, ch := range string(playing) {
		canvasWidth := 128
		margin := 27
		offset := pos * 8
		pixfont.DrawString(c, offset % canvasWidth, margin + offset / canvasWidth * 9, string(ch), color.RGBA{255, 0, 0, 255})
	}
	c.Render()
}

func main() {
	config := &mx.DefaultConfig
	config.Rows = 64
	config.Cols = 64
	config.ChainLength = 2
	config.Brightness = 30

	m, err := mx.NewRGBLedMatrix(config)
	if err != nil {
		panic(err)
	}

	c := mx.NewCanvas(m)
	defer c.Close()
	for {
		render(c)
	}
}
