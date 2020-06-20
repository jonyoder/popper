package popper

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/fortytw2/leaktest"
	"gopkg.in/check.v1"
)

type PopperSuite struct{}

var _ = check.Suite(new(PopperSuite))

func Test(t *testing.T) { check.TestingT(t) }

func (s *PopperSuite) SetUpSuite(c *check.C) {
}

func (s *PopperSuite) TearDownSuite(c *check.C) {
}

func (s *PopperSuite) SetUpTest(c *check.C) {
}

func (s *PopperSuite) TearDownTest(c *check.C) {
}

func (s *PopperSuite) TestNew(c *check.C) {
	defer leaktest.Check(c)

	dir, err := ioutil.TempDir("", "")
	c.Assert(err, check.IsNil)
	defer os.RemoveAll(dir)

	fName := filepath.Join(dir, "test.write")
	fW, err := os.Create(fName)
	c.Assert(err, check.IsNil)
	defer fW.Close()

	// slowly write to file
	done := make(chan struct{})
	go func() {
		for i := 0; i < 1000; i++ {
			fW.WriteString(fmt.Sprintf("this is line %d\n", i))
			time.Sleep(1 * time.Millisecond)
		}
		close(done)
	}()

	// Wait for writing to start
	time.Sleep(100 * time.Millisecond)

	fR, err := os.Open(fName)
	c.Assert(err, check.IsNil)
	defer fR.Close()

	// Read from file. This will be faster than the write,
	// so we will hit an EOF before we are done writing
	// to the file, but we don't know exactly where/when.
	bufR := bufio.NewReader(fR)
	for {
		line, _, err := bufR.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Printf("Error: %s", err)
		}
		log.Printf("%s", string(line))
		time.Sleep(500 * time.Microsecond)
	}

	// Wait for writing to complete
	<-done
}
