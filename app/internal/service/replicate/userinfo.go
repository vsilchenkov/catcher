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

	var res *models.UserInfo
	key := fmt.Sprintf("%s:%s:%s", op, prj.Name, userName)

	useCache := prj.Service.Cache.Use
	if useCache {
		if x, found := s.Cacher.Get(s.Ctx, key); found {
			s.Logger.Debug("Используем кэш UserInfo",
				s.Logger.Op(op),
				s.Logger.Str("key", key))
			res = x.(*models.UserInfo)
		}
	}

	if res == nil {
		var err error
		res, err = g.Get(s.Ctx, userName)
		if err != nil {
			return
		}
		if useCache {
			s.Cacher.Set(s.Ctx, key, res, time.Duration(prj.Service.Cache.Expiration)*time.Minute)
		}

	}

	copier.Copy(&rd.Data.SessionInfo.UserInfo, &res)
}
