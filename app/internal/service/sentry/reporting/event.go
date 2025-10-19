package reporting

import (
	"catcher/app/internal/lib/ioimage"
	"catcher/app/internal/sentry/normalize"
	"fmt"
	"log/slog"

	"github.com/getsentry/sentry-go"
)

const (
	eventLogger  = "1с_service_logger"
	emptyMessage = "Ошибка выполнения"
)

func (r Report) event() (*sentry.Event, error) {

	data := r.Data
	Id := normalize.Id(data.Id)

	event := &sentry.Event{
		Timestamp:   data.Time,
		EventID:     sentry.EventID(Id),
		Message:     r.message(),
		Level:       sentry.LevelError,
		Contexts:    r.contexts(),
		User:        r.user(),
		Breadcrumbs: r.breadcrumbs(),
		Exception:   r.exception(),
		Modules:     r.modules(),
		Release:     normalize.Release(data.AdditionalInfo, data.ConfigInfo.Version),
		Platform:    r.Prj.Sentry.Platform,
		Dist:        r.dist(),
		Environment: r.Prj.Sentry.Environment,
		Attachments: r.attachments(),
		Tags:        r.tags(),
		Logger:      eventLogger,
	}

	if len(event.Exception) == 0 && event.Message == "" {
		event.Message = emptyMessage
	}

	return event, nil

}

func (r Report) message() string {
	return ""
}

func (r Report) dist() string {
	return r.Prj.Id
}

// Contexts Interface
// https://develop.sentry.dev/sdk/data-model/event-payloads/contexts/
func (r Report) contexts() map[string]sentry.Context {

	ctx := make(map[string]sentry.Context)

	clientInfo := r.Data.ClientInfo
	if clientInfo.AppName != "" {
		ctx["app"] = sentry.Context{"app_build": clientInfo.PlatformType,
			"app_version": clientInfo.AppVersion,
			"app_name":    clientInfo.AppName}
	}

	systemInfo := clientInfo.SystemInfo
	if systemInfo.OsVersion != "" {
		ctx["os"] = sentry.Context{"name": systemInfo.OsVersion}
	}

	sessionInfo := r.Data.SessionInfo.UserInfo.SessionInfo

	deviceName := sessionInfo.Device
	if deviceName == "" {
		deviceName = "PC"
	}

	if systemInfo.ClientID != "" {
		ctx["device"] = sentry.Context{
			"name":            deviceName,
			"memory_size":     systemInfo.FullRAM,
			"free_memory":     systemInfo.FreeRAM,
			"cpu_description": systemInfo.Processor,
			"model_id":        systemInfo.ClientID}
	}

	sessionData := sentry.Context{}
	connection := sessionInfo.Connection
	if connection != 0 {
		sessionData["Connection"] = connection
	}
	session := sessionInfo.Session
	if session != 0 {
		sessionData["Session"] = session
	}
	if len(sessionData) != 0 {
		ctx["Session Data"] = sessionData
	}

	return ctx
}

func (r Report) user() sentry.User {

	sessionInfo := r.Data.SessionInfo
	return sentry.User{
		ID:        sessionInfo.UserInfo.Id,
		Username:  sessionInfo.UserName,
		IPAddress: normalize.Ip(sessionInfo.UserInfo.SessionInfo.IP),
	}
}

func (r Report) breadcrumbs() []*sentry.Breadcrumb {

	brc := make([]*sentry.Breadcrumb, 0)
	brc = r.breadcrumbsErrorInfo(brc)

	return brc
}

func (r Report) breadcrumbsErrorInfo(brc []*sentry.Breadcrumb) []*sentry.Breadcrumb {

	// описание ошибки на втором уровне
	// забираем все строки, кроме первой, т.к. первая уже будет в message и сюда заполнится автоматом
	errorInfo := r.Data.ErrorInfo.ApplicationErrorInfo.Errors

	for i := len(errorInfo) - 1; i >= 0; i-- {

		errorData := errorInfo[i] //[2]any
		message := normalize.Message(fmt.Sprintf("%v", errorData[0]))
		if message == "" {
			continue
		}

		сategory := "message"

		brc = append(brc, &sentry.Breadcrumb{
			Type:      "info",
			Category:  сategory,
			Message:   message,
			Level:     sentry.LevelError,
			Timestamp: r.Data.Time,
		})

	}

	return brc

}

func (r Report) modules() map[string]string {

	const emptyVersion string = "0.0.1"

	result := make(map[string]string)

	// Заполняем включенные расширения, кроме постоянных
	extentions := r.Data.ConfigInfo.Extentions
	for _, val := range extentions {

		name := val[0]
		if name == "" || r.Prj.ExistExtention(name) {
			continue
		}

		extention := normalize.SplitString(name, "(", ")")

		version := extention[1]
		if version == "" {
			version = emptyVersion
		}
		result[extention[0]] = version
	}

	return result

}

func (r Report) attachments() []*sentry.Attachment {

	cAttachments := r.Prj.Sentry.Attachments
	if !cAttachments.Use {
		return nil
	}

	filesData := r.Files
	if len(filesData) == 0 {
		return nil
	}

	var result []*sentry.Attachment

	const full = 100
	var percent = full
	if cAttachments.Сompress.Use || cAttachments.Сompress.Percent < percent {
		percent = cAttachments.Сompress.Percent
	}

	for _, v := range filesData {

		var data = v.Data
		if percent < full {
			if cdata, err := ioimage.Compress(data, uint(percent)); err == nil {
				r.Logger.Debug("Сжато вложение",
					r.Logger.Str("Name", v.Name),
					slog.Int("Percent", percent))
				data = cdata
			}
		}

		result = append(result, &sentry.Attachment{
			Filename: v.Name,
			Payload:  data,
		})
	}

	return result
}

func (r Report) tags() map[string]string {

	res := make(map[string]string)

	userInfo := r.Data.SessionInfo.UserInfo

	if userInfo.City != "" {
		res["place.city"] = userInfo.City
	}
	if userInfo.Branch != "" {
		res["place.branch"] = userInfo.Branch
	}
	if userInfo.Position != "" {
		res["user.position"] = userInfo.Position
	}
	if !userInfo.Started.IsZero() {
		res["user.started"] = userInfo.Started.Format("02-01-2006T15:04:05")
	}

	return res
}
