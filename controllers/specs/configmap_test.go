package specs

import (
	"fmt"
	"os"
	"testing"

	"github.com/lsst/qserv-operator/api/v1beta1"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/api/resource"
)

func TestGetReplicationWorkerThread(t *testing.T) {

	tests := []struct {
		name string
		str  string
		want uint
	}{
		{
			"millis-1",
			"1000m",
			2,
		},
		{
			"millis-2",
			"50m",
			2,
		},
		{
			"dec-1",
			"1.5",
			4,
		},
	}

	var nulQuantity resource.Quantity
	got := getReplicationWorkerThread(&nulQuantity)
	assert.Equal(t, uint(16), got)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			quantity := resource.MustParse(tt.str)
			got := getReplicationWorkerThread(&quantity)
			t.Logf("ncore %v", got)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestApplyTemplate(t *testing.T) {
	type args struct {
		str          string
		templateData templateData
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"Existing template",
			args{
				"example.tpl",
				templateData{QstatusMysqldHost: "Value",
					ResultsProtocol: v1beta1.ResultsProtocolTypeHTTP,
				},
			},
			"Test field: Value\nTest field2: HTTP\n",
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

			got, error := applyTemplate(tt.args.str, &tt.args.templateData)
			fmt.Println(error)
			assert.Equal(t, tt.want, got)
		})
	}
}
