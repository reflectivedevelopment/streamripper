package ripper

import "ripper"
import "testing"

func TestSplitter(t *testing.T) {
	s := ripper.Splitter{make(chan ripper.RawBlock, 100), make(chan ripper.SplitBlock, 100)}

	go s.RunSplitter()

	r := ripper.RawBlock{ 0, []byte{0} }

	s.BlocksIn <- r

	r = ripper.RawBlock{ 1, []byte{1} }

	s.BlocksIn <- r

	close(s.BlocksIn)

	o := <- s.BlocksOut

	if o.Id != 0 {
		t.Fail()
	}

	o = <- s.BlocksOut

	if o.Id != 1 {
		t.Fail()
	}

	_, more := <- s.BlocksOut

	if more {
		t.Fail()
	}
}



func TestJoiner(t *testing.T) {
	j := ripper.Joiner{make(chan ripper.SplitBlock, 100), make(chan ripper.RawBlock, 100)}

	go j.RunJoiner()

	r := ripper.SplitBlock{ 1, ripper.RawBlock{ 1, []byte{1} } }

	j.BlocksIn <- r

	r = ripper.SplitBlock{ 0, ripper.RawBlock{ 0, []byte{0} } }

	j.BlocksIn <- r

	close(j.BlocksIn)

	o := <- j.BlocksOut

	if o.Len != 0 {
		t.Fail()
	}

	o = <- j.BlocksOut

	if o.Len != 1 {
		t.Fail()
	}

	_, more := <- j.BlocksOut

	if more {
		t.Fail()
	}
}
