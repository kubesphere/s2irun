
pipeline {
    agent {
        node {
            label 'go'
        }
    }

    parameters {
      string(name:'TAG_NAME',defaultValue: 'latest',description:'')
    }

    environment {
      IMAGE_NAME = 'kubespheredev/s2irun'
    }

    stages {

        stage ('checkout scm') {
            steps {
                checkout(scm)
            }
        }

        stage ('unit test'){
            steps{
                container ('go') {

                    sh '''
mkdir -p /home/jenkins/go/src/github.com/kubesphere
ln -s `pwd` /home/jenkins/go/src/github.com/kubesphere/s2irun
cd /home/jenkins/go/src/github.com/kubesphere/s2irun
make test'''
                }
            }
        }

        stage ('docker push') {
            when{
                branch 'master'
            }

            steps {
                container ('go') {
                    sh "docker build . -t $IMAGE_NAME:$TAG_NAME"
                    withCredentials([usernamePassword(passwordVariable : 'DOCKER_PASSWORD' ,usernameVariable : 'DOCKER_USERNAME' ,credentialsId : "dockerhub-id" ,)]) {
                        sh 'echo "$DOCKER_PASSWORD" | docker login -u "$DOCKER_USERNAME" --password-stdin'
                        sh 'docker push $IMAGE_NAME:$TAG_NAME'
                    }
                }
            }
        }
    }


}
