package services_test

import (
	"testing"
	"time"
	"github.com/stretchr/testify/assert"
	"trading_platform/backend/internal/models"
	"trading_platform/backend/internal/services/simulation"
)

func TestSimulationAccountService(t *testing.T) {
	service := simulation.NewSimulationAccountService()
	
	t.Run("CreateSimulationAccount", func(t *testing.T) {
		// Test valid account creation
		accountData := models.SimulationAccount{
			Name:           "Test Paper Trading Account",
			Description:    "Test account for paper trading",
			InitialBalance: 100000.0,
			Currency:       "USD",
			SimulationType: "PAPER",
		}
		
		account, err := service.CreateSimulationAccount("user123", accountData)
		assert.NoError(t, err)
		assert.NotNil(t, account)
		assert.Equal(t, "Test Paper Trading Account", account.Name)
		assert.Equal(t, "PAPER", account.SimulationType)
		assert.Equal(t, 100000.0, account.InitialBalance)
		assert.Equal(t, 100000.0, account.CurrentBalance)
		assert.True(t, account.IsActive)
		
		// Test invalid account creation - missing name
		invalidAccount := models.SimulationAccount{
			InitialBalance: 100000.0,
			Currency:       "USD",
			SimulationType: "PAPER",
		}
		
		_, err = service.CreateSimulationAccount("user123", invalidAccount)
		assert.Error(t, err)
		
		// Test invalid account creation - negative balance
		invalidAccount = models.SimulationAccount{
			Name:           "Invalid Account",
			InitialBalance: -1000.0,
			Currency:       "USD",
			SimulationType: "PAPER",
		}
		
		_, err = service.CreateSimulationAccount("user123", invalidAccount)
		assert.Error(t, err)
		
		// Test invalid account creation - invalid simulation type
		invalidAccount = models.SimulationAccount{
			Name:           "Invalid Account",
			InitialBalance: 100000.0,
			Currency:       "USD",
			SimulationType: "INVALID",
		}
		
		_, err = service.CreateSimulationAccount("user123", invalidAccount)
		assert.Error(t, err)
	})
	
	t.Run("GetSimulationAccount", func(t *testing.T) {
		account, err := service.GetSimulationAccount("sim123")
		assert.NoError(t, err)
		assert.NotNil(t, account)
		assert.Equal(t, "sim123", account.ID)
		
		_, err = service.GetSimulationAccount("")
		assert.Error(t, err)
	})
	
	t.Run("GetSimulationAccountsByUser", func(t *testing.T) {
		accounts, err := service.GetSimulationAccountsByUser("user123")
		assert.NoError(t, err)
		assert.NotEmpty(t, accounts)
		assert.Equal(t, 2, len(accounts))
		assert.Equal(t, "Paper Trading Account", accounts[0].Name)
		assert.Equal(t, "Backtesting Account", accounts[1].Name)
		
		_, err = service.GetSimulationAccountsByUser("")
		assert.Error(t, err)
	})
	
	t.Run("AddFunds", func(t *testing.T) {
		transaction, err := service.AddFunds("sim123", 10000.0, "Additional deposit")
		assert.NoError(t, err)
		assert.NotNil(t, transaction)
		assert.Equal(t, "DEPOSIT", transaction.Type)
		assert.Equal(t, 10000.0, transaction.Amount)
		assert.Equal(t, "Additional deposit", transaction.Description)
		
		_, err = service.AddFunds("", 10000.0, "Additional deposit")
		assert.Error(t, err)
		
		_, err = service.AddFunds("sim123", -1000.0, "Invalid deposit")
		assert.Error(t, err)
	})
	
	t.Run("WithdrawFunds", func(t *testing.T) {
		transaction, err := service.WithdrawFunds("sim123", 5000.0, "Partial withdrawal")
		assert.NoError(t, err)
		assert.NotNil(t, transaction)
		assert.Equal(t, "WITHDRAWAL", transaction.Type)
		assert.Equal(t, -5000.0, transaction.Amount)
		assert.Equal(t, "Partial withdrawal", transaction.Description)
		
		_, err = service.WithdrawFunds("", 5000.0, "Partial withdrawal")
		assert.Error(t, err)
		
		_, err = service.WithdrawFunds("sim123", -1000.0, "Invalid withdrawal")
		assert.Error(t, err)
	})
	
	t.Run("GetTransactions", func(t *testing.T) {
		startDate := time.Now().Add(-48 * time.Hour)
		endDate := time.Now()
		
		transactions, err := service.GetTransactions("sim123", startDate, endDate)
		assert.NoError(t, err)
		assert.NotEmpty(t, transactions)
		assert.Equal(t, 2, len(transactions))
		assert.Equal(t, "DEPOSIT", transactions[0].Type)
		assert.Equal(t, "P&L", transactions[1].Type)
		
		_, err = service.GetTransactions("", startDate, endDate)
		assert.Error(t, err)
	})
}

func TestVirtualBalanceService(t *testing.T) {
	service := simulation.NewVirtualBalanceService()
	
	t.Run("ProcessOrderImpact", func(t *testing.T) {
		order := models.SimulationOrder{
			Order: models.Order{
				Symbol:    "AAPL",
				Quantity:  100,
				Side:      "BUY",
				OrderType: "MARKET",
			},
			SimulatedFillPrice: 150.25,
			CommissionAmount:   15.03,
		}
		
		transaction, err := service.ProcessOrderImpact("sim123", order)
		assert.NoError(t, err)
		assert.NotNil(t, transaction)
		assert.Equal(t, "P&L", transaction.Type)
		assert.Equal(t, "ORDER", transaction.ReferenceType)
		
		// Test sell order
		sellOrder := models.SimulationOrder{
			Order: models.Order{
				Symbol:    "AAPL",
				Quantity:  100,
				Side:      "SELL",
				OrderType: "MARKET",
			},
			SimulatedFillPrice: 150.25,
			CommissionAmount:   15.03,
		}
		
		transaction, err = service.ProcessOrderImpact("sim123", sellOrder)
		assert.NoError(t, err)
		assert.NotNil(t, transaction)
		assert.Equal(t, "P&L", transaction.Type)
		assert.Greater(t, transaction.Amount, 0.0) // Positive amount for sell orders
		
		_, err = service.ProcessOrderImpact("", order)
		assert.Error(t, err)
	})
	
	t.Run("ProcessPositionUpdate", func(t *testing.T) {
		position := models.SimulationPosition{
			SimulationAccountID:  "sim123",
			Symbol:               "AAPL",
			Quantity:             100,
			Side:                 "BUY",
			SimulatedEntryPrice:  150.25,
			Status:               "CLOSED",
		}
		
		transaction, err := service.ProcessPositionUpdate("sim123", position, 155.50)
		assert.NoError(t, err)
		assert.NotNil(t, transaction)
		assert.Equal(t, "P&L", transaction.Type)
		assert.Equal(t, "POSITION", transaction.ReferenceType)
		assert.Greater(t, transaction.Amount, 0.0) // Profit
		
		// Test open position
		openPosition := models.SimulationPosition{
			SimulationAccountID:  "sim123",
			Symbol:               "AAPL",
			Quantity:             100,
			Side:                 "BUY",
			SimulatedEntryPrice:  150.25,
			Status:               "OPEN",
		}
		
		transaction, err = service.ProcessPositionUpdate("sim123", openPosition, 155.50)
		assert.NoError(t, err)
		assert.Nil(t, transaction) // No transaction for open positions
		
		_, err = service.ProcessPositionUpdate("", position, 155.50)
		assert.Error(t, err)
	})
	
	t.Run("ApplyDividend", func(t *testing.T) {
		transaction, err := service.ApplyDividend("sim123", "AAPL", 0.82, 100)
		assert.NoError(t, err)
		assert.NotNil(t, transaction)
		assert.Equal(t, "DIVIDEND", transaction.Type)
		assert.Equal(t, 82.0, transaction.Amount)
		
		_, err = service.ApplyDividend("", "AAPL", 0.82, 100)
		assert.Error(t, err)
		
		_, err = service.ApplyDividend("sim123", "", 0.82, 100)
		assert.Error(t, err)
		
		_, err = service.ApplyDividend("sim123", "AAPL", -0.82, 100)
		assert.Error(t, err)
		
		_, err = service.ApplyDividend("sim123", "AAPL", 0.82, -100)
		assert.Error(t, err)
	})
	
	t.Run("GetAccountBalance", func(t *testing.T) {
		balance, err := service.GetAccountBalance("sim123")
		assert.NoError(t, err)
		assert.Equal(t, 105000.0, balance)
		
		_, err = service.GetAccountBalance("")
		assert.Error(t, err)
	})
	
	t.Run("GetAccountEquity", func(t *testing.T) {
		equity, err := service.GetAccountEquity("sim123")
		assert.NoError(t, err)
		assert.Equal(t, 108500.0, equity) // 105000 + 3500 unrealized P&L
		
		_, err = service.GetAccountEquity("")
		assert.Error(t, err)
	})
}

func TestSimulationOrderService(t *testing.T) {
	service := simulation.NewSimulationOrderService()
	
	t.Run("CreateOrder", func(t *testing.T) {
		orderData := models.SimulationOrder{
			Order: models.Order{
				UserID:    "user123",
				Symbol:    "AAPL",
				Quantity:  100,
				Side:      "BUY",
				OrderType: "MARKET",
			},
		}
		
		order, err := service.CreateOrder("sim123", orderData)
		assert.NoError(t, err)
		assert.NotNil(t, order)
		assert.Equal(t, "AAPL", order.Symbol)
		assert.Equal(t, 100, order.Quantity)
		assert.Equal(t, "BUY", order.Side)
		assert.Equal(t, "MARKET", order.OrderType)
		assert.Equal(t, "SIM", order.Environment)
		
		// Test invalid order - missing symbol
		invalidOrder := models.SimulationOrder{
			Order: models.Order{
				UserID:    "user123",
				Quantity:  100,
				Side:      "BUY",
				OrderType: "MARKET",
			},
		}
		
		_, err = service.CreateOrder("sim123", invalidOrder)
		assert.Error(t, err)
		
		// Test invalid order - invalid side
		invalidOrder = models.SimulationOrder{
			Order: models.Order{
				UserID:    "user123",
				Symbol:    "AAPL",
				Quantity:  100,
				Side:      "INVALID",
				OrderType: "MARKET",
			},
		}
		
		_, err = service.CreateOrder("sim123", invalidOrder)
		assert.Error(t, err)
		
		_, err = service.CreateOrder("", orderData)
		assert.Error(t, err)
	})
	
	t.Run("GetOrder", func(t *testing.T) {
		order, err := service.GetOrder("order123")
		assert.NoError(t, err)
		assert.NotNil(t, order)
		assert.Equal(t, "order123", order.ID)
		
		_, err = service.GetOrder("")
		assert.Error(t, err)
	})
	
	t.Run("GetOrdersByAccount", func(t *testing.T) {
		orders, err := service.GetOrdersByAccount("sim123")
		assert.NoError(t, err)
		assert.NotEmpty(t, orders)
		assert.Equal(t, 2, len(orders))
		assert.Equal(t, "AAPL", orders[0].Symbol)
		assert.Equal(t, "MSFT", orders[1].Symbol)
		
		_, err = service.GetOrdersByAccount("")
		assert.Error(t, err)
	})
	
	t.Run("CancelOrder", func(t *testing.T) {
		order, err := service.CancelOrder("order123")
		assert.NoError(t, err)
		assert.NotNil(t, order)
		assert.Equal(t, "CANCELLED", order.Status)
		
		_, err = service.CancelOrder("")
		assert.Error(t, err)
	})
	
	t.Run("ModifyOrder", func(t *testing.T) {
		orderData := models.SimulationOrder{
			Order: models.Order{
				Quantity: 200,
				Price:    290.50,
			},
		}
		
		order, err := service.ModifyOrder("order123", orderData)
		assert.NoError(t, err)
		assert.NotNil(t, order)
		assert.Equal(t, 200, order.Quantity)
		assert.Equal(t, 290.50, order.Price)
		
		_, err = service.ModifyOrder("", orderData)
		assert.Error(t, err)
	})
}

func TestMarketSimulationService(t *testing.T) {
	service := simulation.NewMarketSimulationService()
	
	t.Run("GetCurrentMarketPrice", func(t *testing.T) {
		marketData, err := service.GetCurrentMarketPrice("AAPL")
		assert.NoError(t, err)
		assert.NotNil(t, marketData)
		assert.Equal(t, "AAPL", marketData.Symbol)
		assert.True(t, marketData.IsSimulated)
		
		_, err = service.GetCurrentMarketPrice("")
		assert.Error(t, err)
	})
	
	t.Run("GetHistoricalMarketData", func(t *testing.T) {
		startDate := time.Now().Add(-30 * 24 * time.Hour)
		endDate := time.Now()
		
		marketData, err := service.GetHistoricalMarketData("AAPL", startDate, endDate, "1d")
		assert.NoError(t, err)
		assert.NotEmpty(t, marketData)
		assert.Equal(t, "AAPL", marketData[0].Symbol)
		assert.Equal(t, "1d", marketData[0].Timeframe)
		assert.True(t, marketData[0].IsSimulated)
		
		_, err = service.GetHistoricalMarketData("", startDate, endDate, "1d")
		assert.Error(t, err)
		
		_, err = service.GetHistoricalMarketData("AAPL", startDate, endDate, "")
		assert.Error(t, err)
	})
	
	t.Run("SimulateOrderExecution", func(t *testing.T) {
		order := &models.SimulationOrder{
			Order: models.Order{
				Symbol:    "AAPL",
				Quantity:  100,
				Side:      "BUY",
				OrderType: "MARKET",
			},
		}
		
		marketSettings := &models.MarketSettings{
			SlippageModel:   "PERCENTAGE",
			SlippageValue:   0.001,
			LatencyModel:    "FIXED",
			LatencyValue:    100,
			CommissionModel: "PERCENTAGE",
			CommissionValue: 0.001,
		}
		
		err := service.SimulateOrderExecution(order, marketSettings)
		assert.NoError(t, err)
		assert.Equal(t, "FILLED", order.Status)
		assert.Equal(t, 100, order.FilledQty)
		assert.Greater(t, order.SimulatedFillPrice, 0.0)
		assert.Greater(t, order.SlippageAmount, 0.0)
		assert.Equal(t, 100, order.LatencyMs)
		assert.Greater(t, order.CommissionAmount, 0.0)
		
		// Test limit order
		limitOrder := &models.SimulationOrder{
			Order: models.Order{
				Symbol:    "AAPL",
				Quantity:  100,
				Side:      "BUY",
				OrderType: "LIMIT",
				Price:     100.0, // Low price that won't be filled
			},
		}
		
		err = service.SimulateOrderExecution(limitOrder, marketSettings)
		assert.NoError(t, err)
		assert.NotEqual(t, "FILLED", limitOrder.Status) // Should not be filled
		
		err = service.SimulateOrderExecution(nil, marketSettings)
		assert.Error(t, err)
		
		err = service.SimulateOrderExecution(order, nil)
		assert.Error(t, err)
	})
	
	t.Run("SimulateMarketMovement", func(t *testing.T) {
		marketData, err := service.SimulateMarketMovement("AAPL", "1m", 10*time.Minute)
		assert.NoError(t, err)
		assert.NotEmpty(t, marketData)
		assert.Equal(t, 10, len(marketData))
		assert.Equal(t, "AAPL", marketData[0].Symbol)
		assert.Equal(t, "1m", marketData[0].Timeframe)
		assert.True(t, marketData[0].IsSimulated)
		
		_, err = service.SimulateMarketMovement("", "1m", 10*time.Minute)
		assert.Error(t, err)
		
		_, err = service.SimulateMarketMovement("AAPL", "", 10*time.Minute)
		assert.Error(t, err)
	})
	
	t.Run("GetMarketDepth", func(t *testing.T) {
		marketDepth, err := service.GetMarketDepth("AAPL", 5)
		assert.NoError(t, err)
		assert.NotNil(t, marketDepth)
		assert.Equal(t, "AAPL", marketDepth["symbol"])
		assert.NotNil(t, marketDepth["bids"])
		assert.NotNil(t, marketDepth["asks"])
		assert.NotNil(t, marketDepth["spread"])
		
		_, err = service.GetMarketDepth("", 5)
		assert.Error(t, err)
	})
}

func TestBacktestService(t *testing.T) {
	service := simulation.NewBacktestService()
	
	t.Run("CreateBacktestSession", func(t *testing.T) {
		sessionData := models.BacktestSession{
			Name:           "Test Backtest Session",
			Description:    "Test session for backtesting",
			StartDate:      time.Now().Add(-30 * 24 * time.Hour),
			EndDate:        time.Now(),
			Symbols:        []string{"AAPL", "MSFT", "GOOGL"},
			Timeframe:      "1d",
			InitialBalance: 100000.0,
			StrategyID:     "strategy1",
			Parameters:     map[string]interface{}{"param1": 10, "param2": "value"},
		}
		
		session, err := service.CreateBacktestSession("sim123", sessionData)
		assert.NoError(t, err)
		assert.NotNil(t, session)
		assert.Equal(t, "Test Backtest Session", session.Name)
		assert.Equal(t, []string{"AAPL", "MSFT", "GOOGL"}, session.Symbols)
		assert.Equal(t, "1d", session.Timeframe)
		assert.Equal(t, 100000.0, session.InitialBalance)
		assert.Equal(t, 100000.0, session.FinalBalance)
		assert.Equal(t, "PENDING", session.Status)
		
		// Test invalid session - missing name
		invalidSession := models.BacktestSession{
			StartDate:      time.Now().Add(-30 * 24 * time.Hour),
			EndDate:        time.Now(),
			Symbols:        []string{"AAPL", "MSFT", "GOOGL"},
			Timeframe:      "1d",
			InitialBalance: 100000.0,
		}
		
		_, err = service.CreateBacktestSession("sim123", invalidSession)
		assert.Error(t, err)
		
		// Test invalid session - missing symbols
		invalidSession = models.BacktestSession{
			Name:           "Invalid Session",
			StartDate:      time.Now().Add(-30 * 24 * time.Hour),
			EndDate:        time.Now(),
			Timeframe:      "1d",
			InitialBalance: 100000.0,
		}
		
		_, err = service.CreateBacktestSession("sim123", invalidSession)
		assert.Error(t, err)
		
		// Test invalid session - start date after end date
		invalidSession = models.BacktestSession{
			Name:           "Invalid Session",
			StartDate:      time.Now(),
			EndDate:        time.Now().Add(-30 * 24 * time.Hour),
			Symbols:        []string{"AAPL", "MSFT", "GOOGL"},
			Timeframe:      "1d",
			InitialBalance: 100000.0,
		}
		
		_, err = service.CreateBacktestSession("sim123", invalidSession)
		assert.Error(t, err)
		
		_, err = service.CreateBacktestSession("", sessionData)
		assert.Error(t, err)
	})
	
	t.Run("GetBacktestSession", func(t *testing.T) {
		session, err := service.GetBacktestSession("session123")
		assert.NoError(t, err)
		assert.NotNil(t, session)
		assert.Equal(t, "session123", session.ID)
		
		_, err = service.GetBacktestSession("")
		assert.Error(t, err)
	})
	
	t.Run("GetBacktestSessionsByAccount", func(t *testing.T) {
		sessions, err := service.GetBacktestSessionsByAccount("sim123")
		assert.NoError(t, err)
		assert.NotEmpty(t, sessions)
		assert.Equal(t, 2, len(sessions))
		assert.Equal(t, "Moving Average Crossover Backtest", sessions[0].Name)
		assert.Equal(t, "RSI Strategy Backtest", sessions[1].Name)
		
		_, err = service.GetBacktestSessionsByAccount("")
		assert.Error(t, err)
	})
	
	t.Run("RunBacktest", func(t *testing.T) {
		err := service.RunBacktest("session123")
		assert.NoError(t, err)
		
		err = service.RunBacktest("")
		assert.Error(t, err)
	})
	
	t.Run("StopBacktest", func(t *testing.T) {
		err := service.StopBacktest("session123")
		assert.NoError(t, err)
		
		err = service.StopBacktest("")
		assert.Error(t, err)
	})
	
	t.Run("GetBacktestResults", func(t *testing.T) {
		results, err := service.GetBacktestResults("session123")
		assert.NoError(t, err)
		assert.NotEmpty(t, results)
		assert.Equal(t, "session123", results[0].BacktestSessionID)
		
		_, err = service.GetBacktestResults("")
		assert.Error(t, err)
	})
	
	t.Run("GetBacktestPerformanceMetrics", func(t *testing.T) {
		metrics, err := service.GetBacktestPerformanceMetrics("session123")
		assert.NoError(t, err)
		assert.NotNil(t, metrics)
		assert.NotNil(t, metrics["totalReturn"])
		assert.NotNil(t, metrics["sharpeRatio"])
		assert.NotNil(t, metrics["maxDrawdown"])
		assert.NotNil(t, metrics["winRate"])
		
		_, err = service.GetBacktestPerformanceMetrics("")
		assert.Error(t, err)
	})
	
	t.Run("GetBacktestTrades", func(t *testing.T) {
		trades, err := service.GetBacktestTrades("session123")
		assert.NoError(t, err)
		assert.NotEmpty(t, trades)
		assert.True(t, trades[0].IsBacktestOrder)
		assert.NotNil(t, trades[0].BacktestDate)
		
		_, err = service.GetBacktestTrades("")
		assert.Error(t, err)
	})
	
	t.Run("CompareBacktestSessions", func(t *testing.T) {
		comparison, err := service.CompareBacktestSessions([]string{"session1", "session2"})
		assert.NoError(t, err)
		assert.NotNil(t, comparison)
		assert.NotNil(t, comparison["sessions"])
		assert.NotNil(t, comparison["bestPerformer"])
		assert.NotNil(t, comparison["metrics"])
		
		_, err = service.CompareBacktestSessions([]string{})
		assert.Error(t, err)
	})
	
	t.Run("OptimizeStrategy", func(t *testing.T) {
		parameterRanges := map[string]map[string]interface{}{
			"param1": {"min": 10, "max": 20, "step": 1},
			"param2": {"values": []string{"value1", "value2", "value3"}},
			"param3": {"min": 0.01, "max": 0.1, "step": 0.02},
		}
		
		results, err := service.OptimizeStrategy("strategy1", parameterRanges, "sharpeRatio")
		assert.NoError(t, err)
		assert.NotNil(t, results)
		assert.Equal(t, "strategy1", results["strategyID"])
		assert.Equal(t, "sharpeRatio", results["optimizationMetric"])
		assert.NotNil(t, results["optimalParameters"])
		assert.NotNil(t, results["parameterSweep"])
		
		_, err = service.OptimizeStrategy("", parameterRanges, "sharpeRatio")
		assert.Error(t, err)
		
		_, err = service.OptimizeStrategy("strategy1", map[string]map[string]interface{}{}, "sharpeRatio")
		assert.Error(t, err)
		
		_, err = service.OptimizeStrategy("strategy1", parameterRanges, "")
		assert.Error(t, err)
	})
}
