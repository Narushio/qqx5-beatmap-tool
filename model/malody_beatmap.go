package model

type MalodyBeatmap struct {
	Meta  MetaData    `json:"meta"`
	Time  []TimePoint `json:"time"`
	Note  []Note      `json:"note"`
	Extra Extra       `json:"extra"`
}

type MetaData struct {
	Creator    string  `json:"creator"`
	Background string  `json:"background"`
	Version    string  `json:"version"`
	ID         int     `json:"id"`
	Mode       int     `json:"mode"`
	Time       int     `json:"time"`
	Song       Song    `json:"song"`
	ModeExt    ModeExt `json:"mode_ext"`
}

type Song struct {
	Title  string `json:"title"`
	Artist string `json:"artist"`
	ID     int    `json:"id"`
}

type ModeExt struct {
	Column int `json:"column"`
}

type TimePoint struct {
	Beat [3]float64 `json:"beat"`
	BPM  float64    `json:"bpm"`
}

type Note struct {
	Beat    [3]float64  `json:"beat"`
	EndBeat *[3]float64 `json:"endbeat"`
	Column  int         `json:"column"`
	Sound   *string     `json:"sound"`
	Vol     *int        `json:"vol"`
	Offset  *float64    `json:"offset"`
	Type    *int        `json:"type"`
}

type Extra struct {
	Test Test `json:"test"`
}

type Test struct {
	Divide   int `json:"divide"`
	Speed    int `json:"speed"`
	Save     int `json:"save"`
	Lock     int `json:"lock"`
	EditMode int `json:"edit_mode"`
}

func (m *MalodyBeatmap) BpmOffset() *float64 {
	soundNote := Note{}
	for _, n := range m.Note {
		if n.Type != nil && *n.Type != 0 {
			soundNote = n
			break
		}
	}

	return soundNote.Offset
}

func (m *MalodyBeatmap) CalcNoteHitTime(b [3]float64, bpmOffset *float64) float64 {
	// TODO: fix sv songs

	return ms(beat(b)-beat(m.Time[0].Beat), m.Time[0].BPM, bpmOffset)
}

func ms(beats float64, bpm float64, offset *float64) float64 {
	if offset == nil {
		return 1000 * (60 / bpm) * beats
	}

	return 1000*(60/bpm)*beats + *offset
}

func beat(beat [3]float64) float64 {
	return beat[0] + beat[1]/beat[2]
}

func (m *MalodyBeatmap) IsKeyMode() bool {
	return m.Meta.Mode == 0
}

func (m *Note) IsHoldNote() bool {
	return m.EndBeat != nil
}

func (m *Note) QQX5BeatmapTargetTrack(columnCount int) QQX5BeatmapTargetTrackType {
	if columnCount == 4 {
		switch m.Column {
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
		switch m.Column {
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
