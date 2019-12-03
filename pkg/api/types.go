package api

import (
	"errors"
	"fmt"
	"github.com/opencontainers/go-digest"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/docker/distribution/reference"
	"github.com/kubesphere/s2irun/pkg/scm/git"
	utilglog "github.com/kubesphere/s2irun/pkg/utils/glog"
	"github.com/kubesphere/s2irun/pkg/utils/user"
)

var glog = utilglog.StderrLog

// invalidFilenameCharacters contains a list of character we consider malicious
// when injecting the directories into containers.
const invalidFilenameCharacters = `;*?"<>|%#$!+{}&[],"'` + "`"

const (
	// PullAlways means that we always attempt to pull the latest image.
	PullAlways PullPolicy = "always"

	// PullNever means that we never pull an image, but only use a local image.
	PullNever PullPolicy = "never"

	// PullIfNotPresent means that we pull if the image isn't present on disk.
	PullIfNotPresent PullPolicy = "if-not-present"

	// DefaultBuilderPullPolicy specifies the default pull policy to use
	DefaultBuilderPullPolicy = PullIfNotPresent

	// DefaultRuntimeImagePullPolicy specifies the default pull policy to use.
	DefaultRuntimeImagePullPolicy = PullIfNotPresent

	// DefaultPreviousImagePullPolicy specifies policy for pulling the previously
	// build Docker image when doing incremental build
	DefaultPreviousImagePullPolicy = PullIfNotPresent

	// DefaultTag is the image tag, being applied if none is specified.
	DefaultTag = "latest"

	// AnnotationBuildResultKey is the filed in job annotations.
	AnnotationBuildResultKey = "s2iBuildResult"

	// AnnotationBuildSourceKey is the filed in job annotations.
	AnnotationBuildSourceKey = "s2iBuildSource"

	// Command for pull docker image.
	CommandPull = "docker pull "

	// Default source branch.
	DefaultBranch = "master"
)

// Config contains essential fields for performing build.
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

	// Export Push the result image to specify image registry in tag
	Export bool `json:"export,omitempty"`

	// SourceURL is  url of the codes such as https://github.com/a/b.git
	SourceURL string `json:"sourceURL,omitempty"`

	// IsBinaryURL explain the type of SourceURL.
	// If it is IsBinaryURL, it will download the file directly without using git.
	IsBinaryURL bool `json:"isBinaryURL,omitempty"`

	// The RevisionId is a branch name or a SHA-1 hash of every important thing about the commit
	RevisionId string `json:"revisionId,omitempty"`

	// Output build result. If build not in k8s cluster, can not use this field.
	OutputBuildResult bool `json:"outputBuildResult,omitempty"`
}

// DeepCopyInto to implement k8s api requirement
func (c *Config) DeepCopyInto(out *Config) {
	*out = *c

	//slice
	if c.DropCapabilities != nil {
		out.DropCapabilities = make([]string, len(c.DropCapabilities))
		copy(out.DropCapabilities, c.DropCapabilities)
	}
	if c.BuildVolumes != nil {
		out.BuildVolumes = make([]string, len(c.BuildVolumes))
		copy(out.BuildVolumes, c.BuildVolumes)
	}
	if c.AddHost != nil {
		out.AddHost = make([]string, len(c.AddHost))
		copy(out.AddHost, c.AddHost)
	}
	if c.SecurityOpt != nil {
		out.SecurityOpt = make([]string, len(c.SecurityOpt))
		copy(out.SecurityOpt, c.SecurityOpt)
	}

	//pointer
	if c.DockerConfig != nil {
		out.DockerConfig = new(DockerConfig)
		*(out.DockerConfig) = *(c.DockerConfig)
	}
	if c.Source != nil {
		out.Source = new(git.URL)
		*(out.Source) = *(c.Source)
	}
	if c.SourceInfo != nil {
		out.SourceInfo = new(git.SourceInfo)
		*(out.SourceInfo) = *(c.SourceInfo)
	}
	if c.CGroupLimits != nil {
		out.CGroupLimits = new(CGroupLimits)
		*(out.CGroupLimits) = *(c.CGroupLimits)
	}
}

func (c *Config) DeepCopy() *Config {
	out := new(Config)
	c.DeepCopyInto(out)
	return out
}

// EnvironmentSpec specifies a single environment variable.
type EnvironmentSpec struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

// EnvironmentList contains list of environment variables.
type EnvironmentList []EnvironmentSpec

// ProxyConfig holds proxy configuration.
type ProxyConfig struct {
	HTTPProxy  *url.URL
	HTTPSProxy *url.URL
}

// CGroupLimits holds limits used to constrain container resources.
type CGroupLimits struct {
	MemoryLimitBytes int64  `json:"memoryLimitBytes,omitempty"`
	CPUShares        int64  `json:"cpuShares,omitempty"`
	CPUPeriod        int64  `json:"cpuPeriod,omitempty"`
	CPUQuota         int64  `json:"cpuQuota,omitempty"`
	MemorySwap       int64  `json:"memorySwap,omitempty"`
	Parent           string `json:"parent,omitempty"`
}

// VolumeSpec represents a single volume mount point.
type VolumeSpec struct {
	// Source is a reference to the volume source.
	Source string `json:"source,omitempty"`
	// Destination is the path to mount the volume to - absolute or relative.
	Destination string `json:"destination,omitempty"`
	// Keep indicates if the mounted data should be kept in the final image.
	Keep bool `json:"keep,omitempty"`
}

// VolumeList contains list of VolumeSpec.
type VolumeList []VolumeSpec

// DockerConfig contains the configuration for a Docker connection.
type DockerConfig struct {
	// Endpoint is the docker network endpoint or socket
	Endpoint string `json:"endpoint,omitempty"`

	// CertFile is the certificate file path for a TLS connection
	CertFile string `json:"certFile,omitempty"`

	// KeyFile is the key file path for a TLS connection
	KeyFile string `json:"keyFile,omitempty"`

	// CAFile is the certificate authority file path for a TLS connection
	CAFile string `json:"caFile,omitempty"`

	// UseTLS indicates if TLS must be used
	UseTLS bool `json:"useTLS,omitempty"`

	// TLSVerify indicates if TLS peer must be verified
	TLSVerify bool `json:"tlsVerify,omitempty"`
}

// AuthConfig is our abstraction of the Registry authorization information for whatever
// docker client we happen to be based on
type AuthConfig struct {
	Username      string `json:"username,omitempty"`
	Password      string `json:"password,omitempty"`
	Email         string `json:"email,omitempty"`
	ServerAddress string `json:"serverAddress,omitempty"`
}

// ContainerConfig is the abstraction of the docker client provider (formerly go-dockerclient, now either
// engine-api or kube docker client) container.Config type that is leveraged by s2i or origin
type ContainerConfig struct {
	Labels map[string]string
	Env    []string
}

// Image is the abstraction of the docker client provider (formerly go-dockerclient, now either
// engine-api or kube docker client) Image type that is leveraged by s2i or origin
type Image struct {
	ID string
	*ContainerConfig
	Config *ContainerConfig
}

// Result structure contains information from build process.
type Result struct {
	// Success describes whether the build was successful.
	Success bool

	// Messages is a list of messages from build process.
	Messages []string

	// WorkingDir describes temporary directory used for downloading sources, scripts and tar operations.
	WorkingDir string

	// BuildInfo holds information about the result of a build.
	BuildInfo BuildInfo

	// ImageInfo describes resulting image info.
	ResultInfo OutputResultInfo
	// Source info.
	SourceInfo SourceInfo
}

type SourceInfo struct {
	SourceUrl    string `json:"sourceUrl,omitempty"`
	RevisionId   string `json:"revisionId,omitempty"`
	BuilderImage string `json:"builderImage,omitempty"`
	Description  string `json:"description,omitempty"`

	CommitID       string `json:"commitID,omitempty"`
	CommitterName  string `json:"committerName,omitempty"`
	CommitterEmail string `json:"committerEmail,omitempty"`

	BinaryName string `json:"binaryName,omitempty"`
	BinarySize uint64 `json:"binarySize,omitempty"`
}

type OutputResultInfo struct {
	ImageName     string   `json:"imageName,omitempty"`
	ImageID       string   `json:"imageID,omitempty"`
	ImageSize     int64    `json:"imageSize,omitempty"`
	ImageCreated  string   `json:"imageCreated,omitempty"`
	ImageRepoTags []string `json:"imageRepoTags,omitempty"`
	CommandPull   string   `json:"commandPull,omitempty"`
}

// BuildInfo contains information about the build process.
type BuildInfo struct {
	// Stages contains details about each build stage.
	Stages []StageInfo

	// FailureReason is a camel case reason that is used by the machine to reply
	// back to the OpenShift builder with information why any of the steps in the
	// build failed.
	FailureReason FailureReason
}

// StageInfo contains details about a build stage.
type StageInfo struct {
	// Name is the identifier for each build stage.
	Name StageName

	// StartTime identifies when this stage started.
	StartTime time.Time

	// DurationMilliseconds identifies how long this stage ran.
	DurationMilliseconds int64

	// Steps contains details about each build step within a build stage.
	Steps []StepInfo
}

// StageName is the identifier for each build stage.
type StageName string

// Valid StageNames
const (
	// StagePullImages pulls the docker images.
	StagePullImages StageName = "PullImages"

	//StageAssemble runs the assemble steps.
	StageAssemble StageName = "Assemble"

	// StageBuild builds the source.
	StageBuild StageName = "Build"

	// StageCommit commits the container.
	StageCommit StageName = "CommitContainer"

	// StageRetrieve retrieves artifacts.
	StageRetrieve StageName = "RetrieveArtifacts"
)

// StepInfo contains details about a build step.
type StepInfo struct {
	// Name is the identifier for each build step.
	Name StepName

	// StartTime identifies when this step started.
	StartTime time.Time

	// DurationMilliseconds identifies how long this step ran.
	DurationMilliseconds int64
}

// StepName is the identifier for each build step.
type StepName string

// Valid StepNames
const (
	// StepPullBuilderImage pulls the builder image.
	StepPullBuilderImage StepName = "PullBuilderImage"

	// StepPullPreviousImage pulls the previous image for an incremental build.
	StepPullPreviousImage StepName = "PullPreviousImage"

	// StepPullRuntimeImage pull the runtime image.
	StepPullRuntimeImage StepName = "PullRuntimeImage"

	// StepAssembleBuildScripts runs the assemble scripts.
	StepAssembleBuildScripts StepName = "AssembleBuildScripts"

	// StepBuildDockerImage builds the Docker image for layered builds.
	StepBuildDockerImage StepName = "BuildDockerImage"

	// StepCommitContainer commits the container to the builder image.
	StepCommitContainer StepName = "CommitContainer"

	// StepRetrievePreviousArtifacts restores archived artifacts from the previous build.
	StepRetrievePreviousArtifacts StepName = "RetrievePreviousArtifacts"
)

// StepFailureReason holds the type of failure that occurred during the build
// process.
type StepFailureReason string

// StepFailureMessage holds the detailed message of a failure.
type StepFailureMessage string

// FailureReason holds the type of failure that occurred during the build
// process.
type FailureReason struct {
	Reason  StepFailureReason
	Message StepFailureMessage
}

// InstallResult structure describes the result of install operation
type InstallResult struct {
	// Script describes which script this result refers to
	Script string

	// URL describes from where the script was taken
	URL string

	// Downloaded describes if download operation happened, this will be true for
	// external scripts, but false for scripts from inside the image
	Downloaded bool

	// Installed describes if script was installed to upload directory
	Installed bool

	// Error describes last error encountered during install operation
	Error error

	// FailedSources is a list of sources that were attempted but failed
	// when downloading this script
	FailedSources []string
}

// DockerNetworkMode specifies the network mode setting for the docker container
type DockerNetworkMode string

// Image holds information about an image.
type ImageInfo struct {
	Domain string
	Path   string
	Tag    string
	Digest digest.Digest
	named  reference.Named
}

const (
	// DockerNetworkModeHost places the container in the default (host) network namespace.
	DockerNetworkModeHost DockerNetworkMode = "host"
	// DockerNetworkModeBridge instructs docker to create a network namespace for this container connected to the docker0 bridge via a veth-pair.
	DockerNetworkModeBridge DockerNetworkMode = "bridge"
	// DockerNetworkModeContainerPrefix is the string prefix used by NewDockerNetworkModeContainer.
	DockerNetworkModeContainerPrefix string = "container:"
	// DockerNetworkModeNetworkNamespacePrefix is the string prefix used when sharing a namespace from a CRI-O container.
	DockerNetworkModeNetworkNamespacePrefix string = "netns:"
)

// NewDockerNetworkModeContainer creates a DockerNetworkMode value which instructs docker to place the container in the network namespace of an existing container.
// It can be used, for instance, to place the s2i container in the network namespace of the infrastructure container of a k8s pod.
func NewDockerNetworkModeContainer(id string) DockerNetworkMode {
	return DockerNetworkMode(DockerNetworkModeContainerPrefix + id)
}

// PullPolicy specifies a type for the method used to retrieve the Docker image
type PullPolicy string

// String implements the String() function of pflags.Value so this can be used as
// command line parameter.
// This method is really used just to show the default value when printing help.
// It will not default the configuration.
func (p *PullPolicy) String() string {
	if len(string(*p)) == 0 {
		return string(DefaultBuilderPullPolicy)
	}
	return string(*p)
}

// Type implements the Type() function of pflags.Value interface
func (p *PullPolicy) Type() string {
	return "string"
}

// Set implements the Set() function of pflags.Value interface
// The valid options are "always", "never" or "if-not-present"
func (p *PullPolicy) Set(v string) error {
	switch v {
	case "always":
		*p = PullAlways
	case "never":
		*p = PullNever
	case "if-not-present":
		*p = PullIfNotPresent
	default:
		return fmt.Errorf("invalid value %q, valid values are: always, never or if-not-present", v)
	}
	return nil
}

// IsInvalidFilename verifies if the provided filename contains malicious
// characters.
func IsInvalidFilename(name string) bool {
	return strings.ContainsAny(name, invalidFilenameCharacters)
}

// Set implements the Set() function of pflags.Value interface.
// This function parses the string that contains source:destination pair.
// When the destination is not specified, the source get copied into current
// working directory in container.
func (l *VolumeList) Set(value string) error {
	volumes := strings.Split(value, ";")
	newVols := make([]VolumeSpec, len(volumes))
	for i, v := range volumes {
		spec, err := l.parseSpec(v)
		if err != nil {
			return err
		}
		newVols[i] = *spec
	}
	*l = append(*l, newVols...)
	return nil
}

func (l *VolumeList) parseSpec(value string) (*VolumeSpec, error) {
	if len(value) == 0 {
		return nil, errors.New("invalid format, must be source:destination")
	}
	var mount []string
	pos := strings.LastIndex(value, ":")
	if pos == -1 {
		mount = []string{value, ""}
	} else {
		mount = []string{value[:pos], value[pos+1:]}
	}
	mount[0] = strings.Trim(mount[0], `"'`)
	mount[1] = strings.Trim(mount[1], `"'`)
	s := &VolumeSpec{Source: filepath.Clean(mount[0]), Destination: filepath.ToSlash(filepath.Clean(mount[1]))}
	if IsInvalidFilename(s.Source) || IsInvalidFilename(s.Destination) {
		return nil, fmt.Errorf("invalid characters in filename: %q", value)
	}
	return s, nil
}

// String implements the String() function of pflags.Value interface.
func (l *VolumeList) String() string {
	result := []string{}
	for _, i := range *l {
		result = append(result, strings.Join([]string{i.Source, i.Destination}, ":"))
	}
	return strings.Join(result, ",")
}

// Type implements the Type() function of pflags.Value interface.
func (l *VolumeList) Type() string {
	return "string"
}

// Set implements the Set() function of pflags.Value interface.
func (e *EnvironmentList) Set(value string) error {
	parts := strings.SplitN(value, "=", 2)
	if len(parts) != 2 || len(parts[0]) == 0 {
		return fmt.Errorf("invalid environment format %q, must be NAME=VALUE", value)
	}
	if strings.Contains(parts[1], ",") && strings.Contains(parts[1], "=") {
		glog.Warningf("DEPRECATED: Use multiple -e flags to specify multiple environment variables instead of comma (%q)", value)
	}
	*e = append(*e, EnvironmentSpec{
		Name:  strings.TrimSpace(parts[0]),
		Value: strings.TrimSpace(parts[1]),
	})
	return nil
}

// String implements the String() function of pflags.Value interface.
func (e *EnvironmentList) String() string {
	result := []string{}
	for _, i := range *e {
		result = append(result, strings.Join([]string{i.Name, i.Value}, "="))
	}
	return strings.Join(result, ",")
}

// Type implements the Type() function of pflags.Value interface.
func (e *EnvironmentList) Type() string {
	return "string"
}

// AsBinds converts the list of volume definitions to go-dockerclient compatible
// list of bind mounts.
func (l *VolumeList) AsBinds() []string {
	result := make([]string, len(*l))
	for index, v := range *l {
		result[index] = strings.Join([]string{v.Source, v.Destination}, ":")
	}
	return result
}

func Parse(originalName, serverAddress string) (ref string, err error) {

	image, err := parseImage(originalName)
	if err != nil {
		return "", fmt.Errorf("parsing image %q failed: %v", originalName, err)
	}

	// remove schema if required
	if strings.HasPrefix(serverAddress, "http") {
		_,serverAddress,err = Getscheme(serverAddress)
		if err != nil {
			return "", fmt.Errorf("parsing image %q failed: %v", image, err)
		}
		serverAddress = strings.Trim(serverAddress, "//")
	}

	if image.Domain != serverAddress && serverAddress != "" {
		ref = serverAddress + "/" + image.Path + ":" + image.Tag
	} else {
		ref = image.String()
	}

	return ref, err
}

// ParseImage returns an Image struct with all the values filled in for a given image.
// example : localhost:5000/nginx:latest, nginx:perl etc.
func parseImage(image string) (*ImageInfo, error) {

	// Parse the image name and tag.
	named, err := reference.ParseNormalizedNamed(image)

	if err != nil {
		return nil, fmt.Errorf("parsing image %q failed: %v", image, err)
	}
	// Add the latest lag if they did not provide one.
	named = reference.TagNameOnly(named)

	i := &ImageInfo{
		named:  named,
		Domain: reference.Domain(named),
		Path:   reference.Path(named),
	}

	// Add the tag if there was one.
	if tagged, ok := named.(reference.Tagged); ok {
		i.Tag = tagged.Tag()
	}

	return i, nil
}

// Maybe rawurl is of the form scheme:path.
// (Scheme must be [a-zA-Z][a-zA-Z0-9+-.]*)
// If so, return scheme, path; else return "", rawurl.
func Getscheme(rawurl string) (scheme, path string, err error) {
	for i := 0; i < len(rawurl); i++ {
		c := rawurl[i]
		switch {
		case 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z':
		// do nothing
		case '0' <= c && c <= '9' || c == '+' || c == '-' || c == '.':
			if i == 0 {
				return "", rawurl, nil
			}
		case c == ':':
			if i == 0 {
				return "", "", nil
			}
			return rawurl[:i], rawurl[i+1:], nil
		default:
			// we have encountered an invalid character,
			// so there is no valid scheme
			return "", rawurl, nil
		}
	}
	return "", rawurl, nil
}

// String returns the string representation of an image.
func (i *ImageInfo) String() string {
	return i.named.String()
}

// Reference returns either the digest if it is non-empty or the tag for the image.
func (i *ImageInfo) Reference() string {
	if len(i.Digest.String()) > 1 {
		return i.Digest.String()
	}

	return i.Tag
}
