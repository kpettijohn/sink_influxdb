package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/influxdb/influxdb/client"
	"github.com/kelseyhightower/envconfig"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

// Configuration struct for envconfig
type Config struct {
	Host      string
	Port      int
	Database  string
	User      string
	Password  string
	Batchsize int64
	RP        string
}

// Message format from stdin.
//
// KeyName:key:tag1:tag_value:tag2:tag_value|vale|timestamp
// KeyName:cpu:node_type:worker:instance_id:i-52fdg34|30|2015-06-24T05:48:56.865215158Z
type Message struct {
	Key   string
	Value float64
	Tags  map[string]string
	Time  time.Time
}

func main() {

	// Get ENV configuration
	var cfg Config
	err := envconfig.Process("influx", &cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Create influxdb connection
	u, err := url.Parse(fmt.Sprintf("http://%s:%d", cfg.Host, cfg.Port))
	if err != nil {
		log.Fatal(err)
	}

	conf := client.Config{
		URL:      *u,
		Username: cfg.User,
		Password: cfg.Password,
	}

	con, err := client.NewClient(conf)
	if err != nil {
		log.Fatal(err)
	}

	// Create a scanner on stdin
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanLines)
	batch := make([]Message, 0)

	var i int64
	for scanner.Scan() {

		i++
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}

		// Grab a line (string) from stdin
		line := scanner.Text()
		// Make sure the line containes a message
		if len(line) == 0 {
			log.Println("String with no length. Skipping.")
			continue
		}

		// If we hit our batch size write the batch of
		// points/messages to influxdb and start the next batch.
		if i%cfg.Batchsize == 0 {
			pts := CreatePoints(batch)
			WritePoints(con, cfg, pts)
			batch = nil
		}

		message := CreateMessage(line)
		batch = append(batch, message)

		// log.Printf("Key: %s Value: %g", message.Key, message.Value)
	}

	if len(batch) == 0 {
		log.Println("Message with no length. Skipping write and exiting..")
		os.Exit(0)
	}
	// Write any remaining points/messages to influxdb
	pts := CreatePoints(batch)
	WritePoints(con, cfg, pts)
}

func CreateMessage(line string) Message {
	// Remove any leading and trailing white space.
	// Including \n
	line = strings.TrimSpace(line)
	line_split := strings.Split(line, "|")

	if len(line_split) < 2 {
		log.Println("No value. Skipping.")
	}

	tags := TagSplit(line_split[0])

	key := tags["KeyName"]
	delete(tags, "KeyName")
	value, err := strconv.ParseFloat(line_split[1], 32)

	if err != nil {
		log.Println("Unable to parse float32")
	}

	const RFC3339 = "2006-01-02T15:04:05Z07:00"
	timestamp, err := time.Parse(RFC3339, line_split[2])
	if err != nil {
		log.Println("Unable to parse time")
	}

	message := Message{
		Key:   key,
		Value: value,
		Tags:  tags,
		Time:  timestamp,
	}

	return message
}

func TagSplit(t string) map[string]string {
	tagSlice := strings.Split(t, ":")
	tags := make(map[string]string)

	for i, _ := range tagSlice {
		if i%2 == 0 {
			v := i + 1
			if v < len(tagSlice) {
				tags[tagSlice[i]] = tagSlice[v]
			} else {
				continue
			}
		}
	}
	return tags
}

func CreatePoints(batch []Message) []client.Point {
	pts := make([]client.Point, 0)
	for _, m := range batch {
		pt := client.Point{
			Measurement: m.Key,
			Tags:        m.Tags,
			Fields: map[string]interface{}{
				"value": m.Value,
			},
			Time:      m.Time,
			Precision: "s",
		}
		pts = append(pts, pt)
	}
	return pts
}

func WritePoints(con *client.Client, cfg Config, pts []client.Point) {
	bps := client.BatchPoints{
		Points:          pts,
		Database:        cfg.Database,
		RetentionPolicy: cfg.RP,
	}

	message_debug, _ := json.Marshal(bps)
	log.Printf(string(message_debug))

	_, err := con.Write(bps)
	if err != nil {
		log.Fatal(err)
	}
}
