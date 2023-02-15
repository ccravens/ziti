package cmd

import (
	"fmt"
	cmdhelper "github.com/openziti/ziti/ziti/cmd/helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
	"testing"
	"time"
)

/* BEGIN Controller config template structure */

type ControllerConfig struct {
	V            string       `yaml:"v"`
	Db           string       `yaml:"db"`
	Identity     Identity     `yaml:"identity"`
	Ctrl         Ctrl         `yaml:"ctrl"`
	Mgmt         Mgmt         `yaml:"mgmt"`
	HealthChecks HealthChecks `yaml:"healthChecks"`
	Edge         Edge         `yaml:"edge"`
	Web          []Web        `yaml:"web"`
}

type Identity struct {
	Cert        string `yaml:"cert"`
	Server_cert string `yaml:"server_cert"`
	Key         string `yaml:"key"`
	Ca          string `yaml:"ca"`
}

type Ctrl struct {
	Listener string `yaml:"listener"`
}

type Mgmt struct {
	Listener string `yaml:"listener"`
}

type HealthChecks struct {
	BoltCheck BoltCheck `yaml:"boltCheck"`
}

type BoltCheck struct {
	Interval     string `yaml:"interval"`
	Timeout      string `yaml:"timeout"`
	InitialDelay string `yaml:"initialDelay"`
}

type Edge struct {
	Api        Api        `yaml:"api"`
	Enrollment Enrollment `yaml:"enrollment"`
}

type Api struct {
	SessionTimeout string `yaml:"sessionTimeout"`
	Address        string `yaml:"address"`
}

type Enrollment struct {
	SigningCert  SigningCert  `yaml:"signingCert"`
	EdgeIdentity EdgeIdentity `yaml:"edgeIdentity"`
	EdgeRouter   EdgeRouter   `yaml:"edgeRouter"`
}

type SigningCert struct {
	Cert string `yaml:"cert"`
	Key  string `yaml:"key"`
}

type EdgeIdentity struct {
	Duration string `yaml:"duration"`
}

type EdgeRouter struct {
	Duration string `yaml:"duration"`
}

type Web struct {
	Name       string       `yaml:"name"`
	BindPoints []BindPoints `yaml:"bindPoints"`
	Identity   Identity     `yaml:"identity"`
	Options    Options      `yaml:"options"`
	Apis       []Apis       `yaml:"apis"`
}

type BindPoints struct {
	BpInterface string `yaml:"interface"`
	Address     string `yaml:"address"`
}

type Options struct {
	IdleTimeout   string `yaml:"idleTimeout"`
	ReadTimeout   string `yaml:"readTimeout"`
	WriteTimeout  string `yaml:"writeTimeout"`
	MinTLSVersion string `yaml:"minTLSVersion"`
	MaxTLSVersion string `yaml:"maxTLSVersion"`
}

type Apis struct {
	Binding string     `yaml:"binding"`
	Options ApiOptions `yaml:"options"`
}

type ApiOptions struct {
	// Unsure of this format right now
}

/* END Controller config template structure */

var controllerOptions = CreateConfigControllerOptions{}

func TestControllerOutputPathDoesNotExist(t *testing.T) {
	expectedErrorMsg := "stat /IDoNotExist: no such file or directory"

	// Create the options with non-existent path
	options := &CreateConfigControllerOptions{}
	options.Output = "/IDoNotExist/MyController.yaml"

	err := options.run(&ConfigTemplateValues{})

	assert.EqualError(t, err, expectedErrorMsg, "Error does not match, expected %s but got %s", expectedErrorMsg, err)
}

func TestCreateConfigControllerTemplateValues(t *testing.T) {

	// Create and run the CLI command (capture output, otherwise config prints to stdout instead of test results)
	// This must run first, otherwise the addresses used later won't be correct; this command re-allocates the `data` struct
	_ = execCreateConfigControllerCommand(nil, nil)

	expectedNonEmptyStringFields := []string{
		".ZitiHome",
		".Controller.Identity.Cert",
		".Controller.Identity.ServerCert",
		".Controller.Identity.Key",
		".Controller.Identity.Ca",
		".Controller.Ctrl.ListenerAddress",
		".Controller.Ctrl.ListenerPort",
		".Controller.Mgmt.ListenerAddress",
		".Controller.Mgmt.ListenerPort",
		".Controller.EdgeApi.Address",
		".Controller.EdgeApi.Port",
		".Controller.EdgeEnrollment.SigningCert",
		".Controller.EdgeEnrollment.SigningCertKey",
		".Controller.Web.BindPoints.InterfaceAddress",
		".Controller.Web.BindPoints.InterfacePort",
		".Controller.Web.BindPoints.AddressAddress",
		".Controller.Web.BindPoints.AddressPort",
		".Controller.Web.Identity.Ca",
		".Controller.Web.Identity.Key",
		".Controller.Web.Identity.ServerCert",
		".Controller.Web.Identity.Cert",
		".Controller.Web.Options.MinTLSVersion",
		".Controller.Web.Options.MaxTLSVersion",
	}
	expectedNonEmptyStringValues := []*string{
		&data.ZitiHome,
		&data.Controller.Identity.Cert,
		&data.Controller.Identity.ServerCert,
		&data.Controller.Identity.Key,
		&data.Controller.Identity.Ca,
		&data.Controller.Ctrl.ListenerAddress,
		&data.Controller.Ctrl.ListenerPort,
		&data.Controller.Mgmt.ListenerAddress,
		&data.Controller.Mgmt.ListenerPort,
		&data.Controller.EdgeApi.Address,
		&data.Controller.EdgeApi.Port,
		&data.Controller.EdgeEnrollment.SigningCert,
		&data.Controller.EdgeEnrollment.SigningCertKey,
		&data.Controller.Web.BindPoints.InterfaceAddress,
		&data.Controller.Web.BindPoints.InterfacePort,
		&data.Controller.Web.BindPoints.AddressAddress,
		&data.Controller.Web.BindPoints.AddressPort,
		&data.Controller.Web.Identity.Ca,
		&data.Controller.Web.Identity.Key,
		&data.Controller.Web.Identity.ServerCert,
		&data.Controller.Web.Identity.Cert,
		&data.Controller.Web.Options.MinTLSVersion,
		&data.Controller.Web.Options.MaxTLSVersion,
	}

	expectedNonZeroTimeFields := []string{
		".Controller.Ctrl.MinConnectTimeout",
		".Controller.Ctrl.MaxConnectTimeout",
		".Controller.Ctrl.DefaultConnectTimeout",
		".Controller.Mgmt.MinConnectTimeout",
		".Controller.Mgmt.MaxConnectTimeout",
		".Controller.Mgmt.DefaultConnectTimeout",
		".Controller.HealthChecks.Interval",
		".Controller.HealthChecks.Timeout",
		".Controller.HealthChecks.InitialDelay",
		".Controller.EdgeApi.APIActivityUpdateInterval",
		".Controller.EdgeApi.SessionTimeout",
		".Controller.EdgeEnrollment.DefaultEdgeIdentityDuration",
		".Controller.EdgeEnrollment.EdgeIdentityDuration",
		".Controller.EdgeEnrollment.DefaultEdgeRouterDuration",
		".Controller.EdgeEnrollment.EdgeRouterDuration",
		".Controller.Web.Options.IdleTimeout",
		".Controller.Web.Options.ReadTimeout",
		".Controller.Web.Options.WriteTimeout",
	}

	expectedNonZeroTimeValues := []*time.Duration{
		&data.Controller.Ctrl.MinConnectTimeout,
		&data.Controller.Ctrl.MaxConnectTimeout,
		&data.Controller.Ctrl.DefaultConnectTimeout,
		&data.Controller.Mgmt.MinConnectTimeout,
		&data.Controller.Mgmt.MaxConnectTimeout,
		&data.Controller.Mgmt.DefaultConnectTimeout,
		&data.Controller.HealthChecks.Interval,
		&data.Controller.HealthChecks.Timeout,
		&data.Controller.HealthChecks.InitialDelay,
		&data.Controller.EdgeApi.APIActivityUpdateInterval,
		&data.Controller.EdgeApi.SessionTimeout,
		&data.Controller.EdgeEnrollment.DefaultEdgeIdentityDuration,
		&data.Controller.EdgeEnrollment.EdgeIdentityDuration,
		&data.Controller.EdgeEnrollment.DefaultEdgeRouterDuration,
		&data.Controller.EdgeEnrollment.EdgeRouterDuration,
		&data.Controller.Web.Options.IdleTimeout,
		&data.Controller.Web.Options.ReadTimeout,
		&data.Controller.Web.Options.WriteTimeout,
	}

	expectedNonZeroIntFields := []string{
		".Controller.Ctrl.DefaultQueuedConnects",
		".Controller.Ctrl.MinOutstandingConnects",
		".Controller.Ctrl.MaxOutstandingConnects",
		".Controller.Ctrl.DefaultOutstandingConnects",
		".Controller.Mgmt.MinQueuedConnects",
		".Controller.Mgmt.MaxQueuedConnects",
		".Controller.Mgmt.DefaultQueuedConnects",
		".Controller.Mgmt.MinOutstandingConnects",
		".Controller.Mgmt.MaxOutstandingConnects",
		".Controller.Mgmt.DefaultOutstandingConnects",
		".Controller.EdgeApi.APIActivityUpdateBatchSize",
	}

	expectedNonZeroIntValues := []*int{
		&data.Controller.Ctrl.DefaultQueuedConnects,
		&data.Controller.Ctrl.MinOutstandingConnects,
		&data.Controller.Ctrl.MaxOutstandingConnects,
		&data.Controller.Ctrl.DefaultOutstandingConnects,
		&data.Controller.Mgmt.MinQueuedConnects,
		&data.Controller.Mgmt.MaxQueuedConnects,
		&data.Controller.Mgmt.DefaultQueuedConnects,
		&data.Controller.Mgmt.MinOutstandingConnects,
		&data.Controller.Mgmt.MaxOutstandingConnects,
		&data.Controller.Mgmt.DefaultOutstandingConnects,
		&data.Controller.EdgeApi.APIActivityUpdateBatchSize,
	}

	// Check that the expected string template values are not blank
	for field, value := range expectedNonEmptyStringValues {
		assert.NotEqualf(t, "", *value, expectedNonEmptyStringFields[field]+" should be a non-blank value")
	}

	// Check that the expected time.Duration template values are not zero
	for field, value := range expectedNonZeroTimeValues {
		assert.NotZero(t, *value, expectedNonZeroTimeFields[field]+" should be a non-zero value")
	}

	// Check that the expected integer template values are not zero
	for field, value := range expectedNonZeroIntValues {
		assert.NotZero(t, *value, expectedNonZeroIntFields[field]+" should be a non-zero value")
	}
}

func TestCtrlConfigDefaultsWhenUnset(t *testing.T) {
	// Clears template data and unsets all env vars
	clearControllerOptionsAndTemplateData()

	ctrlConfig := execCreateConfigControllerCommand(nil, nil)

	t.Run("TestPKIControllerCert", func(t *testing.T) {
		expectedValue := workingDir + "/" + cmdhelper.HostnameOrNetworkName() + ".cert"

		require.Equal(t, expectedValue, data.Controller.Identity.Cert)
		assert.Equal(t, expectedValue, ctrlConfig.Identity.Cert)
	})

	t.Run("TestWebAdvertisedAddress", func(t *testing.T) {
		expectedValue, _ := os.Hostname()

		require.Equal(t, expectedValue, data.Controller.Web.BindPoints.AddressAddress)
		assert.Equal(t, expectedValue, strings.Split(ctrlConfig.Web[0].BindPoints[0].Address, ":")[0])
	})
}

func TestCtrlConfigDefaultsWhenBlank(t *testing.T) {
	keys := map[string]string{
		"ZITI_PKI_CTRL_CERT":               "",
		"ZITI_CTRL_WEB_ADVERTISED_ADDRESS": "",
	}
	// run the config
	ctrlConfig := execCreateConfigControllerCommand(nil, keys)

	t.Run("TestPKIControllerCert", func(t *testing.T) {
		expectedValue := workingDir + "/" + cmdhelper.HostnameOrNetworkName() + ".cert"

		require.Equal(t, expectedValue, data.Controller.Identity.Cert)
		assert.Equal(t, expectedValue, ctrlConfig.Identity.Cert)
	})

	t.Run("TestWebAdvertisedAddress", func(t *testing.T) {
		expectedValue, _ := os.Hostname()

		require.Equal(t, expectedValue, data.Controller.Web.BindPoints.AddressAddress)
		assert.Equal(t, expectedValue, strings.Split(ctrlConfig.Web[0].BindPoints[0].Address, ":")[0])
	})
}

func TestZitiCtrlIdentityCert(t *testing.T) {
	customValue := "/var/test/custom/path"
	keys := map[string]string{
		"ZITI_PKI_CTRL_CERT": customValue,
	}

	t.Run("TestSetCtrlIdentityCertPath", func(t *testing.T) {
		_ = execCreateConfigControllerCommand(nil, keys)
		assert.Equal(t, customValue, data.Controller.Identity.Cert)
	})

	t.Run("TestSetCtrlIdentityCertPathInConfig", func(t *testing.T) {
		ctrlConfig := execCreateConfigControllerCommand(nil, keys)
		assert.Equal(t, customValue, ctrlConfig.Identity.Cert)
	})
}

func TestAdvertisedBindAddress(t *testing.T) {
	customValue := "123.456.7.8"
	keys := map[string]string{
		"ZITI_CTRL_WEB_ADVERTISED_ADDRESS": customValue,
	}

	t.Run("TestSetAdvertisedBindAddress", func(t *testing.T) {
		_ = execCreateConfigControllerCommand(nil, keys)
		assert.Equal(t, customValue, data.Controller.Web.BindPoints.AddressAddress)
	})

	t.Run("TestSetAdvertisedBindAddressInConfig", func(t *testing.T) {
		configStruct := execCreateConfigControllerCommand(nil, keys)
		assert.Contains(t, customValue, strings.Split(configStruct.Web[0].BindPoints[0].Address, ":")[0])
	})
}

//// Edge Ctrl Listener port should use ZITI_CTRL_WEB_ADVERTISED_PORT if it is set
//func TestListenerAddressWhenEdgeCtrlPortAndListenerHostPortNotSet(t *testing.T) {
//	myPort := "1234"
//	expectedListenerAddress := "0.0.0.0:" + myPort
//
//	// Make sure the related env vars are unset
//	_ = os.Unsetenv("ZITI_CTRL_EDGE_LISTENER_HOST_PORT")
//
//	// Set the edge controller port
//	_ = os.Setenv("ZITI_CTRL_WEB_ADVERTISED_PORT", myPort)
//
//	// Create and run the CLI command (capture output, otherwise config prints to stdout instead of test results)
//	cmd := NewCmdCreateConfigController()
//	_ = captureOutput(func() {
//		_ = cmd.Execute()
//	})
//
//	assert.Equal(t, expectedListenerAddress, data.Controller.Edge.ListenerHostPort)
//}
//
//// Edge Ctrl Listener address and port should always use ZITI_EDGE_CTRL_LISTENER_HOST_PORT value if it is set
//func TestListenerAddressWhenEdgeCtrlPortAndListenerHostPortSet(t *testing.T) {
//	myPort := "1234"
//	expectedListenerAddress := "0.0.0.0:4321" // Expecting a different port even when edge ctrl port is set
//
//	// Set a custom value for the host and port
//	_ = os.Setenv("ZITI_CTRL_EDGE_LISTENER_HOST_PORT", expectedListenerAddress)
//
//	// Set the edge controller port (this should not show up in the end resulting listener address)
//	_ = os.Setenv("ZITI_CTRL_WEB_ADVERTISED_PORT", myPort)
//
//	// Create and run the CLI command (capture output, otherwise config prints to stdout instead of test results)
//	cmd := NewCmdCreateConfigController()
//	_ = captureOutput(func() {
//		_ = cmd.Execute()
//	})
//
//	assert.Equal(t, expectedListenerAddress, data.Controller.Edge.ListenerHostPort)
//}
//
//// Edge Ctrl Advertised Port should update the edge ctrl port to the default when ZITI_CTRL_WEB_ADVERTISED_PORT is not set
//func TestDefaultEdgeCtrlAdvertisedPort(t *testing.T) {
//	expectedPort := "1280" // Expecting the default port of 1280
//
//	// Set a custom value for the host and port
//	_ = os.Unsetenv("ZITI_CTRL_WEB_ADVERTISED_PORT")
//
//	// Create and run the CLI command (capture output, otherwise config prints to stdout instead of test results)
//	cmd := NewCmdCreateConfigController()
//	_ = captureOutput(func() {
//		_ = cmd.Execute()
//	})
//
//	assert.Equal(t, expectedPort, data.Controller.Edge.AdvertisedPort)
//}
//
//// Edge Ctrl Advertised Port should update the edge ctrl port to the custom value when ZITI_CTRL_WEB_ADVERTISED_PORT is set
//func TestEdgeCtrlAdvertisedPortValueWhenSet(t *testing.T) {
//	expectedPort := "1234" // Setting a custom port which is not the default value
//
//	// Set a custom value for the host and port
//	_ = os.Setenv("ZITI_CTRL_WEB_ADVERTISED_PORT", expectedPort)
//
//	// Create and run the CLI command (capture output, otherwise config prints to stdout instead of test results)
//	cmd := NewCmdCreateConfigController()
//	_ = captureOutput(func() {
//		_ = cmd.Execute()
//	})
//
//	assert.Equal(t, expectedPort, data.Controller.Edge.AdvertisedPort)
//}
//
//func TestDefaultEdgeIdentityEnrollmentDuration(t *testing.T) {
//	// Expect the default (3 hours)
//	expectedDuration := time.Duration(180) * time.Minute
//	expectedConfigValue := "180m"
//
//	// Unset the env var so the default is used
//	_ = os.Unsetenv("ZITI_EDGE_IDENTITY_ENROLLMENT_DURATION")
//
//	// Create and run the CLI command (capture output, otherwise config prints to stdout instead of test results)
//	cmd := NewCmdCreateConfigController()
//	configOutput := captureOutput(func() {
//		_ = cmd.Execute()
//	})
//
//	assert.Equal(t, expectedDuration, data.Controller.EdgeIdentityDuration)
//
//	// Expect that the config value is represented correctly
//	configStruct := configToStruct(configOutput)
//	assert.Equal(t, expectedConfigValue, configStruct.Edge.Enrollment.EdgeIdentity.Duration)
//}
//
//func TestEdgeIdentityEnrollmentDurationWhenEnvVarSet(t *testing.T) {
//	expectedDuration := 5 * time.Minute // Setting a custom duration which is not the default value
//	expectedConfigValue := "5m"
//
//	// Set a custom value for the enrollment duration
//	_ = os.Setenv("ZITI_EDGE_IDENTITY_ENROLLMENT_DURATION", fmt.Sprintf("%.0f", expectedDuration.Minutes()))
//
//	// Create and run the CLI command (capture output, otherwise config prints to stdout instead of test results)
//	cmd := NewCmdCreateConfigController()
//	configOutput := captureOutput(func() {
//		_ = cmd.Execute()
//	})
//
//	assert.Equal(t, expectedDuration, data.Controller.EdgeIdentityDuration)
//
//	// Expect that the config value is represented correctly
//	configStruct := configToStruct(configOutput)
//	assert.Equal(t, expectedConfigValue, configStruct.Edge.Enrollment.EdgeIdentity.Duration)
//}
//
//func TestEdgeIdentityEnrollmentDurationWhenEnvVarSetToBlank(t *testing.T) {
//	// Expect the default (3 hours)
//	expectedDuration := time.Duration(180) * time.Minute
//	expectedConfigValue := "180m"
//
//	// Set a custom value for the enrollment duration
//	_ = os.Setenv("ZITI_EDGE_IDENTITY_ENROLLMENT_DURATION", "")
//
//	// Create and run the CLI command (capture output, otherwise config prints to stdout instead of test results)
//	cmd := NewCmdCreateConfigController()
//	configOutput := captureOutput(func() {
//		_ = cmd.Execute()
//	})
//
//	assert.Equal(t, expectedDuration, data.Controller.EdgeIdentityDuration)
//
//	// Expect that the config value is represented correctly
//	configStruct := configToStruct(configOutput)
//	assert.Equal(t, expectedConfigValue, configStruct.Edge.Enrollment.EdgeIdentity.Duration)
//}
//
//func TestEdgeIdentityEnrollmentDurationCLITakesPriority(t *testing.T) {
//	envVarValue := 5 * time.Minute // Setting a custom duration which is not the default value
//	cliValue := "10m"              // Setting a CLI custom duration which is also not the default value
//	expectedConfigValue := "10m"
//
//	// Set a custom value for the enrollment duration
//	_ = os.Setenv("ZITI_EDGE_IDENTITY_ENROLLMENT_DURATION", fmt.Sprintf("%.0f", envVarValue.Minutes()))
//
//	// Create and run the CLI command (capture output, otherwise config prints to stdout instead of test results)
//	cmd := NewCmdCreateConfigController()
//	cmd.SetArgs([]string{"--identityEnrollmentDuration", cliValue})
//	configOutput := captureOutput(func() {
//		_ = cmd.Execute()
//	})
//
//	// Expect that the CLI value was used over the environment variable
//	expectedValue, _ := time.ParseDuration(cliValue)
//	assert.Equal(t, expectedValue, data.Controller.EdgeIdentityDuration)
//
//	// Expect that the config value is represented correctly
//	configStruct := configToStruct(configOutput)
//	assert.Equal(t, expectedConfigValue, configStruct.Edge.Enrollment.EdgeIdentity.Duration)
//}
//
//func TestDefaultEdgeRouterEnrollmentDuration(t *testing.T) {
//	expectedDuration := time.Duration(180) * time.Minute
//	expectedConfigValue := "180m"
//
//	// Unset the env var so the default is used
//	_ = os.Unsetenv("ZITI_EDGE_ROUTER_ENROLLMENT_DURATION")
//
//	// Create and run the CLI command (capture output, otherwise config prints to stdout instead of test results)
//	cmd := NewCmdCreateConfigController()
//	configOutput := captureOutput(func() {
//		_ = cmd.Execute()
//	})
//
//	assert.Equal(t, expectedDuration, data.Controller.EdgeRouterDuration)
//
//	// Expect that the config value is represented correctly
//	configStruct := configToStruct(configOutput)
//	assert.Equal(t, expectedConfigValue, configStruct.Edge.Enrollment.EdgeRouter.Duration)
//}
//
//func TestEdgeRouterEnrollmentDurationWhenEnvVarSet(t *testing.T) {
//	expectedDuration := 5 * time.Minute // Setting a custom duration which is not the default value
//	expectedConfigValue := "5m"
//
//	// Set a custom value for the enrollment duration
//	_ = os.Setenv("ZITI_EDGE_ROUTER_ENROLLMENT_DURATION", fmt.Sprintf("%.0f", expectedDuration.Minutes()))
//
//	// Create and run the CLI command (capture output, otherwise config prints to stdout instead of test results)
//	cmd := NewCmdCreateConfigController()
//	configOutput := captureOutput(func() {
//		_ = cmd.Execute()
//	})
//
//	assert.Equal(t, expectedDuration, data.Controller.EdgeRouterDuration)
//
//	// Expect that the config value is represented correctly
//	configStruct := configToStruct(configOutput)
//	assert.Equal(t, expectedConfigValue, configStruct.Edge.Enrollment.EdgeRouter.Duration)
//}
//
//func TestEdgeRouterEnrollmentDurationWhenEnvVarSetToBlank(t *testing.T) {
//	// Expect the default (3 hours)
//	expectedDuration := time.Duration(180) * time.Minute
//	expectedConfigValue := "180m"
//
//	// Set a custom value for the enrollment duration
//	_ = os.Setenv("ZITI_EDGE_ROUTER_ENROLLMENT_DURATION", "")
//
//	// Create and run the CLI command
//	cmd := NewCmdCreateConfigController()
//	configOutput := captureOutput(func() {
//		_ = cmd.Execute()
//	})
//
//	assert.Equal(t, expectedDuration, data.Controller.EdgeRouterDuration)
//
//	// Expect that the config value is represented correctly
//	configStruct := configToStruct(configOutput)
//	assert.Equal(t, expectedConfigValue, configStruct.Edge.Enrollment.EdgeRouter.Duration)
//}
//
//func TestEdgeRouterEnrollmentDurationCLITakesPriority(t *testing.T) {
//	envVarValue := 5 * time.Minute // Setting a custom duration which is not the default value
//	cliValue := "10m"              // Setting a CLI custom duration which is also not the default value
//	expectedConfigValue := "10m"   // Config value representation should be in minutes
//
//	// Set a custom value for the enrollment duration
//	_ = os.Setenv("ZITI_EDGE_ROUTER_ENROLLMENT_DURATION", fmt.Sprintf("%.0f", envVarValue.Minutes()))
//
//	// Create and run the CLI command (capture output, otherwise config prints to stdout instead of test results)
//	cmd := NewCmdCreateConfigController()
//	cmd.SetArgs([]string{"--routerEnrollmentDuration", cliValue})
//	configOutput := captureOutput(func() {
//		_ = cmd.Execute()
//	})
//
//	// Expect that the CLI value was used over the environment variable
//	expectedValue, _ := time.ParseDuration(cliValue)
//	assert.Equal(t, expectedValue, data.Controller.EdgeRouterDuration)
//
//	// Expect that the config value is represented correctly
//	configStruct := configToStruct(configOutput)
//	assert.Equal(t, expectedConfigValue, configStruct.Edge.Enrollment.EdgeRouter.Duration)
//}
//
//func TestEdgeRouterEnrollmentDurationCLIConvertsToMin(t *testing.T) {
//	cliValue := "1h"             // Setting a CLI custom duration which is also not the default value
//	expectedConfigValue := "60m" // Config value representation should be in minutes
//
//	// Make sure the env var is not set
//	_ = os.Unsetenv("ZITI_EDGE_ROUTER_ENROLLMENT_DURATION")
//
//	// Create and run the CLI command
//	cmd := NewCmdCreateConfigController()
//	cmd.SetArgs([]string{"--routerEnrollmentDuration", cliValue})
//	configOutput := captureOutput(func() {
//		_ = cmd.Execute()
//	})
//
//	// Expect that the CLI value was used over the environment variable
//	expectedValue, _ := time.ParseDuration(cliValue)
//	assert.Equal(t, expectedValue, data.Controller.EdgeRouterDuration)
//
//	// Expect that the config value is represented correctly
//	configStruct := configToStruct(configOutput)
//	assert.Equal(t, expectedConfigValue, configStruct.Edge.Enrollment.EdgeRouter.Duration)
//}

func TestEdgeRouterAndIdentityEnrollmentDurationTogetherCLI(t *testing.T) {
	cliIdentityDurationValue := "1h"
	cliRouterDurationValue := "30m"
	expectedIdentityConfigValue := "60m"
	expectedRouterConfigValue := "30m"

	// Create and run the CLI command
	args := []string{"--routerEnrollmentDuration", cliRouterDurationValue, "--identityEnrollmentDuration", cliIdentityDurationValue}
	configStruct := execCreateConfigControllerCommand(args, nil)

	// Expect that the config values are represented correctly
	assert.Equal(t, expectedIdentityConfigValue, configStruct.Edge.Enrollment.EdgeIdentity.Duration)
	assert.Equal(t, expectedRouterConfigValue, configStruct.Edge.Enrollment.EdgeRouter.Duration)
}

func TestEdgeRouterAndIdentityEnrollmentDurationTogetherEnvVar(t *testing.T) {
	envVarIdentityDurationValue := "120"
	envVarRouterDurationValue := "60"
	expectedIdentityConfigValue := envVarIdentityDurationValue + "m"
	expectedRouterConfigValue := envVarRouterDurationValue + "m"

	// Create and run the CLI command
	keys := map[string]string{
		"ZITI_EDGE_IDENTITY_ENROLLMENT_DURATION": envVarIdentityDurationValue,
		"ZITI_EDGE_ROUTER_ENROLLMENT_DURATION":   envVarRouterDurationValue,
	}
	configStruct := execCreateConfigControllerCommand(nil, keys)

	// Expect that the config values are represented correctly
	assert.Equal(t, expectedIdentityConfigValue, configStruct.Edge.Enrollment.EdgeIdentity.Duration)
	assert.Equal(t, expectedRouterConfigValue, configStruct.Edge.Enrollment.EdgeRouter.Duration)
}

func clearControllerOptionsAndTemplateData() {
	controllerOptions = CreateConfigControllerOptions{}
	data = &ConfigTemplateValues{}

	unsetZitiEnv()
}

func configToStruct(config string) ControllerConfig {
	configStruct := ControllerConfig{}
	err2 := yaml.Unmarshal([]byte(config), &configStruct)
	if err2 != nil {
		fmt.Println(err2)
	}
	return configStruct
}

func execCreateConfigControllerCommand(args []string, keys map[string]string) ControllerConfig {
	// Setup
	clearControllerOptionsAndTemplateData()
	controllerOptions.Output = defaultOutput

	setEnvByMap(keys)
	// Create and run the CLI command (capture output to convert to a template struct)
	cmd := NewCmdCreateConfigController()
	cmd.SetArgs(args)
	configOutput := captureOutput(func() {
		_ = cmd.Execute()
	})

	return configToStruct(configOutput)
}
