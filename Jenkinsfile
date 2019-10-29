
pipeline {
    agent {
        node {
            label 'go'
        }
    }

    environment {
      IMAGE_NAME = 'kubespheredev/s2irun:advanced-2.1.0'
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
                branch 'jenkinsfile'
            }

            steps {
                container ('go') {
                    sh "docker build . -t $IMAGE_NAME"
                    withCredentials([usernamePassword(passwordVariable : 'DOCKER_PASSWORD' ,usernameVariable : 'DOCKER_USERNAME' ,credentialsId : "dockerhub-id" ,)]) {
                        sh 'echo "$DOCKER_PASSWORD" | docker login -u "$DOCKER_USERNAME" --password-stdin'
                        sh 'docker push $IMAGE_NAME'
                    }
                }
            }
        }
    }


}