# S2IRun

[![License](http://img.shields.io/badge/license-apache%20v2-blue.svg)](https://github.com/KubeSphere/s2irun/blob/master/LICENSE)

## 介绍

[S2I](https://github.com/openshift/source-to-image)( source to image )是一款由Openshift开发、自动将代码容器化的工具，通过预置的模板来支持多种语言和框架，诸如Java，Nodejs, python等等，[S2I Operator](https://github.com/kubesphere/s2ioperator)将S2I引入到kubernetes中，S2IRun则是具体执行构建工作的资源对象，也是S2I Operator项目的重要组成部分。

> S2IRun API是以CRD的形式生成的，所以如果想要直接访问原生Https接口，请先参考[kubeernetes API文档](https://kubernetes.io/docs/concepts/overview/kubernetes-api/)， 根据文档中说明构造出需要使用的API

## 快速开始

1. S2IRun是定义的一种CRD，需要与S2I Operator结合使用，S2i Operator的安装使用可参考[S2i Operator安装](https://github.com/kubesphere/s2ioperator#如何安装)。

2. 定义S2IBuilder，S2IBuilder存储了S2IRun运行需要的所有配置信息。每一次生成Docker镜像可以复用这些信息，也可以覆盖一些信息。

   ```
   kubectl apply -f - <<EOF
   apiVersion: devops.kubesphere.io/v1alpha1
   kind: S2iBuilder
   metadata:
       name: s2ibuilder-sample
   spec:
       config:
           displayName: "For Test"
           sourceUrl: "https://github.com/sclorg/django-ex"
           builderImage: centos/python-35-centos7
           imageName: kubesphere/hello-python
           tag: v0.0.1
           builderPullPolicy: if-not-present
   EOF
   ```

   可以通过`kubectl get s2ib` 查看当前所有的S2ibuilder状态

   ```
   kubectl get s2ib
   NAME                RUNCOUNT   LASTRUNSTATE   LASTRUNNAME
   s2ibuilder-sample   2          Successful     s2irun-sample1
   ```

3. 定义并部署s2irun，需要制定builderName和newTag。

   ```
   kubectl apply -f - <<EOF
   apiVersion: devops.kubesphere.io/v1alpha1
   kind: S2iRun
   metadata:
       name: s2irun-sample
   spec:
       builderName: s2ibuilder-sample  //指定使用的s2ibuilder
       newTag: v1  //需要更新的image tag
   EOF
   ```

4. 查看Job运行的Node节点，利用命令docker image ls 查看编译好的镜像。

## 完整S2IRun 配置参数

```
type S2iRunSpec struct {
   //BuilderName specify the name of s2ibuilder, required
   BuilderName string `json:"builderName"`
   
   //BackoffLimit limits the restart count of each s2irun. Default is 0
   BackoffLimit int32 `json:"backoffLimit,omitempty"`
   
   //SecondsAfterFinished if is set and greater than zero, and the job created by s2irun become successful or failed , the job will be auto deleted after SecondsAfterFinished
   SecondsAfterFinished int32 `json:"secondsAfterFinished,omitempty"`
   
   //NewTag override the default tag in its s2ibuilder, image name cannot be changed.
   NewTag string `json:"newTag,omitempty"`
}
```