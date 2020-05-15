package main

import (
	"fmt"
	"strings"
	"time"
)

type Reader interface {
	Read(rc chan string)
}

type Writer interface {
	Write(wc chan string)
}
type LogProcess struct {
	rc chan string
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

func (r *ReadFromFile) Read(rc chan string) {
	line := "message"
	rc <- line
}

func (l *LogProcess) Process()  {
	// Process
	data := <-l.rc
	l.wc <- strings.ToUpper(data)
}

func (w *WriteToInfluxDB) Write(wc chan string)  {
	// write
	fmt.Println(<-wc)
}

func main()  {

	reader := &ReadFromFile{
		path: "/tmp/access.log",
	}

	writer := &WriteToInfluxDB{
		influxDBsn: "username&password..",
	}

	lp := &LogProcess{
		rc : make(chan string),
		wc : make(chan string),
		read: reader,
		write: writer,
	}

	go lp.read.Read(lp.rc)
	go lp.Process()
	go lp.write.Write(lp.wc)

	time.Sleep(time.Second)
}

