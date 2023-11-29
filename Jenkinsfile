pipeline {
    agent any

    parameters {
        booleanParam(name: 'DELETE_OLD_TAGS', defaultValue: true, description: 'Delete old tags')
        booleanParam(name: 'DELETE_STALE_BRANCHES', defaultValue: true, description: 'Delete stale branches')
    }

    stages {
        stage('Delete Old Tags') {

            when {
                expression { params.DELETE_OLD_TAGS }
            }
            
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
                        def tagDateTime = dateFormat.parse(tagDate)
                        // Calculate the days difference
                        def daysDifference = (currentDate.time - tagDateTime.time) / (1000 * 60 * 60 * 24)
                        if (daysDifference > 60) {
                            echo "Deleting old tag: ${tag}"
                            // sh "git tag -d ${tag}"
                            // sh "git push origin :refs/tags/${tag}"
                        }
                    }
                }
            }
        }
        stage('Delete Stale Branches') {

            when {
                expression { params.DELETE_STALE_BRANCHES }
            }
            
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
