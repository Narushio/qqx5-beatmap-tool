package handler

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"

	"github.com/Narushio/qqx5-beatmap-tool/handler/httpdto"
	"github.com/Narushio/qqx5-beatmap-tool/handler/render"
	"github.com/Narushio/qqx5-beatmap-tool/model"
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
		render.Error(ctx, http.StatusInternalServerError, fmt.Errorf("open file: %w", err))
	}

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, f); err != nil {
		render.Error(ctx, http.StatusInternalServerError, fmt.Errorf("copy file: %w", err))
	}

	beatMap := new(model.QQX5BeatmapLevel)
	if err := xml.Unmarshal(buf.Bytes(), beatMap); err != nil {
		render.Error(ctx, http.StatusBadRequest, fmt.Errorf("unmarshalling XML: %w", err))
	}

	render.OK(ctx, httpdto.UpdateXMLBPMResponse{Notes: beatMap.ResetNotesWithBPM(param.BPM).NoteInfo.Normal.ToHTML()})
}

func (h *Convert) OSUToXML(ctx *gin.Context) {
	param := new(httpdto.OSUToXMLParam)
	if err := ctx.Bind(param); err != nil {
		render.Error(ctx, http.StatusBadRequest, fmt.Errorf("binding param: %w", err))
	}

	f, err := param.File.Open()
	if err != nil {
		render.Error(ctx, http.StatusInternalServerError, fmt.Errorf("open file: %w", err))
	}

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, f); err != nil {
		render.Error(ctx, http.StatusInternalServerError, fmt.Errorf("copy file: %w", err))
	}

	beatMap := new(model.OSUBeatMap)
	if err := beatMap.Parse(buf.Bytes()); err != nil {
		render.Error(ctx, http.StatusInternalServerError, fmt.Errorf("parse osu file: %w", err))
	}
}
