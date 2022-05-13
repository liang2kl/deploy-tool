# deploy-tool

由于软件工程课程平台部署体验极差，为了加快部署速度以更快地同步前后端进度，我在个人服务器上使用 docker-compose 完整部署了我们的项目。这个工具的目的是方便从远端部署最新代码，而无需使用 ssh 手动部署。例如，我们想要使用 `dev` 分支更新 `backend`，可以直接发一个 HTTP 请求：

```
curl https://example.com/update/backend/dev
```

然后服务器便会从 remote 拉取 `dev` 上最新的代码，并使用 docker-compose 重新部署。

## Prerequisites

- 需要将软工 GitLab 上的仓库克隆到本地
- 需要设置好访问 remote 的用户名和密码
- 需要配置 docker-compose

其中，`docker-compose.yml` 文件应包含需要更新的所有服务。例如，项目包含前端和后端：

```yml
version: "3.9"
services:
  backend:
    build: ../backend
    ...

  frontend:
    build: ../frontend
    ...
```

另外，对于 Linux 系统，当前用户需要在 `docker` 用户组中，从而可以用非 root 身份运行 docker 命令：

```
sudo usermod -aG docker $USER
```

## Usage

创建 `config.yml`：

```yaml
hostname: '127.0.0.1'
port: 8100
# docker-compose 配置文件路径
docker-compose-file: ../deploy/docker-compose.yml
projects: # docker-compose service 名称到项目 git 目录的映射
  frontend: ../frontend
  backend: ../backend
script: ./update.sh # 更新使用的脚本
interval: 60s # 两次更新之间的最小时间间隔
```

编译（需要 go 1.17）：

```
go mod download
go build -o bot
```

运行：

```
./bot
```

查看 log：

```
curl https://example.com/log/<service_name>
```

发送更新请求：

```
curl https://example.com/update/<service_name>/<branch>
```

接收到请求时，此程序会根据 `projects` 中给出的 git 目录拉取给定 branch，并使用 docker-compose 重新 build 并部署对应的服务。

