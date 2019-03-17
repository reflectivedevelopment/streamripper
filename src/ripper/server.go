package ripper

import "io"
import "encoding/binary"
import "net"
import "os"
import "sync"

func ReadSocketSplitBlock(connId uint64, wg *sync.WaitGroup, conn net.Conn, out chan SplitBlock) {
	defer wg.Done()
	defer conn.Close()

	for {
		var id uint32
		err := binary.Read(conn, binary.LittleEndian, &id)
		if err != nil {
			break;
		}

		var l uint32
		err = binary.Read(conn, binary.LittleEndian, &l)
		if err != nil {
			break;
		}

		buf := make([]byte, l)
		_, err = io.ReadFull(conn, buf)
		if err != nil {
			break;
		}
		r := SplitBlock{ id, RawBlock{ l, buf } }

		out <- r
	}
}

func ReadSocketRawBlock(connId uint64, wg *sync.WaitGroup, conn net.Conn, out chan RawBlock, bufsize uint32) {
	defer wg.Done()
	defer conn.Close()

	for {
		buf := make([]byte, bufsize)
		l, err := conn.Read(buf)
		if err != nil {
			close(out)
			break;
		}
		r := RawBlock{ uint32(l), buf }

		out <- r
	}
}

func WriteSocketSplitBlock(connId uint64, wg *sync.WaitGroup, conn net.Conn, in chan SplitBlock) {
	defer wg.Done()
	defer conn.Close()

	for {
		data, ok := <- in;

		if !ok {
			break
		}

		err := binary.Write(conn, binary.LittleEndian, &data.Id)
		if err != nil {
			close(in)
			break;
		}

		err = binary.Write(conn, binary.LittleEndian, &data.Block.Len)
		if err != nil {
			close(in)
			break;
		}

		_, err = conn.Write(data.Block.Data)
		if err != nil {
			close(in)
			break;
		}
	}
}

func WriteSocketRawBlock(connId uint64, wg *sync.WaitGroup, conn net.Conn, in chan RawBlock) {
	defer wg.Done()
	defer conn.Close()

	for {
		data, ok := <- in;

		if !ok {
			break
		}

		_, err := conn.Write(data.Data)
		if err != nil {
			close(in)
			break;
		}
	}
}

/*
Reads in chan. Closes out chan when complete. 
*/
func RawBlockToStdOut(connId uint64, wg *sync.WaitGroup, in chan RawBlock) {
	defer wg.Done()
	defer os.Stdout.Close()
	defer close(in)

	for {
		data, ok := <- in;

		if !ok {
			break
		}

		_, err := os.Stdout.Write(data.Data)
		if err != nil {
			close(in)
			break;
		}		
	}
}

/*
Write out chan. Ignores in chan. 
*/
func StdInToRawBlock(connId uint64, wg *sync.WaitGroup, out chan RawBlock, bufsize uint32) {
	defer wg.Done()
	defer close(out)

	for {
		buf := make([]byte, bufsize)
		l, err := os.Stdin.Read(buf)
		if err != nil {
			break;
		}
		r := RawBlock{ uint32(l), buf }

		out <- r
	}
}
