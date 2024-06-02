package httpdto

import "mime/multipart"

type OsuToXMLParam struct {
	BPM  float64               `form:"bpm" binding:"required"`
	File *multipart.FileHeader `form:"file" binding:"required"`
}

type OsuToXMLResponse struct {
	Creator string  `json:"creator"`
	Offset  float64 `json:"offset"`
	Notes   string  `json:"notes"`
}

type UpdateXMLBPMParam struct {
	BPM  float64               `form:"bpm" binding:"required"`
	File *multipart.FileHeader `form:"file" binding:"required"`
}

type UpdateXMLBPMResponse struct {
	Notes string `json:"notes"`
}

type MCZToXMLParam struct {
	BPM  float64               `form:"bpm" binding:"required"`
	File *multipart.FileHeader `form:"file" binding:"required"`
}

type MCZToXMLResponse struct {
	Creator string  `json:"creator"`
	Offset  float64 `json:"offset"`
	Title   string  `json:"title"`
	Artist  string  `json:"artist"`
	Notes   string  `json:"notes"`
}
