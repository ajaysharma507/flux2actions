pipeline {
    agent any
    
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
                        // Parse the date using TimeCategory
                        def daysDifference
                        use(groovy.time.TimeCategory) {
                            def tagDateTime = Date.parse("yyyy-MM-dd HH:mm:ss Z", tagDate)
                            def diffInMillis = currentDate - tagDateTime
                            daysDifference = diffInMillis.to(TimeUnit.DAYS)
                        }
                        if (daysDifference > 60) {
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
