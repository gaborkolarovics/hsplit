package writer

import (
	"os"

	"github.com/bobg/hashsplit"
)

// Writer is an io.WriteCloser that splits its input with a hashsplit.Splitter,
// It additionally assembles those chunks into a tree with a hashsplit.TreeBuilder.
type Writer struct {
	spl    *hashsplit.Splitter
	tb     *hashsplit.TreeBuilder
	fanout uint
	prefix string
}

// Option is the type of an option passed to NewWriter.
type Option = func(*Writer)

// NewWriter produces a new Writer writing to the file.
func NewWriter(opts ...Option) *Writer {
	w := &Writer{
		fanout: 8,
		prefix: "",
	}

	var x int32 = 475254

	tb := hashsplit.TreeBuilder{
		F: func(node *hashsplit.TreeBuilderNode) (hashsplit.Node, error) {

			if len(node.Chunks) > 0 {
				x++

				f, err := os.Create(w.prefix + IntToLetters(x))
				if err != nil {
					return nil, err
				}

				for i, chunk := range node.Chunks {
					repr := chunk
					_, err := f.Write(chunk)
					if err != nil {
						return nil, err
					}
					node.Chunks[i] = repr
				}

				err = f.Close()
				if err != nil {
					return nil, err
				}

			}

			return node, nil
		},
	}
	w.tb = &tb

	spl := hashsplit.NewSplitter(func(bytes []byte, level uint) error {
		return tb.Add(bytes, level/w.fanout)
	})

	spl.MinSize = 1024
	spl.SplitBits = 16
	w.spl = spl

	for _, opt := range opts {
		opt(w)
	}
	return w
}

// Write implements io.Writer.
func (w *Writer) Write(inp []byte) (int, error) {
	return w.spl.Write(inp)
}

// Close implements io.Closer.
func (w *Writer) Close() error {
	var err error
	if w.tb == nil {
		return nil
	}
	err = w.spl.Close()
	if err != nil {
		return err
	}

	_, err = w.tb.Root()
	if err != nil {
		return err
	}

	w.tb = nil
	return nil
}

// Bits is an option for NewWriter that changes the number of trailing zero bits in the rolling checksum used to identify chunk boundaries.
// A chunk boundary occurs on average once every 2^n bytes.
// (But the actual median chunk size is the logarithm, base (2^n-1)/(2^n), of 0.5.)
// The default value for n is 16,
// producing a chunk boundary every 65,536 bytes,
// and a median chunk size of 45,426 bytes.
func Bits(n uint) Option {
	return func(w *Writer) {
		w.spl.SplitBits = n
	}
}

// MinSize is an option for NewWriter that sets a lower bound on the size of a chunk.
// No chunk may be smaller than this,
// except for the final one in the input stream.
// The value must be 64 or higher.
// The default is 1024.
func MinSize(n int) Option {
	return func(w *Writer) {
		w.spl.MinSize = n
	}
}

// Fanout is an option for NewWriter that can change the fanout of the nodes in the tree produced.
// The value of n must be 1 or higher.
// The default is 8.
// In a nutshell,
// nodes in the hashsplit tree will tend to have around 2n children each,
// because each chunk's "level" is reduced by dividing it by n.
// For more information see https://pkg.go.dev/github.com/bobg/hashsplit#TreeBuilder.Add.
func Fanout(n uint) Option {
	return func(w *Writer) {
		w.fanout = n
	}
}

// Prefix is an option for NewWriter that sets a file name prefeix to io.writer
func Prefix(s string) Option {
	return func(w *Writer) {
		w.prefix = s
	}
}

// IntToLetters is convert number to alpha to chunk sufix. Like aa, ab, ac..
func IntToLetters(number int32) (letters string) {
	number--
	if firstLetter := number / 26; firstLetter > 0 {
		letters += IntToLetters(firstLetter)
		letters += string('a' + number%26)
	} else {
		letters += string('a' + number)
	}

	return
}
