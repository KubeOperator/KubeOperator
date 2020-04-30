package images

import (
	"ko3-gin/pkg/cluster/adm/constants"
	"reflect"
	"sort"
)

type Components struct {
	ETCD               Image
	CoreDNS            Image
	Pause              Image
	NvidiaDevicePlugin Image
	Keepalived         Image
}

func (c Components) Get(name string) *Image {
	v := reflect.ValueOf(c)
	for i := 0; i < v.NumField(); i++ {
		v, _ := v.Field(i).Interface().(Image)
		if v.Name == name {
			return &v
		}
	}
	return nil
}

var components = Components{
	ETCD:    Image{Name: "etcd", Tag: constants.EtcdVersion},
	CoreDNS: Image{Name: "coredns", Tag: constants.CoreDNSVersion},
	Pause:   Image{Name: "pause", Tag: constants.PauseVersion},
}

func List() []string {
	var items []string

	//for _, version := range K8sVersionsWithV {
	//	for _, name := range []string{"kube-apiserver", "kube-controller-manager", "kube-scheduler", "kube-proxy"} {
	//		items = append(items, Image{Name: name, Tag: version}.BaseName())
	//	}
	//}

	v := reflect.ValueOf(components)
	for i := 0; i < v.NumField(); i++ {
		v, _ := v.Field(i).Interface().(Image)
		items = append(items, v.BaseName())
	}
	sort.Strings(items)
	return items
}

func Get() Components {
	return components
}
