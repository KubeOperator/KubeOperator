package install

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"text/template"

	"github.com/shurcooL/httpfs/vfsutil"
)

type TemplateParameters struct {
	Namespace               string
	TillerNamespace         string
	SSHSecretName           string
	EnableTillerTLS         bool
	TillerTLSCACertContent  string
	TillerTLSCertSecretName string
	HelmVersions            string
	AdditionalArgs          []string
}

func FillInTemplates(params TemplateParameters) (map[string][]byte, error) {
	result := map[string][]byte{}
	err := vfsutil.WalkFiles(templates, "/", func(path string, info os.FileInfo, rs io.ReadSeeker, err error) error {
		if err != nil {
			return fmt.Errorf("cannot walk embedded files: %s", err)
		}
		if info.IsDir() {
			return nil
		}
		manifestTemplateBytes, err := ioutil.ReadAll(rs)
		if err != nil {
			return fmt.Errorf("cannot read embedded file %q: %s", info.Name(), err)
		}
		indent := func(n int, s string) string {
			lines := strings.Split(s, "\n")
			for i, l := range lines {
				lines[i] = strings.Repeat(" ", n) + l
			}
			return strings.Join(lines, "\n")
		}
		manifestTemplate, err := template.New(info.Name()).
			Funcs(template.FuncMap{"indent": indent}).
			Parse(string(manifestTemplateBytes))
		if err != nil {
			return fmt.Errorf("cannot parse embedded file %q: %s", info.Name(), err)
		}
		out := bytes.NewBuffer(nil)
		if err := manifestTemplate.Execute(out, params); err != nil {
			return fmt.Errorf("cannot execute template for embedded file %q: %s", info.Name(), err)
		}
		if len(out.Bytes()) <= 1 { // empty file
			return nil
		}
		result[strings.TrimSuffix(info.Name(), ".tmpl")] = out.Bytes()
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("internal error filling embedded installation templates: %s", err)
	}
	return result, nil
}
