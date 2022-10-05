package util

import (
	"fmt"
	"os"
	"testing"

	"github.com/lsst/qserv-operator/controllers/constants"
	"github.com/stretchr/testify/assert"
)

func TestGetComponentLabels(t *testing.T) {
	type args struct {
		crName string
	}

	crName := "qserv"

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

			dir, err := os.Open(".")
			if err != nil {
				fmt.Println(err)
				return
			}
			files, err := dir.Readdir(0)
			if err != nil {
				fmt.Println(err)
				return
			}
			for _, f := range files {
				fmt.Println(f.Name())
			}

			got := GetComponentLabels(constants.ReplCtl, crName)

			fmt.Println(got)
			assert.Equal(t, tt.want, got)
		})
	}
}
