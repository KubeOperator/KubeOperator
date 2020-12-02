package release

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"

	"github.com/fluxcd/helm-operator/pkg/helm"
	"github.com/ghodss/yaml"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"

	"github.com/fluxcd/helm-operator/pkg/apis/helm.fluxcd.io/v1"
)

// values attempts to compose the final values for the given
// `HelmRelease`. It returns the values as bytes and a checksum,
// or an error in case anything went wrong.
func composeValues(coreV1Client corev1.CoreV1Interface, hr *v1.HelmRelease, chartPath string) (helm.Values, error) {
	result := helm.Values{}

	for _, v := range hr.GetValuesFromSources() {
		var valueFile helm.Values
		ns := hr.Namespace

		switch {
		case v.ConfigMapKeyRef != nil:
			cm := v.ConfigMapKeyRef
			name := cm.Name
			if cm.Namespace != "" {
				ns = cm.Namespace
			}
			key := cm.Key
			if key == "" {
				key = "values.yaml"
			}
			optional := cm.Optional != nil && *cm.Optional
			configMap, err := coreV1Client.ConfigMaps(ns).Get(name, metav1.GetOptions{})
			if err != nil {
				if errors.IsNotFound(err) && optional {
					continue
				}
				return result, err
			}
			d, ok := configMap.Data[key]
			if !ok {
				if optional {
					continue
				}
				return result, fmt.Errorf("could not find key %v in ConfigMap %s/%s", key, ns, name)
			}
			if err := yaml.Unmarshal([]byte(d), &valueFile); err != nil {
				if optional {
					continue
				}
				return result, fmt.Errorf("unable to yaml.Unmarshal %v from %s in ConfigMap %s/%s", d, key, ns, name)
			}
		case v.SecretKeyRef != nil:
			s := v.SecretKeyRef
			name := s.Name
			if s.Namespace != "" {
				ns = s.Namespace
			}
			key := s.Key
			if key == "" {
				key = "values.yaml"
			}
			optional := s.Optional != nil && *s.Optional
			secret, err := coreV1Client.Secrets(ns).Get(name, metav1.GetOptions{})
			if err != nil {
				if errors.IsNotFound(err) && optional {
					continue
				}
				return result, err
			}
			d, ok := secret.Data[key]
			if !ok {
				if optional {
					continue
				}
				return result, fmt.Errorf("could not find key %s in Secret %s/%s", key, ns, name)
			}
			if err := yaml.Unmarshal(d, &valueFile); err != nil {
				return result, fmt.Errorf("unable to yaml.Unmarshal %v from %s in Secret %s/%s", d, key, ns, name)
			}
		case v.ExternalSourceRef != nil:
			es := v.ExternalSourceRef
			u := es.URL
			optional := es.Optional != nil && *es.Optional
			b, err := readURL(u)
			if err != nil {
				if optional {
					continue
				}
				return result, fmt.Errorf("unable to read value file from URL %s", u)
			}
			if err := yaml.Unmarshal(b, &valueFile); err != nil {
				if optional {
					continue
				}
				return result, fmt.Errorf("unable to yaml.Unmarshal %v from URL %s", b, u)
			}
		case v.ChartFileRef != nil:
			cf := v.ChartFileRef
			filePath := cf.Path
			optional := cf.Optional != nil && *cf.Optional
			f, err := readLocalChartFile(filepath.Join(chartPath, filePath))
			if err != nil {
				if optional {
					continue
				}
				return result, fmt.Errorf("unable to read value file from path %s", filePath)
			}
			if err := yaml.Unmarshal(f, &valueFile); err != nil {
				if optional {
					continue
				}
				return result, fmt.Errorf("unable to yaml.Unmarshal %v from path %s", f, filePath)
			}
		}
		result = mergeValues(result, valueFile)
	}

	result = mergeValues(result, hr.Spec.Values)
	return result, nil
}

// readURL attempts to read a file from an HTTP(S) URL.
func readURL(URL string) ([]byte, error) {
	u, err := url.Parse(URL)
	if err != nil {
		return []byte{}, err
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return []byte{}, fmt.Errorf("URL scheme should be HTTP(S), got '%s'", u.Scheme)
	}
	resp, err := http.Get(u.String())
	if err != nil {
		return []byte{}, err
	}
	switch resp.StatusCode {
	case http.StatusOK:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return []byte{}, err
		}
		return body, nil
	default:
		return []byte{}, fmt.Errorf("failed to retrieve file from URL, status '%s (%d)'", resp.Status, resp.StatusCode)
	}
}

// readLocalChartFile attempts to read a file from the chart path.
func readLocalChartFile(filePath string) ([]byte, error) {
	f, err := ioutil.ReadFile(filePath)
	if err != nil {
		return []byte{}, err
	}
	return f, nil
}

// mergeValues merges source and destination map, preferring values
// from the source values. This is slightly adapted from:
// https://github.com/helm/helm/blob/2332b480c9cb70a0d8a85247992d6155fbe82416/cmd/helm/install.go#L359
func mergeValues(dest, src map[string]interface{}) map[string]interface{} {
	for k, v := range src {
		// If the key doesn't exist already, then just set the key to that value
		if _, exists := dest[k]; !exists {
			dest[k] = v
			continue
		}
		nextMap, ok := v.(map[string]interface{})
		// If it isn't another map, overwrite the value
		if !ok {
			dest[k] = v
			continue
		}
		// Edge case: If the key exists in the destination, but isn't a map
		destMap, isMap := dest[k].(map[string]interface{})
		// If the source map has a map for this key, prefer it
		if !isMap {
			dest[k] = v
			continue
		}
		// If we got to this point, it is a map in both, so merge them
		dest[k] = mergeValues(destMap, nextMap)
	}
	return dest
}
