package router

import (
	"github.com/janqii/pusher/server/action"
	"net/http"
)

func ProxyServerRouter(mux map[string]func(http.ResponseWriter, *http.Request)) {
	mux["/subscriber/set"] = action.SetSubscriberAction
	mux["/subscriber/get"] = action.GetSubscriberAction
	mux["/subscriber/del"] = action.DelSubscriberAction
	mux["/message/skip"] = action.SkipMessageAction
}
