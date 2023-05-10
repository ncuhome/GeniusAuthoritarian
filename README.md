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

+ 添加 `GOPRIVATE` 环境变量 `github.com/ncuhome`
+ Goland 设置 => Go => 构建标记 => 自定义标记  => 填入 `dev`
+ 添加运行配置，勾选“使用所有自定义构建标记”
+ 运行一次，此时项目根目录会出现配置文件，填好配置项

调试时：

+ 在 web 目录运行 `pnpm run dev`
+ 在 Goland 启动运行配置项

## :gear: 使用

需要先申请鉴权密钥对，相关功能开发中，暂未开放申请

### 为开放服务非入侵式添加鉴权

目前只支持注入不跨域集群服务，见 [GeniusAuthoritarianGate](https://github.com/ncuhome/GeniusAuthoritarianGate)

将会占用路径 `/login`，cookie `token`

### 前端调用

需要先创建一个页面接收回调信息，回调页面需要处理 Query param `token`，然后用这个 token 请求自己项目后端的登录接口获取正式登录状态

调用示例：

```
window.open('https://v.ncuos.com/?appCode=YourAppCode', '_self')
```

其中 `target` 为前端回调页面 url，登录系统会对域名进行白名单校验，可以附带自定义 path、query 或 hash

### 后端调用

Golang 项目可以选择直接调用 Client [ncuhome/GeniusAuthoritarianClient](https://github.com/ncuhome/GeniusAuthoritarianClient)

先接收到前端传来的 `token`，此为一次性身份校验令牌

然后请求校验接口：

POST `https://v.ncuos.com/api/v1/public/login/verify`

Form:

| key       | type   | required | desc      |
|-----------|--------|----------|-----------|
| token     | string | √        |           |
| appCode   | string | √        |           |
| timeStamp | int64  | √        | unix 时间，秒 |
| signature | string | √        | 请求签名      |

`signature` 的计算方法是，在表单对象中加入 `appSecret`，去掉 `signature`，再将整个对象按键名排序，将键名和键值用 `=` 连接，不同项中间用 `&` 连接之后得到一个字符串，如：`key1=value1&key2=value2`，计算 `sha256` 值。

注意不要把 `appSecret` 当表单值放入请求传出。

成功：

```json5
{
  "code": 0,
  "data": {
    "name": "孙翔宇",
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