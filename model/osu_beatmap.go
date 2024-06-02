package model

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"strconv"
	"strings"
)

// official website https://osu.ppy.sh/wiki/zh/Client/File_formats/osu_(file_format)

type OsuBeatMap struct {
	Metadata     OsuBeatMapMetaData
	TimingPoints []OsuBeatMapTimingPoint
	HitObjects   []OsuBeatMapHitObject
	Difficulty   OsuBeatMapDifficulty
}

type OsuBeatMapDifficulty struct {
	HPDrainRate       float64
	CircleSize        float64
	OverallDifficulty float64
	ApproachRate      float64
	SliderMultiplier  float64
	SliderTickRate    float64
}

type OsuBeatMapMetaData struct {
	Title         string
	TitleUnicode  string
	Artist        string
	ArtistUnicode string
	Creator       string
	Version       string
	Source        string
	Tags          []string
	BeatmapID     int
	BeatmapSetID  int
}

type OsuBeatMapTimingPoint struct {
	Time        int
	BeatLength  float64 // BPM = （1 / BeatLength * 1000 * 60）
	Meter       int
	SampleSet   int
	SampleIndex int
	Volume      int
	Uninherited bool
	Effects     int
}

type OsuBeatMapHitObject struct {
	X         int // floor(x * columnCount / 512)
	Y         int // default 192 for osu!mania
	Time      int
	Type      int
	HitSound  int
	HitSample string // default 0:0:0:0:
}

func (m *OsuBeatMapHitObject) EndTime() int {
	parts := strings.Split(m.HitSample, ":")
	endTime, _ := strconv.Atoi(parts[0])
	return endTime
}

func (m *OsuBeatMapHitObject) IsHoldNote() bool {
	return m.Type == 128
}

func (m *OsuBeatMapHitObject) QQX5BeatmapTargetTrack(columnCount float64) QQX5BeatmapTargetTrackType {
	targetTrack := math.Floor(float64(m.X) * columnCount / 512)

	if columnCount == 4 {
		switch targetTrack {
		case 0:
			return QQX5BeatmapTrackLeft2
		case 1:
			return QQX5BeatmapTrackLeft1
		case 2:
			return QQX5BeatmapTrackRight1
		case 3:
			return QQX5BeatmapTrackRight2
		}
	} else {
		switch targetTrack {
		case 0:
			return QQX5BeatmapTrackLeft2
		case 1:
			return QQX5BeatmapTrackLeft1
		case 2:
			return QQX5BeatmapTrackMiddle
		case 3:
			return QQX5BeatmapTrackRight1
		case 4:
			return QQX5BeatmapTrackRight2
		}
	}

	return ""
}

func (m *OsuBeatMap) Parse(b []byte) error {
	buf := bytes.NewBuffer(b)
	scanner := bufio.NewScanner(buf)

	var currentSection string

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			currentSection = line
			continue
		}

		switch currentSection {
		case "[Metadata]":
			m.parseMetadata(line)
		case "[TimingPoints]":
			m.parseTimingPoints(line)
		case "[Difficulty]":
			m.parseDifficulty(line)
		case "[HitObjects]":
			m.parseHitObjects(line)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("reading buffer: %w", err)
	}

	return nil
}

func (m *OsuBeatMap) parseDifficulty(line string) {
	if line == "" {
		return
	}

	parts := strings.Split(line, ":")
	if len(parts) != 2 {
		return
	}

	key := strings.TrimSpace(parts[0])
	value, _ := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)

	switch key {
	case "HPDrainRate":
		m.Difficulty.HPDrainRate = value
	case "CircleSize":
		m.Difficulty.CircleSize = value
	case "OverallDifficulty":
		m.Difficulty.OverallDifficulty = value
	case "ApproachRate":
		m.Difficulty.ApproachRate = value
	case "SliderMultiplier":
		m.Difficulty.SliderMultiplier = value
	case "SliderTickRate":
		m.Difficulty.SliderTickRate = value
	}
}

func (m *OsuBeatMap) parseMetadata(line string) {
	if line == "" {
		return
	}

	parts := strings.Split(line, ":")
	if len(parts) != 2 {
		return
	}

	key := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])

	switch key {
	case "Title":
		m.Metadata.Title = value
	case "TitleUnicode":
		m.Metadata.TitleUnicode = value
	case "Artist":
		m.Metadata.Artist = value
	case "ArtistUnicode":
		m.Metadata.ArtistUnicode = value
	case "Creator":
		m.Metadata.Creator = value
	case "Version":
		m.Metadata.Version = value
	case "Source":
		m.Metadata.Source = value
	case "Tags":
		m.Metadata.Tags = strings.Split(value, " ")
	case "BeatmapID":
		m.Metadata.BeatmapID, _ = strconv.Atoi(value)
	case "BeatmapSetID":
		m.Metadata.BeatmapSetID, _ = strconv.Atoi(value)
	}
}

func (m *OsuBeatMap) parseTimingPoints(line string) {
	if line == "" {
		return
	}

	parts := strings.Split(line, ",")
	if len(parts) < 8 {
		return
	}

	time, _ := strconv.Atoi(parts[0])
	beatLength, _ := strconv.ParseFloat(parts[1], 64)
	meter, _ := strconv.Atoi(parts[2])
	sampleSet, _ := strconv.Atoi(parts[3])
	sampleIndex, _ := strconv.Atoi(parts[4])
	volume, _ := strconv.Atoi(parts[5])
	uninherited := parts[6] == "1"
	effects, _ := strconv.Atoi(parts[7])

	timingPoint := OsuBeatMapTimingPoint{
		Time:        time,
		BeatLength:  beatLength,
		Meter:       meter,
		SampleSet:   sampleSet,
		SampleIndex: sampleIndex,
		Volume:      volume,
		Uninherited: uninherited,
		Effects:     effects,
	}

	m.TimingPoints = append(m.TimingPoints, timingPoint)
}

func (m *OsuBeatMap) parseHitObjects(line string) {
	if line == "" {
		return
	}

	parts := strings.Split(line, ",")
	if len(parts) < 6 {
		return
	}

	x, _ := strconv.Atoi(parts[0])
	y, _ := strconv.Atoi(parts[1])
	time, _ := strconv.Atoi(parts[2])
	typeVal, _ := strconv.Atoi(parts[3])
	hitSound, _ := strconv.Atoi(parts[4])
	hitSample := parts[5]

	hitObject := OsuBeatMapHitObject{
		X:         x,
		Y:         y,
		Time:      time,
		Type:      typeVal,
		HitSound:  hitSound,
		HitSample: hitSample,
	}

	m.HitObjects = append(m.HitObjects, hitObject)
}
