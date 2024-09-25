package models

type Part struct {
	FileName   string `json:"file_name"`
	PartNumber int    `json:"part_number"`
	Size       int64  `json:"size"`
}
