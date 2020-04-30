package cluster

import (
	"fmt"
	"ko3-gin/pkg/cluster/adm"
	"testing"
	"time"
)

func TestNewCluster(t *testing.T) {
	cluster := &Cluster{
		Name: "test",
	}
	a, err := adm.NewClusterAdm()
	if err != nil {
		t.Error(err)
	}
	for {
		start := time.Now()
		resp, err := a.OnInitialize(*cluster)
		if err != nil {
			t.Error(err)
		}
		cluster = &resp
		condition := cluster.Status.Conditions[len(cluster.Status.Conditions)-1]
		switch condition.Status {
		case ConditionFalse:
			fmt.Printf("OnInitialize.%s [Failed] [%fs] reason: %s message: %s",
				condition.Type, time.Since(start).Seconds(),
				condition.Reason, condition.Message, )
		case ConditionUnknown:
			condition = resp.Status.Conditions[len(resp.Status.Conditions)-2]
			fmt.Printf("OnInitialize.%s [Success] [%fs]", condition.Type, time.Since(start).Seconds())
		case ConditionTrue:
			fmt.Print("all done")
		}
		time.Sleep(5 * time.Second)
	}

}
