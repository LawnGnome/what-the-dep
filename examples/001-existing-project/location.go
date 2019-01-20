package main

import (
	"errors"
	"io"
	"log"
	"strconv"
	"strings"
	"time"

	nmea "github.com/adrianmo/go-nmea"
)

type Location struct {
	current  *Position
	fixes    chan nmea.GPGGA
	input    chan string
	monitors []chan Position
	reader   io.ReadCloser
}

type Position struct {
	Time      time.Time
	Latitude  float64
	Longitude float64
	Altitude  float64
}

func NewLocation(r io.ReadCloser) *Location {
	l := &Location{
		fixes:  make(chan nmea.GPGGA),
		input:  make(chan string),
		reader: r,
	}

	// Wire up the goroutines.
	go l.handleFixes()
	go l.parseInput()
	go l.readInput()

	return l
}

func (l *Location) Current() (Position, error) {
	if l.current != nil {
		return *l.current, nil
	}

	return Position{}, errors.New("No fix available")
}

func (l *Location) Monitor() chan Position {
	monitor := make(chan Position)
	l.monitors = append(l.monitors, monitor)

	return monitor
}

func (l *Location) Shutdown() {
	for _, monitor := range l.monitors {
		close(monitor)
	}

	l.reader.Close()
}

func (l *Location) handleFixes() {
	for {
		fix, ok := <-l.fixes
		if !ok {
			return
		}

		l.update(&fix)
	}
}

func (l *Location) parseInput() {
	for {
		line, ok := <-l.input
		if !ok {
			close(l.fixes)
			return
		}

		sentence, err := nmea.Parse(line)
		// Errors usually indicate just that it's an unhandled type, so
		// we'll ignore them.
		// TODO: logging?
		if err != nil {
			continue
		}

		if sentence.GetSentence().Type == nmea.PrefixGPGGA {
			l.fixes <- sentence.(nmea.GPGGA)
		}
	}
}

func (l *Location) readInput() {
	acc := ""
	buf := make([]byte, 128)
	initial := true

	for {
		n, err := l.reader.Read(buf)
		if err != nil {
			log.Println(err)
			close(l.input)
			return
		}

		acc = acc + string(buf[:n])
		if initial {
			pos := strings.Index(acc, "$G")
			if pos != -1 {
				initial = false
				acc = acc[pos:]
			} else {
				continue
			}
		}

		lines := strings.Split(acc, "\n")
		for i, ss := range lines {
			if i == len(lines)-1 {
				acc = ss
			} else {
				l.input <- strings.TrimSpace(ss)
			}
		}
	}
}

func (l *Location) update(fix *nmea.GPGGA) bool {
	if fix.FixQuality == nmea.Invalid {
		return false
	}

	now := time.Now().UTC()
	t, err := time.Parse("150405", fix.Time)
	if err == nil {
		now = time.Date(now.Year(), now.Month(), now.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), time.UTC)
	}

	altitude, err := strconv.ParseFloat(fix.Altitude, 64)
	if err != nil {
		altitude = 0
	}

	pos := &Position{
		Time:      now,
		Latitude:  float64(fix.Latitude),
		Longitude: float64(fix.Longitude),
		Altitude:  altitude,
	}

	l.current = pos
	for _, monitor := range l.monitors {
		monitor <- *pos
	}

	return true
}
