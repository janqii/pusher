package transport

import (
	"fmt"
)

type PusherManager struct {
}

func (p *PusherManager) Startup() error {
	fmt.Println("PusherManager startup")
	return nil
}

func (p *PusherManager) Shutdown() {
	fmt.Println("PusherManager shutdown")
}
