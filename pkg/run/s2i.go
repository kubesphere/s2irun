package run

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/magicsong/s2irun/pkg/scm/git"

	"github.com/magicsong/s2irun/pkg/api"
	"github.com/magicsong/s2irun/pkg/api/describe"
	"github.com/magicsong/s2irun/pkg/api/validation"
	"github.com/magicsong/s2irun/pkg/build/strategies"
	"github.com/magicsong/s2irun/pkg/docker"
	s2ierr "github.com/magicsong/s2irun/pkg/errors"
	utilglog "github.com/magicsong/s2irun/pkg/utils/glog"
)

const (
	ConfigEnvVariable = "S2I_CONFIG_PATH"
)

var glog = utilglog.StderrLog

// S2I Just run the command
func S2I(cfg *api.Config) error {
	cfg.DockerConfig = docker.GetDefaultDockerConfig()
	if len(cfg.AsDockerfile) > 0 {
		if cfg.RunImage {
			return fmt.Errorf("ERROR: --run cannot be used with --as-dockerfile")
		}
		if len(cfg.RuntimeImage) > 0 {
			return fmt.Errorf("ERROR: --runtime-image cannot be used with --as-dockerfile")
		}
	}
	if cfg.Incremental && len(cfg.RuntimeImage) > 0 {
		return fmt.Errorf("ERROR: Incremental build with runtime image isn't supported")
	}
	//set default image pull policy
	if len(cfg.BuilderPullPolicy) == 0 {
		cfg.BuilderPullPolicy = api.DefaultBuilderPullPolicy
	}
	if len(cfg.PreviousImagePullPolicy) == 0 {
		cfg.PreviousImagePullPolicy = api.DefaultPreviousImagePullPolicy
	}
	if len(cfg.RuntimeImagePullPolicy) == 0 {
		cfg.RuntimeImagePullPolicy = api.DefaultRuntimeImagePullPolicy
	}
	if errs := validation.ValidateConfig(cfg); len(errs) > 0 {
		var buf bytes.Buffer
		for _, e := range errs {
			buf.WriteString("ERROR:")
			buf.WriteString(e.Error())
			buf.WriteString("\n")
		}
		return fmt.Errorf(buf.String())
	}

	client, err := docker.NewEngineAPIClient(cfg.DockerConfig)
	if err != nil {
		return err
	}

	d := docker.New(client, cfg.PullAuthentication, cfg.PushAuthentication)
	err = d.CheckReachable()
	if err != nil {
		return err
	}

	glog.V(2).Infof("\n%s\n", describe.Config(client, cfg))

	builder, _, err := strategies.GetStrategy(client, cfg)
	s2ierr.CheckError(err)
	result, err := builder.Build(cfg)
	if err != nil {
		glog.V(0).Infof("Build failed")
		s2ierr.CheckError(err)
		return err
	} else {
		if len(cfg.AsDockerfile) > 0 {
			glog.V(0).Infof("Application dockerfile generated in %s", cfg.AsDockerfile)
		} else {
			glog.V(0).Infof("Build completed successfully")
		}
	}

	//result.Message store Callback Info
	for _, message := range result.Messages {
		glog.V(1).Infof(message)
	}

	return nil
}

func App() int {
	var apiConfig = new(api.Config)
	path := os.Getenv(ConfigEnvVariable)
	file, err := os.Open(path)
	defer file.Close()
	if os.IsNotExist(err) {
		glog.Errorf("Config file does not exist,please check the path: %s", path)
		return 1
	}

	jsonParser := json.NewDecoder(file)
	err = jsonParser.Decode(apiConfig)
	if err != nil {
		glog.Errorf("There are some errors in config file, please check the error:\n%v", err)
		return 1
	}
	apiConfig.Source, err = git.Parse(apiConfig.SourceURL)
	if err != nil {
		glog.Errorf("SourceURL is illegal, please check the error:\n%v", err)
		return 1
	}
	err = S2I(apiConfig)
	if err != nil {
		glog.Errorf("Build failed, please check the error:\n%v", err)
		return 1
	}
	return 0
}
