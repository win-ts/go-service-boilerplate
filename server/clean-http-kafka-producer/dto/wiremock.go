package dto

// WiremockGetTestResponse represents the response of the wiremock get test
type WiremockGetTestResponse struct {
	Message string `json:"message"`
}

// WiremockGetTestHeader represents the header of the wiremock get test
type WiremockGetTestHeader struct {
	ContentType string
	RequestID   string
}

// ToMap converts the header to a map
func (h WiremockGetTestHeader) ToMap() map[string]string {
	m := make(map[string]string)
	if h.RequestID != "" {
		m["X-Request-ID"] = h.RequestID
	}
	m["Content-Type"] = h.ContentType
	return m
}
