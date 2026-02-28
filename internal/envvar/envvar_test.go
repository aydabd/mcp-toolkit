package envvar

import (
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func resetViper() {
	viper.Reset()
	BindAll()
}

func TestVarKey(t *testing.T) {
	tests := []struct {
		v    Var
		want string
	}{
		{Var{Name: "MCP_ATLASSIAN_URL"}, "atlassian.url"},
		{Var{Name: "MCP_VAULT_ADDR"}, "vault.addr"},
		{Var{Name: "MCP_VERBOSE", Key: "verbose"}, "verbose"},
		{Var{Name: "MCP_CONFIG_DIR"}, "config.dir"},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.want, tt.v.key(), "Name: %s", tt.v.Name)
	}
}

func TestVarGet(t *testing.T) {
	resetViper()

	v := Var{Name: "MCP_TEST_GET", Default: "default_value"}

	// Test default value
	assert.Equal(t, "default_value", v.Get())

	// Test with viper set
	viper.Set(v.key(), "custom_value")
	assert.Equal(t, "custom_value", v.Get())
}

func TestVarGetFromEnv(t *testing.T) {
	resetViper()

	v := Var{Name: "MCP_TEST_ENV", Default: "default"}

	// Set via env var
	os.Setenv("MCP_TEST_ENV", "from_env")
	defer os.Unsetenv("MCP_TEST_ENV")

	// Viper should pick it up
	assert.Equal(t, "from_env", v.Get())
}

func TestVarGetOrError(t *testing.T) {
	resetViper()

	// Required var not set
	v := Var{Name: "MCP_TEST_REQUIRED", Required: true}
	_, err := v.GetOrError()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "required config")

	// Required var set via viper
	viper.Set(v.key(), "value")
	val, err := v.GetOrError()
	require.NoError(t, err)
	assert.Equal(t, "value", val)

	// Optional var not set
	resetViper()
	optVar := Var{Name: "MCP_TEST_OPTIONAL", Required: false}
	val, err = optVar.GetOrError()
	require.NoError(t, err)
	assert.Equal(t, "", val)
}

func TestVarIsSet(t *testing.T) {
	resetViper()

	v := Var{Name: "MCP_TEST_ISSET"}
	assert.False(t, v.IsSet())

	viper.Set(v.key(), "value")
	assert.True(t, v.IsSet())

	// With default
	resetViper()
	vWithDefault := Var{Name: "MCP_TEST_DEFAULT", Default: "def"}
	assert.True(t, vWithDefault.IsSet())
}

func TestVarSet(t *testing.T) {
	resetViper()

	v := Var{Name: "MCP_TEST_SET"}
	v.Set("new_value")
	assert.Equal(t, "new_value", viper.GetString(v.key()))
}

func TestBindAll(t *testing.T) {
	resetViper()

	// Check defaults are set
	assert.Equal(t, "false", viper.GetString("atlassian.read.only"))
	assert.Equal(t, "/usr/local/bin/vault-mcp-server", viper.GetString("vault.binary.path"))
	assert.Equal(t, "docker", viper.GetString("container.runtime"))
}

func TestAll(t *testing.T) {
	vars := All()
	assert.NotEmpty(t, vars)

	names := make(map[string]bool)
	for _, v := range vars {
		names[v.Name] = true
	}
	assert.True(t, names["MCP_ATLASSIAN_URL"])
	assert.True(t, names["MCP_VAULT_ADDR"])
	assert.True(t, names["MCP_VERBOSE"])
}

func TestByCategory(t *testing.T) {
	categories := ByCategory()

	assert.Contains(t, categories, "Atlassian")
	assert.Contains(t, categories, "Vault")
	assert.Contains(t, categories, "Supabase")
	assert.Contains(t, categories, "General")

	atlassianVars := categories["Atlassian"]
	assert.NotEmpty(t, atlassianVars)
}

func TestFormatHelp(t *testing.T) {
	help := FormatHelp()

	assert.Contains(t, help, "Environment Variables:")
	assert.Contains(t, help, "Atlassian:")
	assert.Contains(t, help, "MCP_ATLASSIAN_URL")
	assert.Contains(t, help, "(required)")
}

func TestValidate(t *testing.T) {
	resetViper()

	// Test with missing required var
	requiredVar := Var{Name: "MCP_TEST_VALIDATE", Required: true}
	err := Validate(requiredVar)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing required config")

	// Test with required var set
	viper.Set(requiredVar.key(), "value")
	err = Validate(requiredVar)
	assert.NoError(t, err)

	// Test with optional var
	resetViper()
	optionalVar := Var{Name: "MCP_TEST_OPT", Required: false}
	err = Validate(optionalVar)
	assert.NoError(t, err)

	// Test with default value
	defaultVar := Var{Name: "MCP_TEST_DEF", Required: true, Default: "default"}
	err = Validate(defaultVar)
	assert.NoError(t, err)
}

func TestAtlassianVars(t *testing.T) {
	assert.Equal(t, "MCP_ATLASSIAN_URL", AtlassianURL.Name)
	assert.True(t, AtlassianURL.Required)

	assert.Equal(t, "MCP_ATLASSIAN_API_TOKEN", AtlassianAPIToken.Name)
	assert.True(t, AtlassianAPIToken.Secret)
}

func TestVaultVars(t *testing.T) {
	assert.Equal(t, "MCP_VAULT_ADDR", VaultAddr.Name)
	assert.Equal(t, "/usr/local/bin/vault-mcp-server", VaultBinaryPath.Default)
}

func TestGeneralVars(t *testing.T) {
	assert.Equal(t, "MCP_VERBOSE", Verbose.Name)
	assert.Equal(t, "verbose", Verbose.Key)
	assert.Equal(t, "false", Verbose.Default)
}
