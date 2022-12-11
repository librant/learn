package writev

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/google/vectorio"
)

func TestWritev(t *testing.T) {
	f, err := os.Create("testwritev")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	data := [][]byte{[]byte("foo"), []byte("bar"), []byte("baz")}

	nw, err := vectorio.Writev(f, data)
	f.Seek(0, 0)
	if err != nil {
		t.Errorf("WritevRaw threw error: %s", err)
	}

	if nw != 9 {
		t.Errorf("Length %d of input does not match %d written bytes", 9, nw)
	}

	fromdisk, err := ioutil.ReadAll(f)
	if err != nil {
		t.Errorf("can't read file back, %s", err)
	}
	should := []byte("foobarbaz")
	if bytes.Compare(fromdisk, should) != 0 {
		t.Errorf("contents of file don't match input, %s != %s", fromdisk, should)
	}
	os.Remove("testwritev")
}

func TestWritevSocket(t *testing.T) {
	go func() {
		ln, err := net.Listen("tcp", "127.0.0.1:9999")
		if err != nil {
			t.Errorf("could not listen on 127.0.0.1:9999: %s", err)
		}

		conn, err := ln.Accept()
		if err != nil {
			t.Errorf("could not accept conn: %s", err)
		}
		defer conn.Close()

		buf := make([]byte, 9)
		nr, err := conn.Read(buf)
		if nr != len(buf) {
			t.Errorf("read was wrong length: %d != %d", nr, len(buf))
		}

		good := []byte("foobarbaz")
		if bytes.Compare(buf, good) != 0 {
			t.Errorf("read got wrong data: %s != %s", buf, good)
		} else {
			t.Logf("read got correct data: %s == %s", buf, good)
		}
	}()

	time.Sleep(1 * time.Second)
	addr, _ := net.ResolveTCPAddr("tcp4", "127.0.0.1:9999")
	conn, err := net.DialTCP("tcp4", nil, addr)
	if err != nil {
		t.Errorf("error connecting to 127.0.0.1:9999: %s", err)
	} else {
		t.Logf("connected to server")
	}
	defer conn.Close()
	data := [][]byte{[]byte("foo"), []byte("bar"), []byte("baz")}

	f, err := conn.File()
	if err != nil {
		t.Errorf("could not get file handle for TCP client: %s", err)
	}
	defer f.Close()
	nw, err := vectorio.Writev(f, data)
	f.Seek(0, 0)
	if err != nil {
		t.Errorf("WritevRaw threw error: %s", err)
	}

	if nw != 9 {
		t.Errorf("Length %d of input does not match %d written bytes", len(data), nw)
	} else {
		t.Logf("right number of bytes written")
	}
	time.Sleep(1 * time.Second)
}

func ExampleVectorioCombined() {
	// Create a temp file for demo purposes
	f, err := ioutil.TempFile("", "vectorio")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	data1 := []byte("foobarbaz_buf\n")
	data2 := []byte("barbazfoo_buf\n")

	// This demonstrates the "buffered" form of the library, similar to bufio.
	// w implements io.Writer.
	w, err := vectorio.NewBufferedWritev(f)
	if err != nil {
		panic(err)
	}

	// The simple method is to just write a byte slice;
	// this is converted to a syscall.Iovec and queued
	// for writing
	w.Write(data2)

	// The user can also gain more control with WriteIovec,
	// although this does not have a significant advantage over Write.
	// This is what Write does on your behalf.
	w.WriteIovec(syscall.Iovec{&data1[0], uint64(len(data1))})

	// Flush must be called after writes are complete, to empty out
	// the buffer of pending Iovec.
	// Returns the total number of bytes written, as reported by the underlying syscall.
	nw1, err := w.Flush()
	if err != nil {
		panic(err)
	}

	// One can also write a slice of byte slices ([][]byte) all at once.
	// Note, this usage does *not* implement io.Writer, but if you have a slice
	// of byte slices to write, this is a way to do that without looping.

	multiple := [][]byte{
		[]byte("foobarbaz_slice\n"),
		[]byte("foobazbar_slice\n"),
	}

	// we return the number of bytes written, as reported by the underlying syscall.
	nw2, err := vectorio.Writev(f, multiple)
	if err != nil {
		panic(err)
	}

	fmt.Println("Wrote", nw1+nw2, "bytes to file")
	// Output:
	// Wrote 60 bytes to file
}
