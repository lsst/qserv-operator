package objects

import (
	"fmt"
	"io/ioutil"
	"testing"

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
				templateData{QstatusMysqldHost: "Value"},
			},
			"Test field: Value\n",
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

			got, error := applyTemplate(tt.args.str, &tt.args.templateData)
			fmt.Println(error)
			assert.Equal(t, tt.want, got)
		})
	}
}
