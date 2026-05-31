# API Gateway Config Validator 🛡️

> AI API 中转站 / 网关配置校验工具 | Validate and fix your AI gateway configs

一键检查你的 **One API / New API / relay-server** 配置是否正确，避免启动失败和生产事故。

## 🎯 解决的问题

```
你写了一个 provider YAML 配置 → 启动报错 → 排查半天发现 base_url 写错了
或者 API key 没填，或者 models 数组是空的...
```

**现在不用了。** 一条命令帮你检查所有配置问题。

## ✨ 功能

- ✅ **配置校验** — 检查 server.addr、admin.addr、auth.keys、provider 字段是否完整
- ✅ **目录批量校验** — 一次检查整个 providers/ 目录的所有 YAML
- ✅ **Provider 模板** — 内置 DeepSeek、火山引擎、通义千问、智谱、百度等模板
- ✅ **环境变量检测** — 自动识别 `${VAR}` 语法，提示未设置的变量
- ✅ **多文档 YAML** — 支持 `---` 分割的 YAML 文件
- ✅ **项目初始化** — 一键生成标准配置骨架
- ✅ **单文件二进制** — 下载即用，无需安装任何依赖

## 🚀 快速开始

```bash
# 1. 下载 (Windows/Linux/Mac)
# 见下方"下载"

# 2. 校验配置文件
api-validator validate config.yaml

# 3. 校验整个配置目录
api-validator validate ./providers/

# 4. 生成新项目骨架
api-validator init my-gateway
```

## 📦 下载

| 平台 | 下载 |
|------|------|
| Windows x64 | [api-validator-windows-amd64.exe](./api-validator-windows-amd64.exe) |
| Linux x64 | （需自行编译：`go build -o api-validator .`） |
| macOS x64 | （需自行编译：`go build -o api-validator .`） |
| macOS ARM | （需自行编译：`go build -o api-validator .`） |

## 📖 使用指南

### 校验单个配置文件

```bash
api-validator validate ./config.yaml
```

输出示例：
```
🔍 Config: ./config.yaml
  ✓ PASS

  ✓ No issues found. Config looks good!
```

### 校验发现错误时

```bash
api-validator validate ./bad-config.yaml
```

输出示例：
```
🔍 Config: ./bad-config.yaml
  ✖ FAIL

  ✖ Errors:
    • providers[0](deepseek): base_url is required for type openai

  ⚠ Warnings:
    • providers[0](deepseek): no api_key set
    • providers[0](deepseek): no models configured
```

### 校验整个目录

```bash
api-validator validate ./providers/
```

### 查看 Provider 模板

```bash
api-validator providers
```

内置 6 个常用模板：DeepSeek、火山引擎、通义千问、智谱 GLM、百度文心、Mock。

### 生成新项目

```bash
api-validator init my-project
cd my-project
api-validator validate config.yaml
```

## 🔧 支持的 Provider 类型

| 类型 | 说明 |
|------|------|
| `openai` | OpenAI 兼容的 API (DeepSeek/火山/千问/智谱/百度...) |
| `mock` | 本地 Mock，无需 API Key |

## 📄 配置参考

```yaml
# config.yaml
server:
  addr: ":8080"          # 公开 API 端口
  timeout: "120s"

admin:
  addr: ":8081"          # 管理 API 端口

auth:
  enabled: true
  keys:
    - "sk-your-key"

rate_limit:
  enabled: true
  rate: 10               # 每秒请求数
  capacity: 20           # 突发容量

billing:
  enabled: false
  currency: "CNY"

providers:
  - name: deepseek
    type: openai
    base_url: https://api.deepseek.com
    api_key: ${DEEPSEEK_API_KEY}
    models:
      - deepseek-chat
      - deepseek-reasoner
```

## 🏗️ 构建方法

```bash
git clone https://github.com/.../api-validator.git
cd api-validator
go build -o api-validator .
```

## 📜 许可证

MIT License
