package main

import (
	"context"
	"flag"
	"fmt"
	"ginapp/conf"
	"ginapp/fslib"
	"ginapp/nettool"
	"ginapp/router"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	config := conf.GetConf()
	port := flag.String("p", "4500", "Public Server Port")
	dir := flag.String("d", "", "change directory")
	staticFS := flag.Bool("s", false, "static fs")
	staticBasePath := flag.String("s4", "", "static 404 path, e.g.: /a/404.html")
	usage := flag.Usage
	flag.Usage = func() {
		fmt.Printf("date: %s, version: %s\n", conf.BuildDate, conf.BuildVersion)
		usage()
	}
	flag.Parse()

	// chang directory
	if *dir != "" {
		//home, _ := os.UserHomeDir()
		if err := fslib.Chdir(*dir); err != nil {
			panic(err)
		}
	}

	// gin+port
	engine := gin.New()
	//  https://pkg.go.dev/github.com/gin-gonic/gin#readme-don-t-trust-all-proxies
	engine.SetTrustedProxies([]string{"127.0.0.1"})
	router.Register(engine, *staticFS, *staticBasePath)

	//http.Handle("/", r)
	if false {
		engine.Run(":" + *port)
	} else {
		// set timeout
		publicServer := &http.Server{
			Addr: fmt.Sprintf(":%s", *port),
			// https://ieftimov.com/posts/testing-in-go-test-doubles-by-example/
			// Handler: engine,
			// Handler:      http.TimeoutHandler(http.HandlerFunc(slowHandler), 1*time.Second, "Timeout!\n"),
			/*
				1. 无论handerTimeout还是client主动取消，handler本身还会继续执行,
				2. 不过可以用 case <-c.Request.Context().Done() 或 <- c.Done 判断context 是否关闭： https://github.com/gin-gonic/gin/issues/1452
			*/
			Handler: http.TimeoutHandler(engine, 60*time.Second, "Handler Timeout!\n"),
			// ReadTimeout: the maximum duration for reading the entire request, including the body
			ReadTimeout: config.Http.ReadTimeout,
			// WriteTimeout: the maximum duration before timing out writes of the response
			WriteTimeout: config.Http.WriteTimeout,
			// IdleTimetout: the maximum amount of time to wait for the next request when keep-alive is enabled
			IdleTimeout:       30 * time.Second,
			ReadHeaderTimeout: 2 * time.Second,
			// TLSConfig:         tlsConfig,
			MaxHeaderBytes: 1 << 20,
		}
		// service connections
		go func() {
			if err := publicServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("listen: %s\n", err)
			}
		}()
		log.Printf("curl -D- http://m:%s/status/400 \n", *port)
		log.Printf("curl -s %s:%s/f/w/my/a.log -F file1=@./test.log -F name=alex | jq\n", nettool.GetLocalIP(), *port)
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt)
		<-quit
		log.Println("Shutdown Server ...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := publicServer.Shutdown(ctx); err != nil {
			log.Fatal("Public Server gracefully shutdown:", err)
		}
	}

}
