package util

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/lsst/qserv-operator/controllers/constants"
	"github.com/stretchr/testify/assert"
)

func TestGetComponentLabels(t *testing.T) {
	type args struct {
		cr_name string
	}

	cr_name := "qserv"

	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			"Replication controller",
			args{
				"qserv",
			},
			map[string]string{"app": "qserv", "app.kubernetes.io/managed-by": "qserv-operator", "component": "repl-ctl", "instance": "qserv"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			files, err := ioutil.ReadDir(".")
			if err != nil {
				t.Errorf("%s", err)
			}
			for _, f := range files {
				fmt.Println(f.Name())
			}

			got := GetComponentLabels(constants.ReplCtl, cr_name)

			fmt.Println(got)
			assert.Equal(t, tt.want, got)
		})
	}
}
