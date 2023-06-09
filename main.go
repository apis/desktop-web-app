package main

import (
	"context"
	"desktop-web-app/backend"
	"desktop-web-app/mainWindow"
	"desktop-web-app/staticAssets"
	"embed"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
	"net"
	"net/http"
	"os"
	"time"
)

const localHost = "127.0.0.1"
const anyPort = 0
const uiUrlPrefix = "/ui"
const apiUrlPrefix = "/api"
const embedFsRoot = "frontend/dist"
const defaultUiUrl = "index.html"

//go:embed frontend/dist
var embedFs embed.FS

var log *logrus.Logger

func main() {
	log = &logrus.Logger{
		Out:   os.Stderr,
		Level: logrus.DebugLevel,
		Formatter: &easy.Formatter{
			TimestampFormat: "2006-01-02 15:04:05",
			LogFormat:       "[%lvl%] %time% - %msg%",
		},
	}

	log.Infof("Application starting\n")

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", localHost, anyPort))
	if err != nil {
		log.Fatalf("Server net.Listen() failed: %s\n",
			fmt.Errorf("%w", err))
	}

	ephemeralPort := listener.Addr().(*net.TCPAddr).Port

	log.Infof("Ephemeral port: %d\n", ephemeralPort)

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Mount(apiUrlPrefix, backend.Router(log))
	router.Handle(uiUrlPrefix+"*", staticAssets.Handler(embedFs, embedFsRoot, uiUrlPrefix, defaultUiUrl))

	server := &http.Server{
		Handler: router,
	}

	go func() {
		log.Infof("Server starting\n")

		err := server.Serve(listener)
		if err != nil {
			if err != nil && err != http.ErrServerClosed {
				log.Fatalf("Server ListenAndServe() failed: %s\n",
					fmt.Errorf("%w", err))
			}
		}

		log.Infof("Server stopped\n")
	}()

	mainWindow.Run(localHost, ephemeralPort, uiUrlPrefix, log)

	log.Infof("Server stopping\n")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown() failed: %s\n",
			fmt.Errorf("%w", err))
	}

	log.Infof("Application stopped\n")
}
