package backend

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"net/http"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	"time"
)

var log *logrus.Logger

func Router(externalLog *logrus.Logger) chi.Router {
	log = externalLog
	router := chi.NewRouter()
	router.Get("/get-time", handleGetTime())
	router.Get("/time-event", handleTimeEvent())
	return router
}

func handleGetTime() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		log.Infof("GetTime request from client\n")

		writer.Header().Set("Content-Type", "application/json")

		message, err := json.Marshal(getCurrentTime())
		if err != nil {
			log.Errorf("json.Marshal() failed: %s\n",
				fmt.Errorf("%w", err))
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, err = writer.Write(message)
		if err != nil {
			log.Errorf("writer.Write() failed: %s\n",
				fmt.Errorf("%w", err))
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func getCurrentTime() string {
	return time.Now().Format(time.RFC3339)
}

func handleTimeEvent() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		log.Infof("TimeEvent request from client\n")

		conn, err := websocket.Accept(writer, request, &websocket.AcceptOptions{
			Subprotocols: []string{"time"},
		})
		if err != nil {
			log.Errorf("websocket.Accept() failed: %s\n",
				fmt.Errorf("%w", err))
			return
		}

		defer func() {
			err := conn.Close(websocket.StatusInternalError, "unexpected error")
			if err != nil {
				log.Errorf("conn.Close() failed: %s\n",
					fmt.Errorf("%w", err))
				return
			}

			log.Infof("websocket connection closed\n")
		}()

		for {
			//if websocket.CloseStatus(err) != -1 {
			//	log.Infof("websocket closed\n")
			//	return
			//}

			currentTime := getCurrentTime()

			log.Infof("about to send event message: %s\n", currentTime)

			err = wsjson.Write(request.Context(), conn, currentTime)
			if err != nil {
				log.Errorf("wsjson.Write() failed: %s\n",
					fmt.Errorf("%w", err))
				return
			}

			log.Infof("sent event message: %s\n", currentTime)

			time.Sleep(3 * time.Second)
		}
	}
}
