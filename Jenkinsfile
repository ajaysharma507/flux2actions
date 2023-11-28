pipeline {
    agent any
    
    }
    stages {
        stage('Cleanup Tags') {
            steps {
                script {
                    def oldestDate = new Date() - 60
                    sh "git fetch --tags"
                    def tags = sh(script: 'git tag -l', returnStdout: true).trim().split('\n')
                    for (def tag in tags) {
                        def tagDate = sh(script: "git log -1 --format=%ai ${tag}", returnStdout: true).trim()
                        if (Date.parse("yyyy-MM-dd HH:mm:ss Z", tagDate) < oldestDate) {
                            sh "git tag -d ${tag}"
                            sh "git push origin :refs/tags/${tag}"
                            echo "Deleted tag: ${tag}"
                        }
                    }
                }
            }
        }
        stage('Cleanup Branches') {
            steps {
                script {
                    def oldestDate = new Date() - 60
                    def branches = sh(script: "git for-each-ref --sort=-committerdate refs/heads/ --format='%(refname:short) %(committerdate:short)'", returnStdout: true).trim().split('\n')
                    for (def branch in branches) {
                        def branchName = branch.split()[0]
                        def branchDate = branch.split()[1]
                        if (Date.parse("yyyy-MM-dd", branchDate) < oldestDate) {
                            sh "git branch -D ${branchName}"
                            sh "git push origin --delete ${branchName}"
                            echo "Deleted branch: ${branchName}"
                        }
                    }
                }
            }
        }
    }
