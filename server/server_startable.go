package server

import (
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
	defer httpServer.ShutDown()

	fetcher := &transport.FetcherManager{}
	fetcher.Startup()
	defer fetcher.Shutdown()

	pusher := &transport.PusherManager{}
	pusher.Startup()
	defer pusher.Shutdown()

	log.Println("Pusher is running...")
	wg.Wait()
	log.Println("Pusher is exiting...")

	return nil
}
