package server

import (
	"crypto/tls"
	"encoding/gob"
	"fmt"
	"net"
	"net/http"
	"os"
	"path"

	"ConDetect/backend/i18n"

	"ConDetect/backend/init/app"
	"ConDetect/backend/init/business"

	"ConDetect/backend/cron"
	"ConDetect/backend/init/cache"
	"ConDetect/backend/init/session"
	"ConDetect/backend/init/session/psession"

	"ConDetect/backend/global"
	"ConDetect/backend/init/db"
	"ConDetect/backend/init/hook"
	"ConDetect/backend/init/log"
	"ConDetect/backend/init/migration"
	"ConDetect/backend/init/router"
	"ConDetect/backend/init/validator"
	"ConDetect/backend/init/viper"

	"github.com/gin-gonic/gin"
)

func Start() {
	viper.Init()
	i18n.Init()
	log.Init()
	db.Init()
	migration.Init()
	app.Init()
	validator.Init()
	gob.Register(psession.SessionUser{})
	cache.Init()
	session.Init()
	gin.SetMode("debug")
	cron.Run()
	InitOthers()
	business.Init()
	hook.Init()

	rootRouter := router.Routers()

	tcpItem := "tcp4"
	if global.CONF.System.Ipv6 == "enable" {
		tcpItem = "tcp"
		global.CONF.System.BindAddress = fmt.Sprintf("[%s]", global.CONF.System.BindAddress)
	}
	server := &http.Server{
		Addr:    global.CONF.System.BindAddress + ":" + global.CONF.System.Port,
		Handler: rootRouter,
	}
	ln, err := net.Listen(tcpItem, server.Addr)
	if err != nil {
		panic(err)
	}
	type tcpKeepAliveListener struct {
		*net.TCPListener
	}
	if global.CONF.System.SSL == "enable" {
		certPath := path.Join(global.CONF.System.BaseDir, "condetect/secret/server.crt")
		keyPath := path.Join(global.CONF.System.BaseDir, "condetect/secret/server.key")
		certificate, err := os.ReadFile(certPath)
		if err != nil {
			panic(err)
		}
		key, err := os.ReadFile(keyPath)
		if err != nil {
			panic(err)
		}
		cert, err := tls.X509KeyPair(certificate, key)
		if err != nil {
			panic(err)
		}
		server.TLSConfig = &tls.Config{
			Certificates: []tls.Certificate{cert},
		}
		global.LOG.Infof("listen at https://%s:%s [%s]", global.CONF.System.BindAddress, global.CONF.System.Port, tcpItem)

		if err := server.ServeTLS(tcpKeepAliveListener{ln.(*net.TCPListener)}, certPath, keyPath); err != nil {
			panic(err)
		}
	} else {
		global.LOG.Infof("listen at http://%s:%s [%s]", global.CONF.System.BindAddress, global.CONF.System.Port, tcpItem)
		if err := server.Serve(tcpKeepAliveListener{ln.(*net.TCPListener)}); err != nil {
			panic(err)
		}
	}
}
