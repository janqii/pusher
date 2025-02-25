package server

import (
	"github.com/janqii/pusher/admin"
	"github.com/janqii/pusher/global"
	"github.com/janqii/pusher/server/router"
	"github.com/janqii/pusher/transport"
	"github.com/janqii/pusher/utils"
	"log"
	"net/http"
	"sync"
)

func Startable(cfg *PusherConfig) error {
	wg := new(sync.WaitGroup)

	var (
		err      error
		zkClient *utils.ZK
	)

	if zkClient, err = utils.NewZK(cfg.ZookeeperAddr, cfg.ZookeeperChroot, cfg.ZookeeperTimeout); err != nil {
		log.Printf("init zkClient error: %v", err)
		return err
	}
	defer zkClient.Close()

	if global.KafkaClient, err = global.NewKafkaClient(zkClient); err != nil {
		log.Printf("create kafka client error: %v", err)
		return err
	}
	defer global.KafkaClient.Close()

	var httpMux map[string]func(http.ResponseWriter, *http.Request)
	httpMux = make(map[string]func(http.ResponseWriter, *http.Request))

	fetchManager := &transport.FetchManager{}
	pushManager := &transport.PushManager{}
	global.SubsManager = &admin.SubscribeManager{
		ZkClient:      zkClient,
		ZkChroot:      cfg.ZookeeperChroot,
		SubscriberMap: make(map[string]*admin.Subscriber),
		SubscriberNum: 0,
		FetcherM:      fetchManager,
		PusherM:       pushManager,
		Wg:            wg,
	}

	httpServer := &HttpServer{
		Addr:            ":" + cfg.HttpServerPort,
		Handler:         &HttpHandler{Mux: httpMux},
		ReadTimeout:     cfg.HttpServerReadTimeout,
		WriteTimeout:    cfg.HttpServerWriteTimeout,
		MaxHeaderBytes:  cfg.HttpServerMaxHeaderBytes,
		KeepAliveEnable: cfg.HttpKeepAliveEnabled,
		RouterFunc:      router.ProxyServerRouter,
		Wg:              wg,
		Mux:             httpMux,
	}

	httpServer.Startup()
	defer httpServer.Shutdown()

	fetchManager.Startup()
	defer fetchManager.Shutdown()

	pushManager.Startup()
	defer pushManager.Shutdown()

	global.SubManager.Startup()
	defer subscribeManager.Shutdown()

	log.Println("Pusher is running...")
	wg.Wait()
	log.Println("Pusher is exiting...")

	return nil
}
