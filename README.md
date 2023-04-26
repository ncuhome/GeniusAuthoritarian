# Genius Authoritarian

统一鉴权中心，提供一站式连携登录服务

目前的工作原理为，自动从飞书同步通讯录，使用电话号码作为身份标记。暂不支持从钉钉同步或混合同步。

## :wrench: 项目完整本地调试

首先准备环境：  `mysql`、`redis`、`ldap`

ldap 需要配置 `DOMAIN` 为 `ncuos.com`，`ORGANISATION` 为 `NCUHOME`

预先准备：

+ 在 web 目录中运行 `pnpm run build` 以生成站点静态文件。这么做不是为了生成用于调试的站点，只是为了让 go embed 不报错
+ 在 Goland 运行配置添加环境变量 `DEV_MODE`，值为 `TRUE`
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