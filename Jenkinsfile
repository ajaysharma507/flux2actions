
pipeline {
    agent any
    
    triggers {
        cron('0 0 * * *') // Run every day at midnight
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
                        def tagDateParsed = new Date(tagDate)
                        
                        if (tagDateParsed < oldestDate) {
                            sh "git tag -d ${tag}"
                            sh "git push origin :refs/tags/${tag}"
                            echo "Deleted tag: ${tag}"
                        }
                    }
                }
            }
        }

    }
}
