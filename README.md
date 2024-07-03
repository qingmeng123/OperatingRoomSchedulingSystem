
<div align="center">

# 手术室排程系统


</div>



## 实现功能

1. 用户管理权限设置
2. 用户密码加密存储
3. 手术室资源管理系统 
4. 手术室排程缓存处理
5. 排程信息甘特图展示
6. 容器化部署
7. 跨域 cors 设置


## 技术栈

- golang
    - Gin 
    - gorm
    - go-redis
    - jwt-go
    - scrypt
    - go-ini
- JavaScript
    - vue
    - vue cli
    - vue router
    - element-plus
    - axios
    - gantt-schedule-timeline-calendar
- MySQL version:8.0.21

## 项目预览




## 目录结构

```shell
├─  .gitignore
│  go.mod // 项目依赖
│  go.sum
│  LICENSE
│  README.md
│  Dockerfile
│  docker-compose.yml        
├─api         
├─cache  //缓存
├─cmd   //项目启动入口
├─config // 项目配置入口   
├─dao  // 数据库备份文件（初始化）
├─deploy 
├─model // 数据模型层
├─service
├─tool       
└─web // 前端开发源码
    
```

## 运行&&部署

1. 克隆项目
```shell
git clone https://github.com/qingmeng123/OperatingRoomSchedulingSystem.git
```

2. 转到下面文件夹下

```shell
cd yourPath/OperatingRoomSchedulingSystem
```

3. 安装依赖

```
go mod tidy
```

4. 初始化项目配置config.ini

```ini
./config/config.ini

[server]
# debug 开发模式，release 生产模式
AppMode = debug
HttpPort = :8080
JwtKey = duryun
PageSize=10
EndTime=23

[database]
Db = mysql
DbHost = 127.0.0.1
DbPort = 3306
DbUser = root
DbPassWord =@XUEHUI.
DbName = scheduling_system
```

5. 在deploy中将sql文件导入数据库

   推荐navicat或者其他sql管理工具导入

6. 启动项目

```shell
 go run main.go
```

此时，项目启动，你可以访问页面

```shell
首页
http://localhost:3000

```

