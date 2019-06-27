# S2IRun

[![License](http://img.shields.io/badge/license-apache%20v2-blue.svg)](https://github.com/KubeSphere/s2irun/blob/master/LICENSE)

## 介绍

[S2I](https://github.com/openshift/source-to-image)( source to image )是一款由Openshift开发、自动将代码容器化的工具，通过预置的模板来支持多种语言和框架，诸如Java，Nodejs, python等等。S2IRun可以将S2I容器化，以便可以在任何容器环境中运行，也可以更好的运行在kubernetes中。

S2IRun主要的工作就是根据用户配置的参数，去执行拉取代码、生成执行文件、构建镜像、推送镜像等等相关操作。对于用户而言，只需提供代码仓库地址等相关配置参数即可，也不需要编写Dockerfile，从源码到最终镜像之间的所有操作都由S2IRun自动完成。



## 快速开始

1. Fork示例代码**[devops-java-sample](https://github.com/kubesphere/devops-java-sample)**至个人仓库

2. 编写配置文件，并放置自定义目录，配置文件参考如下。

   ```json
   {
     "builderImage": "kubesphere/java-8-centos7:advanced-2.0.0",   // builderImage基础镜像
     "pushAuthentication": {
       "username": "USERNAME_REPLACE",
       "password": "PASSWORD_REPLACE"
     },		// 镜像仓库账户信息，例如dockerhub账号
     "imageName": "USERNAME_REPLACE/s2i",
     "tag": "USERNAME_REPLACE/s2i:tag",
     "builderPullPolicy": "if-not-present",    // builder镜像拉取策略
     "export": true,
     "sourceUrl": "https://GIT_USERNAME_REPLACE:GIT_PASSWORD_REPLACE@github.com/GIT_USERNAME_REPLACE/devops-java-sample.git",		// 代码仓库地址
   }
   ```

3. 设置环境变量`S2I_CONFIG_PATH`，为配置文件的路径，例如`export S2I_CONFIG_PATH=/etc/data/config.json`。

4. 执行命令`go run cmd/main.go --v=4 --logtostderr=true`，S2IRun则开始运行。拉取builderImage，拉取代码，生成运行文件、构建runtimeImage等步骤。

5. 最后看到log中输出`Build completed successfully`，及构建完成。

## 完整S2IRun 配置参数

```
type Config struct {
	// DisplayName is a result image display-name label. This defaults to the
	// output image name.
	DisplayName string `json:"displayName,omitempty"`

	// Description is a result image description label. The default is no
	// description.
	Description string `json:"description,omitempty"`

	// BuilderImage describes which image is used for building the result images.
	BuilderImage string `json:"builderImage,omitempty"`

	// BuilderImageVersion provides optional version information about the builder image.
	BuilderImageVersion string `json:"builderImageVersion,omitempty"`

	// BuilderBaseImageVersion provides optional version information about the builder base image.
	BuilderBaseImageVersion string `json:"builderBaseImageVersion,omitempty"`

	// RuntimeImage specifies the image that will be a base for resulting image
	// and will be used for running an application. By default, BuilderImage is
	// used for building and running, but the latter may be overridden.
	RuntimeImage string `json:"runtimeImage,omitempty"`

	// RuntimeImagePullPolicy specifies when to pull a runtime image.
	RuntimeImagePullPolicy PullPolicy `json:"runtimeImagePullPolicy,omitempty"`

	// RuntimeAuthentication holds the authentication information for pulling the
	// runtime Docker images from private repositories.
	RuntimeAuthentication AuthConfig `json:"runtimeAuthentication,omitempty"`

	// RuntimeArtifacts specifies a list of source/destination pairs that will
	// be copied from builder to a runtime image. Source can be a file or
	// directory. Destination must be a directory. Regardless whether it
	// is an absolute or relative path, it will be placed into image's WORKDIR.
	// Destination also can be empty or equals to ".", in this case it just
	// refers to a root of WORKDIR.
	// In case it's empty, S2I will try to get this list from
	// io.openshift.s2i.assemble-input-files label on a RuntimeImage.
	RuntimeArtifacts VolumeList `json:"runtimeArtifacts,omitempty"`

	// DockerConfig describes how to access host docker daemon.
	DockerConfig *DockerConfig `json:"dockerConfig,omitempty"`

	// PullAuthentication holds the authentication information for pulling the
	// Docker images from private repositories
	PullAuthentication AuthConfig `json:"pullAuthentication,omitempty"`

	// PullAuthentication holds the authentication information for pulling the
	// Docker images from private repositories
	PushAuthentication AuthConfig `json:"pushAuthentication,omitempty"`

	// IncrementalAuthentication holds the authentication information for pulling the
	// previous image from private repositories
	IncrementalAuthentication AuthConfig `json:"incrementalAuthentication,omitempty"`

	// DockerNetworkMode is used to set the docker network setting to --net=container:<id>
	// when the builder is invoked from a container.
	DockerNetworkMode DockerNetworkMode `json:"dockerNetworkMode,omitempty"`

	// PreserveWorkingDir describes if working directory should be left after processing.
	PreserveWorkingDir bool `json:"preserveWorkingDir,omitempty"`

	// IgnoreSubmodules determines whether we will attempt to pull in submodules
	// (via --recursive or submodule init)
	IgnoreSubmodules bool `json:"ignoreSubmodules,omitempty"`

	// Source URL describing the location of sources used to build the result image.
	Source *git.URL `json:"source,omitempty"`

	// Tag is a result image tag name.
	Tag string `json:"tag,omitempty"`

	// BuilderPullPolicy specifies when to pull the builder image
	BuilderPullPolicy PullPolicy `json:"builderPullPolicy,omitempty"`

	// PreviousImagePullPolicy specifies when to pull the previously build image
	// when doing incremental build
	PreviousImagePullPolicy PullPolicy `json:"previousImagePullPolicy,omitempty"`

	// Incremental describes whether to try to perform incremental build.
	Incremental bool `json:"incremental,omitempty"`

	// IncrementalFromTag sets an alternative image tag to look for existing
	// artifacts. Tag is used by default if this is not set.
	IncrementalFromTag string `json:"incrementalFromTag,omitempty"`

	// RemovePreviousImage describes if previous image should be removed after successful build.
	// This applies only to incremental builds.
	RemovePreviousImage bool `json:"removePreviousImage,omitempty"`

	// Environment is a map of environment variables to be passed to the image.
	Environment EnvironmentList `json:"environment,omitempty"`

	// LabelNamespace provides the namespace under which the labels will be generated.
	LabelNamespace string `json:"labelNamespace,omitempty"`

	// CallbackURL is a URL which is called upon successful build to inform about that fact.
	CallbackURL string `json:"callbackURL,omitempty"`

	// ScriptsURL is a URL describing where to fetch the S2I scripts from during build process.
	// This url can be a reference within the builder image if the scheme is specified as image://
	ScriptsURL string `json:"scriptsURL,omitempty"`

	// Destination specifies a location where the untar operation will place its artifacts.
	Destination string `json:"destination,omitempty"`

	// WorkingDir describes temporary directory used for downloading sources, scripts and tar operations.
	WorkingDir string `json:"workingDir,omitempty"`

	// WorkingSourceDir describes the subdirectory off of WorkingDir set up during the repo download
	// that is later used as the root for ignore processing
	WorkingSourceDir string `json:"workingSourceDir,omitempty"`

	// LayeredBuild describes if this is build which layered scripts and sources on top of BuilderImage.
	LayeredBuild bool `json:"layeredBuild,omitempty"`

	// Operate quietly. Progress and assemble script output are not reported, only fatal errors.
	// (default: false).
	Quiet bool `json:"quiet,omitempty"`

	// ForceCopy results in only the file SCM plugin being used (i.e. no `git clone`); allows for empty directories to be included
	// in resulting image (since git does not support that).
	// (default: false).
	ForceCopy bool `json:"forceCopy,omitempty"`

	// Specify a relative directory inside the application repository that should
	// be used as a root directory for the application.
	ContextDir string `json:"contextDir,omitempty"`

	// AllowedUIDs is a list of user ranges of users allowed to run the builder image.
	// If a range is specified and the builder (or runtime) image uses a non-numeric
	// user or a user that is outside the specified range, then the build fails.
	AllowedUIDs user.RangeList `json:"allowedUIDs,omitempty"`

	// AssembleUser specifies the user to run the assemble script in container
	AssembleUser string `json:"assembleUser,omitempty"`

	// RunImage will trigger a "docker run ..." invocation of the produced image so the user
	// can see if it operates as he would expect
	RunImage bool `json:"runImage,omitempty"`

	// Usage allows for properly shortcircuiting s2i logic when `s2i usage` is invoked
	Usage bool `json:"usage,omitempty"`

	// Injections specifies a list source/destination folders that are injected to
	// the container that runs assemble.
	// All files we inject will be truncated after the assemble script finishes.
	Injections VolumeList `json:"injections,omitempty"`

	// CGroupLimits describes the cgroups limits that will be applied to any containers
	// run by s2i.
	CGroupLimits *CGroupLimits `json:"cGroupLimits,omitempty"`

	// DropCapabilities contains a list of capabilities to drop when executing containers
	DropCapabilities []string `json:"dropCapabilities,omitempty"`

	// ScriptDownloadProxyConfig optionally specifies the http and https proxy
	// to use when downloading scripts
	ScriptDownloadProxyConfig *ProxyConfig `json:"scriptDownloadProxyConfig,omitempty"`

	// ExcludeRegExp contains a string representation of the regular expression desired for
	// deciding which files to exclude from the tar stream
	ExcludeRegExp string `json:"excludeRegExp,omitempty"`

	// BlockOnBuild prevents s2i from performing a docker build operation
	// if one is necessary to execute ONBUILD commands, or to layer source code into
	// the container for images that don't have a tar binary available, if the
	// image contains ONBUILD commands that would be executed.
	BlockOnBuild bool `json:"blockOnBuild,omitempty"`

	// HasOnBuild will be set to true if the builder image contains ONBUILD instructions
	HasOnBuild bool `json:"hasOnBuild,omitempty"`

	// BuildVolumes specifies a list of volumes to mount to container running the
	// build.
	BuildVolumes []string `json:"buildVolumes,omitempty"`

	// Labels specify labels and their values to be applied to the resulting image. Label keys
	// must have non-zero length. The labels defined here override generated labels in case
	// they have the same name.
	Labels map[string]string `json:"labels,omitempty"`

	// SourceInfo provides the info about the source to be built rather than relying
	// on the Downloader to retrieve it.
	SourceInfo *git.SourceInfo `json:"sourceInfo,omitempty"`

	// SecurityOpt are passed as options to the docker containers launched by s2i.
	SecurityOpt []string `json:"securityOpt,omitempty"`

	// KeepSymlinks indicates to copy symlinks as symlinks. Default behavior is to follow
	// symlinks and copy files by content.
	KeepSymlinks bool `json:"keepSymlinks,omitempty"`

	// AsDockerfile indicates the path where the Dockerfile should be written instead of building
	// a new image.
	AsDockerfile string `json:"asDockerfile,omitempty"`

	// ImageWorkDir is the default working directory for the builder image.
	ImageWorkDir string `json:"imageWorkDir,omitempty"`

	// ImageScriptsURL is the default location to find the assemble/run scripts for a builder image.
	// This url can be a reference within the builder image if the scheme is specified as image://
	ImageScriptsURL string `json:"imageScriptsURL,omitempty"`

	// AddHost Add a line to /etc/hosts for test purpose or private use in LAN. Its format is host:IP,muliple hosts can be added  by using multiple --add-host
	AddHost []string `json:"addHost,omitempty"`

	//Export Push the result image to specify image registry in tag
	Export bool `json:"export,omitempty"`

	SourceURL string `json:"sourceURL,omitempty"`
}
```