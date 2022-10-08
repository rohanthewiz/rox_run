package main

import (
	"fmt"
	"log"

	"github.com/rohanthewiz/rox"
	"github.com/valyala/fasthttp"
)

func main() {
	r := rox.New(rox.Options{
		Verbose: true,
		Port:    "3020",
		TLS: rox.TLSOpts{
			UseTLS:   false,
			CertFile: "/etc/letsencrypt/live/mysite.com/cert.pem",
			KeyFile:  "/etc/letsencrypt/live/mysite.com/privkey.pem",
		},
	})

	// var customHdlr fasthttp.RequestHandler = func(ctx *fasthttp.RequestCtx) {
	// 	ctx.Response.Header.Add("Content-Type", "text/html")
	// 	_, _ = ctx.WriteString("Yo. It's not found")
	// }
	// r.Options.CustomNotFoundHandler = &customHdlr

	// Logging middleware
	r.Use(
		func(ctx *fasthttp.RequestCtx) (ok bool) {
			log.Printf("MW:: Requested path: %s\n", ctx.Path())
			return true
		},
		fasthttp.StatusServiceUnavailable, // 503
	)

	// Auth middleware
	r.Use(
		func(ctx *fasthttp.RequestCtx) (ok bool) {
			authed := true // pretend we got a good response from our auth check
			if !authed {
				return false
			}
			log.Printf("MW:: You are authorized for: %s\n", ctx.Path())
			return true
		},
		fasthttp.StatusUnauthorized,
	)

	// Add routes for static files
	r.AddStaticFilesRoute("/images/", "artifacts/images", 1)
	r.AddStaticFilesRoute("/css/", "artifacts/css", 1)
	// r.AddStaticFilesRoute("/.well-known/acme-challenge/", "certs", 0) // great for letsEncrypt!

	r.Get("/", func(ctx *fasthttp.RequestCtx, params rox.Params) {
		println("We are here :-)")

		ctx.Response.Header.Add("Content-Type", "text/html")
		_, _ = ctx.WriteString(`
<img src="/images/tiger.jpg" style="float: right; width: 400px; height: 400px"/>
<h2>Hello there! Rox here.</h2>
`)
		_, _ = ctx.WriteString("<h3>Hello there! Susie here.</h3>")
	})

	r.Get("/hail/:name/mood/:whatmood", func(ctx *fasthttp.RequestCtx, params rox.Params) {
		ctx.Response.Header.Add("Content-Type", "text/html")
		name := params.ByName("name")
		mood := params.ByName("whatmood")

		if mood == "happy" {
			_, _ = ctx.WriteString(fmt.Sprintf("<h2>Hello there %s! Welcome!</h2>", name))
		} else {
			_, _ = ctx.WriteString(fmt.Sprintf("<h2>Hey %s. Go Away!</h2>", name))
		}
		_, _ = ctx.WriteString(`<img src="/images/tiger.jpg" style="float: right"/>`)
	})

	r.Get("/greet/:name", func(ctx *fasthttp.RequestCtx, params rox.Params) {
		ctx.Response.Header.Add("Content-Type", "text/html")
		_, _ = ctx.WriteString("Hey " + params.ByName("name") + "!")
	})
	r.Get("/greet/city", func(ctx *fasthttp.RequestCtx, params rox.Params) {
		ctx.Response.Header.Add("Content-Type", "text/html")
		_, _ = ctx.WriteString("Hey big city!")
	})
	r.Get("/greet/city/street", func(ctx *fasthttp.RequestCtx, params rox.Params) {
		ctx.Response.Header.Add("Content-Type", "text/html")
		_, _ = ctx.WriteString("Hey big city street!")
	})

	r.Serve()
}
