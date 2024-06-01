package httpdto

import "mime/multipart"

type OSUToXMLParam struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
}

type UpdateXMLBPMParam struct {
	BPM  float64               `form:"bpm" binding:"required"`
	File *multipart.FileHeader `form:"file" binding:"required"`
}

type UpdateXMLBPMResponse struct {
	Notes string `json:"notes"`
}
