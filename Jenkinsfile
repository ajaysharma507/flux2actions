pipeline {
    
    agent { 
        label 'graviton-medium'
    }
    // triggers {
    // Run once a week every Saturday
    // cron('30 0 * * 6')
    // }
    options {
		buildDiscarder(logRotator(numToKeepStr: '100'))
		disableConcurrentBuilds()
		timestamps()
	}
    environment {
        GIT_TOKEN               = credentials('github-pat-jenkins')
        GIT_ORG                 = 'git@github.com:polarsecurity'
        GIT_SSH                 = credentials('jenkins_ssh_git')
    }

    stages {
        
        stage('Delete Stale Branches') {

            steps {
                script {
                        echo "Deleting stale branches"
                        deleteStale('polarsecurity', 100)
                        }
                    } 
                }
           }
    }

def deleteStale(organization, limit)  {
    // Set the credentials for authenticating with the GitHub API
    def credentials = [
        username: 'jenkins',
        password: "${GIT_TOKEN}"
    ]

    // Set the URL for the GitHub API
    def apiUrl = 'https://api.github.com'

    // Set the URL for the organization's repositories
    def reposUrl = "${apiUrl}/orgs/${organization}/repos"

    // Set the authentication headers
    def authHeaders = [
        "Authorization": "Basic " + "${credentials.username}:${credentials.password}".getBytes().encodeBase64().toString(),
        "Content-Type": "application/json"
    ]
	
    // Set the parameters for the API request
    def params = [
        "type": "all",
        "per_page": "${limit}",
        //"archived": "false"
    ]

    // Send the API request and store the response
    def response = new URL(reposUrl).getText(requestProperties: authHeaders, query: params)

    // Parse the JSON response and print the names of the repositories
    def json = new groovy.json.JsonSlurperClassic().parseText(response)
    

    json.each { repo ->
        global_prettyPrintWithHeaderAndFooter header: "Backing up Git repos", body: "${repo.name}", lineLenght: 20
        def repoUrl = "git@github.com:polarsecurity/${repo.name}.git"

	// Get a list of all remote branches
        def branches = sh(script: 'git ls-remote --heads origin | cut -f 2 | sed \'s,refs/heads/,,\'', returnStdout: true).trim().split('\n')

            // Iterate through branches
            branches.each { branch ->
               // Get the branch's last commit date
                def lastCommitDate = sh(script: "git log -1 --format=%ct origin/${branch}", returnStdout: true).trim().toLong() * 1000

                // Calculate 6 months in milliseconds
                def sixMonthsAgo = System.currentTimeMillis() - (6 * 30 * 24 * 60 * 60 * 1000)

	 }
     }
}	
	
	
