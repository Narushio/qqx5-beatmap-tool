package usecase

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"

	"github.com/gin-gonic/gin"

	"github.com/Narushio/qqx5-beatmap-tool/form"
	"github.com/Narushio/qqx5-beatmap-tool/model"
)

type FileUseCase struct{}

func (f *FileUseCase) UpdateXMLBPM(ctx *gin.Context, param *form.UpdateXMLBPMParam) (*string, error) {
	xmlFile, err := param.File.Open()
	if err != nil {
		return nil, fmt.Errorf("opening file")
	}
	defer xmlFile.Close()

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, xmlFile); err != nil {
		return nil, fmt.Errorf("copying file")
	}

	level := &model.Level{}
	if err := xml.Unmarshal(buf.Bytes(), level); err != nil {
		return nil, fmt.Errorf("unmarshalling XML")
	}

	level = level.ResetNotesWithBPM(param.BPM)
	noteStr := level.NoteInfo.Normal.ToHTML()
	return &noteStr, nil
}
