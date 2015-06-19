package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/influxdb/influxdb/client"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	Host      = "influxdb"
	Port      = 8086
	Database  = "test"
	BatchSize = 20
)

type Message struct {
	Key   string
	Value float64
}

func main() {

	u, err := url.Parse(fmt.Sprintf("http://%s:%d", Host, Port))
	if err != nil {
		log.Fatal(err)
	}

	conf := client.Config{
		URL:      *u,
		Username: os.Getenv("INFLUX_USER"),
		Password: os.Getenv("INFLUX_PWD"),
	}
	con, err := client.NewClient(conf)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanLines)
	batch := make([]Message, 0)

	var i int64
	for scanner.Scan() {

		i++
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}

		line := fmt.Sprint(scanner.Text())
		if len(line) == 0 {
			log.Println("String with no length. Skipping.")
			continue
		}

		if i%BatchSize == 0 || len(line) == 0 {
			writePoints(con, batch)
			batch = nil
		}

		line = strings.TrimSpace(line)
		line_split := strings.Split(line, "|")

		if len(line_split) < 2 {
			log.Println("No value. Skipping.")
			log.Println(line_split)
			continue
		}

		key := line_split[0]
		value, err := strconv.ParseFloat(line_split[1], 32)

		if err != nil {
			log.Println("Unable to parse float32")
			continue
		}

		message := Message{
			Key:   key,
			Value: value,
		}

		if err != nil {
			log.Println("Unable to read string from stdin")
		}

		batch = append(batch, message)

		log.Printf("Key: %s Value: %g", message.Key, message.Value)
	}
	writePoints(con, batch)
}

func writePoints(con *client.Client, batch []Message) {
	pts := make([]client.Point, 0)
	region := "us-west-1b"
	node_type := "api"

	for _, m := range batch {
		pt := client.Point{
			Measurement: m.Key,
			Tags: map[string]string{
				"region":    region,
				"node_type": node_type,
			},
			Fields: map[string]interface{}{
				"value": m.Value,
			},
			Time:      time.Now(),
			Precision: "s",
		}
		pts = append(pts, pt)
	}

	bps := client.BatchPoints{
		Points:          pts,
		Database:        Database,
		RetentionPolicy: "default",
	}

	message_debug, _ := json.Marshal(bps)
	log.Printf(string(message_debug))

	_, err := con.Write(bps)
	if err != nil {
		log.Fatal(err)
	}
}
