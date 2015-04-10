package action

import (
	"github.com/janqii/pusher/serializer"
	"io"
	"log"
	"net/http"
	"strings"
)

func SkipMessageAction(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	msgConverter := strings.Join(query["format"], "")
	if msgConverter == "" {
		msgConverter = "json"
	}

	s := serializer.Serializer{Converter: msgConverter}

	echo2client(w, s, nil)
}

func echo2client(w http.ResponseWriter, s serializer.Serializer, e error) {
	b, e := s.Marshal(map[string]interface{}{
		"errno":  0,
		"errmsg": "ok",
		"data":   nil,
	})
	if e != nil {
		log.Printf("marshal http response error, %v", e)
	} else {
		io.WriteString(w, string(b))
	}

	return
}
