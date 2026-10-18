# TarsDify

MCP integration for Dify.

## Changelog

### v1.0.1
- Version bump with stability improvements

### v1.0.2
- JSON 配置文件支持（config.json）
- 环境变量与配置文件优先级合并
- 配置验证与错误提示

### v1.0.3
- 结构化请求日志中间件
- 自动追踪 ID（X-Trace-ID）
- 请求耗时与状态码记录

### v1.0.4
- 连接池优化（MaxIdleConns / MaxIdleConnsPerHost）
- 指数退避重试机制（默认3次）
- Dial 超时与 KeepAlive 调优

### v1.0.5
- 速率限制（令牌桶，/invoke 默认60/秒）
- 安全头增强（CSP, HSTS, Referrer-Policy）

### v1.0.6
- 深度健康检查（upstream 不可达时返回 degraded）
- /metrics 端点（请求计数 + 错误计数 + 运行时间）

### v1.0.7
- 请求体大小限制（默认 1MB）
- JSON Content-Type 校验中间件
- 标准化错误响应格式

### v1.0.8
- /version 端点（版本号 + 构建时间 + Git commit）
- 构建元数据支持

### v1.0.9
- 多环境配置支持（config.dev.json / config.staging.json / config.prod.json）
- LoadForEnv 自动按环境选择配置文件

### v1.1.0
- 自动更新检查 stub（/update 端点）
- GitHub Release 版本比对框架
- 当前版本与最新版本对比
