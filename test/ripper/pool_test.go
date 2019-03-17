package ripper

import "bytes"
import "ripper"
import "strings"
import "testing"

func TestSplitter(t *testing.T) {
	s := ripper.Splitter{}
	s.Blocksize = 42
	s.BlocksIn = make(chan ripper.RawBlock, 100)
	s.BlocksOut = make(chan ripper.SplitBlock, 100)

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
	j := ripper.Joiner{}
	j.Blocksize = 42
	j.BlocksIn = make(chan ripper.SplitBlock, 100)
	j.BlocksOut = make(chan ripper.RawBlock, 100)

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

func TestIpsum(t *testing.T) {
	raw := `Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.`

	r := strings.NewReader(raw)

	s := ripper.Splitter{}
	s.Blocksize = 5
	s.BlocksIn = make(chan ripper.RawBlock, 100)
	s.BlocksOut = make(chan ripper.SplitBlock, 100)

	go s.RunSplitter()

	j := ripper.Joiner{}
	j.Blocksize = 5
	j.BlocksIn = s.BlocksOut
	j.BlocksOut = make(chan ripper.RawBlock, 100)

	go j.RunJoiner()

	s.AddIn(r)

	buf := new(bytes.Buffer)

	j.AddOut(buf)
	j.WaitOut.Wait()

	if buf.String() != raw {
		t.Fail()
	}
}