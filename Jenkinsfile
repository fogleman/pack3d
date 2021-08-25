pipeline {
    agent {
        label 'docker'
    }

    stages {
        stage('Test single model (pass if not fails)') {
            steps {
                sh 'bin/pack3d --json_file=tests/jenkins_tests/input_jenkins_test_1.json --filename=tests/jenkins_tests/output_jenkins_test_1'
                sh 'rm tests/jenkins_tests/output_jenkins_test_1.json'
            }
        }

        stage('Test multiple models and co-packing (pass if not fails)') {
            steps {
                sh 'bin/pack3d --json_file=tests/jenkins_tests/input_jenkins_test_2.json --filename=tests/jenkins_tests/output_jenkins_test_2'
                sh 'rm tests/jenkins_tests/output_jenkins_test_2.json'
            }
        }
    }
}
