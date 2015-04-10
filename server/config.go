package server

import (
	"errors"
	"flag"
	"fmt"
	"github.com/janqii/pusher/internal/version"
	"os"
	"strings"
	"time"
)

type PusherConfig struct {
	ID             int
	HttpServerPort string
	PrintVersion   bool

	HttpServerReadTimeout    time.Duration
	HttpServerWriteTimeout   time.Duration
	HttpServerMaxHeaderBytes int
	HttpKeepAliveEnabled     bool

	ZookeeperAddr    []string
	ZookeeperChroot  string
	ZookeeperTimeout time.Duration
}

func NewPusherConfig() (*PusherConfig, error) {
	cfg := new(PusherConfig)
	if err := cfg.Parse(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (cfg *PusherConfig) Parse() error {
	id := flag.Int("id", 0, "pusher id")

	httpServerPort := flag.String("http_port", "", "http server port")
	needPrintVersion := flag.Int("version", 0, "print version")

	httpServerReadTimeout := flag.Int64("http_server_read_timeout", 5000, "http server read timeout")
	httpServerWriteTimeout := flag.Int64("http_server_write_timeout", 5000, "http server write timeout")
	httpServerMaxHeaderBytes := flag.Int("http_server_max_header_bytes", 1<<20, "http server max header bytes")
	httpKeepAliveEnabled := flag.Int("http_keep_alive", 0, "http keep alive enable")

	zookeeperAddr := flag.String("zookeeper_addr", "", "zookeeper address")
	zookeeperChroot := flag.String("zookeeper_chroot", "", "zookeeper chroot")
	zookeeperTimeout := flag.Int64("zookeeper_timeout", 1, "zookeeper connect timeout")

	flag.Parse()

	cfg.PrintVersion = (*needPrintVersion > 0)
	if cfg.PrintVersion {
		fmt.Println(version.String("pusher"))
		os.Exit(1)
	}

	if *id <= 0 {
		return errors.New("id nil")
	}
	if *httpServerPort == "" {
		return errors.New("http_port nil")
	}
	if *zookeeperAddr == "" {
		return errors.New("zookeeper_addr nil")
	}
	if *zookeeperChroot == "" {
		return errors.New("zookeeper_chroot nil")
	}

	cfg.ID = *id
	cfg.HttpServerPort = *httpServerPort

	cfg.HttpServerReadTimeout = time.Duration(*httpServerReadTimeout) * time.Millisecond
	cfg.HttpServerWriteTimeout = time.Duration(*httpServerWriteTimeout) * time.Millisecond
	cfg.HttpServerMaxHeaderBytes = *httpServerMaxHeaderBytes
	cfg.HttpKeepAliveEnabled = (*httpKeepAliveEnabled > 0)

	cfg.ZookeeperAddr = strings.Split(*zookeeperAddr, ",")
	cfg.ZookeeperChroot = *zookeeperChroot
	cfg.ZookeeperTimeout = time.Duration(*zookeeperTimeout) * time.Second

	return nil
}
