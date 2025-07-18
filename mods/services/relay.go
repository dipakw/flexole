package services

import (
	"context"
	"io"
	"net"
	"sync"
)

func relay(ctx context.Context, src, dst net.Conn) {
	var wg sync.WaitGroup
	wg.Add(2)

	// Channel to signal when either connection is closed or context is done
	done := make(chan struct{}, 1)

	// Copy from src to dst
	go func() {
		defer wg.Done()
		defer dst.Close() // Ensure dst is closed when done

		_, err := io.Copy(dst, src)

		if err != nil {
			select {
			case done <- struct{}{}:
			default:
			}
		}
	}()

	// Copy from dst to src
	go func() {
		defer wg.Done()
		defer src.Close() // Ensure src is closed when done

		_, err := io.Copy(src, dst)

		if err != nil {
			select {
			case done <- struct{}{}:
			default:
			}
		}
	}()

	// Close both connections on either copy error or context cancellation
	go func() {
		select {
		case <-done:
		case <-ctx.Done():
		}
		src.Close()
		dst.Close()
	}()

	// Wait for both copy operations to complete
	wg.Wait()
}
