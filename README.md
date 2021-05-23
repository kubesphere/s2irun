## S2IRun
[![Gitpod ready-to-code](https://img.shields.io/badge/Gitpod-ready--to--code-blue?logo=gitpod)](https://gitpod.io/#https://github.com/kubesphere/s2irun)
[![License](http://img.shields.io/badge/license-apache%20v2-blue.svg)](https://github.com/KubeSphere/s2irun/blob/master/LICENSE)  [![Go Report Card](https://goreportcard.com/badge/github.com/kubesphere/s2ioperator)](https://goreportcard.com/report/github.com/kubesphere/s2irun)  [![S2IRun release](https://img.shields.io/github/release/kubesphere/s2irun.svg?color=release&label=release&logo=release&logoColor=release)](https://github.com/kubesphere/s2irun/releases/tag/v0.0.3)
![Docker Pulls](https://img.shields.io/docker/pulls/kubesphere/s2irun)

Source to image is a command toolkit  and workflow for building reproducible container images from source code. S2IRun is the instance that executes the s2i build process, and can be containerized so that it can run in any container environment. S2IRun produces ready-to-run images by injecting source code into a container image and letting the container prepare that source code for execution. 

You've just built and run a new container image from source code in a git repository, no Dockerfile necessary.

## Quick start

#### Prerequisites

- [golang](https://golang.org/dl/) environment
- [docker](https://docs.docker.com/install/) version 17.03+.

#### Start

1. Create your own fork of [devops-java-sample](https://github.com/kubesphere/devops-java-sample), this is a web project with java8.

2. Clone [S2IRun repo](https://github.com/kubesphere/s2irun) to your machine

3. Write a config file to set parameters used in build:

   ```json
   {
     "builderImage": "kubesphere/java-8-centos7:v2.1.0",
     "pushAuthentication": {
       "username": "USERNAME_REPLACE",
       "password": "PASSWORD_REPLACE"
     },
     "imageName": "USERNAME_REPLACE/s2irun-sample",
     "tag": "USERNAME_REPLACE/s2irun-sample:tag",
     "builderPullPolicy": "if-not-present",
     "export": true,
     "sourceUrl": "https://github.com/GIT_USERNAME_REPLACE/devops-java-sample.git",	                                  
   }
   ```

4. To build image in right,  you'll need to set environment variables:

   ```shell
   export S2I_CONFIG_PATH=PATH/TO/config.json
   ```

5. Run following command to start build.

   ```
   go run cmd/main.go --v=4 --logtostderr=true
   ```

   Process details:

   ​	1. It will pull image that defind in filed `builderImage`. 

   ​	2. clone source code.

   ​	3. start build process

   ​	4. after build completed, it will push read-to-run image

6. Finally, it will push image to your docker registry, please check it.

## About more 

- See [CONTRIBUTING](https://github.com/kubesphere/kubesphere/blob/master/docs/en/guides/Development-workflow.md) for an overview of our processes
- Guidance for [S2I builder image](docs/builder_image.md) creators
