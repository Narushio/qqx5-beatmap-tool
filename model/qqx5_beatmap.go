package model

import (
	"encoding/xml"
	"fmt"
	"math"
	"strings"
)

type (
	QQX5BeatmapNoteType        string
	QQX5BeatmapSectionType     string
	QQX5BeatmapTargetTrackType string
)

const (
	QQX5BeatmapShortNote QQX5BeatmapNoteType = "short"
	QQX5BeatmapLongNote  QQX5BeatmapNoteType = "long"
	QQX5BeatmapSlipNote  QQX5BeatmapNoteType = "slip"

	QQX5BeatmapPreviousSection QQX5BeatmapSectionType = "previous"
	QQX5BeatmapBeginSection    QQX5BeatmapSectionType = "begin"
	QQX5BeatmapNoteSection     QQX5BeatmapSectionType = "note"
	QQX5BeatmapShowTimeSection QQX5BeatmapSectionType = "showtime"

	QQX5BeatmapTrackLeft2  QQX5BeatmapTargetTrackType = "Left2"
	QQX5BeatmapTrackLeft1  QQX5BeatmapTargetTrackType = "Left1"
	QQX5BeatmapTrackMiddle QQX5BeatmapTargetTrackType = "Middle"
	QQX5BeatmapTrackRight1 QQX5BeatmapTargetTrackType = "Right1"
	QQX5BeatmapTrackRight2 QQX5BeatmapTargetTrackType = "Right2"

	QQX5BeatmapIndent         string  = "&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;"
	QQX5BeatmapToleranceRange float64 = 100.0
)

type QQX5BeatmapLevel struct {
	XMLName           xml.Name                     `xml:"Level"`
	LevelInfo         QQX5BeatmapLevelInfo         `xml:"LevelInfo"`
	MusicInfo         QQX5BeatmapMusicInfo         `xml:"MusicInfo"`
	SectionSeq        QQX5BeatmapSectionSeq        `xml:"SectionSeq"`
	IndicatorResetPos QQX5BeatmapIndicatorResetPos `xml:"IndicatorResetPos"`
	NoteInfo          QQX5BeatmapNoteInfo          `xml:"NoteInfo"`
	ActionSeq         QQX5BeatmapActionSeq         `xml:"ActionSeq"`
	CameraSeq         QQX5BeatmapCameraSeq         `xml:"CameraSeq"`
	DancerSort        QQX5BeatmapDancerSort        `xml:"DancerSort"`
	StageEffectSeq    QQX5BeatmapStageEffectSeq    `xml:"StageEffectSeq"`
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

			if resultMS-ms > QQX5BeatmapToleranceRange {
				bar -= 1
			} else if resultMS-ms < -QQX5BeatmapToleranceRange {
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

func (l *QQX5BeatmapLevel) ResetNotesWithBPM(bpm float64) *QQX5BeatmapLevel {
	oldBarMS, oldPosMS := calculateBarAndPosMS(l.LevelInfo.BPM)
	newBarMS, newPosMS := calculateBarAndPosMS(bpm)

	resetNotes := func(notes []*QQX5BeatmapNote) {
		for _, n := range notes {
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

type QQX5BeatmapLevelInfo struct {
	BPM             float64 `xml:"BPM"`
	BeatPerBar      int     `xml:"BeatPerBar"`
	BeatLen         int     `xml:"BeatLen"`
	EnterTimeAdjust float64 `xml:"EnterTimeAdjust"`
	NotePreShow     float64 `xml:"NotePreShow"`
	LevelTime       int     `xml:"LevelTime"`
	BarAmount       int     `xml:"BarAmount"`
	BeginBarLen     int     `xml:"BeginBarLen"`
	IsFourTrack     bool    `xml:"IsFourTrack"`
	TrackCount      int     `xml:"TrackCount"`
	LevelPreTime    int     `xml:"LevelPreTime"`
	Star            int     `xml:"Star"`
}

type QQX5BeatmapMusicInfo struct {
	Author   string `xml:"Author"`
	Title    string `xml:"Title"`
	Artist   string `xml:"Artist"`
	FilePath string `xml:"FilePath"`
}

type QQX5BeatmapSectionSeq struct {
	Sections []*QQX5BeatmapSection `xml:"Section"`
}

type QQX5BeatmapSection struct {
	Type     string `xml:"type,attr"`
	StartBar int    `xml:"startbar,attr"`
	EndBar   int    `xml:"endbar,attr"`
	Mark     string `xml:"mark,attr"`
	Param1   string `xml:"param1,attr"`
}

type QQX5BeatmapIndicatorResetPos struct {
	PosNum int `xml:"PosNum,attr"`
}

type QQX5BeatmapNoteInfo struct {
	Normal QQX5BeatmapNormal `xml:"Normal"`
}

type QQX5BeatmapNormal struct {
	Notes        []*QQX5BeatmapNote        `xml:"Note"`
	CombineNotes []*QQX5BeatmapCombineNote `xml:"CombineNote"`
}

func (n *QQX5BeatmapNormal) ParseFromOsuBeatMap(osuBeatMap *OSUBeatMap, bpm float64) {
	var (
		barTime = 60000 / bpm * 4
		posTime = barTime / 32
	)

	for _, ho := range osuBeatMap.HitObjects {
		bar, pos := calculateBarAndPos(float64(ho.Time), barTime, posTime)

		note := &QQX5BeatmapNote{
			Bar:         bar,
			Pos:         pos,
			TargetTrack: ho.QQX5BeatmapTargetTrack(osuBeatMap.Difficulty.CircleSize),
			NoteType:    QQX5BeatmapShortNote,
		}

		n.Notes = append(n.Notes, note)
	}
}

func (n *QQX5BeatmapNormal) ToHTML() string {
	var noteStr string
	for _, n := range n.Notes {
		noteStr += n.ToHTML()
	}
	for _, cn := range n.CombineNotes {
		noteStr += cn.ToHTML()
	}
	return fmt.Sprintf("&lt;Normal&gt;<br>" + noteStr + "&lt;/Normal&gt;<br>")
}

type QQX5BeatmapNote struct {
	Bar         int                        `xml:"Bar,attr"`
	Pos         int                        `xml:"Pos,attr"`
	FromTrack   *string                    `xml:"from_track,attr,omitempty"`
	TargetTrack QQX5BeatmapTargetTrackType `xml:"target_track,attr"`
	EndTrack    *string                    `xml:"end_track,attr,omitempty"`
	NoteType    QQX5BeatmapNoteType        `xml:"note_type,attr"`
	EndBar      *int                       `xml:"EndBar,attr,omitempty"`
	EndPos      *int                       `xml:"EndPos,attr,omitempty"`
}

func (n *QQX5BeatmapNote) ToHTML() string {
	var builder strings.Builder

	builder.WriteString(QQX5BeatmapIndent + "&lt;Note ")
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

func (n *QQX5BeatmapNote) ToMilliseconds(barMS, posMS float64) (float64, float64) {
	hitMS := float64(n.Bar-1)*barMS + float64(n.Pos)*posMS

	if n.IsLongNote() {
		return hitMS, float64(*n.EndBar-1)*barMS + float64(*n.EndPos)*posMS
	}

	return hitMS, 0
}

func (n *QQX5BeatmapNote) IsLongNote() bool {
	return n.NoteType == QQX5BeatmapLongNote
}

type QQX5BeatmapCombineNote struct {
	Notes []*QQX5BeatmapNote `xml:"Note"`
}

func (c *QQX5BeatmapCombineNote) ToHTML() string {
	var noteStr string
	for _, n := range c.Notes {
		noteStr += QQX5BeatmapIndent + n.ToHTML()
	}

	return fmt.Sprintf(QQX5BeatmapIndent + "&lt;CombineNote&gt;<br>" + noteStr + QQX5BeatmapIndent + "&lt;/CombineNote&gt;<br>")
}

type QQX5BeatmapActionSeq struct {
	Type        string                   `xml:"type,attr"`
	ActionLists []*QQX5BeatmapActionList `xml:"ActionList"`
}

type QQX5BeatmapActionList struct {
	StartBar int    `xml:"start_bar,attr"`
	DanceLen int    `xml:"dance_len,attr"`
	SeqLen   int    `xml:"seq_len,attr"`
	Level    int    `xml:"level,attr"`
	Type     string `xml:"type,attr"`
}

type QQX5BeatmapCameraSeq struct {
	Cameras []*QQX5BeatmapCamera `xml:"Camera"`
}

type QQX5BeatmapCamera struct {
	Name   string `xml:"name,attr"`
	Bar    int    `xml:"bar,attr"`
	Pos    int    `xml:"pos,attr"`
	EndBar int    `xml:"end_bar,attr"`
	EndPos int    `xml:"end_pos,attr"`
}

type QQX5BeatmapDancerSort struct {
	Bars []int `xml:"Bar"`
}

type QQX5BeatmapStageEffectSeq struct {
	Effects []*QQX5BeatmapEffect `xml:"effect"`
}

type QQX5BeatmapEffect struct {
	Name   string `xml:"name,attr"`
	Bar    int    `xml:"bar,attr"`
	Length int    `xml:"length,attr"`
}
