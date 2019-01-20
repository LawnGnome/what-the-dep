package main

import (
	"flag"
	"log"

	"github.com/aecerbot/protocols/src/protocols"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	serial "gopkg.in/bugst/go-serial.v1"
	goredis "gopkg.in/redis.v4"
)

var (
	device string
	key    string
	redis  string
)

func init() {
	flag.StringVar(&device, "device", "/dev/gps", "GPS serial device")
	flag.StringVar(&key, "key", "position", "Redis key")
	flag.StringVar(&redis, "redis", "redis:6379", "Redis host and port")
}

func main() {
	flag.Parse()

	cfg := &serial.Mode{
		BaudRate: 4800,
	}
	s, err := serial.Open(device, cfg)
	if err != nil {
		log.Fatal(err)
	}

	client := goredis.NewClient(&goredis.Options{
		Addr: redis,
	})

	location := NewLocation(s)
	c := location.Monitor()
	for {
		pos, ok := <-c
		if !ok {
			log.Println("Channel closed; exiting")
		}

		log.Println(pos)

		time, err := ptypes.TimestampProto(pos.Time)
		if err != nil {
			log.Println(err)
			continue
		}

		msg, err := proto.Marshal(&protocols.Position{
			Time:      time,
			Latitude:  pos.Latitude,
			Longitude: pos.Longitude,
			Altitude:  pos.Altitude,
		})
		if err != nil {
			log.Println(err)
			continue
		}

		if err := client.Set(key, msg, 0); err != nil {
			log.Println(err)
			continue
		}
	}
}
