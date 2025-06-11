package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// StockData 股票数据结构
type StockData struct {
	Timestamp int64   `json:"timestamp"`
	Symbol    string  `json:"symbol"`
	Open      float64 `json:"open"`
	High      float64 `json:"high"`
	Low       float64 `json:"low"`
	Close     float64 `json:"close"`
	Volume    int64   `json:"volume"`
}

// StockResponse API响应结构
type StockResponse struct {
	Symbol    string      `json:"symbol"`
	Data      []StockData `json:"data"`
	Timestamp int64       `json:"timestamp"`
}

// WebSocketMessage WebSocket消息结构
type WebSocketMessage struct {
	Type string    `json:"type"`
	Data StockData `json:"data"`
}

// StockManager 股票数据管理器
type StockManager struct {
	mu           sync.RWMutex
	stockData    map[string][]StockData
	clients      map[*websocket.Conn]string
	basePrices   map[string]float64
	upgrader     websocket.Upgrader
	updateTicker *time.Ticker
}

// NewStockManager 创建新的股票管理器
func NewStockManager() *StockManager {
	sm := &StockManager{
		stockData:  make(map[string][]StockData),
		clients:    make(map[*websocket.Conn]string),
		basePrices: make(map[string]float64),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // 允许跨域
			},
		},
	}

	// 初始化基础价格
	sm.basePrices["AAPL"] = 150.0
	sm.basePrices["GOOGL"] = 2800.0
	sm.basePrices["MSFT"] = 330.0
	sm.basePrices["TSLA"] = 200.0
	sm.basePrices["AMZN"] = 3200.0

	// 生成初始数据
	sm.generateInitialData()

	return sm
}

// generateInitialData 生成初始历史数据
func (sm *StockManager) generateInitialData() {
	symbols := []string{"AAPL", "GOOGL", "MSFT", "TSLA", "AMZN"}
	now := time.Now()

	for _, symbol := range symbols {
		basePrice := sm.basePrices[symbol]
		data := make([]StockData, 100)

		for i := range 100 {
			timestamp := now.Add(time.Duration(-100+i) * time.Minute)

			// 生成价格变动
			change := (rand.Float64() - 0.5) * basePrice * 0.02
			basePrice = math.Max(basePrice+change, basePrice*0.9)

			open := basePrice
			high := open + rand.Float64()*open*0.01
			low := open - rand.Float64()*open*0.01
			close := low + rand.Float64()*(high-low)
			volume := rand.Int63n(900000) + 100000

			data[i] = StockData{
				Timestamp: timestamp.UnixMilli(),
				Symbol:    symbol,
				Open:      math.Round(open*100) / 100,
				High:      math.Round(high*100) / 100,
				Low:       math.Round(low*100) / 100,
				Close:     math.Round(close*100) / 100,
				Volume:    volume,
			}

			basePrice = close
		}

		sm.stockData[symbol] = data
		sm.basePrices[symbol] = basePrice
	}
}

// generateNewData 生成新的实时数据
func (sm *StockManager) generateNewData(symbol string) StockData {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	lastData := sm.stockData[symbol][len(sm.stockData[symbol])-1]
	now := time.Now()

	// 生成价格变动
	change := (rand.Float64() - 0.5) * lastData.Close * 0.01
	newPrice := math.Max(lastData.Close+change, lastData.Close*0.99)

	newData := StockData{
		Timestamp: now.UnixMilli(),
		Symbol:    symbol,
		Open:      lastData.Close,
		High:      math.Max(lastData.Close, newPrice),
		Low:       math.Min(lastData.Close, newPrice),
		Close:     math.Round(newPrice*100) / 100,
		Volume:    rand.Int63n(500000) + 50000,
	}

	// 添加新数据
	sm.stockData[symbol] = append(sm.stockData[symbol], newData)

	// 保持最多200个数据点
	if len(sm.stockData[symbol]) > 200 {
		sm.stockData[symbol] = sm.stockData[symbol][1:]
	}

	return newData
}

// startRealTimeUpdate 启动实时数据更新
func (sm *StockManager) startRealTimeUpdate() {
	sm.updateTicker = time.NewTicker(10 * time.Second)
	go func() {
		for range sm.updateTicker.C {
			symbols := []string{"AAPL", "GOOGL", "MSFT", "TSLA", "AMZN"}
			for _, symbol := range symbols {
				newData := sm.generateNewData(symbol)
				sm.broadcastToClients(symbol, newData)
			}
		}
	}()
}

// broadcastToClients 向订阅的客户端广播数据
func (sm *StockManager) broadcastToClients(symbol string, data StockData) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	message := WebSocketMessage{
		Type: "update",
		Data: data,
	}

	for client, clientSymbol := range sm.clients {
		if clientSymbol == symbol {
			err := client.WriteJSON(message)
			if err != nil {
				log.Printf("WebSocket write error: %v", err)
				client.Close()
				delete(sm.clients, client)
			}
		}
	}
}

// getStockDataHandler 获取股票历史数据
func (sm *StockManager) getStockDataHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	symbol := vars["symbol"]

	// 设置CORS头
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "OPTIONS" {
		return
	}

	sm.mu.RLock()
	data, exists := sm.stockData[symbol]
	sm.mu.RUnlock()

	if !exists {
		http.Error(w, "Symbol not found", http.StatusNotFound)
		return
	}

	// 获取查询参数
	limitStr := r.URL.Query().Get("limit")
	limit := 100
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	// 限制返回的数据量
	start := 0
	if len(data) > limit {
		start = len(data) - limit
	}

	response := StockResponse{
		Symbol:    symbol,
		Data:      data[start:],
		Timestamp: time.Now().UnixMilli(),
	}

	json.NewEncoder(w).Encode(response)
}

// websocketHandler WebSocket处理器
func (sm *StockManager) websocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := sm.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer conn.Close()

	// 获取订阅的股票符号
	symbol := r.URL.Query().Get("symbol")
	if symbol == "" {
		symbol = "SHSE.000001" // 默认订阅上证指数
	}

	sm.mu.Lock()
	sm.clients[conn] = symbol
	sm.mu.Unlock()

	log.Printf("Client connected for symbol: %s", symbol)

	// 发送初始数据
	sm.mu.RLock()
	if data, exists := sm.stockData[symbol]; exists && len(data) > 0 {
		lastData := data[len(data)-1]
		message := WebSocketMessage{
			Type: "initial",
			Data: lastData,
		}
		conn.WriteJSON(message)
	}
	sm.mu.RUnlock()

	// 监听客户端消息
	for {
		var msg map[string]any
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("WebSocket read error: %v", err)
			break
		}

		// 处理订阅变更
		if newSymbol, ok := msg["symbol"].(string); ok {
			sm.mu.Lock()
			sm.clients[conn] = newSymbol
			sm.mu.Unlock()
			log.Printf("Client switched to symbol: %s", newSymbol)

			// 发送新股票的最新数据
			sm.mu.RLock()
			if data, exists := sm.stockData[newSymbol]; exists && len(data) > 0 {
				lastData := data[len(data)-1]
				message := WebSocketMessage{
					Type: "initial",
					Data: lastData,
				}
				conn.WriteJSON(message)
			}
			sm.mu.RUnlock()
		}
	}

	// 清理连接
	sm.mu.Lock()
	delete(sm.clients, conn)
	sm.mu.Unlock()
	log.Printf("Client disconnected")
}

// getSymbolsHandler 获取所有可用的股票符号
func (sm *StockManager) getSymbolsHandler(w http.ResponseWriter, r *http.Request) {
	// 设置CORS头
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "OPTIONS" {
		return
	}

	symbols := []map[string]any{
		{"symbol": "AAPL", "name": "Apple Inc.", "price": sm.basePrices["AAPL"]},
		{"symbol": "GOOGL", "name": "Alphabet Inc.", "price": sm.basePrices["GOOGL"]},
		{"symbol": "MSFT", "name": "Microsoft Corporation", "price": sm.basePrices["MSFT"]},
		{"symbol": "TSLA", "name": "Tesla, Inc.", "price": sm.basePrices["TSLA"]},
		{"symbol": "AMZN", "name": "Amazon.com, Inc.", "price": sm.basePrices["AMZN"]},
	}

	json.NewEncoder(w).Encode(symbols)
}

// healthHandler 健康检查处理器
func (sm *StockManager) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	status := map[string]any{
		"status":    "ok",
		"timestamp": time.Now().UnixMilli(),
		"symbols":   len(sm.stockData),
		"clients":   len(sm.clients),
	}

	json.NewEncoder(w).Encode(status)
}

// 配置结构体
type Config struct {
	API struct {
		Port  int    `toml:"port"`
		gmapi string `toml:"gmapi"`
		gmcsv string `toml:"gmcsv"`
	} `toml:"api"`
}

var cfg Config

func main() {
	// 读取配置文件
	if _, err := toml.DecodeFile("cfg.toml", &cfg); err != nil {
		fmt.Println("Error loading config file:", err)
		return
	}
	fmt.Println(` -=> Loading params from: cfg.toml`)
	// 初始化随机种子
	// rand.Seed(time.Now().UnixNano())

	// 创建股票管理器
	stockManager := NewStockManager()

	// 启动实时数据更新
	stockManager.startRealTimeUpdate()

	// 创建路由
	r := mux.NewRouter()

	// API路由
	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/stocks/{symbol}", stockManager.getStockDataHandler).Methods("GET", "OPTIONS")
	api.HandleFunc("/symbols", stockManager.getSymbolsHandler).Methods("GET", "OPTIONS")
	api.HandleFunc("/health", stockManager.healthHandler).Methods("GET")

	// WebSocket路由
	r.HandleFunc("/ws", stockManager.websocketHandler)

	// 静态文件服务（可选）
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	// 启动服务器
	// port := ":8080"
	port := fmt.Sprintf(":%d", cfg.API.Port)
	fmt.Printf("股票数据服务器启动在端口 %s\n", port)
	fmt.Println("API端点:")
	fmt.Println("  GET  /api/v1/stocks/{symbol} - 获取股票历史数据")
	fmt.Println("  GET  /api/v1/symbols - 获取所有股票符号")
	fmt.Println("  GET  /api/v1/health - 健康检查")
	fmt.Println("  WS   /ws?symbol={symbol} - WebSocket实时数据")
	fmt.Println()
	fmt.Println("示例请求:")
	fmt.Printf("  curl http://localhost%s/api/v1/stocks/AAPL\n", port)
	fmt.Printf("  curl http://localhost%s/api/v1/symbols\n", port)
	// fmt.Println("  curl http://localhost:8080/api/v1/symbols")

	log.Fatal(http.ListenAndServe(port, r))
}
