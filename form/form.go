package form

import "mime/multipart"

type UpdateXMLBPMParam struct {
	BPM  float64               `form:"bpm" binding:"required"`
	File *multipart.FileHeader `form:"file" binding:"required"`
}
