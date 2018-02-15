package main

import (
	"log"
	"os"

	"ireul.com/web"
)

func main() {
	log.SetPrefix("[saito] ")

	pm := NewPackageManager(PackageManagerOption{
		Domain:       os.Getenv("DOMAIN"),
		Token:        os.Getenv("GITHUB_TOKEN"),
		Organization: os.Getenv("GITHUB_ORG"),
	})
	pm.StartTicking()

	m := web.New()
	m.Use(web.Logger())
	m.Use(web.Recovery())
	m.Use(web.Static("public", web.StaticOptions{BinFS: true}))
	m.Use(web.Renderer(web.RenderOptions{BinFS: true}))

	m.Get("/", func(ctx *web.Context) {
		ctx.Data["Domain"] = pm.Option.Domain
		ctx.Data["Packages"] = pm.List()
		ctx.HTML(200, "index")
	})

	m.Get("/:name/?*", func(ctx *web.Context) {
		ctx.Data["Domain"] = pm.Option.Domain
		p := pm.Get(ctx.Params(":name"))
		if len(p.Name) == 0 {
			ctx.HTML(404, "404")
		} else {
			ctx.Data["Package"] = p
			ctx.HTML(200, "package")
		}
	})

	m.Run(os.Getenv("HOST"), os.Getenv("PORT"))
}
