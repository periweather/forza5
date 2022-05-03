package forza5

import (
	"context"
	"fmt"
	"io"
	"net"
	"time"
)

const maxBufferSize = 1024

func Server(ctx context.Context, address string) (err error) {

	conn, err := net.ListenPacket("udp", address)
	if err != nil {
		return
	}

	defer conn.Close()

	doneChan := make(chan error, 1)
	buffer := make([]byte, maxBufferSize)

	go func() {
		count := 0
		for {
			n, addr, err := conn.ReadFrom(buffer)
			if err != nil {
				doneChan <- err
				return
			}

			fmt.Printf("Packet received (Server) %d: bytes=%d from=%s buffer=%x\n", count, n, addr.String(), buffer[:n])

			//fmt.Printf("Buffer (Server): ", buffer)

			deadline := time.Now().Add(15 * time.Second)
			err = conn.SetWriteDeadline(deadline)
			if err != nil {
				doneChan <- err
				return
			}

			//n, err = conn.WriteTo(buffer[:n], addr)
			//if err != nil {
			//	doneChan <- err
			//	return
			//}
			//
			//fmt.Printf("Packet written (Server): bytes=%d to=%s\n", n, addr.String())

			count++
		}
	}()

	select {
	case <-ctx.Done():
		fmt.Println("Cancelled")
		err = ctx.Err()
	case err = <-doneChan:
	}

	return nil
}

func Client(ctx context.Context, address string, reader io.Reader) (err error) {
	raddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return
	}

	conn, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		return
	}

	defer conn.Close()

	doneChan := make(chan error, 1)

	go func() {
		n, err := io.Copy(conn, reader)
		if err != nil {
			doneChan <- err
			return
		}

		fmt.Printf("Packet Written (Client): bytes=%d\n", n)

		buffer := make([]byte, maxBufferSize)
		deadline := time.Now().Add(15 * time.Second)
		err = conn.SetReadDeadline(deadline)
		if err != nil {
			doneChan <- err
			return
		}

		nRead, addr, err := conn.ReadFrom(buffer)
		if err != nil {
			doneChan <- err
			return
		}

		fmt.Printf("Packet Received (Client): bytes=%d from=%s\n", nRead, addr.String())
	}()

	select {
	case <-ctx.Done():
		fmt.Println("Cancelled")
		err = ctx.Err()
	case err = <-doneChan:
	}

	return nil
}
