package outputresult

import (
	"encoding/json"
	dockertypes "github.com/docker/docker/api/types"
	"github.com/kubesphere/s2irun/pkg/api"
	utilglog "github.com/kubesphere/s2irun/pkg/utils/glog"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/retry"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

var glog = utilglog.StderrLog

func OutputResult(builderConfig *api.Config, imageInspect *dockertypes.ImageInspect, result *api.Result) {
	// build result info.
	result.ResultInfo.ImageID = imageInspect.ID
	result.ResultInfo.ImageCreated = imageInspect.Created
	result.ResultInfo.ImageRepoTags = imageInspect.RepoTags
	result.ResultInfo.ImageName = builderConfig.Tag
	result.ResultInfo.CommandPull = api.CommandPull + builderConfig.Tag
	result.ResultInfo.ImageSize = imageInspect.Size

	// build source info.
	result.SourceInfo.SourceUrl = builderConfig.SourceURL
	result.SourceInfo.BuilderImage = builderConfig.BuilderImage
	result.SourceInfo.Description = builderConfig.Description
	if builderConfig.RevisionId == "" {
		builderConfig.RevisionId = api.DefaultBranch
	}

	if builderConfig.IsBinaryURL == true {
		result.SourceInfo.BinaryName = builderConfig.SourceInfo.BinaryName
		result.SourceInfo.BinarySize = builderConfig.SourceInfo.BinarySize
	} else {
		result.SourceInfo.RevisionId = builderConfig.RevisionId
		result.SourceInfo.CommitID = builderConfig.SourceInfo.CommitID
		result.SourceInfo.CommitterName = builderConfig.SourceInfo.CommitterName
		result.SourceInfo.CommitterEmail = builderConfig.SourceInfo.CommitterEmail
	}

	// add Annotation
	err := addBuildResultToAnnotation(result)
	if err != nil {
		glog.V(1).Info("Output build result failed, reason: %s.", err)
	}
}

func addBuildResultToAnnotation(buildResult *api.Result) error {
	namespace := os.Getenv(api.S2iRunNamespace)
	jobName := os.Getenv(api.S2iRunJobName)
	if namespace == "" || jobName == "" {
		glog.Warning("failed to get env S2iRunJobName and S2iRunNamespace")
		return nil
	}
	cfg, err := config.GetConfig()
	if err != nil {
		glog.Errorf("failed to get k8s config file, reason: %s", err)
		return err
	}
	k8sClient := kubernetes.NewForConfigOrDie(cfg)

	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		job, err := k8sClient.BatchV1().Jobs(namespace).Get(jobName, v1.GetOptions{})
		if err != nil {
			glog.Errorf("failed to get job %s in namespace %s, reason: %s", jobName, namespace, err)
			return err
		}

		//update job annotations
		result, _ := json.Marshal(buildResult.ResultInfo)
		source, _ := json.Marshal(buildResult.SourceInfo)
		job.Annotations[api.AnnotationBuildResultKey] = string(result)
		job.Annotations[api.AnnotationBuildSourceKey] = string(source)

		_, err = k8sClient.BatchV1().Jobs(namespace).Update(job)
		return err
	})

	return retryErr
}
