package i18n

import (
	"github.com/KubeOperator/KubeOperator/pkg/i18n/localizations"
)

func Tr(key string, replace map[string]string) string {
	l := localizations.New("zh", "en")
	key = "messages." + key
	if replace != nil {
		rep := localizations.Replacements{}
		for k, v := range replace {
			rep[k] = v
		}
		return l.GetWithLocale("zh", key, &rep)
	} else {
		return l.GetWithLocale("zh", key)
	}

}
