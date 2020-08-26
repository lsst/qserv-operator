package qserv

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
				log.Error(err, "toto")
			}
			for _, f := range files {
				fmt.Println(f.Name())
			}

			got, error := applyTemplate(tt.args.str, tt.args.templateData)
			fmt.Println(error)
			assert.Equal(t, tt.want, got)
		})
	}
}
