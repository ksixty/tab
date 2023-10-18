package main

import (
	"image"
	"image/color"
	"image/draw"
	"os"
	"fmt"
	"time"
	"strconv"

	mx "github.com/mcuadros/go-rpi-rgb-led-matrix"

	//"github.com/pbnjay/pixfont"
	"k60.in/go/tab/cozette"
	"k60.in/go/tab/cozetteb"
)

func check(b []byte, e error, default_value []byte) {
	if e != nil {
		b = default_value
	}
}

func weekdayRu(t time.Time) (int, string, bool) {
	weekday := int(t.Weekday())
	names := []string {
		"воскресенье", 
		"понедельник",
		"вторник",
		"среда",
		"четверг",
		"пятница",
		"суббота",
	}
	return weekday, names[weekday], weekday % 6 == 0
}

func render(c *mx.Canvas) {
	now := time.Now()
	date_str := fmt.Sprintf("%s", now.Format("2006-01-02"))
	weekday, weekday_str, isWeekend := weekdayRu(now)
	h_str := now.Format("15")
	m_str := now.Format("04")
	s_str := now.Format("05")
	weather, err := os.ReadFile("/mnt/tmp/weather")
	check(weather, err, []byte{})

	//clock1 := color.RGBA{76,  0,  0, 255}
	clock2 := color.RGBA{64, 16,  0, 255}
	clock3 := color.RGBA{60, 42,  0, 255}
	gray  := color.RGBA{13, 13, 13, 255}
	//gray2 := color.RGBA{13, 13, 13, 255}
	cyan  := color.RGBA{10, 15, 20, 255}
	//blue  := color.RGBA{20, 40, 255, 100}
	red   := color.RGBA{24, 4, 0, 100}

	dateColor := gray
	if isWeekend {
		dateColor = red
	}

	draw.Draw(c, c.Bounds(), &image.Uniform{color.RGBA{0, 0, 0, 255}}, image.ZP, draw.Src)

	//draw.Draw(c, image.Rect(7, 0, 22, 50), &image.Uniform{clock2}, image.ZP, draw.Src)

	mar_l := 5
	mar_b := 3
	cozetteb.Font.DrawString(c, mar_l+2,  0+mar_b,  h_str, clock2)
	cozetteb.Font.DrawString(c, mar_l+2, 12+mar_b,  m_str, clock2)
	cozetteb.Font.DrawString(c, mar_l+2, 24+mar_b,  s_str, clock3)

	cozette.Font.DrawString(c, mar_l+2+18,  0+mar_b,  date_str, gray)
	cozette.Font.DrawString(c, mar_l+2+18, 12+mar_b, weekday_str, dateColor)
	cozette.Font.DrawString(c, mar_l+2+18, 24+mar_b, string(weather), cyan)

	cozette.Font.DrawString(c, mar_l+2+18, 7+mar_b,  "-----  ", color.RGBA{3, 3, 3, 255})
	cozette.Font.DrawString(c, mar_l+2+18, 7+mar_b,  "     --", color.RGBA{8, 0, 0, 255})

	currentDayOverlay := "";
	for i := 0; i < 7; i++ {
		if i == ((weekday - 1) % 7) {
			currentDayOverlay += "-"
		} else {
			currentDayOverlay += " "
		}
	}
	if !isWeekend {
		dateColor = color.RGBA{20, 20, 20, 255}
	}
	cozette.Font.DrawString(c, mar_l+2+18, 7+mar_b, currentDayOverlay, dateColor)

	c.Render()
}

func main() {
	config := &mx.DefaultConfig
	config.Rows = 64
	config.Cols = 64
	config.ChainLength = 2
	config.Brightness = 20
	config.PWMLSBNanoseconds = 170

	m, err := mx.NewRGBLedMatrix(config)
	if err != nil {
		panic(err)
	}

	c := mx.NewCanvas(m)
	defer c.Close()
	for {
        brightnessStr, err := os.ReadFile("/mnt/tmp/brightness")
        check(brightnessStr, err, []byte{})
        brightness, err := strconv.ParseInt(string(brightnessStr), 10, 32)
        if err != nil {
            brightness = 1
        }
		render(c, brightness)
		time.Sleep(100 * time.Millisecond)
	}
}
