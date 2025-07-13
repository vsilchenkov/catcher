package models

import (
	"catcher/app/internal/config"
	"time"
)

// https://its.1c.ru/db/v8327doc#bookmark:dev:TI000002558
type Repport struct {
	Time       time.Time `json:"time"`
	Id         string    `json:"id"`
	ClientInfo struct {
		PlatformType string `json:"platformType"`
		AppVersion   string `json:"appVersion"`
		AppName      string `json:"appName"`
		SystemInfo   struct {
			OsVersion string `json:"osVersion"`
			FullRAM   int64  `json:"fullRAM"`
			FreeRAM   int64  `json:"freeRAM"`
			Processor string `json:"processor"`
			Useragent string `json:"useragent"`
			ClientID  string `json:"clientID"`
		} `json:"systemInfo"`
	} `json:"clientInfo"`
	SessionInfo struct {
		UserName                           string   `json:"userName"`
		DataSeparation                     string   `json:"dataSeparation"`
		PlatformInterfaceLanguageCode      string   `json:"platformInterfaceLanguageCode"`
		ConfigurationInterfaceLanguageCode string   `json:"configurationInterfaceLanguageCode"`
		LocaleCode                         string   `json:"localeCode"`
		UserInfo                           UserInfo // Добавляются сервисом userinfo
	} `json:"sessionInfo"`
	InfoBaseInfo struct {
		LocaleCode string `json:"localeCode"`
	} `json:"infoBaseInfo"`
	ServerInfo struct {
		PlatformType string `json:"platformType"`
		AppVersion   string `json:"appVersion"`
		Dbms         string `json:"dbms"`
	} `json:"serverInfo"`
	ConfigInfo struct {
		Name               string      `json:"name"`
		Description        string      `json:"description"`
		Version            string      `json:"version"`
		CompatibilityMode  string      `json:"compatibilityMode"`
		Hash               string      `json:"hash"`
		ChangeEnabled      bool        `json:"changeEnabled"`
		Extentions         []extention `json:"extentions"`
		DisabledExtentions []extention `json:"disabledExtentions"`
	} `json:"configInfo"`
	ErrorInfo struct {
		UserDescription string `json:"userDescription"`
		SystemErrorInfo struct {
			ClientStack     string `json:"clientStack"`
			ClientStackHash string `json:"clientStackHash"`
			ServerStack     string `json:"serverStack"`
			ServerStackHash string `json:"serverStackHash"`
			SystemCrash     bool   `json:"systemCrash"`
		} `json:"systemErrorInfo"`
		ApplicationErrorInfo struct {
			Errors    []errorsHeap `json:"errors"`
			Stack     []stackHeap  `json:"stack"`
			StackHash string       `json:"stackHash"`
		} `json:"applicationErrorInfo"`
	} `json:"errorInfo"`
	Screenshot struct {
		File string `json:"file"`
	} `json:"screenshot"`
	AdditionalInfo  string   `json:"additionalInfo"`
	AdditionalData  string   `json:"additionalData"`
	AdditionalFiles []string `json:"additionalFiles"`
	Dump            struct {
		TypeDump        string `json:"type"`
		File            string `json:"file"`
		ReasonForNoDump struct {
			GenericFailure        string `json:"genericFailure"`
			UserRefused           string `json:"userRefused"`
			InsufficientResources string `json:"insufficientResources"`
		} `json:"reasonForNoDump"`
	} `json:"dump"`
}

type extention [2]string
type errorsHeap [2]any
type stackHeap [3]any

type FileData struct {
	Name string
	Data []byte
}

type RepportData struct {
	ID          string
	Prj         config.Project
	Data        *Repport
	Files       []FileData
	Src         string
	SrcDirFiles string
}

func (r Repport) ProjectByConfig(c *config.Config) (config.Project, error) {

	var prj config.Project
	var err error

	additionalInfo := r.AdditionalInfo
	if additionalInfo != "" {
		prj, err = c.ProjectById(additionalInfo)
	} else {
		prj, err = c.ProjectByName(r.ConfigInfo.Name)
	}

	return prj, err
}

type UserInfo struct {
	Id          string
	City        string
	Branch      string
	Position    string
	Started     time.Time
	SessionInfo struct {
		IP         string
		Device     string
		Session    int
		Connection int
	}
}

func (u UserInfo) Empty() bool {
	return u == UserInfo{}
}
