package service

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/fluxcd/flux/pkg/cluster"
	"github.com/fluxcd/flux/pkg/cluster/kubernetes"
	"github.com/fluxcd/flux/pkg/event"
	"github.com/fluxcd/flux/pkg/git"
	"github.com/fluxcd/flux/pkg/manifests"
	"github.com/fluxcd/flux/pkg/resource"
	"github.com/fluxcd/flux/pkg/ssh"
	fluxsync "github.com/fluxcd/flux/pkg/sync"
	helmopclient "github.com/fluxcd/helm-operator/pkg/client/clientset/versioned"
	"github.com/pkg/errors"
	crd "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	k8sclientdynamic "k8s.io/client-go/dynamic"
	k8sclient "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"os/exec"
	"sync"
	"time"
)

type Cluster struct {
	cluster.Cluster
	KoCluster model.Cluster
	Manifests manifests.Manifests
}

func NewCluster(clusterName string) (*Cluster, error) {
	service := NewClusterService()
	koCluster, err := service.Get(clusterName)
	if err != nil {
		return nil, err
	}
	endpoint, err := service.GetApiServerEndpoint(clusterName)
	if err != nil {
		return nil, err
	}
	secrets, err := service.GetSecrets(clusterName)
	if err != nil {
		return nil, err
	}
	restClientConfig := rest.Config{
		Host:            fmt.Sprintf("https://%s:%d", endpoint.Address, endpoint.Port),
		QPS:             50.0,
		Burst:           100,
		BearerToken:     secrets.KubernetesToken,
		TLSClientConfig: rest.TLSClientConfig{Insecure: true},
	}

	var k8s cluster.Cluster
	var k8sManifests manifests.Manifests
	{
		clientset, err := k8sclient.NewForConfig(&restClientConfig)
		if err != nil {
			return nil, err
		}
		dynamicClientset, err := k8sclientdynamic.NewForConfig(&restClientConfig)
		if err != nil {
			return nil, err
		}

		hrClientset, err := helmopclient.NewForConfig(&restClientConfig)
		if err != nil {
			return nil, err
		}

		crdClient, err := crd.NewForConfig(&restClientConfig)
		if err != nil {
			return nil, err
		}
		discoClientset := kubernetes.MakeCachedDiscovery(clientset.Discovery(), crdClient, make(chan struct{}))

		sshKeyRing := ssh.NewNopSSHKeyRing()

		kubectl, err := exec.LookPath("kubectl")
		client := kubernetes.MakeClusterClientset(clientset, dynamicClientset, hrClientset, discoClientset)
		kubectlApplier := kubernetes.NewKubectl(kubectl, &restClientConfig)
		allowedNamespaces := make(map[string]struct{})
		k8sInst := kubernetes.NewCluster(client, kubectlApplier, sshKeyRing, nil, allowedNamespaces, nil, []string{})
		if err := k8sInst.Ping(); err != nil {
			return nil, err
		}
		namespacer, err := kubernetes.NewNamespacer(discoClientset, constant.DefaultNamespace)
		if err != nil {
			return nil, err
		}
		k8sManifests = kubernetes.NewManifests(namespacer, nil)
	}
	return &Cluster{
		Cluster:   k8s,
		KoCluster: koCluster.Cluster,
		Manifests: k8sManifests,
	}, nil

}

type MultiClusterRepositorySync struct {
	Repo       *model.MultiClusterRepository
	GitRepo    *git.Repo
	GitTimeout time.Duration
	Clusters   []Cluster
}

func NewMultiClusterRepositorySync(repository *model.MultiClusterRepository, clusterNames []string) (*MultiClusterRepositorySync, error) {
	var ms MultiClusterRepositorySync
	ms.Repo = repository
	ms.GitTimeout = 5 * time.Minute
	for _, c := range clusterNames {
		cls, err := NewCluster(c)
		if err != nil {
			return nil, err
		}
		ms.Clusters = append(ms.Clusters, *cls)
	}
	gitRemote := git.Remote{URL: repository.Source}
	repo := git.NewRepo(gitRemote, git.PollInterval(repository.SyncInterval), git.Timeout(ms.GitTimeout), git.Branch(repository.Branch), git.IsReadOnly(true))
	ms.GitRepo = repo
	return &ms, nil
}

func (m *MultiClusterRepositorySync) Sync() {
	interval := (time.Now().UnixNano() - m.Repo.LastSyncTime.UnixNano()) / 1e6
	if interval < m.Repo.SyncInterval || m.Repo.Status == constant.StatusPending {
		return
	}
	m.Repo.Status = constant.StatusRunning
	m.Repo.LastSyncTime = time.Now()
	db.DB.Save(m.Repo)

	err := m.Repo.Pull()
	if err != nil {
		log.Error(err)
		return
	}
	newSyncHead, err := m.GitRepo.BranchHead(context.TODO())
	if err != nil {
		log.Error(err)
		return
	}
	if m.Repo.LastSyncHead == newSyncHead {
		return
	}

	var syncLog model.MultiClusterSyncLog
	syncLog.Status = constant.StatusRunning
	syncLog.MultiClusterRepositoryID = m.Repo.ID
	db.DB.Create(&syncLog)
	hash := makeGitConfigHash(m.GitRepo.Origin(), m.Repo.Branch)
	wg := &sync.WaitGroup{}
	for _, c := range m.Clusters {
		c := c
		go func() {
			wg.Add(1)
			var clusterSyncLog model.MultiClusterSyncClusterLog
			clusterSyncLog.MultiClusterSyncLogID = syncLog.ID
			clusterSyncLog.Status = constant.StatusRunning
			db.DB.Create(&clusterSyncLog)
			store, clean, err := m.getManifestStoreByRevision(context.TODO(), c.Manifests, newSyncHead)
			if err != nil {
				log.Error(err)
				clusterSyncLog.Status = constant.StatusFailed
				clusterSyncLog.Message = err.Error()
				db.DB.Save(&syncLog)
				return
			}
			resourceMap, errEvents, err := doSync(context.TODO(), store, c.Cluster, hash)
			if err != nil {
				log.Error(err)
				clusterSyncLog.Status = constant.StatusFailed
				clusterSyncLog.Message = err.Error()
				db.DB.Save(&syncLog)
				return
			}
			clusterSyncLog.Status = constant.StatusSuccess
			db.DB.Save(&syncLog)

			for _, v := range resourceMap {
				var resourceLog model.MultiClusterSyncClusterResourceLog
				resourceLog.MultiClusterSyncClusterLogID = clusterSyncLog.ID
				resourceLog.ResourceName = v.ResourceID().String()
				resourceLog.SourceFile = v.Source()
				resourceLog.Status, resourceLog.Message = func() (string, string) {
					for i := range errEvents {
						if errEvents[i].ID.String() == v.ResourceID().String() {
							return constant.StatusFailed, errEvents[i].Error
						}
					}
					return constant.StatusSuccess, ""
				}()
				db.DB.Create(&resourceLog)
			}
			wg.Done()
			defer clean()
		}()
	}
	wg.Wait()
	syncLog.Status = constant.StatusSuccess
	db.DB.Save(&syncLog)
	if err := refresh(context.TODO(), m.GitTimeout, m.GitRepo); err != nil {
		log.Error(err)
		return
	}

}

func makeGitConfigHash(remote git.Remote, branch string) string {
	urlbit := remote.SafeURL()
	pathshash := sha256.New()
	pathshash.Write([]byte(urlbit))
	pathshash.Write([]byte(branch))
	for _, path := range []string{} {
		pathshash.Write([]byte(path))
	}
	return base64.RawURLEncoding.EncodeToString(pathshash.Sum(nil))
}

func doSync(ctx context.Context, manifestsStore manifests.Store, clus cluster.Cluster, syncSetName string) (map[string]resource.Resource, []event.ResourceError, error) {
	resources, err := manifestsStore.GetAllResourcesByID(ctx)
	if err != nil {
		return nil, nil, errors.Wrap(err, "loading resources from repo")
	}

	var resourceErrors []event.ResourceError
	if err := fluxsync.Sync(syncSetName, resources, clus); err != nil {
		switch syncerr := err.(type) {
		case cluster.SyncError:
			for _, e := range syncerr {
				resourceErrors = append(resourceErrors, event.ResourceError{
					ID:    e.ResourceID,
					Path:  e.Source,
					Error: e.Error.Error(),
				})
			}
		default:
			return nil, nil, err
		}
	}
	return resources, resourceErrors, nil
}

func (m *MultiClusterRepositorySync) getManifestStoreByRevision(ctx context.Context, man manifests.Manifests, revision string) (store manifests.Store, cleanupClone func(), err error) {
	clone, cleanupClone, err := m.cloneRepo(ctx, revision)
	if err != nil {
		return nil, nil, errors.Wrap(err, "cloning repo")
	}

	store, err = m.getManifestStore(clone, man)
	return store, cleanupClone, err
}

func (m *MultiClusterRepositorySync) cloneRepo(ctx context.Context, revision string) (clone *git.Export, cleanup func(), err error) {
	ctxGitOp, cancel := context.WithTimeout(ctx, m.GitTimeout)
	defer cancel()
	clone, err = m.GitRepo.Export(ctxGitOp, revision)
	if err != nil {
		return nil, nil, err
	}

	cleanup = func() {
		if err := clone.Clean(); err != nil {
			log.Error(err)
		}
	}

	return clone, cleanup, nil
}

type repo interface {
	Dir() string
}

func (m *MultiClusterRepositorySync) getManifestStore(r repo, man manifests.Manifests) (manifests.Store, error) {
	absPaths := git.MakeAbsolutePaths(r, []string{})
	return manifests.NewRawFiles(r.Dir(), absPaths, man), nil
}

func refresh(ctx context.Context, timeout time.Duration, repo *git.Repo) error {
	ctxGitOp, cancel := context.WithTimeout(ctx, timeout)
	err := repo.Refresh(ctxGitOp)
	cancel()
	return err
}
