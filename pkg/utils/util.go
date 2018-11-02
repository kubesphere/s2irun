package utils

import (
	"github.com/docker/docker/api/types/container"

	utilglog "github.com/magicsong/s2irun/pkg/utils/glog"
)

var glog = utilglog.StderrLog

// SafeForLoggingContainerConfig returns a copy of the container.Config object
// with sensitive information (proxy environment variables containing credentials)
// redacted.
func SafeForLoggingContainerConfig(config *container.Config) *container.Config {
	strippedEnv := SafeForLoggingEnv(config.Env)
	newConfig := *config
	newConfig.Env = strippedEnv
	return &newConfig
}
