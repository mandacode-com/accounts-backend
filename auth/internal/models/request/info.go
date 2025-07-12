package reqmodels

type RequestInfo struct {
	IP        string `json:"ip"`
	UserAgent string `json:"user_agent"`
	DeviceID  string `json:"device_id"`
	Location  string `json:"location"`
	IsMobile  bool   `json:"is_mobile"`
	IsWeb     bool   `json:"is_web"`
}
