node('master') {
  stage('Unit Tests') {
    git url: "https://github.com/bobbydeveaux/micro-auth-worker.git"
  }
  stage('Build Bin') {
    sh "go get -v -d ./..."
    sh "CGO_ENABLED=0 GOOS=linux go build -o micro-auth-worker ."
  }
  stage('Build Image') {
    sh "oc start-build micro-auth-worker --from-file=. --follow"
  }
  stage('Deploy') {
    openshiftDeploy depCfg: 'micro-auth-worker', namespace: 'fbac'
    openshiftVerifyDeployment depCfg: 'micro-auth-worker', replicaCount: 1, verifyReplicaCount: true
  }
}