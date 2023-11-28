pipeline {
    agent any
    environment {
        GIT_REPO_URL = 'https://github.com/ajaysharma507/flux2actions.git'
    }
    stages {
        stage('Delete Old Tags') {
            steps {
                script {
                    // Clone the Git repository
                    checkout([$class: 'Git', branches: [[name: '*/master']], credentialsId: 'your-git-credentials-id', url: GIT_REPO_URL])
                    // Get a list of all tags
                    def tags = sh(script: 'git tag -l', returnStdout: true).trim().split('\n')
                    // Get the current date
                    def currentDate = new Date()
                    // Iterate through each tag and delete tags older than 60 days
                    tags.each { tag ->
                        def tagDate = sh(script: "git log -1 --format=%ai ${tag}", returnStdout: true).trim()
                        def daysDifference = currentDate - Date.parse("yyyy-MM-dd HH:mm:ss Z", tagDate)
                        def days = daysDifference.toDays()
                        if (days > 60) {
                            echo "Deleting old tag: ${tag}"
                            sh "git tag -d ${tag}"
                            sh "git push origin :refs/tags/${tag}"
                        }
                    }
                }
            }
        }
    }
}
