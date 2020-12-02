/*

This package is for maintaining the link between `HelmRelease`
resources and the Helm releases to which they
correspond. Specifically,

 1. updating the `HelmRelease` status based on the progress of
   syncing, and the state of the associated Helm release; and,

 2. attributing each resource in a Helm release (under our control) to
 the associated `HelmRelease`.

*/
package status

import (
	"time"

	"github.com/go-kit/kit/log"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	kube "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/retry"

	v1 "github.com/fluxcd/helm-operator/pkg/apis/helm.fluxcd.io/v1"
	ifclientset "github.com/fluxcd/helm-operator/pkg/client/clientset/versioned"
	v1client "github.com/fluxcd/helm-operator/pkg/client/clientset/versioned/typed/helm.fluxcd.io/v1"
	iflister "github.com/fluxcd/helm-operator/pkg/client/listers/helm.fluxcd.io/v1"
	"github.com/fluxcd/helm-operator/pkg/helm"
)

type Updater struct {
	hrClient           ifclientset.Interface
	hrLister           iflister.HelmReleaseLister
	kube               kube.Interface
	helmClients        *helm.Clients
	defaultHelmVersion string
}

func New(hrClient ifclientset.Interface, hrLister iflister.HelmReleaseLister, helmClients *helm.Clients, defaultHelmVersion string) *Updater {
	return &Updater{
		hrClient:           hrClient,
		hrLister:           hrLister,
		helmClients:        helmClients,
		defaultHelmVersion: defaultHelmVersion,
	}
}

func (u *Updater) Loop(stop <-chan struct{}, interval time.Duration, logger log.Logger) {
	ticker := time.NewTicker(interval)
	var logErr error

bail:
	for {
		select {
		case <-stop:
			break bail
		case <-ticker.C:
		}
		list, err := u.hrLister.List(labels.Everything())
		if err != nil {
			logErr = err
			break bail
		}
		for _, hr := range list {
			nsHrClient := u.hrClient.HelmV1().HelmReleases(hr.Namespace)
			releaseName := hr.GetReleaseName()
			c, ok := u.helmClients.Load(hr.GetHelmVersion(u.defaultHelmVersion))
			// If we are unable to get the client, we do not care why
			if !ok {
				continue
			}
			rel, _ := c.Get(hr.GetReleaseName(), helm.GetOptions{Namespace: hr.GetTargetNamespace()})
			// If we are unable to get the status, we do not care why
			if rel == nil {
				continue
			}
			if err := SetReleaseStatus(nsHrClient, hr, releaseName, rel.Info.Status.String()); err != nil {
				logger.Log("namespace", hr.Namespace, "resource", hr.Name, "err", err)
				continue
			}
		}
	}

	ticker.Stop()
	logger.Log("loop", "stopping", "err", logErr)
}

// SetReleaseStatus updates the status of the HelmRelease to the given
// release name and/or release status.
func SetReleaseStatus(client v1client.HelmReleaseInterface, hr *v1.HelmRelease,
	releaseName, releaseStatus string) error {

	firstTry := true
	err := retry.RetryOnConflict(retry.DefaultBackoff, func() (err error) {
		if !firstTry {
			var getErr error
			hr, getErr = client.Get(hr.Name, metav1.GetOptions{})
			if getErr != nil {
				return getErr
			}
		}

		if hr.Status.ReleaseName == releaseName && hr.Status.ReleaseStatus == releaseStatus {
			return
		}

		cHr := hr.DeepCopy()
		cHr.Status.ReleaseName = releaseName
		cHr.Status.ReleaseStatus = releaseStatus

		_, err = client.UpdateStatus(cHr)
		firstTry = false
		return
	})
	return err
}

// SetReleaseRevision updates the revision in the status of the HelmRelease
// to the given revision, and sets the current revision as the previous one.
func SetReleaseRevision(client v1client.HelmReleaseInterface, hr *v1.HelmRelease, revision string) error {

	firstTry := true
	err := retry.RetryOnConflict(retry.DefaultBackoff, func() (err error) {
		if !firstTry {
			var getErr error
			hr, getErr = client.Get(hr.Name, metav1.GetOptions{})
			if getErr != nil {
				return getErr
			}
		}

		if revision == "" || hr.Status.Revision == revision {
			return
		}

		cHr := hr.DeepCopy()
		cHr.Status.Revision = revision

		_, err = client.UpdateStatus(cHr)
		firstTry = false
		return
	})
	return err
}

// SetObservedGeneration updates the observed generation status of the
// HelmRelease to the given generation.
func SetObservedGeneration(client v1client.HelmReleaseInterface, hr *v1.HelmRelease, generation int64) error {
	firstTry := true
	err := retry.RetryOnConflict(retry.DefaultBackoff, func() (err error) {
		if !firstTry {
			var getErr error
			hr, getErr = client.Get(hr.Name, metav1.GetOptions{})
			if getErr != nil {
				return getErr
			}
		}

		if hr.Status.ObservedGeneration >= generation {
			return
		}

		cHr := hr.DeepCopy()
		cHr.Status.ObservedGeneration = generation

		_, err = client.UpdateStatus(cHr)
		firstTry = false
		return
	})
	return err
}

// HasSynced returns if the HelmRelease has been processed by the
// controller.
func HasSynced(hr v1.HelmRelease) bool {
	return hr.Status.ObservedGeneration >= hr.Generation
}

// HasRolledBack returns if the current generation of the HelmRelease
// has been rolled back.
func HasRolledBack(hr v1.HelmRelease) bool {
	if !HasSynced(hr) {
		return false
	}

	rolledBack := GetCondition(hr.Status, v1.HelmReleaseRolledBack)
	if rolledBack == nil {
		return false
	}

	return rolledBack.Status == corev1.ConditionTrue
}

// ShouldRetryUpgrade returns if the upgrade of a rolled back release should
// be retried.
func ShouldRetryUpgrade(hr v1.HelmRelease) bool {
	if !hr.Spec.Rollback.Retry {
		return false
	}
	return hr.Spec.Rollback.GetMaxRetries() == 0 || hr.Status.RollbackCount <= hr.Spec.Rollback.GetMaxRetries()
}
