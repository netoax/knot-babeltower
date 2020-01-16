pipeline {
    agent { docker { image 'golang' } }
    stages {
        stage('install dependencies') {
            steps {
                sh 'sudo apt-get install python-pip'
                sh 'make tools'
                sh 'make deps'
            }
        }
        stage('run linter') {
            steps {
                sh 'make lint'
            }
        }
        stage('run tests') {
            steps {
                sh 'make tests'
            }
        }
    }
}