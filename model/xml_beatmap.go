package model

import (
	"encoding/xml"
	"fmt"
	"math"
	"strings"
)

type NoteType string

type SectionType string

const (
	ShortNote       NoteType    = "short"
	LongNote        NoteType    = "long"
	SlipNote        NoteType    = "slip"
	PreviousSection SectionType = "previous"
	BeginSection    SectionType = "begin"
	NoteSection     SectionType = "note"
	ShowTimeSection SectionType = "showtime"
	Indent          string      = "&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;"
	ToleranceRange  float64     = 100.0
)

type Level struct {
	XMLName           xml.Name          `xml:"Level"`
	LevelInfo         LevelInfo         `xml:"LevelInfo"`
	MusicInfo         MusicInfo         `xml:"MusicInfo"`
	SectionSeq        SectionSeq        `xml:"SectionSeq"`
	IndicatorResetPos IndicatorResetPos `xml:"IndicatorResetPos"`
	NoteInfo          NoteInfo          `xml:"NoteInfo"`
	ActionSeq         ActionSeq         `xml:"ActionSeq"`
	CameraSeq         CameraSeq         `xml:"CameraSeq"`
	DancerSort        DancerSort        `xml:"DancerSort"`
	StageEffectSeq    StageEffectSeq    `xml:"StageEffectSeq"`
}

func calculateBarAndPosMS(bpm float64) (float64, float64) {
	barMS := 60000 / bpm * 4
	return barMS, barMS / 64
}

func calculateBarAndPos(ms, barMS, posMS float64) (int, int) {
	bar := math.Floor(ms / barMS)
	pos := math.Round(math.Mod(ms, barMS) / posMS)
	if math.Mod(pos, 2) != 0 {
		if math.Mod(ms/barMS, posMS) > math.Round(math.Mod(ms/barMS, posMS)) {
			pos += 1
		} else {
			pos -= 1
		}
	}

	// Fix error bar and pos
	if pos == 64 {
		pos = 0
	}

	{
		for {
			resultMS := bar*barMS + pos*posMS

			if resultMS-ms > ToleranceRange {
				bar -= 1
			} else if resultMS-ms < -ToleranceRange {
				bar += 1
			} else {
				break
			}
		}
	}

	// QQX5XML Beatmap Editor Default append 1 bar
	bar += 1

	return int(bar), int(pos)
}

func (l *Level) ResetNotesWithBPM(bpm float64) *Level {
	oldBarMS, oldPosMS := calculateBarAndPosMS(l.LevelInfo.BPM)
	newBarMS, newPosMS := calculateBarAndPosMS(bpm)

	resetNotes := func(notes []*Note) {
		for _, n := range notes {
			if n.Bar == 125 && n.Pos == 0 && n.NoteType == SlipNote {
				fmt.Println("debug")
			}
			if n.Bar == 109 && n.Pos == 32 && n.IsLongNote() {
				fmt.Println("debug")
			}
			if n.Bar == 28 && n.Pos == 16 && n.IsLongNote() {
				fmt.Println("debug")
			}
			hitMS, releaseMS := n.ToMilliseconds(oldBarMS, oldPosMS)
			n.Bar, n.Pos = calculateBarAndPos(hitMS, newBarMS, newPosMS)
			if n.IsLongNote() {
				*n.EndBar, *n.EndPos = calculateBarAndPos(releaseMS, newBarMS, newPosMS)
			}
		}
	}

	resetNotes(l.NoteInfo.Normal.Notes)

	for _, cn := range l.NoteInfo.Normal.CombineNotes {
		resetNotes(cn.Notes)
	}

	return l
}

type LevelInfo struct {
	BPM             float64 `xml:"BPM"`
	BeatPerBar      int     `xml:"BeatPerBar"`
	BeatLen         int     `xml:"BeatLen"`
	EnterTimeAdjust float64 `xml:"EnterTimeAdjust"`
	NotePreShow     int     `xml:"NotePreShow"`
	LevelTime       int     `xml:"LevelTime"`
	BarAmount       int     `xml:"BarAmount"`
	BeginBarLen     int     `xml:"BeginBarLen"`
	IsFourTrack     bool    `xml:"IsFourTrack"`
	TrackCount      int     `xml:"TrackCount"`
	LevelPreTime    int     `xml:"LevelPreTime"`
	Star            int     `xml:"Star"`
}

type MusicInfo struct {
	Author   string `xml:"Author"`
	Title    string `xml:"Title"`
	Artist   string `xml:"Artist"`
	FilePath string `xml:"FilePath"`
}

type SectionSeq struct {
	Sections []*Section `xml:"Section"`
}

type Section struct {
	Type     string `xml:"type,attr"`
	StartBar int    `xml:"startbar,attr"`
	EndBar   int    `xml:"endbar,attr"`
	Mark     string `xml:"mark,attr"`
	Param1   string `xml:"param1,attr"`
}

type IndicatorResetPos struct {
	PosNum int `xml:"PosNum,attr"`
}

type NoteInfo struct {
	Normal Normal `xml:"Normal"`
}

type Normal struct {
	Notes        []*Note        `xml:"Note"`
	CombineNotes []*CombineNote `xml:"CombineNote"`
}

func (n *Normal) ToHTML() string {
	var noteStr string
	for _, n := range n.Notes {
		noteStr += n.ToHTML()
	}
	for _, cn := range n.CombineNotes {
		noteStr += cn.ToHTML()
	}
	return fmt.Sprintf("&lt;Normal&gt;<br>" + noteStr + "&lt;/Normal&gt;<br>")
}

type Note struct {
	Bar         int      `xml:"Bar,attr"`
	Pos         int      `xml:"Pos,attr"`
	FromTrack   *string  `xml:"from_track,attr,omitempty"`
	TargetTrack string   `xml:"target_track,attr"`
	EndTrack    *string  `xml:"end_track,attr,omitempty"`
	NoteType    NoteType `xml:"note_type,attr"`
	EndBar      *int     `xml:"EndBar,attr,omitempty"`
	EndPos      *int     `xml:"EndPos,attr,omitempty"`
}

func (n *Note) ToHTML() string {
	var builder strings.Builder

	builder.WriteString(Indent + "&lt;Note ")
	builder.WriteString(fmt.Sprintf("Bar=\"%d\" Pos=\"%d\" ", n.Bar, n.Pos))

	if n.FromTrack != nil {
		builder.WriteString(fmt.Sprintf("from_track=\"%s\" ", *n.FromTrack))
	}

	builder.WriteString(fmt.Sprintf("target_track=\"%s\" ", n.TargetTrack))

	if n.FromTrack == nil {
		builder.WriteString(fmt.Sprintf("end_track=\"%s\" ", *n.EndTrack))
	}

	builder.WriteString(fmt.Sprintf("note_type=\"%s\" ", n.NoteType))

	if n.IsLongNote() {
		builder.WriteString(fmt.Sprintf("EndBar=\"%d\" EndPos=\"%d\" ", *n.EndBar, *n.EndPos))
	}

	builder.WriteString("/&gt;<br>")

	return builder.String()
}

func (n *Note) ToMilliseconds(barMS, posMS float64) (float64, float64) {
	hitMS := float64(n.Bar-1)*barMS + float64(n.Pos)*posMS

	if n.IsLongNote() {
		return hitMS, float64(*n.EndBar-1)*barMS + float64(*n.EndPos)*posMS
	}

	return hitMS, 0
}

func (n *Note) IsLongNote() bool {
	return n.NoteType == LongNote
}

type CombineNote struct {
	Notes []*Note `xml:"Note"`
}

func (c *CombineNote) ToHTML() string {
	var noteStr string
	for _, n := range c.Notes {
		noteStr += Indent + n.ToHTML()
	}

	return fmt.Sprintf(Indent + "&lt;CombineNote&gt;<br>" + noteStr + Indent + "&lt;/CombineNote&gt;<br>")
}

type ActionSeq struct {
	Type        string        `xml:"type,attr"`
	ActionLists []*ActionList `xml:"ActionList"`
}

type ActionList struct {
	StartBar int    `xml:"start_bar,attr"`
	DanceLen int    `xml:"dance_len,attr"`
	SeqLen   int    `xml:"seq_len,attr"`
	Level    int    `xml:"level,attr"`
	Type     string `xml:"type,attr"`
}

type CameraSeq struct {
	Cameras []*Camera `xml:"Camera"`
}

type Camera struct {
	Name   string `xml:"name,attr"`
	Bar    int    `xml:"bar,attr"`
	Pos    int    `xml:"pos,attr"`
	EndBar int    `xml:"end_bar,attr"`
	EndPos int    `xml:"end_pos,attr"`
}

type DancerSort struct {
	Bars []int `xml:"Bar"`
}

type StageEffectSeq struct {
	Effects []*Effect `xml:"effect"`
}

type Effect struct {
	Name   string `xml:"name,attr"`
	Bar    int    `xml:"bar,attr"`
	Length int    `xml:"length,attr"`
}
