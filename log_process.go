package main

type LogProcess struct {
	path string // read file path
	influxDBsn string  // influx data source
}

func (l *LogProcess) ReadFromFile()  {
	// read
}

func (l *LogProcess) Process()  {
	// Process
}

func (l *LogProcess) WriteToInfluxDB()  {
	// write
}

func main()  {
	lp := &LogProcess{
		path: "/tmp/access.log",
		influxDBsn: "username&password..",
	}

	go lp.ReadFromFile()
	go lp.Process()
	go lp.WriteToInfluxDB()

}

