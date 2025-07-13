package config

import (
	"catcher/app/internal/git/gitlab"
	"slices"

	"github.com/jinzhu/copier"
)

type Project struct {
	Name    string `yaml:"Name" binding:"required"`
	Id      string `yaml:"ID" binding:"required"`
	Service struct {
		Use         bool   `yaml:"Use"`
		Url         string `yaml:"Url"`
		IimeOut     int    `yaml:"IimeOut"`
		Credintials struct {
			UserName string `yaml:"UserName"`
			Password string `yaml:"Password"`
		} `yaml:"Credintials"`
		Cache     cache `yaml:"Cache"`
		Exeptions struct {
			Use   bool  `yaml:"Use"`
			Cache cache `yaml:"Cache"`
		} `yaml:"Exeptions"`
		Test struct {
			UserName string `yaml:"UserName"`
		} `yaml:"Test"`
	} `yaml:"Service"`
	release string `yaml:"Release"`
	Sentry  struct {
		Dsn           string `yaml:"Dsn" binding:"required"`
		Environment   string `yaml:"Environment" binding:"required"`
		Platform      string `yaml:"Platform" binding:"required"`
		ContextAround struct {
			Use      bool  `yaml:"Use"`
			Quantity int   `yaml:"Quantity"`
			Cache    cache `yaml:"Cache"`
		} `yaml:"ContextAround"`
		Attachments struct {
			Use      bool `yaml:"Use" binding:"required"`
			Сompress struct {
				Use     bool `yaml:"Use"`
				Percent int  `yaml:"Percent" validate:"required,min=1,max=100"`
			} `yaml:"Сompress"`
		} `yaml:"Attachments"`
		SendingCache cache `yaml:"SendingCache"`
	} `yaml:"Sentry" binding:"required"`
	Git struct {
		Use            bool   `yaml:"Use"`
		Url            string `yaml:"Url"`
		Path           string `yaml:"Path"`
		Token          string `yaml:"Token"`
		Branch         string `yaml:"Branch"`
		SourceCodeRoot string `yaml:"SourceCodeRoot"`
	} `yaml:"Git"`
	Extentions []string `yaml:"Extentions"`
}

type cache struct {
	Use        bool `yaml:"Use"`
	Expiration int  `yaml:"Expiration" validate:"required,min=1"`
}

func (p *Project) SetRelease(release string) {
	p.release = release
}

func (p Project) Release() string {
	return p.release
}

func (p Project) ExistExtention(input string) bool {
	return slices.Contains(p.Extentions, input)
}

type gitter interface {
	GetFileContent(filePath string) (*string, error)
}

func (p Project) GetGit() (gitter, error) {

	if !p.Git.Use {
		return nil, nil
	}

	var prjGit gitlab.Project
	copier.Copy(&prjGit, &p.Git)
	return gitlab.Get(prjGit)

}
