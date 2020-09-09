package main

import (
	"fmt"
	"github.com/dollarkillerx/erguotou"
	"io/ioutil"
	"log"
)

func main() {
	app := erguotou.New()

	app.Get("/file/:path", func(ctx *erguotou.Context) {
		val, ex := ctx.PathValueString("path")
		if !ex {
			ctx.Json(404, erguotou.H{"message": "404"})
			return
		}

		file, err := ioutil.ReadFile(fmt.Sprintf("file/%s", val))
		if err != nil {
			ctx.Json(404, err)
			return
		}

		ctx.Write(200, file)
	})

	app.Post("/file/:path", func(ctx *erguotou.Context) {
		val, ex := ctx.PathValueString("path")
		if !ex {
			ctx.Json(404, erguotou.H{"message": "404"})
			return
		}
		fileBody := ctx.Ctx.Request.Body()
		if err := ioutil.WriteFile(fmt.Sprintf("file/%s", val), fileBody, 00666); err != nil {
			ctx.Json(500, err)
			return
		}
		ctx.Json(200, erguotou.H{"message": "success"})
	})

	if err := app.Run(erguotou.SetHost("0.0.0.0:8089"), erguotou.SetDebug(true)); err != nil {
		log.Fatalln(err)
	}
}
