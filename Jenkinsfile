pipeline {
    agent {
        label 'docker'
    }

    stages {
        stage('Test single model (pass if not fails)') {
            steps {
                sh 'bin/pack3d 20 40 30 2.0 pack3d_transforms 1 tests/stl/logo.stl'
                sh 'rm pack3d_transforms.json'
            }
        }

        stage('Test multiple models (pass if not fails)') {
            steps {
                sh 'bin/pack3d 20 40 30 2.0 pack3d_transforms 1 tests/stl/cube.stl 1 tests/stl/logo.stl 1 tests/stl/corner.stl'
                sh 'rm pack3d_transforms.json'
            }
        }
    }
}