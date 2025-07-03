package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type StockDB struct {
	db *sql.DB
}

// 优化后的表结构创建
func (s *StockDB) CreateStockTables(stockCode string) error {
	tables := []string{
		// K线日频表 - 使用自增ID作为主键，时间字段设置唯一索引
		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS kline_daily_%s (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            trade_date DATE NOT NULL UNIQUE,  -- 唯一约束，防止重复数据
            open_price DECIMAL(10,3) NOT NULL,
            high_price DECIMAL(10,3) NOT NULL,
            low_price DECIMAL(10,3) NOT NULL,
            close_price DECIMAL(10,3) NOT NULL,
            volume BIGINT DEFAULT 0,
            amount DECIMAL(15,2) DEFAULT 0,
            turnover_rate DECIMAL(5,2),
            price_change DECIMAL(10,3),
            price_change_pct DECIMAL(5,2),
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
            updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
        )`, stockCode),

		// 为日K表创建索引
		fmt.Sprintf(`CREATE INDEX IF NOT EXISTS idx_daily_%s_date ON kline_daily_%s (trade_date DESC)`, stockCode, stockCode),
		fmt.Sprintf(`CREATE INDEX IF NOT EXISTS idx_daily_%s_volume ON kline_daily_%s (volume DESC)`, stockCode, stockCode),

		// K线1分钟表 - 时间精度到分钟
		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS kline_1min_%s (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            trade_time DATETIME NOT NULL UNIQUE,  -- 精确到分钟的唯一时间戳
            open_price DECIMAL(10,3) NOT NULL,
            high_price DECIMAL(10,3) NOT NULL,
            low_price DECIMAL(10,3) NOT NULL,
            close_price DECIMAL(10,3) NOT NULL,
            volume BIGINT DEFAULT 0,
            amount DECIMAL(15,2) DEFAULT 0,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP
        )`, stockCode),

		// 为1分钟表创建复合索引
		fmt.Sprintf(`CREATE INDEX IF NOT EXISTS idx_1min_%s_time ON kline_1min_%s (trade_time DESC)`, stockCode, stockCode),
		fmt.Sprintf(`CREATE INDEX IF NOT EXISTS idx_1min_%s_date_time ON kline_1min_%s (date(trade_time), trade_time)`, stockCode, stockCode),

		// 财务衍生数据表
		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS financial_derived_%s (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            trade_date DATE NOT NULL UNIQUE,
            pe_ratio DECIMAL(8,2),
            pb_ratio DECIMAL(8,2),
            ps_ratio DECIMAL(8,2),
            market_cap DECIMAL(15,2),
            total_market_cap DECIMAL(15,2),
            turnover_rate DECIMAL(5,2),
            dividend_yield DECIMAL(5,2),
            roe DECIMAL(5,2),
            roa DECIMAL(5,2),
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
            updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
        )`, stockCode),

		// 为财务表创建索引
		fmt.Sprintf(`CREATE INDEX IF NOT EXISTS idx_financial_%s_date ON financial_derived_%s (trade_date DESC)`, stockCode, stockCode),
	}

	for _, tableSQL := range tables {
		if _, err := s.db.Exec(tableSQL); err != nil {
			return fmt.Errorf("执行SQL失败: %v, SQL: %s", err, tableSQL)
		}
	}

	return nil
}

// 数据结构定义
type KlineDaily struct {
	ID             int64     `json:"id" db:"id"`
	TradeDate      time.Time `json:"trade_date" db:"trade_date"`
	OpenPrice      float64   `json:"open_price" db:"open_price"`
	HighPrice      float64   `json:"high_price" db:"high_price"`
	LowPrice       float64   `json:"low_price" db:"low_price"`
	ClosePrice     float64   `json:"close_price" db:"close_price"`
	Volume         int64     `json:"volume" db:"volume"`
	Amount         float64   `json:"amount" db:"amount"`
	TurnoverRate   float64   `json:"turnover_rate" db:"turnover_rate"`
	PriceChange    float64   `json:"price_change" db:"price_change"`
	PriceChangePct float64   `json:"price_change_pct" db:"price_change_pct"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

// 使用UPSERT语法插入或更新日K数据（SQLite 3.24+支持）
func (s *StockDB) UpsertDailyKline(stockCode string, data *KlineDaily) error {
	query := fmt.Sprintf(`INSERT INTO kline_daily_%s 
        (trade_date, open_price, high_price, low_price, close_price, volume, amount, 
         turnover_rate, price_change, price_change_pct, updated_at) 
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
        ON CONFLICT(trade_date) DO UPDATE SET
            open_price = excluded.open_price,
            high_price = excluded.high_price,
            low_price = excluded.low_price,
            close_price = excluded.close_price,
            volume = excluded.volume,
            amount = excluded.amount,
            turnover_rate = excluded.turnover_rate,
            price_change = excluded.price_change,
            price_change_pct = excluded.price_change_pct,
            updated_at = CURRENT_TIMESTAMP`, stockCode)

	_, err := s.db.Exec(query, data.TradeDate, data.OpenPrice, data.HighPrice,
		data.LowPrice, data.ClosePrice, data.Volume, data.Amount,
		data.TurnoverRate, data.PriceChange, data.PriceChangePct)
	return err
}

// 高效的批量插入（使用事务）
func (s *StockDB) BatchUpsertDailyKline(stockCode string, dataList []*KlineDaily) error {
	if len(dataList) == 0 {
		return nil
	}

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := fmt.Sprintf(`INSERT INTO kline_daily_%s 
        (trade_date, open_price, high_price, low_price, close_price, volume, amount, 
         turnover_rate, price_change, price_change_pct, updated_at) 
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
        ON CONFLICT(trade_date) DO UPDATE SET
            open_price = excluded.open_price,
            high_price = excluded.high_price,
            low_price = excluded.low_price,
            close_price = excluded.close_price,
            volume = excluded.volume,
            amount = excluded.amount,
            turnover_rate = excluded.turnover_rate,
            price_change = excluded.price_change,
            price_change_pct = excluded.price_change_pct,
            updated_at = CURRENT_TIMESTAMP`, stockCode)

	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, data := range dataList {
		_, err = stmt.Exec(data.TradeDate, data.OpenPrice, data.HighPrice,
			data.LowPrice, data.ClosePrice, data.Volume, data.Amount,
			data.TurnoverRate, data.PriceChange, data.PriceChangePct)
		if err != nil {
			return fmt.Errorf("批量插入失败，股票代码: %s, 日期: %v, 错误: %v",
				stockCode, data.TradeDate, err)
		}
	}

	return tx.Commit()
}

// 按时间范围查询（利用索引优化）
func (s *StockDB) GetDailyKlineByDateRange(stockCode string, startDate, endDate time.Time) ([]*KlineDaily, error) {
	query := fmt.Sprintf(`SELECT id, trade_date, open_price, high_price, low_price, 
        close_price, volume, amount, turnover_rate, price_change, price_change_pct,
        created_at, updated_at 
        FROM kline_daily_%s 
        WHERE trade_date BETWEEN ? AND ? 
        ORDER BY trade_date ASC`, stockCode)

	rows, err := s.db.Query(query, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*KlineDaily
	for rows.Next() {
		data := &KlineDaily{}
		err := rows.Scan(&data.ID, &data.TradeDate, &data.OpenPrice, &data.HighPrice,
			&data.LowPrice, &data.ClosePrice, &data.Volume, &data.Amount,
			&data.TurnoverRate, &data.PriceChange, &data.PriceChangePct,
			&data.CreatedAt, &data.UpdatedAt)
		if err != nil {
			return nil, err
		}
		results = append(results, data)
	}

	return results, nil
}

// 获取最新N条记录（利用索引快速查询）
func (s *StockDB) GetLatestDailyKline(stockCode string, limit int) ([]*KlineDaily, error) {
	query := fmt.Sprintf(`SELECT id, trade_date, open_price, high_price, low_price, 
        close_price, volume, amount, turnover_rate, price_change, price_change_pct,
        created_at, updated_at 
        FROM kline_daily_%s 
        ORDER BY trade_date DESC 
        LIMIT ?`, stockCode)

	rows, err := s.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*KlineDaily
	for rows.Next() {
		data := &KlineDaily{}
		err := rows.Scan(&data.ID, &data.TradeDate, &data.OpenPrice, &data.HighPrice,
			&data.LowPrice, &data.ClosePrice, &data.Volume, &data.Amount,
			&data.TurnoverRate, &data.PriceChange, &data.PriceChangePct,
			&data.CreatedAt, &data.UpdatedAt)
		if err != nil {
			return nil, err
		}
		results = append(results, data)
	}

	return results, nil
}

// 可以追踪数据的生命周期
func (s *StockDB) GetDataModificationHistory(stockCode string, days int) error {
	query := fmt.Sprintf(`
        SELECT trade_date, created_at, updated_at,
               CASE 
                   WHEN created_at = updated_at THEN '新增'
                   ELSE '修改'
               END as operation_type
        FROM kline_daily_%s 
        WHERE created_at >= datetime('now', '-%d days')
        ORDER BY created_at DESC`, stockCode, days)

	rows, err := s.db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	fmt.Printf("股票 %s 近%d天的数据变更记录:\n", stockCode, days)
	for rows.Next() {
		var tradeDate, createdAt, updatedAt, opType string
		rows.Scan(&tradeDate, &createdAt, &updatedAt, &opType)
		fmt.Printf("日期: %s, 创建: %s, 更新: %s, 操作: %s\n",
			tradeDate, createdAt, updatedAt, opType)
	}

	return nil
}

// 检测异常更新（比如某个交易日的数据被频繁修改）
func (s *StockDB) DetectAbnormalUpdates(stockCode string) ([]*AbnormalRecord, error) {
	query := fmt.Sprintf(`
        SELECT trade_date, created_at, updated_at,
               (julianday(updated_at) - julianday(created_at)) * 24 * 60 as minutes_diff
        FROM kline_daily_%s 
        WHERE created_at != updated_at 
        AND (julianday(updated_at) - julianday(created_at)) * 24 * 60 > 30  -- 创建后30分钟内被修改
        ORDER BY updated_at DESC`, stockCode)

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var abnormalRecords []*AbnormalRecord
	for rows.Next() {
		record := &AbnormalRecord{}
		err := rows.Scan(&record.TradeDate, &record.CreatedAt,
			&record.UpdatedAt, &record.MinutesDiff)
		if err != nil {
			return nil, err
		}
		abnormalRecords = append(abnormalRecords, record)
	}

	return abnormalRecords, nil
}

type AbnormalRecord struct {
	TradeDate   string  `json:"trade_date"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
	MinutesDiff float64 `json:"minutes_diff"`
}

// 获取指定时间后发生变更的数据（用于数据同步）
func (s *StockDB) GetChangedDataSince(stockCode string, since time.Time) ([]*KlineDaily, error) {
	query := fmt.Sprintf(`
        SELECT id, trade_date, open_price, high_price, low_price, close_price,
               volume, amount, created_at, updated_at
        FROM kline_daily_%s 
        WHERE updated_at > ?  -- 关键：基于updated_at做增量同步
        ORDER BY updated_at ASC`, stockCode)

	rows, err := s.db.Query(query, since)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*KlineDaily
	for rows.Next() {
		data := &KlineDaily{}
		err := rows.Scan(&data.ID, &data.TradeDate, &data.OpenPrice,
			&data.HighPrice, &data.LowPrice, &data.ClosePrice,
			&data.Volume, &data.Amount, &data.CreatedAt, &data.UpdatedAt)
		if err != nil {
			return nil, err
		}
		results = append(results, data)
	}

	return results, nil
}

// 增量同步示例
// func (s *StockDB) SyncToRemoteSystem(stockCode string) error {
// 	// 获取上次同步时间
// 	lastSyncTime, err := s.getLastSyncTime(stockCode)
// 	if err != nil {
// 		return err
// 	}

// 	// 获取变更数据
// 	changedData, err := s.GetChangedDataSince(stockCode, lastSyncTime)
// 	if err != nil {
// 		return err
// 	}

// 	if len(changedData) == 0 {
// 		fmt.Printf("股票 %s 没有新的变更数据\n", stockCode)
// 		return nil
// 	}

// 	// 同步到远程系统
// 	for _, data := range changedData {
// 		if err := s.sendToRemoteSystem(data); err != nil {
// 			return fmt.Errorf("同步数据失败: %v", err)
// 		}
// 	}

// 	// 更新同步时间
// 	return s.updateLastSyncTime(stockCode, time.Now())
// }

// 创建历史表用于数据版本控制
func (s *StockDB) CreateHistoryTable(stockCode string) error {
	query := fmt.Sprintf(`
        CREATE TABLE IF NOT EXISTS kline_daily_%s_history (
            history_id INTEGER PRIMARY KEY AUTOINCREMENT,
            original_id INTEGER NOT NULL,
            trade_date DATE NOT NULL,
            open_price DECIMAL(10,3),
            high_price DECIMAL(10,3),
            low_price DECIMAL(10,3),
            close_price DECIMAL(10,3),
            volume BIGINT,
            amount DECIMAL(15,2),
            operation_type VARCHAR(10) NOT NULL,  -- INSERT, UPDATE, DELETE
            operation_time DATETIME DEFAULT CURRENT_TIMESTAMP,
            operator VARCHAR(50)
        )`, stockCode)

	_, err := s.db.Exec(query)
	return err
}

// 在更新数据前，先保存历史版本
func (s *StockDB) UpdateWithHistory(stockCode string, data *KlineDaily, operator string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. 查询原始数据
	var original KlineDaily
	selectQuery := fmt.Sprintf(`
        SELECT id, trade_date, open_price, high_price, low_price, close_price, volume, amount
        FROM kline_daily_%s WHERE trade_date = ?`, stockCode)

	err = tx.QueryRow(selectQuery, data.TradeDate).Scan(
		&original.ID, &original.TradeDate, &original.OpenPrice,
		&original.HighPrice, &original.LowPrice, &original.ClosePrice,
		&original.Volume, &original.Amount)

	if err != nil && err != sql.ErrNoRows {
		return err
	}

	// 2. 如果存在原始数据,保存到历史表
	if err != sql.ErrNoRows {
		historyQuery := fmt.Sprintf(`
            INSERT INTO kline_daily_%s_history 
            (original_id, trade_date, open_price, high_price, low_price, close_price, 
             volume, amount, operation_type, operator)
            VALUES (?, ?, ?, ?, ?, ?, ?, ?, 'UPDATE', ?)`, stockCode)

		_, err = tx.Exec(historyQuery, original.ID, original.TradeDate,
			original.OpenPrice, original.HighPrice, original.LowPrice,
			original.ClosePrice, original.Volume, original.Amount, operator)
		if err != nil {
			return err
		}
	}

	// 3. 更新主表数据
	updateQuery := fmt.Sprintf(`
        INSERT INTO kline_daily_%s 
        (trade_date, open_price, high_price, low_price, close_price, volume, amount, updated_at)
        VALUES (?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
        ON CONFLICT(trade_date) DO UPDATE SET
            open_price = excluded.open_price,
            high_price = excluded.high_price,
            low_price = excluded.low_price,
            close_price = excluded.close_price,
            volume = excluded.volume,
            amount = excluded.amount,
            updated_at = CURRENT_TIMESTAMP`, stockCode)

	_, err = tx.Exec(updateQuery, data.TradeDate, data.OpenPrice,
		data.HighPrice, data.LowPrice, data.ClosePrice, data.Volume, data.Amount)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// 监控数据更新频率，识别热点数据
func (s *StockDB) AnalyzeUpdatePatterns(stockCode string, days int) error {
	query := fmt.Sprintf(`
        SELECT 
            date(updated_at) as update_date,
            COUNT(*) as update_count,
            COUNT(DISTINCT trade_date) as affected_trading_days
        FROM kline_daily_%s 
        WHERE updated_at >= datetime('now', '-%d days')
        AND created_at != updated_at  -- 只统计真正的更新操作
        GROUP BY date(updated_at)
        ORDER BY update_count DESC`, stockCode, days)

	rows, err := s.db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	fmt.Printf("股票 %s 近%d天的数据更新模式分析:\n", stockCode, days)
	for rows.Next() {
		var updateDate string
		var updateCount, affectedDays int
		rows.Scan(&updateDate, &updateCount, &affectedDays)
		fmt.Printf("日期: %s, 更新次数: %d, 影响交易日: %d\n",
			updateDate, updateCount, affectedDays)
	}

	return nil
}

// 清理长时间未更新的临时数据
func (s *StockDB) CleanupStaleData(stockCode string, days int) error {
	query := fmt.Sprintf(`
        DELETE FROM kline_daily_%s 
        WHERE volume = 0 
        AND amount = 0 
        AND updated_at < datetime('now', '-%d days')`, stockCode, days)

	result, err := s.db.Exec(query)
	if err != nil {
		return err
	}

	affected, _ := result.RowsAffected()
	fmt.Printf("清理了股票 %s 的 %d 条过期数据\n", stockCode, affected)

	return nil
}

// 检查数据是否存在（基于时间戳）
func (s *StockDB) ExistsByDate(stockCode string, tradeDate time.Time) (bool, error) {
	query := fmt.Sprintf(`SELECT COUNT(1) FROM kline_daily_%s WHERE trade_date = ?`, stockCode)

	var count int
	err := s.db.QueryRow(query, tradeDate).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// 数据库优化相关方法
func (s *StockDB) OptimizeDatabase() error {
	// 执行VACUUM清理数据库
	if _, err := s.db.Exec("VACUUM"); err != nil {
		return fmt.Errorf("VACUUM失败: %v", err)
	}

	// 分析表统计信息，优化查询计划
	if _, err := s.db.Exec("ANALYZE"); err != nil {
		return fmt.Errorf("ANALYZE失败: %v", err)
	}

	return nil
}

// 获取表的统计信息
func (s *StockDB) GetTableStats(stockCode string) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	tables := []string{"kline_daily", "kline_1min", "financial_derived"}

	for _, table := range tables {
		tableName := fmt.Sprintf("%s_%s", table, stockCode)

		// 获取记录数
		var count int64
		query := fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName)
		err := s.db.QueryRow(query).Scan(&count)
		if err != nil {
			// 表可能不存在，跳过
			continue
		}

		stats[tableName+"_count"] = count

		// 获取最早和最新日期
		if table == "kline_daily" || table == "financial_derived" {
			var minDate, maxDate sql.NullTime
			dateQuery := fmt.Sprintf("SELECT MIN(trade_date), MAX(trade_date) FROM %s", tableName)
			err = s.db.QueryRow(dateQuery).Scan(&minDate, &maxDate)
			if err == nil {
				if minDate.Valid {
					stats[tableName+"_min_date"] = minDate.Time
				}
				if maxDate.Valid {
					stats[tableName+"_max_date"] = maxDate.Time
				}
			}
		}
	}

	return stats, nil
}

// 使用示例
func main() {
	db, err := sql.Open("sqlite3", "stock_optimized.db?cache=shared&mode=rwc&_journal_mode=WAL")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// 设置连接池参数
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Hour)

	stockDB := &StockDB{db: db}

	// 创建表
	if err := stockDB.CreateStockTables("000001"); err != nil {
		panic(err)
	}

	// 插入测试数据
	testData := &KlineDaily{
		TradeDate:      time.Now().Truncate(24 * time.Hour),
		OpenPrice:      10.50,
		HighPrice:      11.20,
		LowPrice:       10.30,
		ClosePrice:     11.00,
		Volume:         1000000,
		Amount:         10800000.00,
		TurnoverRate:   5.2,
		PriceChange:    0.50,
		PriceChangePct: 4.76,
	}

	if err := stockDB.UpsertDailyKline("000001", testData); err != nil {
		fmt.Printf("插入数据失败: %v\n", err)
	} else {
		fmt.Println("数据插入成功")
	}

	// 查询统计信息
	stats, err := stockDB.GetTableStats("000001")
	if err != nil {
		fmt.Printf("获取统计信息失败: %v\n", err)
	} else {
		fmt.Printf("表统计信息: %+v\n", stats)
	}
}
