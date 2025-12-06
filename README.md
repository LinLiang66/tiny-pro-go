# Tiny Pro Go 

## 说明

此项目是 TinyPro 的基于GoLang的后端服务，支持在线配置菜单、路由、国际化，支持页签模式、多级菜单，支持丰富的模板类型，支持多种构建工具，功能强大、开箱即用！

## 快速上手

### 依赖安装

请选择任何一个你喜欢的包管理工具进行安装，这里使用 `go mod`：

```
bash
go mod tidy
```
### 启动开发环境

```
bash
go run main.go
```
## 目录结构

```
config            # 配置文件
controller        # 控制器
entity            # 实体类
impl              # 业务逻辑实现
middleware        # 中间件
routes            # 路由配置
service           # 服务层
utils             # 工具类
.gitignore        # Git 忽略文件
go.mod            # Go 模块文件
go.sum            # Go 模块校验文件
main.go           # 主程序入口
```
## 二次开发指南

1. 修改 `config/config.dev.yaml`为`config/config.yaml` 并文件中的数据库连接信息
2. 在 `src/controller` 目录下添加新的控制器
3. 在 `src/entity/dto` 目录下定义新的数据传输对象
4. 在 `src/impl` 目录下实现业务逻辑
5. 在 `src/routes` 目录下配置路由

## 遇到困难?

加官方小助手微信 opentiny-official，加入技术交流群