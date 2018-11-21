package ripper

import "ripper"
import "testing"

func TestSplitter(t *testing.T) {
	s := ripper.Splitter{make(chan ripper.RawBlock), make(chan ripper.SplitBlock)}

	go s.RunSplitter()

}



