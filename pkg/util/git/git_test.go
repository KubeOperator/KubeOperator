package git

import (
	"log"
	"testing"

	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

// func generateTree(root *dto.TreeNode) error {
// 	children, err := ioutil.ReadDir(root.Path)
// 	if err != nil {
// 		return err
// 	}
// 	for _, c := range children {
// 		if c.Name() != ".git" {
// 			var node dto.TreeNode
// 			if !c.IsDir() {
// 				node.Name = c.Name()
// 				node.Dir = false
// 				node.Path = path.Join(root.Path, c.Name())
// 			} else {
// 				node.Name = c.Name()
// 				node.Dir = true
// 				node.Path = path.Join(root.Path, c.Name())
// 				node.Children = make([]*dto.TreeNode, 0)
// 				err := generateTree(&node)
// 				if err != nil {
// 					return err
// 				}
// 			}
// 			root.Children = append(root.Children, &node)
// 			sort.Sort(root.Children)
// 		}
// 	}
// 	return nil
// }

func TestCloneRepository(t *testing.T) {
	err := CloneRepository("https://github.com/KubeOperator/MultiClusterRepositoryExample.git", "/var/ko/data/test", "master", &http.BasicAuth{
		Username: "Aaron3S",
		Password: "scydeai251",
	})
	if err != nil {
		log.Fatal(err)
	}
}

func TestUpdateRepository(t *testing.T) {
	err := UpdateRepository("/var/ko/data/test", "master", &http.BasicAuth{
		Username: "Aaron3S",
		Password: "scydeai251",
	})
	if err != nil {
		log.Fatal(err)
	}
}
