@Library('shared-lib@main')_
import groovy.json.JsonSlurperClassic
import groovy.json.JsonOutput

pipeline {
    agent any
    
    // agent { 
    //     label 'graviton-medium'
    // }
    // triggers {
    //     //Run once a week every Saturday
    // cron('30 0 * * 6')
    // }
    options {
		buildDiscarder(logRotator(numToKeepStr: '100'))
		disableConcurrentBuilds()
		timestamps()
	}
    environment {
        GIT_TOKEN               = credentials('github-pat-jenkins')
        GIT_ORG                 = 'https://github.com/polarsecurity'
        GIT_SSH                 = credentials('jenkins_ssh_git')
    }

    stages {
        
        stage('Delete Stale Branches') {

            steps {
                script {

                    // Get a list of all remote branches
                    def branches = sh(script: 'git ls-remote --heads origin | cut -f 2 | sed \'s,refs/heads/,,\'', returnStdout: true).trim().split('\n')

                    // Iterate through branches
                    branches.each { branch ->
                        // Get the branch's last commit date
                        def lastCommitDate = sh(script: "git log -1 --format=%ct origin/${branch}", returnStdout: true).trim().toLong() * 1000

                        // Calculate 6 months in milliseconds
                        def sixMonthsAgo = System.currentTimeMillis() - (6 * 30 * 24 * 60 * 60 * 1000)

                        // If the branch is older than 6 months, delete it
                        if (lastCommitDate < sixMonthsAgo) {
                            // sh "git push origin --delete ${branch}"
                            echo "Deleted stale branch: ${branch}"
                        }
                    }
                }
            }
        }
    }
}
