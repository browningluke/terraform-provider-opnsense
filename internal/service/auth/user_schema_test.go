package auth

import (
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
	opnauth "github.com/browningluke/opnsense-go/pkg/auth"
	"github.com/browningluke/terraform-provider-opnsense/internal/tools"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestConvertUserSchemaToStruct(t *testing.T) {
	model := &userResourceModel{
		Enabled:          types.BoolValue(true),
		Name:             types.StringValue("alice"),
		FullName:         types.StringValue("Alice Example"),
		Email:            types.StringValue("alice@example.com"),
		Comment:          types.StringValue(""),
		Shell:            types.StringValue("/bin/sh"),
		Groups:           tools.StringSliceToSet([]string{"admins"}),
		Privileges:       tools.StringSliceToSet([]string{"page-all"}),
		AuthorizedKeys:   tools.StringSliceToSet([]string{"ssh-ed25519 AAAA alice@example.com"}),
		ScramblePassword: types.BoolValue(true),
		Password:         types.StringNull(),
	}

	user, err := convertUserSchemaToStruct(model, []string{"1999"})
	require.NoError(t, err)
	require.Equal(t, "alice", user.Name)
	require.Equal(t, "0", user.Disabled)
	require.Equal(t, "1", user.ScrambledPassword)
	require.Equal(t, api.SelectedMap("/bin/sh"), user.Shell)
	require.Equal(t, api.SelectedMapList{"1999"}, user.GroupMemberships)
	require.Equal(t, "ssh-ed25519 AAAA alice@example.com", user.AuthorizedKeys)
}

func TestConvertUserSchemaToStruct_PasswordDisablesScramble(t *testing.T) {
	model := &userResourceModel{
		Enabled:          types.BoolValue(true),
		Name:             types.StringValue("alice"),
		FullName:         types.StringValue("Alice Example"),
		Email:            types.StringValue("alice@example.com"),
		Comment:          types.StringValue(""),
		Shell:            types.StringValue("/bin/sh"),
		Groups:           tools.StringSliceToSet([]string{}),
		Privileges:       tools.StringSliceToSet([]string{}),
		AuthorizedKeys:   tools.StringSliceToSet([]string{}),
		ScramblePassword: types.BoolValue(true),
		Password:         types.StringValue("s3cret"),
	}

	user, err := convertUserSchemaToStruct(model, []string{})
	require.NoError(t, err)
	require.Equal(t, "s3cret", user.Password)
	require.Equal(t, "0", user.ScrambledPassword)
}

func TestConvertUserStructToSchema(t *testing.T) {
	user := &opnauth.User{
		Name:             "alice",
		Disabled:         "0",
		Fullname:         "Alice Example",
		Email:            "alice@example.com",
		Shell:            api.SelectedMap("/bin/sh"),
		Privilege:        api.SelectedMapList{"page-all"},
		AuthorizedKeys:   "ssh-ed25519 AAAA alice@example.com\nssh-ed25519 BBBB other",
		GroupMemberships: api.SelectedMapList{"1999"},
	}

	model, err := convertUserStructToSchema(user, []string{"admins"})
	require.NoError(t, err)
	require.True(t, model.Enabled.ValueBool())
	require.Equal(t, "alice", model.Name.ValueString())
	require.Equal(t, "/bin/sh", model.Shell.ValueString())
	require.ElementsMatch(t, []string{"admins"}, tools.SetToStringSlice(model.Groups))
	require.ElementsMatch(t, []string{"page-all"}, tools.SetToStringSlice(model.Privileges))
	require.ElementsMatch(t, []string{
		"ssh-ed25519 AAAA alice@example.com",
		"ssh-ed25519 BBBB other",
	}, tools.SetToStringSlice(model.AuthorizedKeys))
}
