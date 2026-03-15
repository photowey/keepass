# `keepass`

[English](./README.md) | 中文

`keepass` 是一个使用 Go 编写的本地优先命令行密码管理器。

设计目标只有三个：

- `安全优先`：用户只需要记住一个 `master password`
- `简单好用`：命令短、参数少、默认交互式提示
- `快速查询`：先精确匹配，再做唯一前缀匹配

## 它管理什么

每条记录由唯一的 `alias` 标识，并可包含：

- `username`
- `password`
- `uri`
- `note`
- `tags`

例如：

- `github` -> `abc`
- `gitea` -> `abc`

同一个用户名可以出现在多个网站；真正区分条目的是 `alias`。

## 安全模型

- Vault 文件：`~/.keepass/default.kp`
- 配置文件：`~/.keepass/keepass.config.json`
- 配置文件只保存非敏感设置
- 密码条目统一保存在加密 vault 中
- vault 文件必须带 `format_version`
- 不支持的文件版本直接失败，不猜测、不降级
- 初始化和解锁都必须输入 `master password`
- 明文密码不会落盘

## 快速开始

初始化密码库：

```bash
keepass init
```

添加记录：

```bash
keepass add github abc --uri https://github.com --note "personal" --tag code
keepass add gitea abc --uri https://gitea.example.com --note "work" --tag code
```

添加时：

- 如果你手动输入账号密码，CLI 会要求二次确认
- 如果你留空，系统会自动生成密码

列出记录：

```bash
keepass list
keepass list --tag code
```

查看摘要：

```bash
keepass get github
keepass get gith
```

显式查看明文密码：

```bash
keepass get gith --reveal
```

更新与删除：

```bash
keepass update github
keepass delete github
```

查看当前配置：

```bash
keepass config
```

## 别名解析规则

查询规则固定为：

1. 先做精确匹配
2. 再做唯一前缀匹配
3. 如果前缀有歧义，直接报错

例如：

- `keepass get github` -> 精确命中
- `keepass get gith` -> 唯一前缀命中
- `keepass get gi` -> 如果同时存在 `github` 和 `gitea`，则直接失败

## 自动生成密码

如果你不手动输入账号密码，`keepass` 会使用安全随机源自动生成密码。

默认生成规则位于 `~/.keepass/keepass.config.json`：

```json
{
  "version": 1,
  "vault": {
    "path": "~/.keepass/default.kp",
    "format_version": 1
  },
  "security": {
    "argon2id": {
      "time": 3,
      "memory_kib": 65536,
      "threads": 4,
      "key_length": 32
    }
  },
  "password_generator": {
    "default_length": 21,
    "alphabet": "ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz23456789-_"
  }
}
```

## 测试

当前测试覆盖包括：

- 配置校验测试
- 密码生成器测试
- vault 格式与加解密测试
- 管理器规则测试
- 命令级闭环测试
- vault 解码 fuzz 入口

执行：

```bash
GOCACHE=/tmp/go-cache go test ./...
```
