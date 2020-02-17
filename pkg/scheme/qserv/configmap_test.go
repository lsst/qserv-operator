package qserv

import (
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
		want bool
	}{
		{
			"String in slice returns true",
			args{
				"info",
				templateData{QstatusMysqldHost: "CZAR_TEST"},
			},
			true,
		},
		{
			"String not in slice returns false",
			args{
				"error",
				templateData{QstatusMysqldHost: "CZAR_TEST2"},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := applyTemplate(tt.args.str, tt.args.templateData)
			assert.Equal(t, tt.want, got)
		})
	}
}
