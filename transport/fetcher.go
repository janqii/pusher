package transport

import (
	"fmt"
)

type FetchManager struct {
}

func (f *FetchManager) Startup() error {
	fmt.Println("FetchManager startup")
	return nil
}

func (f *FetchManager) Shutdown() {
	fmt.Println("FetchManager shutdown")
}
