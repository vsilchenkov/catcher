package models

type RegistryInfo struct {
	NeedSendReport bool   `json:"needSendReport"`
	UserMessage    string `json:"userMessage"`
	DumpType       int    `json:"dumpType"`
}

type RegistryInput struct {
	ConfigName                         string   `json:"configName" binding:"required"`
	ConfigHash                         string   `json:"configHash"`
	ConfigVersion                      string   `json:"configVersion"`
	AppStackHash                       string   `json:"appStackHash"`
	ClientStackHash                    string   `json:"clientStackHash"`
	PlatformType                       string   `json:"platformType"`
	AppName                            string   `json:"appName"`
	AppVersion                         string   `json:"appVersion"`
	PlatformInterfaceLanguageCode      string   `json:"platformInterfaceLanguageCode"`
	ConfigurationInterfaceLanguageCode string   `json:"configurationInterfaceLanguageCode"`
	ErrorCategories                    []string `json:"ErrorCategories"`
	ClientID                           string   `json:"clientID"`
	ReportID                           string   `json:"reportID"`
}

type RegistryPushReportInput struct {
	ID   string
	Data []byte
}
type RegistryPushReportResult struct {
	ID      string
	EventID *EventID
}

func NewRegistryInput() RegistryInput {
	return RegistryInput{
		ConfigName: "none",
	}
}
