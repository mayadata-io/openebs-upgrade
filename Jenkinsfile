def ORG = "mayadataio"
def REPO = "openebs-upgrade"

pipeline {
    agent any
    stages {
        stage('Build Image') {
            steps {
                script {
                  GIT_SHA = sh(
                            returnStdout: true,
                            script: "git log -n 1 --pretty=format:'%h'"
                            ).trim()
                  echo "Building Go Code"
                  sh "docker build -t ${ORG}/${REPO}:ci-${GIT_SHA} ."
                }
            }
        }
        stage('Push Image') {
            steps {
                script {
		             withCredentials([usernamePassword( credentialsId: 'docke_cred', usernameVariable: 'USERNAME', passwordVariable: 'PASSWORD')]) {
                            if (env.BRANCH_NAME == 'master')  {
                               echo "Pushing the image with the tag..."
                               sh "docker login -u${USERNAME} -p${PASSWORD} "
			                         sh "docker push ${ORG}/${REPO}:ci-${GIT_SHA}"
                            } else {
			                   echo "WARNING: Not pushing Image"
                        }
                    }
                }
            }
        }
    }
    post {
        always {
            echo 'This will always run'
            deleteDir()
        }
        success {
            script{
              withCredentials([string(credentialsId: 'OPENEBS_UPGRADE_CODECOV_TOKEN', variable: 'CODECOV_TOKEN')]) {
                CODECOV_BASH = sh(returnStdout: true, script: "curl -s https://codecov.io/bash -o .codecov")
		sh 'chmod +x .codecov'
		sh "./.codecov -t ${CODECOV_TOKEN}"
              }
            }
            echo 'This will run only if successful'
            slackSend channel: '#jenkins-builds',
                   color: 'good',
                   message: "The pipeline ${currentBuild.fullDisplayName} completed successfully :dance: :thumbsup: "
        }
        failure {
            echo 'This will run only if failed'
            slackSend channel: '#jenkins-builds',
                  color: 'RED',
                  message: "The pipeline ${currentBuild.fullDisplayName} failed. :scream_cat: :japanese_goblin: "
        }
        unstable {
            echo 'This will run only if the run was marked as unstable'
            slackSend channel: '#jenkins-builds',
                   color: 'good',
                   message: "The pipeline ${currentBuild.fullDisplayName} is unstable :scream_cat: :japanese_goblin: "
        }
        changed {
/*            slackSend channel: '#jenkins-builds',
                   color: 'good',
                   message: "Build ${currentBuild.fullDisplayName} is now stable :dance: :thumbsup: "
            echo 'This will run only if the state of the Pipeline has changed'
*/            echo 'For example, if the Pipeline was previously failing but is now successful'
        }
    }
}
