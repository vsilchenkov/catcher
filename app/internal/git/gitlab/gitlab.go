package gitlab

import (
	"sync"

	"github.com/cockroachdb/errors"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

type Git struct {
	*gitlab.Client
	Project Project
}

type Project struct {
	Url    string
	Path   string
	Token  string
	Branch string
}

type gitssmap map[Project]Git

var gits gitssmap
var gitsMu sync.RWMutex

func init() {
	gits = make(gitssmap)
}

func New(prj Project) (*Git, error) {

	g, err := gitlab.NewClient(prj.Token, gitlab.WithBaseURL(prj.Url))
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create client")
	}

	git := Git{g, prj}

	gitsMu.Lock()
	gits[prj] = git
	gitsMu.Unlock()

	return &git, nil

}

func Get(prj Project) (*Git, error) {

	gitsMu.RLock()
	git, ok := gits[prj]
	gitsMu.RUnlock()

	if  ok {
		return &git, nil
	}

	return New(prj)
}

func (g Git) GetFileContent(filePath string) (*string, error) {

	gf := &gitlab.GetRawFileOptions{
		Ref: gitlab.Ptr(g.Project.Branch),
	}

	file, _, err := g.RepositoryFiles.GetRawFile(g.Project.Path, filePath, gf)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get content")
	}

	content := string(file)
	return &content, nil

}
