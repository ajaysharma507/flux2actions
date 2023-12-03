@Library('shared-lib@main')_
import java.text.SimpleDateFormat

pipeline {
    agent { 
        label 'graviton-medium'
    }
    triggers {
        //Run once a week every Saturday
    cron('30 0 * * 6')
    }
    options {
		buildDiscarder(logRotator(numToKeepStr: '100'))
		disableConcurrentBuilds()
		timestamps()
	}
    environment {
        AWS_ACCESS_KEY_ID       = credentials('AwsAccesskeyID-Devops')
        AWS_SECRET_ACCESS_KEY   = credentials('AwsSecretAccessKey-Devops')
        GIT_TOKEN               = credentials('github-pat-jenkins')
        GIT_ORG                 = 'https://github.com/polarsecurity'
        GIT_SSH                 = credentials('jenkins_ssh_git')
    }
    stages {
        stage ('Git backup'){
            steps {
                script {
                        echo "Backing up Git repos"
                        deleteStale('polarsecurity', 200)
                    }
                }
            }
        }

    post {
        failure  {
            slackSend channel: '#devops_offshore',color: 'bad', tokenCredentialId: 'polar-slack',message: "FAILURE: #${BUILD_NUMBER} build FAILED!\n${RUN_DISPLAY_URL}"
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

    // Calculate the date six months ago
    def currentTimeMillis = System.currentTimeMillis()
    def sixMonthsInMillis = 6L * 29 * 24 * 60 * 60 * 1000
    def sixMonthsAgo = currentTimeMillis - sixMonthsInMillis


    // Send the API request and store the response
    def response = new URL(reposUrl).getText(requestProperties: authHeaders, query: params)

    // Parse the JSON response and print the names of the repositories
    def json = new groovy.json.JsonSlurperClassic().parseText(response)

    json.each { repo ->
        global_prettyPrintWithHeaderAndFooter header: "Backing up Git repos", body: "${repo.name}", lineLenght: 20
        def repoUrl = "git@github.com:polarsecurity/${repo.name}.git"

	// Initialize variables for pagination
	def page = 1
	def branches = []
	while (true) {
	// Get branches for the current page
	def branchesUrl = "https://api.github.com/repos/${organization}/${repo.name}/branches?page=${page}"
	def branchesResponse = new URL(branchesUrl).getText(requestProperties: authHeaders, query: params)
        def branchesjson = new groovy.json.JsonSlurperClassic().parseText(branchesResponse)
	if (branchesjson.size() == 0) {
		break
	}
	branches.addAll(branchesjson)
	page++
	}

        // Iterate over each branch and delete if it's older than the threshold
         branches.each { branch ->

         def lastCommitUrl = "https://api.github.com/repos/${organization}/${repo.name}/commits/${branch.commit.sha}"
         def lastCommitResponse = new URL(lastCommitUrl).getText(requestProperties: authHeaders, query: params)
         def lastCommit = new groovy.json.JsonSlurperClassic().parseText(lastCommitResponse)

	  // Extract the committer and date from the last commit
	  def commitCommitter = lastCommit.commit.committer
	  def dateString = commitCommitter ? commitCommitter.date : null
	  // Use SimpleDateFormat to parse the date string
	  def sdf = new SimpleDateFormat("yyyy-MM-dd'T'HH:mm:ss'Z'")
	  def branchDate = dateString ? sdf.parse(dateString)?.time : null

	    if (branchDate != null && branchDate < sixMonthsAgo) {
	        // Convert the Date object to a string for echo
	        def branchDateString = sdf.format(branchDate)
	        echo "Deleted branch ${branch.name} from ${repo.name}. Last commit date: ${branchDateString}"
                def sixMonthsAgoString = sdf.format(sixMonthsAgo)
                echo (sixMonthsAgoString)
            }

        }
    }
}
