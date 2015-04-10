package action

import (
	"github.com/janqii/pusher/serializer"
	//"log"
	"net/http"
	"strings"
)

func SetSubscriberAction(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	msgConverter := strings.Join(query["format"], "")
	if msgConverter == "" {
		msgConverter = "json"
	}

	s := serializer.Serializer{Converter: msgConverter}

	echo2client(w, s, nil)
}

func GetSubscriberAction(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	msgConverter := strings.Join(query["format"], "")
	if msgConverter == "" {
		msgConverter = "json"
	}

	s := serializer.Serializer{Converter: msgConverter}

	echo2client(w, s, nil)
}

func DelSubscriberAction(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	msgConverter := strings.Join(query["format"], "")
	if msgConverter == "" {
		msgConverter = "json"
	}

	s := serializer.Serializer{Converter: msgConverter}

	echo2client(w, s, nil)
}
