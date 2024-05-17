# Genius Authoritarian

统一授权系统，提供一站式连携登录服务

目前的工作原理为，自动从飞书同步通讯录，使用电话号码作为身份标记。暂不支持从钉钉同步或混合同步。

本项目旨在解决：

1. 更换办公软件会导致增加大量迁移工作
2. 旧 ldap 已无法访问，无法管理，我们需要新的全自动 ldap 管理系统
3. 现在没有一个可用的导航或工作台，进入任何控制面板都要口口相传域名
4. 希望能一并解决研发的系统需要挨个分发账号的问题

## :gear: 使用

需要先申请鉴权密钥对，请后端同学分别在 [预发布版后台 (v.ncuhome.club)](https://v.ncuhome.club) 和 [生产版后台 (v.ncuos.com)](https://v.ncuos.com) 创建相应应用用于测试和上线。应用一经创建就会显示在导航栏，请谨慎操作。创建应用时回调地址可以带自定义参数

具体开发文档请转到 [wiki](https://github.com/ncuhome/GeniusAuthoritarian/wiki)

## :star2: 主要特性

+ 支持三种粒度的 Token 撤销，时间复杂度为 O(1)：单个 Token、某用户所有 Token、某应用所有 Token
+ 高性能分布式 Token 解析：在支持 Token 撤销的基础上，应用校验统一授权签发的 Token 时不用和统一授权后端通信
+ 高可用分布式部署：统一授权后端支持高可用分布式部署
+ 双边 SSL 认证：应用与统一授权后端使用 rpc 通信时支持使用 mTLS 进行双边认证
+ 先进的无密码认证：在第三方 SSO 的基础上增加基于 Totp 的二步校验以及标准通行密钥认证
+ ed25519 密钥对分发：支持自动为研发同学生成 ed25519 SSH 密钥对，以及离职自动销毁等

## :wrench: 项目完整本地调试

前后端可分开调试，前端嵌入并没有任何注入

首先准备环境：`mysql`、`redis`，建议映射预发布环境

构建 go 程序时添加 `-tags="dev"` 可以切换到调试模式

预先准备：

+ 添加 `GOPRIVATE` 环境变量 `github.com/ncuhome`
+ （可选） 设置 => Go => 构建标记 => 自定义标记  => 填入 `dev`
+ 添加运行配置，勾选 “使用所有自定义构建标记” 或在 ”Go 工具实参“ 栏填入 `-tags="dev"`。运行种类为目录 `cmd\core`
+ 运行一次，此时项目根目录会出现配置文件，填好配置项

调试时：

+ 在 web 目录运行 `pnpm run dev`
+ 在 Goland 启动运行配置项

## :children_crossing: 项目子程序说明

Core 程序为 standalone 程序，包含所有登录逻辑、接口和功能，支持多实例并行运行。编程时请注意对多实例并行进行兼容，善用 mysql 锁、redis 锁与 redis 订阅/发布。api 为 80 端口，ssh rpc 为 81 端口，refreshToken 为 82 端口

sshDev 为 ssh server（rpc client），通过 rpc 从 core 获取用户名与公钥以此建立 ssh 服务，可以部署在任意地方，ssh 端口为默认端口
