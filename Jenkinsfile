pipeline {
    agent {
        label 'docker'
    }

    stages {
        stage('Test single model') {
            steps {
                sh 'bin/pack3d 20 40 30 2.0 pack3d_transforms 1 tests/stl/logo.stl'

                sh 'cat pack3d_transforms.json | md5sum'
                sh 'cat tests/fixtures/expected_transformation_logo.json | md5sum'
                sh 'rm pack3d_transforms.json'
            }
        }

        stage('Test multiple models') {
            steps {
                sh 'bin/pack3d 20 40 30 2.0 pack3d_transforms 1 tests/stl/cube.stl 1 tests/stl/logo.stl 1 tests/stl/corner.stl'

                sh 'cat pack3d_transforms.json'
                sh 'rm pack3d_transforms.json'
            }
        }
    }
}