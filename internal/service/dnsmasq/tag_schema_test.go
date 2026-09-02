package dnsmasq

import (
	"testing"

	"github.com/browningluke/opnsense-go/pkg/dnsmasq"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestConvertDnsmasqTagToStruct(t *testing.T) {
	tests := []struct {
		name     string
		input    *tagResourceModel
		expected *dnsmasq.Tag
	}{
		{
			name: "full tag test",
			input: &tagResourceModel{
				Tag: types.StringValue("myTagName"),
				Id:  types.StringValue("f1c0fd09-0652-4ee0-9e5c-9e2e55dae883"),
			},
			expected: &dnsmasq.Tag{
				Tag: "myTagName",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := convertTagSchemaToStruct(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected.Tag, result.Tag)
		})
	}
}

func TestConvertTagStructToSchema(t *testing.T) {
	tests := []struct {
		name     string
		input    *dnsmasq.Tag
		expected *tagResourceModel
	}{
		{
			name: "basic PSK conversion",
			input: &dnsmasq.Tag{
				Tag: "testTagName",
			},
			expected: &tagResourceModel{
				Tag: types.StringValue("testTagName"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := convertTagStructToSchema(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected.Tag, result.Tag)
		})
	}
}
