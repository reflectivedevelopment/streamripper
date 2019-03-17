package ripper

import "encoding/binary"
import "io"
import "sync"

type RawBlock struct {
	Len uint32
	Data []byte
}

type SplitBlock struct {
	Id uint32
	Block RawBlock
}

type Splitter struct {
	BlocksIn chan RawBlock
	BlocksOut chan SplitBlock
	Blocksize uint32
	WaitIn sync.WaitGroup
	WaitOut sync.WaitGroup
}

type Joiner struct {
	BlocksIn chan SplitBlock
	BlocksOut chan RawBlock
	Blocksize uint32 
	WaitIn sync.WaitGroup
	WaitOut sync.WaitGroup
}

func (s *Splitter) RunSplitter() {
	var i uint32 = 0

	for {
		rb, more := <- s.BlocksIn
		if !more {
			close(s.BlocksOut)
			break;
		}
		sb := SplitBlock{i, rb}
		s.BlocksOut <- sb
		i ++
	}
}

func (j *Joiner) RunJoiner() {
	var i uint32 = 0
	var m map[uint32]RawBlock = make(map[uint32]RawBlock)
	for {
		sb, more := <- j.BlocksIn
		// This may leave some items in the map, but this is fine as the incoming connection
		// aborted leaving us in a state that we cannot fulfill joining the records together
		if !more {
			close(j.BlocksOut)
			break;
		}
		/*
		TODO: it is possible that if a record is lost that the map will grow forever, 
		but what can we do about this?
		*/
		m[sb.Id] = sb.Block

		for {
			_, ok := m[i]
			if !ok {
				break
			}
			if ok {
				rb := m[i]
				delete(m, i)
				j.BlocksOut <- rb
				i ++
			}
		}
	}
}

func (s *Splitter) AddIn(in io.Reader) {
	s.WaitIn.Add(1)

	go func() {
		defer s.WaitIn.Done()
		defer close(s.BlocksIn)

		for {
			buf := make([]byte, s.Blocksize)
			l, res := in.Read(buf)
			if (res != nil) {
				break;
			}
			if (l > 0) {
				r := RawBlock{ uint32(l), buf[0:l] }
				s.BlocksIn <- r
			}
		}
	} ()
}

func (s *Splitter) AddOut(out io.Writer) {
	s.WaitOut.Add(1)

	go func() {
		defer s.WaitOut.Done()

		for 
		{
			d, ok := <- s.BlocksOut;

			if !ok {
				break
			}
				
			err := binary.Write(out, binary.LittleEndian, d.Id)
			if err != nil {
				break;
			}
			err = binary.Write(out, binary.LittleEndian, d.Block.Len)
			if err != nil {
				break;
			}
			_, err = out.Write(d.Block.Data)
			if err != nil {
				break;
			}
		}
	} ()
}

func (j *Joiner) AddIn(in io.Reader) {
	j.WaitIn.Add(1)

	go func() {
		defer j.WaitIn.Done()
		
		for {
			var id uint32
			err := binary.Read(in, binary.LittleEndian, &id)
			if err != nil {
				break;
			}

			var l uint32
			err = binary.Read(in, binary.LittleEndian, &l)
			if err != nil {
				break;
			}

			buf := make([]byte, l)
			_, err = io.ReadFull(in, buf)
			if err != nil {
				break;
			}
			r := SplitBlock{ id, RawBlock{ l, buf } }

			j.BlocksIn <- r
		}

	} ()
}

func (j *Joiner) AddOut(out io.Writer) {
	j.WaitOut.Add(1)

	go func() {
		defer j.WaitOut.Done()

		for {
			data, ok := <- j.BlocksOut;

			if !ok {
					break
			}

			_, err := out.Write(data.Data)
			if err != nil {
					break;
			}

		}
	} ()
}