package transport

import (
	"fmt"
)

type PushManager struct {
}

func (p *PushManager) Startup() error {
	fmt.Println("PushManager startup")
	return nil
}

func (p *PushManager) Shutdown() {
	fmt.Println("PushManager shutdown")
}
