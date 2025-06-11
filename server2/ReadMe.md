# 股票行情实时监控系统

一个基于Vue3 + KLineChart前端和Go语言后端的股票行情实时监控系统，支持WebSocket实时数据推送。

## 系统架构

- **前端**: Vue3 + KLineChart + WebSocket客户端
- **后端**: Go + Gorilla WebSocket + RESTful API
- **数据更新**: 每10秒自动更新股票数据
- **支持股票**: AAPL, GOOGL, MSFT, TSLA, AMZN

## 功能特性

✅ 实时K线图显示  
✅ 多股票切换  
✅ WebSocket实时数据推送  
✅ 响应式设计  
✅ 连接状态显示  
✅ 错误处理和自动重连  
✅ 模拟数据生成  
✅ CORS跨域支持  

## 快速开始

### 1. 后端部署

#### 安装依赖
```bash
# 创建项目目录
mkdir stock-monitor
cd stock-monitor

# 初始化Go模块
go mod init stock-monitor

# 安装依赖
go get github.com/gorilla/mux
go get github.com/gorilla/websocket
```

#### 运行后端服务
```bash
# 将go代码保存为main.go
go run main.go
```

服务器将在 `http://localhost:8080` 启动

#### API端点

| 方法 | 端点 | 描述 |
|------|------|------|
| GET | `/api/v1/stocks/{symbol}` | 获取股票历史数据 |
| GET | `/api/v1/symbols` | 获取所有可用股票符号 |
| GET | `/api/v1/health` | 健康检查 |
| WS | `/ws?symbol={symbol}` | WebSocket实时数据 |

#### 示例请求
```bash
# 获取苹果股票数据
curl http://localhost:8080/api/v1/stocks/AAPL

# 获取所有股票符号
curl http://localhost:8080/api/v1/symbols

# 健康检查
curl http://localhost:8080/api/v1/health
```

### 2. 前端部署

#### 方式一：直接运行HTML文件
将Vue前端代码保存为 `index.html`，直接在浏览器中打开即可。

#### 方式二：使用静态服务器
```bash
# 使用Python启动静态服务器
python -m http.server 3000

# 或使用Node.js的http-server
npx http-server -p 3000
```

然后访问 `http://localhost:3000`

### 3. 项目结构
```
stock-monitor/
├── main.go              # Go后端服务器
├── go.mod              # Go依赖管理
├── index.html          # Vue3前端页面
├── static/             # 静态文件目录（可选）
└── README.md           # 项目说明
```

## 配置说明

### 后端配置
- **端口**: 默认8080，可在main.go中修改
- **更新频率**: 每10秒，可在startRealTimeUpdate函数中调整
- **数据保持量**: 最多200个数据点，可调整
- **支持的股票**: 在basePrices map中配置

### 前端配置
- **API地址**: 在JavaScript中的API_BASE变量修改
- **WebSocket地址**: 在WS_BASE变量修改
- **图表配置**: 在initChart函数中自定义样式

## 数据格式

### 股票数据结构
```json
{
  "timestamp": 1640995200000,
  "symbol": "AAPL",
  "open": 150.00,
  "high": 152.00,
  "low": 149.00,
  "close": 151.50,
  "volume": 1000000
}
```

### WebSocket消息格式
```json
{
  "type": "update",
  "data": {
    "timestamp": 1640995200000,
    "symbol": "AAPL",
    "open": 150.00,
    "high": 152.00,
    "low": 149.00,
    "close": 151.50,
    "volume": 1000000
  }
}
```

## 开发指南

### 添加新股票
1. 在Go后端的`basePrices` map中添加股票符号和基础价格
2. 在`generateInitialData`函数的symbols数组中添加符号
3. 在前端的股票选择器中添加选项

### 自定义更新频率
修改Go后端中的ticker时间：
```go
sm.updateTicker = time.NewTicker(5 * time.Second) // 改为5秒更新
```

### 集成真实股票API
替换`generateNewData`函数，调用真实的股票数据API：
```go
func (sm *StockManager) fetchRealStockData(symbol string) StockData {
    // 调用真实API，如Alpha Vantage, Yahoo Finance等
    // ...
}
```

## 故障排除

### 常见问题

1. **WebSocket连接失败**
   - 检查后端服务是否正常运行
   - 确认防火墙设置
   - 检查浏览器控制台错误信息

2. **CORS错误**
   - 后端已配置CORS，如仍有问题请检查浏览器设置
   - 确保API地址正确

3. **数据不更新**
   - 检查WebSocket连接状态
   - 查看浏览器网络面板
   - 确认后端定时器正常工作

4. **图表显示异常**
   - 检查KLineChart库是否正确加载
   - 确认数据格式正确
   - 查看浏览器控制台错误

### 日志查看
后端日志会显示：
- WebSocket连接/断开信息
- 客户端订阅股票变更
- 错误信息

前端在浏览器控制台显示：
- WebSocket连接状态
- 数据更新信息
- 错误信息

## 扩展建议

1. **数据持久化**: 使用Redis或数据库存储历史数据
2. **用户认证**: 添加用户登录和权限管理
3. **多时间段**: 支持分钟、小时、日线切换
4. **技术指标**: 集成MA、MACD、RSI等技术指标
5. **价格预警**: 添加价格提醒功能
6. **移动端适配**: 优化移动设备显示效果

## 许可证

MIT License

## 贡献

欢迎提交Issue和Pull Request来改进项目！