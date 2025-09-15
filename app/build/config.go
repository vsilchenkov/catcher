package build

type Option struct {
	VersionInfo
	Version     string
	ProjectName string
	WorkingDir  string
	Interactive bool
}

type VersionInfo struct {
	FixedFileInfo  FixedFileInfo  `json:"FixedFileInfo"`
	StringFileInfo StringFileInfo `json:"StringFileInfo"`
	IconPath       string         `json:"IconPath"`
}
type FileVersion struct {
	Major int `json:"Major"`
	Minor int `json:"Minor"`
	Patch int `json:"Patch"`
	Build int `json:"Build"`
}

type FixedFileInfo struct {
	FileVersion FileVersion `json:"FileVersion"`
}

type StringFileInfo struct {
	ProductVersion   string `json:"ProductVersion"`
	CompanyName      string `json:"CompanyName"`
	FileDescription  string `json:"FileDescription"`
	InternalName     string `json:"InternalName"`
	LegalCopyright   string `json:"LegalCopyright"`
	LegalTrademarks  string `json:"LegalTrademarks"`
	OriginalFilename string `json:"OriginalFilename"`
	PrivateBuild     string `json:"PrivateBuild"`
	ProductName      string `json:"ProductName"`
	SpecialBuild     string `json:"SpecialBuild"`
	Comments         string `json:"Comments"`
}
