# deploy-tool

由于软件工程课程平台部署体验极差，为了加快部署速度以更快地同步前后端进度，我在个人服务器上使用 docker-compose 完整部署了我们的项目。这个工具的目的是方便从远端部署最新代码，而无需手动 ssh。例如，我们想要更新 `backend` 的 `dev` 分支，可以直接发一个 HTTP 请求：

```
curl https://example.com/update/backend/dev
```

然后服务器便会从 remote 拉取 `dev` 上最新的代码，并使用 docker-compose 重新部署。

## Usage

创建 `config.yml`：

```yaml
hostname: '127.0.0.1'
port: 8100
docker-compose-dir: ../deploy # docker-compose.yml 所在目录
projects: # docker-compose service 名称到项目 git 目录的映射
  service: ../service
  push-service: ../service
  test-service: ../service
  frontend: ../frontend
  backend: ../backend
script: ./update.sh # 更新使用的脚本
interval: 60s # 两次更新之间的最小时间间隔
```

接收到请求时，此程序会根据 `projects` 中给出的 git 目录拉取给定 branch，并使用 docker-compose 重新 build 并部署对应的服务。

编译：

```
go mod download
go build -o bot
```

运行：

```
./bot
```
