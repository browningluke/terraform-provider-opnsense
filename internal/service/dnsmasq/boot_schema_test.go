package dnsmasq

import (
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/dnsmasq"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestConvertDnsmasqBootToStruct(t *testing.T) {
	tests := []struct {
		name     string
		input    *bootResourceModel
		expected *dnsmasq.Boot
	}{
		{
			name: "basic boot conversion",
			input: &bootResourceModel{
				Interface:     types.StringValue("lan"),
				FileName:      types.StringValue("pxelinux.0"),
				ServerName:    types.StringValue("bootserver"),
				ServerAddress: types.StringValue("192.168.1.10"),
				Description:   types.StringValue("Test boot configuration"),
			},
			expected: &dnsmasq.Boot{
				Interface:     api.SelectedMap("lan"),
				Filename:      "pxelinux.0",
				Servername:    "bootserver",
				ServerAddress: "192.168.1.10",
				Description:   "Test boot configuration",
			},
		},
		{
			name: "boot conversion with tags",
			input: &bootResourceModel{
				Interface: types.StringValue("lan"),
				Tag: types.SetValueMust(types.StringType, []attr.Value{
					types.StringValue("tag1"),
					types.StringValue("tag2"),
				}),
				FileName:      types.StringValue("grubx64.efi"),
				ServerName:    types.StringValue("bootserver"),
				ServerAddress: types.StringValue("192.168.1.20"),
				Description:   types.StringValue("Boot with tags"),
			},
			expected: &dnsmasq.Boot{
				Interface:     api.SelectedMap("lan"),
				Tag:           api.SelectedMapList([]string{"tag1", "tag2"}),
				Filename:      "grubx64.efi",
				Servername:    "bootserver",
				ServerAddress: "192.168.1.20",
				Description:   "Boot with tags",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := convertBootSchemaToStruct(tt.input)
			assert.NoError(t, err)

			assert.Equal(t, tt.expected.Filename, result.Filename)
			assert.Equal(t, tt.expected.Servername, result.Servername)
			assert.Equal(t, tt.expected.ServerAddress, result.ServerAddress)
			assert.Equal(t, tt.expected.Description, result.Description)
			assert.Equal(t, tt.expected.Interface.String(), result.Interface.String())

			if tt.expected.Tag == nil {
				assert.Nil(t, result.Tag)
			} else {
				assert.Equal(t, tt.expected.Tag, result.Tag)
			}
		})
	}
}

func TestConvertDnsmasqBootStructToSchema(t *testing.T) {
	tests := []struct {
		name     string
		input    *dnsmasq.Boot
		expected *bootResourceModel
	}{
		{
			name: "basic boot conversion",
			input: &dnsmasq.Boot{
				Interface:     api.SelectedMap("lan"),
				Filename:      "pxelinux.0",
				Servername:    "bootserver",
				ServerAddress: "192.168.1.10",
				Description:   "Test boot configuration",
			},
			expected: &bootResourceModel{
				Interface:     types.StringValue("lan"),
				Tag:           types.SetNull(types.StringType),
				FileName:      types.StringValue("pxelinux.0"),
				ServerName:    types.StringValue("bootserver"),
				ServerAddress: types.StringValue("192.168.1.10"),
				Description:   types.StringValue("Test boot configuration"),
			},
		},
		{
			name: "boot conversion with tags",
			input: &dnsmasq.Boot{
				Interface:     api.SelectedMap("opt1"),
				Tag:           api.SelectedMapList([]string{"tag1", "tag2"}),
				Filename:      "grubx64.efi",
				Servername:    "bootserver2",
				ServerAddress: "192.168.1.20",
				Description:   "Boot with tags",
			},
			expected: &bootResourceModel{
				Interface:     types.StringValue("opt1"),
				FileName:      types.StringValue("grubx64.efi"),
				ServerName:    types.StringValue("bootserver2"),
				ServerAddress: types.StringValue("192.168.1.20"),
				Description:   types.StringValue("Boot with tags"),
			},
		},
		{
			name: "empty interface conversion",
			input: &dnsmasq.Boot{
				Filename:      "pxelinux.0",
				Servername:    "bootserver",
				ServerAddress: "192.168.1.10",
				Description:   "Boot without interface",
			},
			expected: &bootResourceModel{
				Interface:     types.StringNull(),
				FileName:      types.StringValue("pxelinux.0"),
				ServerName:    types.StringValue("bootserver"),
				ServerAddress: types.StringValue("192.168.1.10"),
				Description:   types.StringValue("Boot without interface"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := convertBootStructToSchema(tt.input)
			assert.NoError(t, err)

			assert.Equal(t, tt.expected.Interface, result.Interface)
			assert.Equal(t, tt.expected.FileName, result.FileName)
			assert.Equal(t, tt.expected.ServerName, result.ServerName)
			assert.Equal(t, tt.expected.ServerAddress, result.ServerAddress)
			assert.Equal(t, tt.expected.Description, result.Description)

			if tt.input.Tag == nil {
				assert.True(t, result.Tag.IsNull())
			} else {
				assert.Equal(t, types.SetValueMust(types.StringType, []attr.Value{
					types.StringValue("tag1"),
					types.StringValue("tag2"),
				}), result.Tag)
			}
		})
	}
}
