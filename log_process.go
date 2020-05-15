package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

type Reader interface {
	Read(rc chan []byte)
}

type Writer interface {
	Write(wc chan string)
}
type LogProcess struct {
	rc chan []byte
	wc chan string
	read Reader
	write Writer
}

type ReadFromFile struct {
	path string
}

type WriteToInfluxDB struct {
	influxDBsn string
}

func (r *ReadFromFile) Read(rc chan []byte) {
	// open file
	f, err := os.Open(r.path)
	if err !=nil {
		panic(fmt.Sprintf("open file error: %s", err.Error()))
	}

	// read lines from the last line
	f.Seek(0, 2)
	rd := bufio.NewReader(f)

	for {
		line, err := rd.ReadBytes('\n')
		if err == io.EOF {
			time.Sleep(500*time.Microsecond)
			continue
		} else if err != nil {
			panic(fmt.Sprintf("read bytes error: %s", err))
		}
		// remove line breaker
		rc <- line[:len(line)-1]
	}
}

func (l *LogProcess) Process()  {
	// Process
	for data := range l.rc {
		l.wc <- strings.ToUpper(string(data))
	}
}

func (w *WriteToInfluxDB) Write(wc chan string)  {
	// write
	for data := range wc{
		fmt.Println(data)
	}
}

func main()  {

	reader := &ReadFromFile{
		path: "./access.log",
	}

	writer := &WriteToInfluxDB{
		influxDBsn: "username&password..",
	}

	lp := &LogProcess{
		rc : make(chan []byte),
		wc : make(chan string),
		read: reader,
		write: writer,
	}

	go lp.read.Read(lp.rc)
	go lp.Process()
	go lp.write.Write(lp.wc)

	time.Sleep(30 * time.Second)
}

