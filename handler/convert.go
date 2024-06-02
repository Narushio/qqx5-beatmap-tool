package handler

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/Narushio/qqx5-beatmap-tool/handler/httpdto"
	"github.com/Narushio/qqx5-beatmap-tool/handler/render"
	"github.com/Narushio/qqx5-beatmap-tool/model"
	"github.com/Narushio/qqx5-beatmap-tool/utils/extractor"
	"github.com/gin-gonic/gin"
)

type Convert struct{}

func NewConvert() *Convert {
	return &Convert{}
}

func (h *Convert) UpdateXMLBPM(ctx *gin.Context) {
	param := new(httpdto.UpdateXMLBPMParam)
	if err := ctx.Bind(param); err != nil {
		render.Error(ctx, http.StatusBadRequest, fmt.Errorf("binding param: %w", err))
	}

	f, err := param.File.Open()
	if err != nil {
		render.Error(ctx, http.StatusInternalServerError, fmt.Errorf("opening xml file: %w", err))
	}
	defer f.Close()

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, f); err != nil {
		render.Error(ctx, http.StatusInternalServerError, fmt.Errorf("copying xml file: %w", err))
	}

	beatMap := new(model.QQX5BeatmapLevel)
	if err := xml.Unmarshal(buf.Bytes(), beatMap); err != nil {
		render.Error(ctx, http.StatusBadRequest, fmt.Errorf("unmarshalling xml: %w", err))
	}

	render.OK(ctx, httpdto.UpdateXMLBPMResponse{
		Notes: beatMap.ResetNotesWithBPM(param.BPM).NoteInfo.Normal.ToHTML(),
	})
}

func (h *Convert) OSUToXML(ctx *gin.Context) {
	param := new(httpdto.OsuToXMLParam)
	if err := ctx.Bind(param); err != nil {
		render.Error(ctx, http.StatusBadRequest, fmt.Errorf("binding param: %w", err))
	}

	f, err := param.File.Open()
	if err != nil {
		render.Error(ctx, http.StatusInternalServerError, fmt.Errorf("opening osu file: %w", err))
	}
	defer f.Close()

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, f); err != nil {
		render.Error(ctx, http.StatusInternalServerError, fmt.Errorf("copying osu file: %w", err))
	}

	beatMap := new(model.OsuBeatMap)
	if err := beatMap.Parse(buf.Bytes()); err != nil {
		render.Error(ctx, http.StatusInternalServerError, fmt.Errorf("parsing osu file: %w", err))
	}

	notes := new(model.QQX5BeatmapNormal)
	notes.ParseFromOsuBeatMap(beatMap, param.BPM)
	offset := beatMap.HitObjects[0].Time

	render.OK(ctx, httpdto.OsuToXMLResponse{
		Offset:  float64(offset),
		Creator: beatMap.Metadata.Creator,
		Notes:   notes.ToHTML(),
	})
}

func (h *Convert) MCZToXML(ctx *gin.Context) {
	param := new(httpdto.MCZToXMLParam)
	if err := ctx.Bind(param); err != nil {
		render.Error(ctx, http.StatusBadRequest, fmt.Errorf("binding param: %w", err))
	}

	f, err := param.File.Open()
	if err != nil {
		render.Error(ctx, http.StatusInternalServerError, fmt.Errorf("opening mcz file: %w", err))
	}
	defer f.Close()

	mcBytes, err := extractor.MCZ(f, ".mc")
	if err != nil {
		render.Error(ctx, http.StatusInternalServerError, fmt.Errorf("extracting mcz file: %w", err))
	}

	var beatMap *model.MalodyBeatmap
	if err := json.Unmarshal(mcBytes, &beatMap); err != nil {
		render.Error(ctx, http.StatusInternalServerError, fmt.Errorf("unmarshaling json: %w", err))
	}

	if !beatMap.IsKeyMode() {
		render.Error(ctx, http.StatusBadRequest, errors.New("unsupported convert mode"))
	}

	notes := new(model.QQX5BeatmapNormal)
	notes.ParseFromMalodyBeatmap(beatMap, param.BPM)
	offset := beatMap.CalcNoteHitTime(beatMap.Note[0].Beat, beatMap.BpmOffset())

	render.OK(ctx, httpdto.MCZToXMLResponse{
		Offset:  offset,
		Creator: beatMap.Meta.Creator,
		Title:   beatMap.Meta.Song.Title,
		Artist:  beatMap.Meta.Song.Artist,
		Notes:   notes.ToHTML(),
	})
}
