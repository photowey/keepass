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

- `github` -> `hellopass`
- `gitea` -> `hellopass`

同一个用户名可以出现在多个网站；真正区分条目的是 `alias`。

## 安全模型

- Vault 文件：`~/.keepass/keepass.kp`
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
keepass add github hellopass --uri https://github.com --note "personal" --tag code
keepass add gitea hellopass --uri https://gitea.example.com --note "work" --tag code
```

添加时：

- 如果你手动输入账号密码，CLI 会要求二次确认
- 如果你留空，系统会自动生成密码

列出记录：

```bash
keepass list
keepass list --tag code
keepass list --json
```

查看摘要：

```bash
keepass get github
keepass get gith
```

显式查看明文密码：

```bash
keepass get gith --reveal
keepass get gith --json
keepass get gith --json --reveal
keepass get gith --copy
keepass get gith --copy --copy-timeout 0
```

更新与删除：

```bash
keepass update github
keepass delete github
```

数据迁移与恢复：

```bash
keepass export --path ./entries.json
keepass import --path ./entries.json --conflict overwrite
keepass backup --path ./backup-bundle
keepass restore --path ./backup-bundle --force
```

凭据审计与轮换：

```bash
keepass audit --json
keepass rotate github --generate
```

使用当前 Argon2 参数重写 vault：

```bash
keepass rehash
```

查看当前配置：

```bash
keepass config
keepass config --json
```

审计本地 vault 健康状态：

```bash
keepass doctor
keepass doctor --json
```

## 自动化提示（非交互模式）

当 stdin 不是 TTY（脚本/CI/管道）时，为避免命令卡住，部分命令不会发起交互提示：

- `keepass add` 必须通过参数提供 `alias` 和 `username`（不会再交互询问）。
- `keepass update` 必须显式提供变更参数，例如 `--username`、`--password`、`--clear-uri`、`--clear-note`。
- `keepass delete` 在非交互模式下必须显式传入 `--yes`。
- 你也可以通过 `--non-interactive` 在 TTY 场景强制启用该行为。

## 退出码（Exit Codes）

- `1`：通用错误
- `2`：用法/参数错误
- `3`：未初始化（缺少 config/vault）
- `4`：解锁失败（master password 错误）

## Shell 自动补全

生成补全脚本：

```bash
keepass completion bash
keepass completion zsh
keepass completion fish
keepass completion powershell
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

生成器实际上支持任意非空字符表。

默认字符表刻意不包含大部分特殊符号，主要是为了兼容更多网站、Shell 场景以及手动输入。
如果你的使用环境明确要求更多符号，可以切换 `password_generator.preset`，或者直接自定义 `password_generator.alphabet`。

内置预设：

- `compatible`
  - 默认值，优先兼容更多网站和 Shell 场景
- `symbols`
  - 增加一组适中的特殊符号
- `strict-high-entropy`
  - 使用更大的混合字符表和更多符号

默认生成规则位于 `~/.keepass/keepass.config.json`：

```json
{
  "version": 1,
  "vault": {
    "path": "~/.keepass/keepass.kp",
    "format_version": 1
  },
  "security": {
    "argon2id": {
      "time": 3,
      "memory_kib": 262144,
      "threads": 4,
      "key_length": 32
    }
  },
  "password_generator": {
    "default_length": 21,
    "preset": "compatible"
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

## 发布完整性

GitHub Release 产物包含：

- `SHA256SUMS.txt` 中的逐文件 SHA256 校验和
- 用于构建来源证明的 GitHub Artifact Attestation
