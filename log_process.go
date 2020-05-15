package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	_ "github.com/influxdata/influxdb1-client" // this is important because of the bug in go mod
	client "github.com/influxdata/influxdb1-client/v2"
)

type Reader interface {
	Read(rc chan []byte)
}

type Writer interface {
	Write(wc chan *Message)
}
type LogProcess struct {
	rc chan []byte
	wc chan *Message
	read Reader
	write Writer
}

type ReadFromFile struct {
	path string
}

type WriteToInfluxDB struct {
	influxDBsn string
}

type Message struct {
	TimeLocal 						time.Time
	BytesSent 						int
	Path, Method, Scheme, Status 	string
	UpstreamTime, RequestTime 		float64
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
	// 100.97.120.0 - - [08/Jan/2016:10:40:18 +0800] http "GET / HTTP/1.0" 200 612 "-" "KeepAliveClient" "-" 1.005 1.854

	r := regexp.MustCompile(`([\d\.]+)\s+([^ \[]+)\s+([^ \[]+)\s+\[([^\]]+)\]\s+([a-z]+)\s+\"([^"]+)\"\s(\d{3})\s+
								(\d+)\s+\"([^"]+)\"\s+\"(.*?)\"([\d\.-]+)\"\s+([\d\.-]+)\s+([\d\.-]+)`)

	for data := range l.rc {
		sub := r.FindStringSubmatch(string(data))
		if len(sub) != 14 {
			log.Println("find submatch string fail:", string(data))
			continue
		}

		message := &Message{}
		loc, _ := time.LoadLocation("Asia/Shanghai")
		t, err := time.ParseInLocation("02/Jan/2006:15:04:05 +0000", sub[4], loc )
		if err != nil {
			log.Println("parse in location fail:", err.Error(), sub[4])
		}

		bytesSent, _ := strconv.Atoi(sub[8])

		reqSli := strings.Split(sub[6], " ")
		if len(reqSli) != 3 {
			log.Println("strings split fail", sub[6])
			continue
		}

		u, err := url.Parse(reqSli[1])
		if err != nil {
			log.Println("url parse fail", err)
		}

		upstreamTime, _ := strconv.ParseFloat(sub[12], 64)
		requestTime, _ := strconv.ParseFloat(sub[13], 64)

		message.TimeLocal = t
		message.Method = reqSli[0]
		message.BytesSent = bytesSent
		message.Path = u.Path
		message.Scheme = sub[5]
		message.Status = sub[7]
		message.UpstreamTime = upstreamTime
		message.RequestTime = requestTime

		l.wc <- message
	}
}

func (w *WriteToInfluxDB) Write(wc chan *Message)  {
	// Create a new HTTPCLient
	infSli := strings.Split(w.influxDBsn, "@")
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: infSli[0],
		Username: infSli[1],
		Password: infSli[2],
	})
	if err != nil {
		fmt.Println("Error creating InfluxDB Client: ", err.Error())
	}
	defer c.Close()

	for v := range wc {
		// Create a new point batch
		batchPoints, err := client.NewBatchPoints(client.BatchPointsConfig{
			Database: infSli[3],
			Precision: infSli[4],
		})
		if err != nil {
			log.Fatal(err)
		}

		// Create a point and add to batch
		tags := map[string]string{
			"Path": v.Path,
			"Method": v.Method,
			"Scheme": v.Scheme,
			"Status": v.Status,
		}

		fields := map[string]interface{}{
			"UpstreamTime": v.UpstreamTime,
			"RequestTime": v.RequestTime,
			"BytesSent" : v.BytesSent,
		}

		point, err := client.NewPoint("nginx_log", tags, fields, v.TimeLocal)
		if err != nil {
			log.Fatal(err)
		}
		batchPoints.AddPoint(point)

		// Write the batch
		if err := c.Write(batchPoints); err != nil {
			log.Fatal(err)
		}

		log.Println("write to db success!")
	}

}

func main()  {
	var path, influxDBsn string

	flag.StringVar(&path, "path",
		"./access.log",
		"Read file path" )

	flag.StringVar(&influxDBsn,
		"influxDBsn",
		"http://127.0.0.1:8086@yng@yngpass@yng@s",
		"influx data source" )
	flag.Parse()

	reader := &ReadFromFile{
		path: path,
	}

	writer := &WriteToInfluxDB{
		influxDBsn: influxDBsn,
	}

	lp := &LogProcess{
		rc : make(chan []byte),
		wc : make(chan *Message),
		read: reader,
		write: writer,
	}

	go lp.read.Read(lp.rc)
	go lp.Process()
	go lp.write.Write(lp.wc)

	time.Sleep(30 * time.Second)
}

