package status

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/retry"

	"github.com/fluxcd/helm-operator/pkg/apis/helm.fluxcd.io/v1"
	v1client "github.com/fluxcd/helm-operator/pkg/client/clientset/versioned/typed/helm.fluxcd.io/v1"
)

// NewCondition creates a new HelmReleaseCondition.
func NewCondition(conditionType v1.HelmReleaseConditionType, status corev1.ConditionStatus,
	reason, message string) v1.HelmReleaseCondition {

	return v1.HelmReleaseCondition{
		Type:               conditionType,
		Status:             status,
		LastUpdateTime:     metav1.Now(),
		LastTransitionTime: metav1.Now(),
		Reason:             reason,
		Message:            message,
	}
}

// GetCondition returns the condition with the given type.
func GetCondition(status v1.HelmReleaseStatus, conditionType v1.HelmReleaseConditionType) *v1.HelmReleaseCondition {

	for i := range status.Conditions {
		c := status.Conditions[i]
		if c.Type == conditionType {
			return &c
		}
	}
	return nil
}

// SetCondition updates the HelmRelease to include the given condition.
func SetCondition(client v1client.HelmReleaseInterface, hr *v1.HelmRelease, condition v1.HelmReleaseCondition) error {

	firstTry := true
	err := retry.RetryOnConflict(retry.DefaultBackoff, func() (err error) {
		if !firstTry {
			var getErr error
			hr, getErr = client.Get(hr.Name, metav1.GetOptions{})
			if getErr != nil {
				return getErr
			}
		}

		cHr := hr.DeepCopy()
		currCondition := GetCondition(cHr.Status, condition.Type)
		if currCondition != nil && currCondition.Status == condition.Status {
			condition.LastTransitionTime = currCondition.LastTransitionTime
		}
		newConditions := filterOutCondition(cHr.Status.Conditions, condition.Type)
		cHr.Status.Conditions = append(newConditions, condition)

		switch {
		case condition.Type == v1.HelmReleaseReleased && condition.Status == corev1.ConditionTrue:
			cHr.Status.RollbackCount = 0
		case condition.Type == v1.HelmReleaseRolledBack && condition.Status == corev1.ConditionTrue:
			cHr.Status.RollbackCount = hr.Status.RollbackCount + 1
		}

		_, err = client.UpdateStatus(cHr)
		firstTry = false
		return
	})
	return err
}

// UnsetCondition updates the HelmRelease to exclude the given condition.
func UnsetCondition(client v1client.HelmReleaseInterface,
	hr *v1.HelmRelease, conditionType v1.HelmReleaseConditionType) error {

	firstTry := true
	err := retry.RetryOnConflict(retry.DefaultBackoff, func() (err error) {
		if !firstTry {
			var getErr error
			hr, getErr = client.Get(hr.Name, metav1.GetOptions{})
			if getErr != nil {
				return getErr
			}
		}

		if GetCondition(hr.Status, conditionType) == nil {
			return
		}

		cHr := hr.DeepCopy()
		cHr.Status.Conditions = filterOutCondition(cHr.Status.Conditions, conditionType)

		switch {
		case conditionType == v1.HelmReleaseRolledBack:
			cHr.Status.RollbackCount = 0
		}

		_, err = client.UpdateStatus(cHr)
		firstTry = false
		return
	})
	return err
}

// filterOutCondition returns a new slice of conditions without the
// conditions of the given type.
func filterOutCondition(conditions []v1.HelmReleaseCondition,
	conditionType v1.HelmReleaseConditionType) []v1.HelmReleaseCondition {

	var newConditions []v1.HelmReleaseCondition
	for _, c := range conditions {
		if c.Type == conditionType {
			continue
		}
		newConditions = append(newConditions, c)
	}
	return newConditions
}
