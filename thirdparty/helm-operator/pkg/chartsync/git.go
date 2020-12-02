package chartsync

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/fluxcd/flux/pkg/git"
	"github.com/go-kit/kit/log"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	corev1client "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/cache"

	"github.com/fluxcd/helm-operator/pkg/apis/helm.fluxcd.io/v1"
	lister "github.com/fluxcd/helm-operator/pkg/client/listers/helm.fluxcd.io/v1"
)

// Various (final) errors.
var (
	ErrReleasesForMirror = errors.New("failed to get HelmRelease resources for mirror")
	ErrNoMirror          = errors.New("no existing git mirror found")
	ErrMirrorSync        = errors.New("failed syncing git mirror")
)

// ReleaseQueue is an add-only `workqueue.RateLimitingInterface`
type ReleaseQueue interface {
	AddRateLimited(item interface{})
}

// GitConfig holds the configuration for git operations.
type GitConfig struct {
	GitTimeout      time.Duration
	GitPollInterval time.Duration
	GitDefaultRef   string
}

// GitChartSync syncs `sourceRef`s with their mirrors, and queues
// updates for `v1.HelmRelease`s the sync changes are relevant for.
type GitChartSync struct {
	logger log.Logger
	config GitConfig

	coreV1Client corev1client.CoreV1Interface
	lister       lister.HelmReleaseLister

	mirrors *git.Mirrors

	releaseSourcesMu   sync.RWMutex
	releaseSourcesByID map[string]sourceRef

	releaseQueue ReleaseQueue
}

// sourceRef is used for book keeping, so that we know when a
// signal we receive from a mirror is an actual update for a
// release, and if the source we hold is still the one referred
// to in the `v1.HelmRelease`.
type sourceRef struct {
	mirror string
	remote string
	ref    string
	head   string
}

// forHelmRelease returns true if the given `v1.HelmRelease`s
// `v1.GitChartSource` matches the sourceRef.
func (c sourceRef) forHelmRelease(hr *v1.HelmRelease, defaultGitRef string) bool {
	if hr == nil || hr.Spec.GitChartSource == nil {
		return false
	}

	// reject git source if URL and path are missing
	if hr.Spec.GitURL == "" || hr.Spec.Path == "" {
		return false
	}

	return c.mirror == mirrorName(hr) && c.remote == hr.Spec.GitURL && c.ref == hr.Spec.GitChartSource.RefOrDefault(defaultGitRef)
}

func NewGitChartSync(logger log.Logger,
	coreV1Client corev1client.CoreV1Interface, lister lister.HelmReleaseLister, cfg GitConfig, queue ReleaseQueue) *GitChartSync {

	return &GitChartSync{
		logger:             logger,
		config:             cfg,
		coreV1Client:       coreV1Client,
		lister:             lister,
		mirrors:            git.NewMirrors(),
		releaseSourcesByID: make(map[string]sourceRef),
		releaseQueue:       queue,
	}
}

// Run starts the mirroring of git repositories, and processes mirror
// changes on signal, scheduling a release for a `HelmRelease` resource
// when the update is relevant to the release.
func (c *GitChartSync) Run(stopCh <-chan struct{}, errCh chan error, wg *sync.WaitGroup) {
	c.logger.Log("info", "starting sync of git chart sources")

	wg.Add(1)
	go func() {
		for {
			select {
			case changed := <-c.mirrors.Changes():
				for mirrorName := range changed {
					repo, ok := c.mirrors.Get(mirrorName)

					hrs, err := c.helmReleasesForMirror(mirrorName)
					if err != nil {
						c.logger.Log("error", ErrReleasesForMirror.Error(), "mirror", mirrorName, "err", err)
						continue
					}

					// We received a signal from a no longer existing
					// mirror.
					if !ok {
						if len(hrs) == 0 {
							// If there are no references to it either,
							// just continue with the next mirror...
							continue
						}

						c.logger.Log("warning", ErrNoMirror.Error(), "mirror", mirrorName)
						for _, hr := range hrs {
							c.maybeMirror(mirrorName, hr.Spec.GitChartSource, hr.Namespace)
						}
						// Wait for the signal from the newly requested mirror...
						continue
					}

					// We received a signal from a mirror, but no
					// resource refers to it anymore.
					if ok && len(hrs) == 0 {
						// Garbage collect the mirror.
						c.mirrors.StopOne(mirrorName)
						continue
					}

					c.processChangedMirror(mirrorName, repo, hrs)
				}
			case <-stopCh:
				c.logger.Log("info", "stopping sync of git chart sources")
				c.mirrors.StopAllAndWait()
				wg.Done()
				return
			}
		}
	}()
}

// GetMirrorCopy returns a newly exported copy of the git mirror at the
// recorded HEAD and a string with the HEAD commit hash, or an error.
func (c *GitChartSync) GetMirrorCopy(hr *v1.HelmRelease) (*git.Export, string, error) {
	mirror := mirrorName(hr)
	repo, ok := c.mirrors.Get(mirror)
	if !ok {
		// We did not find a mirror; request one, return, and wait for
		// signal.
		c.maybeMirror(mirror, hr.Spec.GitChartSource, hr.Namespace)
		return nil, "", ChartNotReadyError{ErrNoMirror}
	}

	s, ok, err := c.sync(hr, mirror, repo)
	if err != nil {
		return nil, "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), c.config.GitTimeout)
	defer cancel()
	export, err := repo.Export(ctx, s.head)
	if err != nil {
		return nil, "", ChartUnavailableError{err}
	}

	return export, s.head, nil
}

// Delete cleans up the source reference for the given `v1.HelmRelease`,
// this includes the mirror if there is no reference to it from sources.
// It returns a boolean indicating a successful removal (`true` if so,
// `false` otherwise).
func (c *GitChartSync) Delete(hr *v1.HelmRelease) bool {
	c.releaseSourcesMu.Lock()
	defer c.releaseSourcesMu.Unlock()

	// Attempt to get the source from store, and delete it if found.
	source, ok := c.releaseSourcesByID[hr.ResourceID().String()]
	if ok {
		delete(c.releaseSourcesByID, hr.ResourceID().String())

		if hrs, err := c.helmReleasesForMirror(source.mirror); err == nil && len(hrs) == 0 {
			// The mirror is no longer in use by any source;
			// stop and delete the mirror.
			c.mirrors.StopOne(source.mirror)
		}
	}
	return ok
}

// SyncMirrors instructs all git mirrors to sync from their respective
// upstreams.
func (c *GitChartSync) SyncMirrors() {
	c.logger.Log("info", "starting sync of git mirrors")
	for _, err := range c.mirrors.RefreshAll(c.config.GitTimeout) {
		c.logger.Log("error", ErrMirrorSync.Error(), "err", err)
	}
	c.logger.Log("info", "finished syncing git mirror")
}

// processChangedMirror syncs all given `v1.HelmRelease`s with the
// mirror we received a change signal for and schedules a release,
// but only if the sync indicated the change was relevant.
func (c *GitChartSync) processChangedMirror(mirror string, repo *git.Repo, hrs []*v1.HelmRelease) {
	for _, hr := range hrs {
		if _, ok, _ := c.sync(hr, mirror, repo); ok {
			cacheKey, err := cache.MetaNamespaceKeyFunc(hr.GetObjectMeta())
			if err != nil {
				continue // this should never happen
			}
			// Schedule release sync by adding it to the queue.
			c.releaseQueue.AddRateLimited(cacheKey)
		}
	}
}

// sync synchronizes the record we have for the given `v1.HelmRelease`
// with the given mirror. It always updates the HEAD record in the
// `sourceRef`, but only returns `true` if the update was relevant for
// the release (e.g. a change in git the chart source path, or a new
// record). In case of failure it returns an error.
func (c *GitChartSync) sync(hr *v1.HelmRelease, mirrorName string, repo *git.Repo) (sourceRef, bool, error) {
	source := hr.Spec.GitChartSource
	if source == nil {
		return sourceRef{}, false, nil
	}

	if status, err := repo.Status(); status != git.RepoReady {
		return sourceRef{}, false, ChartNotReadyError{err}
	}

	c.releaseSourcesMu.RLock()
	s, ok := c.releaseSourcesByID[hr.ResourceID().String()]
	c.releaseSourcesMu.RUnlock()

	var changed bool
	if !ok || !s.forHelmRelease(hr, c.config.GitDefaultRef) {
		s = sourceRef{mirror: mirrorName, remote: source.GitURL, ref: source.RefOrDefault(c.config.GitDefaultRef)}
		changed = true
	}

	ctx, cancel := context.WithTimeout(context.Background(), c.config.GitTimeout)
	head, err := repo.Revision(ctx, s.ref)
	cancel()
	if err != nil {
		return sourceRef{}, false, ChartUnavailableError{err}
	}

	if !changed {
		// If the head still equals to what is in our books, there are no changes.
		if s.head == head {
			return s, false, nil
		}

		// Check if the mirror has seen commits in paths we are interested in for
		// this release.
		ctx, cancel = context.WithTimeout(context.Background(), c.config.GitTimeout)
		commits, err := repo.CommitsBetween(ctx, s.head, head, source.Path)
		cancel()
		if err != nil {
			return sourceRef{}, false, ChartUnavailableError{err}
		}
		changed = len(commits) > 0
	}

	// Update the HEAD reference
	s.head = head

	c.releaseSourcesMu.Lock()
	c.releaseSourcesByID[hr.ResourceID().String()] = s
	c.releaseSourcesMu.Unlock()

	return s, changed, nil
}

// maybeMirror requests a new mirror for the given remote. The return value
// indicates whether the repo was already present (`true` if so,
// `false` otherwise).
func (c *GitChartSync) maybeMirror(mirrorName string, source *v1.GitChartSource, namespace string) bool {
	gitURL := source.GitURL
	var err error

	if gitURL, err = c.addAuthForHTTPS(gitURL, source.SecretRef, namespace); err != nil {
		c.logger.Log("error", GitAuthError{err}.Error())
		return false
	}

	ok := c.mirrors.Mirror(
		mirrorName, git.Remote{URL: gitURL}, git.Timeout(c.config.GitTimeout),
		git.PollInterval(c.config.GitPollInterval), git.ReadOnly)
	if !ok {
		c.logger.Log("info", "started mirroring new remote", "remote", source.GitURL, "mirror", mirrorName)
	}
	return ok
}

// helmReleasesForMirror returns a slice of `HelmRelease`s that make
// use of the given mirror.
func (c *GitChartSync) helmReleasesForMirror(mirror string) ([]*v1.HelmRelease, error) {
	hrs, err := c.lister.List(labels.Everything())
	if err != nil {
		return nil, err
	}
	mHrs := make([]*v1.HelmRelease, 0)
	for _, hr := range hrs {
		if m := mirrorName(hr); m == "" || m != mirror {
			continue
		}
		mHrs = append(mHrs, hr.DeepCopy()) // to prevent modifying the (shared) lister store
	}
	return mHrs, nil
}

// mirrorName returns the name of the mirror for the given
// `v1.HelmRelease`.
func mirrorName(hr *v1.HelmRelease) string {
	if hr != nil && hr.Spec.GitChartSource != nil {
		if hr.Spec.GitChartSource.SecretRef == nil {
			return hr.Spec.GitURL
		}
		return fmt.Sprintf("%s/%s/%s", hr.GetNamespace(), hr.Spec.GitChartSource.SecretRef.Name, hr.Spec.GitURL)
	}
	return ""
}

// addAuthForHTTPS will attempt to add basic auth credentials from the
// given secretRef to the given gitURL and return the result, but only
// if the scheme of the URL is HTTPS. In case of a failure it returns
// an error.
func (c *GitChartSync) addAuthForHTTPS(gitURL string, secretRef *corev1.LocalObjectReference, namespace string) (string, error) {
	if secretRef == nil {
		return gitURL, nil
	}

	modifiedURL, err := url.Parse(strings.ToLower(gitURL))
	if err != nil {
		return "", err
	}

	if modifiedURL.Scheme != "https" {
		return gitURL, nil
	}

	username, password, err := c.getAuthFromSecret(secretRef, namespace)
	if err != nil {
		return "", err
	}

	modifiedURL.User = url.UserPassword(username, password)

	return modifiedURL.String(), nil
}

// getAuthFromSecret resolve the given `secretRef` from the given namespace
// using the core v1 secrets client, and return the username and password.
// If this errors, or the secret does not contain the expected keys, an
// error is returned.
func (c *GitChartSync) getAuthFromSecret(secretRef *corev1.LocalObjectReference, ns string) (string, string, error) {
	secretName := secretRef.Name

	secret, err := c.coreV1Client.Secrets(ns).Get(secretName, metav1.GetOptions{})
	if err != nil {
		return "", "", err
	}

	d, ok := secret.Data["username"]
	if !ok {
		return "", "", fmt.Errorf("could not find username key in secret %s/%s", ns, secretName)
	}
	username := string(d)

	d, ok = secret.Data["password"]
	if !ok {
		return "", "", fmt.Errorf("could not find password key in secret %s/%s", ns, secretName)
	}
	password := string(d)

	return username, password, nil
}
