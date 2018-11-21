package ripper


import "log"


type RawBlock struct {
	Len int32
	Data []byte
}

type SplitBlock struct {
	Id int32
	Block RawBlock
}

type Splitter struct {
	BlocksIn chan RawBlock
	BlocksOut chan SplitBlock
}

type Joiner struct {
	BlocksIn chan SplitBlock
	BlocksOut chan RawBlock
}

func (s Splitter) RunSplitter() {
	var i int32 = 0

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

func (j Joiner) RunJoiner() {
	var i int32 = 0
	var m map[int32]RawBlock = make(map[int32]RawBlock)
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

