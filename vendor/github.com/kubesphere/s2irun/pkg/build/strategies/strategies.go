package strategies

import (
	"time"

	"github.com/kubesphere/s2irun/pkg/api"
	"github.com/kubesphere/s2irun/pkg/build"
	"github.com/kubesphere/s2irun/pkg/build/strategies/dockerfile"
	"github.com/kubesphere/s2irun/pkg/build/strategies/onbuild"
	"github.com/kubesphere/s2irun/pkg/build/strategies/sti"
	"github.com/kubesphere/s2irun/pkg/docker"
	"github.com/kubesphere/s2irun/pkg/utils/fs"
	utilstatus "github.com/kubesphere/s2irun/pkg/utils/status"
)

// GetStrategy decides what build strategy will be used for the STI build.
// TODO: deprecated, use Strategy() instead
func GetStrategy(client docker.Client, config *api.Config) (build.Builder, api.BuildInfo, error) {
	return Strategy(client, config, build.Overrides{})
}

// Strategy creates the appropriate build strategy for the provided config, using
// the overrides provided. Not all strategies support all overrides.
func Strategy(client docker.Client, config *api.Config, overrides build.Overrides) (build.Builder, api.BuildInfo, error) {
	var builder build.Builder
	var buildInfo api.BuildInfo
	var err error

	fileSystem := fs.NewFileSystem()

	startTime := time.Now()

	if len(config.AsDockerfile) != 0 {
		builder, err = dockerfile.New(config, fileSystem)
		if err != nil {
			buildInfo.FailureReason = utilstatus.NewFailureReason(
				utilstatus.ReasonGenericS2IBuildFailed,
				utilstatus.ReasonMessageGenericS2iBuildFailed,
			)
			return nil, buildInfo, err
		}
		return builder, buildInfo, nil
	}

	dkr := docker.New(client, config.PullAuthentication, config.PushAuthentication)
	image, err := docker.GetBuilderImage(dkr, config)
	buildInfo.Stages = api.RecordStageAndStepInfo(buildInfo.Stages, api.StagePullImages, api.StepPullBuilderImage, startTime, time.Now())
	if err != nil {
		buildInfo.FailureReason = utilstatus.NewFailureReason(
			utilstatus.ReasonPullBuilderImageFailed,
			utilstatus.ReasonMessagePullBuilderImageFailed,
		)
		return nil, buildInfo, err
	}
	config.HasOnBuild = image.OnBuild

	if config.AssembleUser, err = docker.GetAssembleUser(dkr, config); err != nil {
		buildInfo.FailureReason = utilstatus.NewFailureReason(
			utilstatus.ReasonPullBuilderImageFailed,
			utilstatus.ReasonMessagePullBuilderImageFailed,
		)
		return nil, buildInfo, err
	}

	// if we're blocking onbuild, just do a normal s2i build flow
	// which won't do a docker build and invoke the onbuild commands
	if image.OnBuild && !config.BlockOnBuild {
		builder, err = onbuild.New(client, config, fileSystem, overrides)
		if err != nil {
			buildInfo.FailureReason = utilstatus.NewFailureReason(
				utilstatus.ReasonGenericS2IBuildFailed,
				utilstatus.ReasonMessageGenericS2iBuildFailed,
			)
			return nil, buildInfo, err
		}
		return builder, buildInfo, nil
	}

	builder, err = sti.New(client, config, fileSystem, overrides)
	if err != nil {
		buildInfo.FailureReason = utilstatus.NewFailureReason(
			utilstatus.ReasonGenericS2IBuildFailed,
			utilstatus.ReasonMessageGenericS2iBuildFailed,
		)
		return nil, buildInfo, err
	}
	return builder, buildInfo, err
}
