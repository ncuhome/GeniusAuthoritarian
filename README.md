# Genius Authoritarian

统一鉴权中心，提供一站式连携登录服务

## :wrench: 程序调试

首先需要运行  `mysql`、`redis`、`ldap`

预先准备：

+ 在 web 目录中运行 `pnpm run build` 以生成站点静态文件。这么做不是为了生成用于调试的站点，只是为了让 go embed 不报错
+ 在 Goland 运行配置添加环境变量 `DEV_MODE`，值为 `TRUE`
+ 运行一次，此时项目根目录会出现配置文件，填好配置项

调试时：

+ 在 web 目录运行 `pnpm run dev`
+ 在 Goland 启动运行配置项