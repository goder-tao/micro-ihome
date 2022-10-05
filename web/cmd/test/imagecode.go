package main

import (
	"image/color"
	"image/png"
	"log"
	"net/http"
	"github.com/afocus/captcha"
)


func main() {
	c := captcha.New()

	if err := c.SetFont("/home/tao/Data/Software/project/go/project/IHome/web/assert/comic.ttf"); err != nil {
		log.Fatalln("SetFont: ", err)
	}

	c.SetSize(128, 64)

	c.SetDisturbance(captcha.MEDIUM)

	c.SetFrontColor(color.RGBA{255, 255, 255, 255})
	c.SetBkgColor(color.RGBA{255, 0, 0, 255}, color.RGBA{0, 0, 255, 255}, color.RGBA{0, 153, 0, 255})

	http.HandleFunc("/r", func(w http.ResponseWriter, r *http.Request) {
		img,str := c.Create(4,captcha.ALL)
		png.Encode(w, img)
		println(str)
	})

	http.ListenAndServe(":8085", nil)
}