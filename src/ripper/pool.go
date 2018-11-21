package ripper


import "log"


type RawBlock struct {
	Len int32
	Data []byte
}

type SplitBlock struct {
	Id int32
	BlocksOut RawBlock
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

}

func (j Joiner) RunJointer() {

}

func Test() {
	log.Println("Hello World!")
}
