# Genius Authoritarian

统一鉴权中心，提供一站式连携登录服务

目前的工作原理为，自动从飞书同步通讯录，使用电话号码作为身份标记。暂不支持从钉钉同步或混合同步。

本项目旨在解决：

1. 更换办公软件会导致增加大量迁移工作
2. 旧 ldap 已无法访问，无法管理，我们需要新的全自动 ldap 管理系统
3. 现在没有一个可用的导航或工作台，进入任何控制面板都要口口相传域名
4. 希望能一并解决研发的系统需要挨个分发账号的问题

## :wrench: 项目完整本地调试

前后端可分开调试，前端嵌入并没有任何注入

首先准备环境：`mysql`、`redis`，建议映射预发布环境

构建 go 程序时添加 `-tags dev` 可以切换到调试模式

预先准备：

+ Goland 设置 => Go => 构建标记 => 自定义标记  => 填入 `dev`
+ 添加运行配置，勾选“使用所有自定义构建标记”
+ 运行一次，此时项目根目录会出现配置文件，填好配置项

调试时：

+ 在 web 目录运行 `pnpm run dev`
+ 在 Goland 启动运行配置项

## :gear: 使用

### 为开放服务非入侵式添加鉴权

目前只支持注入不跨域集群服务，见 [GeniusAuthoritarianGate](https://github.com/ncuhome/GeniusAuthoritarianGate)

将会占用路径 `/login`，cookie `token`

### 前端调用

需要先创建一个页面接收回调信息，回调页面需要处理 Query param `token`，然后用这个 token 请求自己项目后端的登录接口拿第二个 `token`

调用示例：

```
window.open('https://v.ncuos.com/?target=https://example.ncuos.com/login', '_self')
```

其中 `target` 为前端回调页面 url，登录系统会对域名进行白名单校验，可以附带自定义 path、query 或 hash

### 后端调用

Golang 项目可以选择直接调用 Client [ncuhome/GeniusAuthoritarianClient](https://github.com/ncuhome/GeniusAuthoritarianClient)

先接收到前端传来的 `token`，此为一次性身份校验令牌

然后请求校验接口：

POST `https://v.ncuos.com/api/v1/public/login/verify`

Form:

| key    | type     | required |
|--------|----------|----------|
| token  | string   | √        |
| groups | []string | x        |

当 groups 为空时，接口会返回用户所在所有组。当 groups 不为空时，接口返回匹配结果，无匹配值将会导致登陆失败

groups 的值参考 [departments.go](/pkg/departments/departments.go)

成功：

```json5
{
  "code": 0,
  "data": {
    "groups": [
      "研发",
      "中心"
    ]
  }
}
```

失败：

```json5
{
    "code": 5,
    "msg": "身份校验失败，权限不足"
}
```

## :card_index_dividers: TODO

+ [x] ~~对接飞书登录~~
+ [x] ~~对接钉钉登录~~
+ [x] ~~添加个人资料管理~~
+ [ ] 添加服务鉴权密钥对
+ [ ] 添加导航支持
+ [ ] 接入 ldap