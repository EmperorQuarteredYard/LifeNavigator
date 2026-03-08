# 项目简介

~~本项目系人生规划系统，是作者在用了Notion，小青帐等工具后发现他们的定位不符合我的需求，故欲自己开发这样一款应用以：\
统一管理目标追踪、时间与金钱预算分配，并在将来实现跨平台服务的应用~~\
**这个项目就是一坨**
# 项目技术栈
- 多平台适配：flutter框架
- 服务端：gin+MySQL+JWT \
# 如何开始
在/backend下新建一个docker-compose.yml，内容如下所示：
```yaml
version: '3.8'

services:
  web:
    build:
      context: .
    ports:
      - "5083:5083"               # 映射宿主机端口到容器端口（根据你的应用端口调整）
    environment:
      - DB_HOST=mysql              # MySQL 服务名
      - DB_PORT=3306
      - DB_USER=root
      - DB_PASSWORD=CCEl8shKS1Byu0QFbSwy0   # 请修改为强密码
      - DB_NAME=MySQLDatabase               # 数据库名
      - JWT_ACCESS_SECRET=S4f7XnZh5zha9dLFbbjj1TpT77LBft8aDc4DOj34hluLaEZaG8pzERXKuZQQMqCTEZlLiZgtv_RAIUXvyzKFkjjhbA5Nr9yRV_JuibeIl9jChwzC7Of8FGg7HHXdiTTf #这里是随机生成的字符可以不用管
      - JWT_REFRESH_SECRET=yONJ7qBqz9mDGSufSQZs30bT5h2lumRBtaK_M8XD766NkSU73xK8uWeat1MfpUmWPTFo_Xrg7ybKHConVxYcyfOrOgboWdrhBdnYfBe6bUoxLmed9OBCKddymfgVU8t8 #这里是随机生成的字符可以不用管
      - JWT_ISSUER=LifeNavigatorDemo
      - JWT_ACCESS_TTL=7200000000000      #2小时
      - JWT_REFRESH_TTL=259200000000000   #3天
      - GLM_TOKEN=这里硬编码了GLM的Token
    depends_on:
      - mysql
    networks:
      - app-network

  mysql:
    image: mysql:9.6
#    ports:
#      - "3306:3306"                # 可选：暴露到宿主机方便调试
    environment:
      MYSQL_ROOT_PASSWORD: CCEl8shKS1Byu0QFbSwy0   # root 密码
      MYSQL_DATABASE: LifeNavigatorData                 # 数据库
      MYSQL_USER: user                       # 普通用户
      MYSQL_PASSWORD: xPBPH4UzbQ           # 普通用户密码
    volumes:
      - mysql-data:/var/lib/mysql            # 持久化数据
    networks:
      - app-network
    command:
      - --character-set-server=utf8mb4
      - --collation-server=utf8mb4_unicode_ci
    healthcheck:
      test: [ "CMD", "mysqladmin", "ping", "-h", "localhost", "-u", "root", "-p$$MYSQL_ROOT_PASSWORD" ]
      interval: 5s       # 每 5 秒检查一次
      timeout: 3s        # 每次检查超时 3 秒
      retries: 5         # 连续失败 5 次后标记为 unhealthy
      start_period: 10s  # 启动后等待 10 秒再开始检查（给 MySQL 初始化时间）

networks:
  app-network:
    driver: bridge

volumes:
  mysql-data:
```
在/backend下运行 
```bash
docker-compose up --build
```
即可