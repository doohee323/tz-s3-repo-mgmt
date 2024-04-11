# tz-s3-repo-mgmt

-. Golang

```
- go
    download and install 
    https://golang.org
    
    mkdir -p ~/go
    cd ~/go

    mkdir bin pkg src
    mkdir -p src/bitbucket.org
    mkdir -p src/github.com/doohee323

    vi ~/.bash_profile
    export GOROOT=/usr/local/opt/go/libexec
    export GOPATH=~/go
    export PATH=$GOPATH/bin:.:$PATH
    source .bash_profile
```

-. Env

```
vi .env 
    DB_DRIVER=mysql
    DB_HOST=xxx
    DB_USER=xxx
    DB_PASSWORD=xxx
    DB_NAME=devops
    DB_PORT=3306
    SECRET_KEY=yourSecretString
    TOKEN_LIFESPAN=1
    S3REPOS=unity-repo.xxx,flutter-repo.xxx
    AWS_REGIONS=xxx

vi config
    [default]
    region = xxx
    output = json

vi credentials
    [default]
    aws_access_key_id = xxx
    aws_secret_access_key = xxx
    
check docker/s3-repo-mysql/sql/init.sql

```

-. build:

```
# checkout
    cd $GOPATH/src/github.com/doohee323
    git clone https://github.com/doohee323/tz-s3-repo-mgmt.git
    cd tz-s3-repo-mgmt

# add modules
    go get github.com/gin-gonic/gin
    go get golang.org/x/crypto/bcrypt
    go get github.com/dgrijalva/jwt-go
    go get gorm.io/gorm
    go get gorm.io/driver/mysql
    go get github.com/joho/godotenv
    
    export GOPRIVATE="github.com/doohee323"
    go mod init github.com/doohee323/tz-s3-repo-mgmt
    go mod tidy
    go run main.go
    go get -u github.com/gin-gonic/gin
    go get -u github.com/adrg/frontmatter
    go get -u github.com/dgrijalva/jwt-go
    go get -u github.com/gin-contrib/cors

# build
    go clean --cache
    go build
    ./tz-s3-repo-mgmt

# run on docker
    cd docker
    vi s3-repo/Dockerfile
    vi s3-repo-mysql/Dockerfile
    docker-compose -f docker-compose.yml up
    
open http://localhost:8080
```

-. Test:

```
# run Vue alone
    cd statics
    npm i
    npm run build
    npm run serve

open http://localhost:8080 or 
open http://localhost:8081
```
