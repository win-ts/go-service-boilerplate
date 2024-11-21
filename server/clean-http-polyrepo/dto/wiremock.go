package dto

type WiremockGetTestResponse struct {
	Message string `json:"message"`
}

type WiremockGetTestHeader struct {
	ContentType string
	RequestID   string
}

func (h WiremockGetTestHeader) ToMap() map[string]string {
	m := make(map[string]string)
	if h.RequestID != "" {
		m["X-Request-ID"] = h.RequestID
	}
	m["Content-Type"] = h.ContentType
	return m
}
