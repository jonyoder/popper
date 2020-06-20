package popper

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/fortytw2/leaktest"
	"gopkg.in/check.v1"
)

type PopperSuite struct {
	dir string
}

var _ = check.Suite(new(PopperSuite))

func Test(t *testing.T) { check.TestingT(t) }

func (s *PopperSuite) SetUpSuite(c *check.C) {
	var err error
	s.dir, err = ioutil.TempDir("", "")
	c.Assert(err, check.IsNil)
}

func (s *PopperSuite) TearDownSuite(c *check.C) {
	c.Assert(os.RemoveAll(s.dir), check.IsNil)
}

func (s *PopperSuite) SetUpTest(c *check.C) {
}

func (s *PopperSuite) TearDownTest(c *check.C) {
}

type SlowReader struct {
	reader io.Reader
	delay  time.Duration
}

func (s *SlowReader) Read(p []byte) (n int, err error) {
	time.Sleep(s.delay)
	return s.reader.Read(p)
}

func (s *PopperSuite) makeSourceFile(c *check.C) string {
	// Make a test file for the slow reader
	fSource := filepath.Join(s.dir, "test.source")
	func() {
		fS, err := os.Create(fSource)
		c.Assert(err, check.IsNil)
		defer fS.Close()
		for i := 0; i < 1000; i++ {
			fS.WriteString(fmt.Sprintf("Here is %3d\n", i))
		}
	}()
	return fSource
}
func (s *PopperSuite) TestNew(c *check.C) {
	defer leaktest.Check(c)

	// Make a test file for the slow reader, and
	// open it for reading. Normally, this reader
	// would be something like a response Body.
	fSource := s.makeSourceFile(c)
	fReadSource, err := os.Open(fSource)
	c.Assert(err, check.IsNil)
	defer fReadSource.Close()

	// Open a new destination file to write to
	fDestName := filepath.Join(s.dir, "test.write")
	fDest, err := os.Create(fDestName)
	c.Assert(err, check.IsNil)
	defer fDest.Close()

	// Create a new popper
	poppingReader := NewPopper(fReadSource)

	// Slowly read and copy to the new file. Each read will be
	// delayed by 1ms. Use a small buffer so we need to read in
	// many chunks
	done := make(chan struct{})
	go func() {
		_, err = io.CopyBuffer(fDest, &SlowReader{
			reader: poppingReader,
			delay:  time.Millisecond,
		}, make([]byte, 48))
		c.Assert(err, check.IsNil)
		close(done)
	}()

	// Wait a bit for writing to start
	time.Sleep(100 * time.Millisecond)

	// Open destination file before we complete writing to it
	fR, err := os.Open(fDestName)
	c.Assert(err, check.IsNil)
	defer fR.Close()

	// Get a new reader from popper
	pReader := poppingReader.NewReader(fR)

	// Read from popper.
	b, err := ioutil.ReadAll(pReader)
	c.Assert(err, check.IsNil)
	lines := strings.Split(string(b), "\n")
	lastLine := lines[len(lines)-2]
	log.Printf("Last line: %s", lastLine)

	// Wait for writing to complete
	now := time.Now()
	log.Printf("Done reading from file")
	<-done
	log.Printf("Done copying to file %dms later", time.Now().Sub(now).Milliseconds())
}
