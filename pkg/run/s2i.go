package run

import (
	"bytes"
	"fmt"

	"github.com/MagicSong/s2irun/pkg/api"
	"github.com/MagicSong/s2irun/pkg/api/describe"
	"github.com/MagicSong/s2irun/pkg/api/validation"
	"github.com/MagicSong/s2irun/pkg/build/strategies"
	"github.com/MagicSong/s2irun/pkg/docker"
	s2ierr "github.com/MagicSong/s2irun/pkg/errors"
	utilglog "github.com/MagicSong/s2irun/pkg/utils/glog"
)

var glog = utilglog.StderrLog

// S2I Just run the command
func S2I(cfg *api.Config) error {
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

	d := docker.New(client, cfg.PullAuthentication)
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
