package provider

import (
	"io"
	"sync"
)

func Download(wg *sync.WaitGroup, fileId string, resultChan chan io.Reader, errChan chan error) {
	// TODO: impl
}