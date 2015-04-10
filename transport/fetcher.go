package transport

import (
	"fmt"
)

type FetcherManager struct {
}

func (f *FetcherManager) Startup() error {
	fmt.Println("Fetcher startup")
	return nil
}

func (f *FetcherManager) Shutdown() {
	fmt.Println("Fetcher shutdown")
}
