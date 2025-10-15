# golang restful api的一个示例
## 去除了业务相关内容

## cmd目录是执行服务启动/模型生成/定时任务的入口
## config目录存放配置/编译发布相关内容
## internal目录为服务业务逻辑
## pkg是框架相关内容,可以独立为一个仓库
## tests目录是tdd的体现,存放服务中的各种测试

### 业务编写中,主要依赖appContext,可以专注于业务逻辑,appContext中包含了外部资源操作,数据库,redis,文件等.