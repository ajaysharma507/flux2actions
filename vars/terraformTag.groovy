#!/usr/bin/env groovy

def call(Map config = [:]) {
     gitRepo       = config.get('gitRepo')
     environment   = config.get('environment')
      

    //set git tags + changelog by envar and push
    wrap([$class: 'BuildUser']) {
	sh "git config --global user.email \"jenkins@polar.security\""
	sh "git config --global user.name \"Jenkins\""
    }           
    TAG = "${environment}-latest"
    //print the tag that will be used for visibality
    //println "tagging with - ${TAG}"

    //delete, replace and create tag in git with ref to code
    withCredentials([usernamePassword(credentialsId: 'polardev-jenkins-integration', passwordVariable: 'GIT_PASSWORD', usernameVariable: 'GIT_USERNAME')]) {
        sh ("git push https://${GIT_USERNAME}:${GIT_PASSWORD}@github.com/polarsecurity/${gitRepo}.git :refs/tags/${TAG}")
        // sh("git tag -fa ${TAG} -m \"${TAG} ${RELEASE_NOTES}\"")
        sh("git tag -fa ${TAG} -m \"${TAG} ${GIT_COMMIT}\"")
        sh("git push https://${GIT_USERNAME}:${GIT_PASSWORD}@github.com/polarsecurity/${gitRepo}.git --tags")
    }                   
}
