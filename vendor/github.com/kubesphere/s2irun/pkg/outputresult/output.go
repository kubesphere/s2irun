package outputresult

import (
	"encoding/json"
	"os"
	"time"

	"github.com/kubesphere/s2irun/pkg/api"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/retry"
	"sigs.k8s.io/controller-runtime/pkg/client/config"

	dockertypes "github.com/docker/docker/api/types"
	utilglog "github.com/kubesphere/s2irun/pkg/utils/glog"
)

var (
	Retry = wait.Backoff{
		Steps:    10,
		Duration: 10 * time.Millisecond,
		Factor:   1.0,
		Jitter:   0.1,
	}
	glog = utilglog.StderrLog
)

func OutputResult(builderConfig *api.Config, imageInspect *dockertypes.ImageInspect, result *api.Result) *api.Result {
	// build result info.
	result.ResultInfo.ImageID = imageInspect.ID
	result.ResultInfo.ImageCreated = imageInspect.Created
	result.ResultInfo.ImageRepoTags = imageInspect.RepoTags
	result.ResultInfo.CommandPull = api.CommandPull + builderConfig.Tag
	result.ResultInfo.ImageSize = imageInspect.Size

	// build source info.
	if builderConfig.IsBinaryURL == true {
		result.SourceInfo.BinaryName = builderConfig.SourceInfo.BinaryName
		result.SourceInfo.BinarySize = builderConfig.SourceInfo.BinarySize
	} else {
		result.SourceInfo.CommitID = builderConfig.SourceInfo.CommitID
		result.SourceInfo.CommitterName = builderConfig.SourceInfo.CommitterName
		result.SourceInfo.CommitterEmail = builderConfig.SourceInfo.CommitterEmail
	}

	return result
}

func AddBuildResultToAnnotation(buildResult *api.Result) error {
	namespace := os.Getenv("POD_NAMESPACE")
	podName := os.Getenv("POD_NAME")
	if namespace == "" || podName == "" {
		glog.Warning("failed to get env S2iRun PodName and S2iRun Namespace")
		return nil
	}
	cfg, err := config.GetConfig()
	if err != nil {
		glog.Errorf("failed to get k8s config file, reason: %s", err)
		return err
	}
	k8sClient := kubernetes.NewForConfigOrDie(cfg)

	retryErr := retry.RetryOnConflict(Retry, func() error {
		pod, err := k8sClient.CoreV1().Pods(namespace).Get(podName, v1.GetOptions{})
		if err != nil {
			glog.Errorf("failed to get pod %s in namespace %s, reason: %s", podName, namespace, err)
			return err
		}

		//update pod annotations
		result, _ := json.Marshal(buildResult.ResultInfo)
		source, _ := json.Marshal(buildResult.SourceInfo)
		pod.Annotations = map[string]string{
			api.AnnotationBuildResultKey: string(result),
			api.AnnotationBuildSourceKey: string(source),
		}

		_, err = k8sClient.CoreV1().Pods(namespace).Update(pod)
		return err
	})

	return retryErr
}
