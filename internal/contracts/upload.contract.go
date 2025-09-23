package contracts

type UploadRequest struct {
	Folder string `form:"folder" json:"folder"`
	File   string `form:"file" json:"file" binding:"required"`
}

type UploadResponse struct {
	URL      string `json:"url"`
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
	Provider string `json:"provider"`
}

type FileInfo struct {
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
	URL      string `json:"url"`
	Provider string `json:"provider"`
}

type RetrieveRequest struct {
	Filename string `uri:"filename" binding:"required"`
	Folder   string `query:"folder"`
}

type ErrorResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
