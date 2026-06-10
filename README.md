# LLM Tracer 🚀

`LLM Tracer` 是一个轻量级、跨平台的本地大模型代理拦截与对话追踪记录工具。它能无感地拦截发往 OpenAI、Anthropic 等主流服务商的 API 请求（包括普通接口与流式 SSE 响应），将请求流、思维链（Thinking Process）及 Token 消耗等数据持久化到本地 SQLite 数据库中，并通过内置的极美暗黑风 Web 控制台进行可视化分析与级联调试。

---

## 🌟 核心特性

- **轻量单二进制部署**：前端 React 构建产物被完全嵌入 Go 二进制中，双击即可直接运行，包体积仅约 15MB。
- **系统托盘常驻 (Desktop Native)**：程序常驻系统托盘，提供一键打开控制台、重启代理及优雅退出等功能，启动时自动拉起浏览器展示控制台。
- **配置与数据库隐藏化 (`~/.llm_tracer`)**：所有配置文件（`config.json`）与本地 SQLite 数据库（`llm_tracer.db`）统一重定向至系统家目录的隐藏文件夹中，权限稳健、便于维护。
- **端口冲突自探测与Fallback**：动态探测代理接口（默认 `:1238`）与控制台接口（默认 `:56129`），若端口被占用则自动后移分配空闲端口，支持实例多开不冲突。
- **多路 Subagent 级联追踪**：支持追踪复杂的 Agent 级联调用（如子代理的多路并发请求），根据上下文指纹和语义匹配，将子请求完美归档至主干会话中。
- **内置一键配置指南**：控制台配置面板内置 Python/Node.js/环境变量 等多语言客户端接入指引，支持一键复制代码段。

---

## 📂 项目结构

```text
llm_tracer/
├── main.go            # 代理转发服务、SSE解析流扫描、系统托盘管理与 Web 控制台托管
├── parser.go          # OpenAI / Responses / Anthropic 协议解析与会话指纹提取
├── db.go              # SQLite 数据库读写、归并统计与级联关系恢复
├── main_test.go       # 离线集成测试与级联回归测试用例
├── Makefile           # 一键编译、构建与运行指令集合
├── pack_mac.sh        # 一键打包 macOS 桌面端应用包 (.app) 自动化脚本
├── logo.png           # 专为此工具生成的 App 高清原画图标（用于 macOS 打包）
├── static/            # 前端 React 打包的静态资源目录（Go embed 静态载入）
└── frontend/          # React + Vite + TailwindCSS 独立前端控制台源码
```

---

## 🛠️ 构建与运行

### 准备工作
确保本地已安装：
- **Go** (1.20 或更高版本，建议 1.26+)
- **Node.js & npm** (用于前端 React 编译)

### 1. 一键构建前端与后端
在项目根目录下，直接运行以下命令：
```bash
make build
```
这会依次执行：
1. `npm run build` (将前端产物编译至 `static/` 目录)
2. `go build` (编译生成 `llm_tracer_web` 单二进制文件)

### 2. 运行代理服务
```bash
./llm_tracer_web
```
启动后：
- 默认代理服务将监听在 `http://localhost:1238`。
- 默认控制台服务将监听在 `http://localhost:56129`，并会自动在默认浏览器中拉起控制台。
- 常驻系统托盘将出现在系统顶部菜单栏（macOS）或右下角任务栏（Windows/Linux）。

> [!NOTE]
> 如果您想自定义端口，可以通过命令行参数覆盖：
> `./llm_tracer_web -listen :8080 -console :3000`

### 3. 运行集成测试
运行后端的回归测试：
```bash
make test
```

### 4. macOS 原生桌面应用打包 (.app)
如果您在 macOS 上运行，并希望获得完美的系统原生体验（**双击直接静默运行、不弹出终端命令框，且拥有独立的 Dock 及状态栏图标**），可运行项目内置的自动打包脚本：
```bash
chmod +x pack_mac.sh
./pack_mac.sh
```
运行后，当前目录将自动生成 **`LLMTracer.app`** 原生应用程序包：
- **高清视网膜图标**：打包脚本会自动读取根目录下的 `logo.png` 原画，利用 macOS 系统内置的 `sips` 与 `iconutil` 图像引擎无损制作符合 Mac 各种分辨率标准的 `icon.icns` 图标集，注入 App 的资源中。
- **开箱即用运行**：双击生成的 `LLMTracer.app` 即可在后台静默运行（终端不弹窗），并自动在浏览器中唤起控制台。
- **系统集成**：您可以将 `LLMTracer.app` 拖入系统的 `/Applications`（应用程序）目录，即可直接在 Launchpad（启动台）或 Spotlight 搜索中一键拉起。

---

## 🚀 客户端接入指引

要通过该代理拦截和记录您的 LLM 请求，仅需将您开发库的 `baseURL` 重定向到本地代理即可（ApiKey 可以填写任意非空字符串）。

### 方式一：设置系统环境变量 (推荐)
最无感的配置方式，能够使所有默认读取环境变量的 SDK 自动生效：
```bash
export OPENAI_BASE_URL=http://localhost:1238/v1
export ANTHROPIC_BASE_URL=http://localhost:1238/v1
```

### 方式二：Python (openai-python)
```python
import openai

# 将 base_url 指向 LLM Tracer 本地代理
openai.base_url = "http://localhost:1238/v1"
# apiKey 可填任意非空字符串
openai.api_key = "anything"

response = openai.chat.completions.create(
    model="gpt-4o",
    messages=[{"role": "user", "content": "你好，这是一次测试！"}]
)
print(response.choices[0].message.content)
```

### 方式三：Node.js / TypeScript
```javascript
import OpenAI from 'openai';

const openai = new OpenAI({
  baseURL: 'http://localhost:1238/v1',
  apiKey: 'anything'
});

const completion = await openai.chat.completions.create({
  model: 'gpt-4o',
  messages: [{ role: 'user', content: 'Hello!' }],
});
```

---

## ⚖️ 授权协议

本项目采用 [MIT License](LICENSE) 开源。
