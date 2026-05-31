# API 配置校验器 🛡️

> **还在为写错 YAML 配置排查半天吗？一条命令搞定。**

[![Release](https://img.shields.io/github/v/release/18296023612/api-validator)](https://github.com/18296023612/api-validator/releases)
[![Go](https://img.shields.io/badge/Go-1.21%2B-blue)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-green)](LICENSE)

---

## 📺 为什么需要这个工具？

```
你写了一个 provider 配置：

  providers:
    - name: deepseek
      type: openai
      models:
        - deepseek-chat

启动 → 报错 → 排查 → 发现 base_url 忘写了 😤

又或者：

  api_key: ${SOME_VAR}

结果 SOME_VAR 忘了设环境变量，服务拿着空字符串去调 API → 401
```

**这个工具就是来解决这个问题的。** 在你启动服务之前，先跑一条命令检查配置有没有问题。

适用场景：**One API / New API / 自建 AI 中转站**

---

## ✨ 功能一览

| 功能 | 说明 |
|------|------|
| ✅ **配置校验** | 检查 server.addr / admin.addr / auth.keys / provider 字段是否完整 |
| ✅ **目录批量校验** | 一次检查整个 `providers/` 目录 |
| ✅ **Provider 模板** | 内置 DeepSeek / 火山引擎 / 通义千问 / 智谱 / 百度 / Mock |
| ✅ **环境变量检测** | 自动发现 `${VAR}` 语法，提示哪些变量没设 |
| ✅ **多文档 YAML** | 支持 `---` 分隔的复杂文件 |
| ✅ **项目初始化** | `api-validator init my-project` 一键生成标准配置 |
| ✅ **零依赖** | 单文件二进制，下载即用 |

---

## 🚀 5 秒上手

### 1️⃣ 下载

从 [Releases](https://github.com/18296023612/api-validator/releases) 下载最新版：

| 平台 | 下载 |
|------|------|
| **Windows x64** | [api-validator-windows-amd64.exe](https://github.com/18296023612/api-validator/releases/download/v1.0.0/api-validator-windows-amd64.exe) |
| **Linux x64** | `go build -o api-validator .`（一条命令编译） |
| **macOS** | `go build -o api-validator .`（一条命令编译） |

### 2️⃣ 校验配置文件

```bash
# 校验单个文件
api-validator validate config.yaml

# 🔍 Config: ./config.yaml
#   ✓ PASS
#   ✓ No issues found. Config looks good!

# 校验目录
api-validator validate ./providers/
```

### 3️⃣ 发现错误时

```bash
api-validator validate ./bad-config.yaml

# 🔍 Config: ./bad-config.yaml
#   ✖ FAIL
#   ✖ Errors:
#     • providers[0](deepseek): base_url is required for type openai
#   ⚠ Warnings:
#     • providers[0](deepseek): no api_key set
#     • providers[0](deepseek): no models configured
```

红色错误 ❌ = 必须修复，黄色警告 ⚠ = 建议处理。

### 4️⃣ 生成新项目骨架

```bash
api-validator init my-gateway
cd my-gateway
api-validator validate config.yaml   # 一键校验
```

---

## 💡 使用技巧

### 嵌入 CI/CD 流水线

```yaml
# GitHub Actions 示例
- name: 校验配置
  run: |
    ./api-validator validate config/providers/
```

### 配合 relay-server / One API 使用

```bash
# 写配置 → 校验 → 启动
api-validator validate ./providers/ && relay-server
```

---

## 🔧 支持的 Provider 类型

| 类型 | 说明 |
|------|------|
| `openai` | OpenAI 兼容格式（DeepSeek / 火山 / 千问 / 智谱 / 百度...） |
| `mock` | 本地 Mock 测试，无需 API Key |

## 📄 配置参考（标准模板）

```yaml
server:
  addr: ":8080"
  timeout: "120s"

admin:
  addr: ":8081"

auth:
  enabled: true
  keys:
    - "sk-your-key"

rate_limit:
  enabled: true
  rate: 10
  capacity: 20

providers:
  - name: deepseek
    type: openai
    base_url: https://api.deepseek.com
    api_key: ${DEEPSEEK_API_KEY}
    models:
      - deepseek-chat
      - deepseek-reasoner
```

---

## 🏗️ 自行编译

```bash
git clone https://github.com/18296023612/api-validator.git
cd api-validator
go build -o api-validator .
```

需要 Go 1.21 或更高版本。

---

## 🧪 测试

```bash
go test -v ./...
```

11 个单元测试全部通过 ✅

---

## 📜 许可证

MIT License

---

## ⭐ 支持这个项目

如果这个工具帮到了你，欢迎：

- ⭐ **Star 这个仓库**（让更多人看到）
- 🐛 **提 Issue**（报 bug 或建议新功能）
- ☕ **请我喝咖啡**（扫下方赞赏码）

> 国内开发不易，一个 Star 就是最大的支持 🙏
