package replicate

import (
	"catcher/app/internal/config"
	"catcher/app/internal/models"
	"context"
	"fmt"
	"time"

	"github.com/jinzhu/copier"
)

type Geter interface {
	Get(ctx context.Context, input string) (*models.UserInfo, error)
}

func (s Service) AddUserInfo(g Geter, prj config.Project, rd *models.RepportData) {

	const op = "userInfo.AddUserInfo"

	var userName string
	if s.Config.UseDebug() && prj.Service.Test.UserName != "" {
		userName = prj.Service.Test.UserName
	} else {
		userName = rd.Data.SessionInfo.UserName
	}

	res := models.UserInfo{}
	key := fmt.Sprintf("%s:%s:%s", prj.Id, op, userName)

	useCache := prj.Service.Cache.Use
	found := false
	if useCache {
		var err error
		found, err = s.Cacher.Get(s.Ctx, key, &res)
		if found {
			s.Logger.Debug("Используем кэш UserInfo",
				s.Logger.Op(op),
				s.Logger.Str("key", key))
		} else if err != nil {
			s.Logger.Error("Ошибка получения значения из кэша",
				s.Logger.Op(op),
				s.Logger.Str("key", key),
				s.Logger.Err(err))
		}
	}

	if !found {
		ptr, err := g.Get(s.Ctx, userName)
		if err != nil {
			return
		}
		if ptr != nil {
			res = *ptr
			if useCache {
				// Проверка: не кэшировать пустой UserInfo
				if res != (models.UserInfo{}) {
					s.Cacher.Set(s.Ctx, key, res, time.Duration(prj.Service.Cache.Expiration)*time.Minute)
				} else {
					s.Logger.Warn("Попытка добавить в кэш пустой UserInfo",
						s.Logger.Op(op),
						s.Logger.Str("key", key))
				}
			}
		}
	}

	copier.Copy(&rd.Data.SessionInfo.UserInfo, &res)
}
