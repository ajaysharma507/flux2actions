pipeline {
    agent any
    environment {
        GIT_REPO_URL = 'https://your.git.repo/url.git'
    }
    stages {
        stage('Delete Old Tags') {
            steps {
                script {
                    def tags = sh(script: 'git tag -l', returnStdout: true).trim().split('\n')
                    // Get the current date
                    def currentDate = new Date()
                    // Iterate through each tag and delete tags older than 60 days
                    tags.each { tag ->
                        def tagDate = sh(script: "git log -1 --format=%ai ${tag}", returnStdout: true).trim()
                        // Parse the date using SimpleDateFormat
                        def dateFormat = new java.text.SimpleDateFormat("yyyy-MM-dd HH:mm:ss Z")
                        def dateObj = dateFormat.parse(tagDate)
                        // Calculate the days difference
                        def daysDifference = currentDate - dateObj
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







