package git

import (
	"fmt"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"k8s.io/kubernetes/third_party/forked/etcd237/pkg/fileutil"
)

type Credential struct {
	Username string
	Password string
}

func CloneRepository(url string, path string, branch string, auth transport.AuthMethod) error {
	exists := fileutil.Exist(path)
	if exists {
		return fmt.Errorf("path: %s already exists", path)
	}
	r, err := git.PlainInit(path, false)
	if err != nil {
		return err
	}
	if _, err = r.CreateRemote(&config.RemoteConfig{
		Name:  "origin",
		URLs:  []string{url},
		Fetch: []config.RefSpec{config.RefSpec(fmt.Sprintf("+refs/heads/%s:refs/remotes/origin/%s", branch, branch))},
	}); err != nil {
		return err
	}
	workTree, err := r.Worktree()
	if err != nil {
		return err
	}
	if err := workTree.Pull(&git.PullOptions{RemoteName: "origin", Auth: auth, ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", branch))}); err != nil {
		return err
	}
	return err
}

func UpdateRepository(path string, branch string, auth transport.AuthMethod) error {
	r, err := git.PlainOpen(path)
	if err != nil {
		return err
	}
	workTree, err := r.Worktree()
	if err != nil {
		return err
	}
	if err := workTree.Pull(&git.PullOptions{RemoteName: "origin", Auth: auth, ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", branch))}); err != nil && err != git.NoErrAlreadyUpToDate {
		return err
	}
	return nil
}

func PushRepository(path string, auth transport.AuthMethod) error {
	r, err := git.PlainOpen(path)
	if err != nil {
		return err
	}
	workTree, err := r.Worktree()
	if err != nil {
		return err
	}
	err = workTree.AddWithOptions(&git.AddOptions{
		Path: path,
		All:  true,
	})
	if err != nil {
		return err
	}
	_, err = workTree.Commit("Update by KubeOperator", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "KubeOperator",
			Email: "",
			When:  time.Now(),
		},
	})
	if err != nil {
		return err
	}
	err = r.Push(&git.PushOptions{
		RemoteName: "origin",
		Auth:       auth,
		Force:      true,
	})
	if err != nil {
		return err
	}
	return nil
}
