package main

import (
	"fmt"
	"io"
	"os"
	"runtime"

	"github.com/gaborkolarovics/hsplit/pkg/writer"
	"github.com/integrii/flaggy"
)

var (
	version = "unversioned"
	commit  string
	date    string

	fileName string
	prefix   string = "x"
	bits     uint   = 16
	minSize  int    = 1024
	fanout   uint   = 8
)

func main() {
	info := fmt.Sprintf(
		"%s\t[OS: %s; Arch: %s; commit: %s; date: %s]",
		version,
		runtime.GOOS,
		runtime.GOARCH,
		commit,
		date,
	)

	flaggy.SetName("hsplit")
	flaggy.SetDescription("split a file into pieces with rolling hash")
	flaggy.DefaultParser.AdditionalHelpPrepend = "\r\nDESCRIPTION\r\n\t" +
		"Output pieces of FILE to PREFIXaaaa, PREFIXaaab, ...; default trailing zero bits in the rolling checksum is 16, and default PREFIX is 'x'."

	flaggy.AddPositionalValue(&fileName, "FILE", 1, true, "filename_desc")
	flaggy.AddPositionalValue(&prefix, "PREFIX", 2, false, "prefix_desc")
	flaggy.UInt(&bits, "b", "bits", "SplitBits is the number of trailing zero bits in the rolling checksum required to produce a chunk.")
	flaggy.Int(&minSize, "m", "minSize", "MinSize is the minimum chunk size.")
	flaggy.UInt(&fanout, "f", "fanout", "Fan-out of the nodes in the tree produced.")

	flaggy.SetVersion(info)
	flaggy.Parse()

	f, err := os.Open(fileName)
	if err != nil {
		check(err)
	}
	defer f.Close()

	w := writer.NewWriter(writer.Prefix(prefix), writer.Bits(bits), writer.Fanout(fanout), writer.MinSize(minSize))
	_, err = io.Copy(w, f)
	if err != nil {
		check(err)
	}
	w.Close()

}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
