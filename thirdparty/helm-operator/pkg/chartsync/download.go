package chartsync

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	helmfluxv1 "github.com/fluxcd/helm-operator/pkg/apis/helm.fluxcd.io/v1"
	"github.com/fluxcd/helm-operator/pkg/helm"
)

// EnsureChartFetched returns the path to a downloaded chart, fetching
// it first if necessary. It returns the (expected) path to the chart,
// a boolean indicating a fetch, and either an error or nil.
func EnsureChartFetched(client helm.Client, base string, source *helmfluxv1.RepoChartSource) (string, bool, error) {
	repoPath, filename, err := makeChartPath(base, client.Version(), source)
	if err != nil {
		return "", false, ChartUnavailableError{err}
	}
	chartPath := filepath.Join(repoPath, filename)
	stat, err := os.Stat(chartPath)
	switch {
	case os.IsNotExist(err):
		chartPath, err = downloadChart(client, repoPath, source)
		if err != nil {
			return chartPath, false, ChartUnavailableError{err}
		}
		return chartPath, true, nil
	case err != nil:
		return chartPath, false, ChartUnavailableError{err}
	case stat.IsDir():
		return chartPath, false, ChartUnavailableError{errors.New("path to chart exists but is a directory")}
	}
	return chartPath, false, nil
}

// makeChartPath gives the expected filesystem location for a chart,
// without testing whether the file exists or not.
func makeChartPath(base string, clientVersion string, source *helmfluxv1.RepoChartSource) (string, string, error) {
	// We don't need to obscure the location of the charts in the
	// filesystem; but we do need a stable, filesystem-friendly path
	// to them that is based on the URL and the client version.
	repoPath := filepath.Join(base, clientVersion, base64.URLEncoding.EncodeToString([]byte(source.CleanRepoURL())))
	if err := os.MkdirAll(repoPath, 00750); err != nil {
		return "", "", err
	}
	filename := fmt.Sprintf("%s-%s.tgz", source.Name, source.Version)
	return repoPath, filename, nil
}

// downloadChart attempts to pull a chart tarball, given the name,
// version and repo URL in `source`, and the path to write the file
// to in `destFolder`.
func downloadChart(helm helm.Client, destFolder string, source *helmfluxv1.RepoChartSource) (string, error) {
	return helm.PullWithRepoURL(source.RepoURL, source.Name, source.Version, destFolder)
}
