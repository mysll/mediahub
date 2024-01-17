package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"mediahub/internal/conf"
	"net/http"
	"time"
)

var httpSrv *http.Server

func serve() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.LoggerWithWriter(log.StandardLogger().Out), gin.RecoveryWithWriter(log.StandardLogger().Out))
	initRouter(r)
	httpBase := fmt.Sprintf("%s:%d", conf.GetConfig().Server.Address, conf.GetConfig().Server.HttpPort)
	log.Infof("start HTTP server @ %s", httpBase)
	httpSrv = &http.Server{Addr: httpBase, Handler: r}
	go func() {
		err := httpSrv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("failed to start http: %s", err.Error())
		}
	}()
}

func shutdown() {
	if httpSrv != nil {
		log.Infof("closing HTTP server")
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		if err := httpSrv.Shutdown(ctx); err != nil {
			log.Fatal("HTTP server shutdown err: ", err)
		}

	}
}
