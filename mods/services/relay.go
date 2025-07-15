package services

import (
	"io"
	"net"
	"sync"
)

func relay(src, dst net.Conn) {
	var wg sync.WaitGroup
	wg.Add(2)

	// Channel to signal when either connection is closed
	done := make(chan struct{})

	// Copy from src to dst
	go func() {
		defer wg.Done()
		defer dst.Close() // Ensure dst is closed when done

		_, err := io.Copy(dst, src)

		if err != nil {
			// Log error if needed, but don't panic
			select {
			case done <- struct{}{}: // Signal other goroutine to close
			default: // Channel might be closed or full
			}
		}
	}()

	// Copy from dst to src
	go func() {
		defer wg.Done()
		defer src.Close() // Ensure src is closed when done

		_, err := io.Copy(src, dst)

		if err != nil {
			// Log error if needed, but don't panic
			select {
			case done <- struct{}{}: // Signal other goroutine to close
			default: // Channel might be closed or full
			}
		}
	}()

	// Wait for one direction to complete, then close both connections
	go func() {
		<-done
		src.Close()
		dst.Close()
	}()

	// Wait for both copy operations to complete
	wg.Wait()
}
