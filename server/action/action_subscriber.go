package action

import (
	"github.com/janqii/pusher/admin"
	"github.com/janqii/pusher/serializer"
	"log"
	"net/http"
	"strings"
)

type AddSbRequest struct {
	Name   string
	Config admin.SubscriberConfig
}
type AddSbResponse struct {
	Errno  int
	Errmsg string
}

func AddSubscriberAction(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("read request body error")
		return
	}

	resData := AddSbResponse{
		Errno:  0,
		Errmsg: "ok",
	}

	var reqData AddSbRequest
	if err := json.Unmarshal(body, &reqData); err != nil {
		log.Printf("Unmarshal HttpRequest error: %v", err)
		resData = AddSbResponse{
			Errno:  -1,
			Errmsg: "Unmarshal HttpRequest error",
		}
	} else if err = global.SubManager.AddItem(reqData.Name, reqData.Config); err != nil {
		log.Printf("global.SubManager.AddItem(reqData.Name, reqData.Config) error: %v", err)
		resData = AddSubResponse{
			Errno:  -1,
			Errmsg: "SubscribeManager Additem error",
		}
	}
}

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
