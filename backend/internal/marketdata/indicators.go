package marketdata

import (
	"errors"
	"fmt"
	"math"
)

// IndicatorType represents the type of technical indicator
type IndicatorType string

const (
	// Moving Averages
	SMA  IndicatorType = "SMA"  // Simple Moving Average
	EMA  IndicatorType = "EMA"  // Exponential Moving Average
	WMA  IndicatorType = "WMA"  // Weighted Moving Average
	DEMA IndicatorType = "DEMA" // Double Exponential Moving Average
	TEMA IndicatorType = "TEMA" // Triple Exponential Moving Average

	// Oscillators
	RSI   IndicatorType = "RSI"   // Relative Strength Index
	MACD  IndicatorType = "MACD"  // Moving Average Convergence Divergence
	STOCH IndicatorType = "STOCH" // Stochastic Oscillator
	CCI   IndicatorType = "CCI"   // Commodity Channel Index
	ADX   IndicatorType = "ADX"   // Average Directional Index

	// Volatility
	BOLL   IndicatorType = "BOLL"   // Bollinger Bands
	ATR    IndicatorType = "ATR"    // Average True Range
	STDDEV IndicatorType = "STDDEV" // Standard Deviation

	// Volume
	OBV  IndicatorType = "OBV"  // On-Balance Volume
	VWAP IndicatorType = "VWAP" // Volume Weighted Average Price
	AD   IndicatorType = "AD"   // Accumulation/Distribution Line

	// Trend
	PSAR   IndicatorType = "PSAR"   // Parabolic SAR
	ICHIMOKU IndicatorType = "ICHIMOKU" // Ichimoku Cloud
)

// IndicatorResult represents the result of a technical indicator calculation
type IndicatorResult struct {
	Type       IndicatorType           `json:"type"`
	Symbol     string                  `json:"symbol"`
	Interval   string                  `json:"interval"`
	Parameters map[string]interface{}  `json:"parameters"`
	Values     []map[string]interface{} `json:"values"`
}

// TechnicalIndicator defines the interface for technical indicators
type TechnicalIndicator interface {
	Calculate(data []OHLCV, params map[string]interface{}) (IndicatorResult, error)
	GetDefaultParameters() map[string]interface{}
	Validate(params map[string]interface{}) error
}

// IndicatorLibrary provides a collection of technical indicators
type IndicatorLibrary struct {
	indicators map[IndicatorType]TechnicalIndicator
}

// NewIndicatorLibrary creates a new indicator library
func NewIndicatorLibrary() *IndicatorLibrary {
	lib := &IndicatorLibrary{
		indicators: make(map[IndicatorType]TechnicalIndicator),
	}

	// Register indicators
	lib.RegisterIndicator(SMA, NewSMAIndicator())
	lib.RegisterIndicator(EMA, NewEMAIndicator())
	lib.RegisterIndicator(RSI, NewRSIIndicator())
	lib.RegisterIndicator(MACD, NewMACDIndicator())
	lib.RegisterIndicator(BOLL, NewBollingerBandsIndicator())
	lib.RegisterIndicator(ATR, NewATRIndicator())
	lib.RegisterIndicator(STOCH, NewStochasticIndicator())
	lib.RegisterIndicator(OBV, NewOBVIndicator())
	lib.RegisterIndicator(ADX, NewADXIndicator())

	return lib
}

// RegisterIndicator registers a technical indicator
func (l *IndicatorLibrary) RegisterIndicator(indicatorType IndicatorType, indicator TechnicalIndicator) {
	l.indicators[indicatorType] = indicator
}

// GetIndicator gets a technical indicator by type
func (l *IndicatorLibrary) GetIndicator(indicatorType IndicatorType) (TechnicalIndicator, error) {
	indicator, ok := l.indicators[indicatorType]
	if !ok {
		return nil, fmt.Errorf("indicator not found: %s", indicatorType)
	}
	return indicator, nil
}

// Calculate calculates a technical indicator
func (l *IndicatorLibrary) Calculate(
	indicatorType IndicatorType,
	data []OHLCV,
	params map[string]interface{},
) (IndicatorResult, error) {
	indicator, err := l.GetIndicator(indicatorType)
	if err != nil {
		return IndicatorResult{}, err
	}

	// Validate parameters
	if err := indicator.Validate(params); err != nil {
		return IndicatorResult{}, err
	}

	// Calculate indicator
	return indicator.Calculate(data, params)
}

// GetAvailableIndicators gets a list of available indicators
func (l *IndicatorLibrary) GetAvailableIndicators() []IndicatorType {
	indicators := make([]IndicatorType, 0, len(l.indicators))
	for indicator := range l.indicators {
		indicators = append(indicators, indicator)
	}
	return indicators
}

// GetDefaultParameters gets the default parameters for an indicator
func (l *IndicatorLibrary) GetDefaultParameters(indicatorType IndicatorType) (map[string]interface{}, error) {
	indicator, err := l.GetIndicator(indicatorType)
	if err != nil {
		return nil, err
	}
	return indicator.GetDefaultParameters(), nil
}

// SMAIndicator implements the Simple Moving Average indicator
type SMAIndicator struct{}

// NewSMAIndicator creates a new SMA indicator
func NewSMAIndicator() *SMAIndicator {
	return &SMAIndicator{}
}

// Calculate calculates the SMA indicator
func (i *SMAIndicator) Calculate(data []OHLCV, params map[string]interface{}) (IndicatorResult, error) {
	// Get parameters
	period := int(params["period"].(float64))
	priceType := params["price"].(string)

	// Validate data
	if len(data) < period {
		return IndicatorResult{}, fmt.Errorf("not enough data points for SMA calculation (need at least %d)", period)
	}

	// Calculate SMA
	result := IndicatorResult{
		Type:       SMA,
		Symbol:     data[0].Symbol,
		Interval:   data[0].Interval,
		Parameters: params,
		Values:     make([]map[string]interface{}, 0, len(data)-period+1),
	}

	for i := period - 1; i < len(data); i++ {
		sum := 0.0
		for j := 0; j < period; j++ {
			price := getPriceValue(data[i-j], priceType)
			sum += price
		}
		sma := sum / float64(period)

		result.Values = append(result.Values, map[string]interface{}{
			"timestamp": data[i].Timestamp,
			"value":     sma,
		})
	}

	return result, nil
}

// GetDefaultParameters gets the default parameters for the SMA indicator
func (i *SMAIndicator) GetDefaultParameters() map[string]interface{} {
	return map[string]interface{}{
		"period": float64(14),
		"price":  "close",
	}
}

// Validate validates the parameters for the SMA indicator
func (i *SMAIndicator) Validate(params map[string]interface{}) error {
	// Check required parameters
	if _, ok := params["period"]; !ok {
		params["period"] = i.GetDefaultParameters()["period"]
	}
	if _, ok := params["price"]; !ok {
		params["price"] = i.GetDefaultParameters()["price"]
	}

	// Validate period
	period, ok := params["period"].(float64)
	if !ok {
		return errors.New("period must be a number")
	}
	if period < 1 {
		return errors.New("period must be greater than 0")
	}

	// Validate price type
	priceType, ok := params["price"].(string)
	if !ok {
		return errors.New("price must be a string")
	}
	if !isValidPriceType(priceType) {
		return errors.New("invalid price type (must be open, high, low, close, or volume)")
	}

	return nil
}

// EMAIndicator implements the Exponential Moving Average indicator
type EMAIndicator struct{}

// NewEMAIndicator creates a new EMA indicator
func NewEMAIndicator() *EMAIndicator {
	return &EMAIndicator{}
}

// Calculate calculates the EMA indicator
func (i *EMAIndicator) Calculate(data []OHLCV, params map[string]interface{}) (IndicatorResult, error) {
	// Get parameters
	period := int(params["period"].(float64))
	priceType := params["price"].(string)

	// Validate data
	if len(data) < period {
		return IndicatorResult{}, fmt.Errorf("not enough data points for EMA calculation (need at least %d)", period)
	}

	// Calculate EMA
	result := IndicatorResult{
		Type:       EMA,
		Symbol:     data[0].Symbol,
		Interval:   data[0].Interval,
		Parameters: params,
		Values:     make([]map[string]interface{}, 0, len(data)-period+1),
	}

	// Calculate multiplier
	multiplier := 2.0 / (float64(period) + 1.0)

	// Calculate first EMA as SMA
	sum := 0.0
	for i := 0; i < period; i++ {
		price := getPriceValue(data[i], priceType)
		sum += price
	}
	ema := sum / float64(period)

	result.Values = append(result.Values, map[string]interface{}{
		"timestamp": data[period-1].Timestamp,
		"value":     ema,
	})

	// Calculate remaining EMAs
	for i := period; i < len(data); i++ {
		price := getPriceValue(data[i], priceType)
		ema = (price - ema) * multiplier + ema

		result.Values = append(result.Values, map[string]interface{}{
			"timestamp": data[i].Timestamp,
			"value":     ema,
		})
	}

	return result, nil
}

// GetDefaultParameters gets the default parameters for the EMA indicator
func (i *EMAIndicator) GetDefaultParameters() map[string]interface{} {
	return map[string]interface{}{
		"period": float64(14),
		"price":  "close",
	}
}

// Validate validates the parameters for the EMA indicator
func (i *EMAIndicator) Validate(params map[string]interface{}) error {
	// Check required parameters
	if _, ok := params["period"]; !ok {
		params["period"] = i.GetDefaultParameters()["period"]
	}
	if _, ok := params["price"]; !ok {
		params["price"] = i.GetDefaultParameters()["price"]
	}

	// Validate period
	period, ok := params["period"].(float64)
	if !ok {
		return errors.New("period must be a number")
	}
	if period < 1 {
		return errors.New("period must be greater than 0")
	}

	// Validate price type
	priceType, ok := params["price"].(string)
	if !ok {
		return errors.New("price must be a string")
	}
	if !isValidPriceType(priceType) {
		return errors.New("invalid price type (must be open, high, low, close, or volume)")
	}

	return nil
}

// RSIIndicator implements the Relative Strength Index indicator
type RSIIndicator struct{}

// NewRSIIndicator creates a new RSI indicator
func NewRSIIndicator() *RSIIndicator {
	return &RSIIndicator{}
}

// Calculate calculates the RSI indicator
func (i *RSIIndicator) Calculate(data []OHLCV, params map[string]interface{}) (IndicatorResult, error) {
	// Get parameters
	period := int(params["period"].(float64))
	priceType := params["price"].(string)

	// Validate data
	if len(data) < period + 1 {
		return IndicatorResult{}, fmt.Errorf("not enough data points for RSI calculation (need at least %d)", period+1)
	}

	// Calculate RSI
	result := IndicatorResult{
		Type:       RSI,
		Symbol:     data[0].Symbol,
		Interval:   data[0].Interval,
		Parameters: params,
		Values:     make([]map[string]interface{}, 0, len(data)-period),
	}

	// Calculate price changes
	changes := make([]float64, len(data)-1)
	for i := 1; i < len(data); i++ {
		currentPrice := getPriceValue(data[i], priceType)
		previousPrice := getPriceValue(data[i-1], priceType)
		changes[i-1] = currentPrice - previousPrice
	}

	// Calculate first average gain and loss
	var sumGain, sumLoss float64
	for i := 0; i < period; i++ {
		if changes[i] > 0 {
			sumGain += changes[i]
		} else {
			sumLoss += -changes[i]
		}
	}
	avgGain := sumGain / float64(period)
	avgLoss := sumLoss / float64(period)

	// Calculate first RSI
	var rs, rsi float64
	if avgLoss == 0 {
		rsi = 100
	} else {
		rs = avgGain / avgLoss
		rsi = 100 - (100 / (1 + rs))
	}

	result.Values = append(result.Values, map[string]interface{}{
		"timestamp": data[period].Timestamp,
		"value":     rsi,
	})

	// Calculate remaining RSIs
	for i := period; i < len(changes); i++ {
		// Update average gain and loss
		if changes[i] > 0 {
			avgGain = (avgGain*float64(period-1) + changes[i]) / float64(period)
			avgLoss = (avgLoss*float64(period-1)) / float64(period)
		} else {
			avgGain = (avgGain*float64(period-1)) / float64(period)
			avgLoss = (avgLoss*float64(period-1) - changes[i]) / float64(period)
		}

		// Calculate RSI
		if avgLoss == 0 {
			rsi = 100
		} else {
			rs = avgGain / avgLoss
			rsi = 100 - (100 / (1 + rs))
		}

		result.Values = append(result.Values, map[string]interface{}{
			"timestamp": data[i+1].Timestamp,
			"value":     rsi,
		})
	}

	return result, nil
}

// GetDefaultParameters gets the default parameters for the RSI indicator
func (i *RSIIndicator) GetDefaultParameters() map[string]interface{} {
	return map[string]interface{}{
		"period": float64(14),
		"price":  "close",
	}
}

// Validate validates the parameters for the RSI indicator
func (i *RSIIndicator) Validate(params map[string]interface{}) error {
	// Check required parameters
	if _, ok := params["period"]; !ok {
		params["period"] = i.GetDefaultParameters()["period"]
	}
	if _, ok := params["price"]; !ok {
		params["price"] = i.GetDefaultParameters()["price"]
	}

	// Validate period
	period, ok := params["period"].(float64)
	if !ok {
		return errors.New("period must be a number")
	}
	if period < 1 {
		return errors.New("period must be greater than 0")
	}

	// Validate price type
	priceType, ok := params["price"].(string)
	if !ok {
		return errors.New("price must be a string")
	}
	if !isValidPriceType(priceType) {
		return errors.New("invalid price type (must be open, high, low, close, or volume)")
	}

	return nil
}

// MACDIndicator implements the Moving Average Convergence Divergence indicator
type MACDIndicator struct{}

// NewMACDIndicator creates a new MACD indicator
func NewMACDIndicator() *MACDIndicator {
	return &MACDIndicator{}
}

// Calculate calculates the MACD indicator
func (i *MACDIndicator) Calculate(data []OHLCV, params map[string]interface{}) (IndicatorResult, error) {
	// Get parameters
	fastPeriod := int(params["fastPeriod"].(float64))
	slowPeriod := int(params["slowPeriod"].(float64))
	signalPeriod := int(params["signalPeriod"].(float64))
	priceType := params["price"].(string)

	// Validate data
	minDataPoints := slowPeriod + signalPeriod
	if len(data) < minDataPoints {
		return IndicatorResult{}, fmt.Errorf("not enough data points for MACD calculation (need at least %d)", minDataPoints)
	}

	// Calculate MACD
	result := IndicatorResult{
		Type:       MACD,
		Symbol:     data[0].Symbol,
		Interval:   data[0].Interval,
		Parameters: params,
		Values:     make([]map[string]interface{}, 0, len(data)-minDataPoints+1),
	}

	// Calculate fast EMA
	fastEMA := calculateEMA(data, fastPeriod, priceType)

	// Calculate slow EMA
	slowEMA := calculateEMA(data, slowPeriod, priceType)

	// Calculate MACD line (fast EMA - slow EMA)
	macdLine := make([]float64, len(fastEMA))
	for i := 0; i < len(macdLine); i++ {
		if i < slowPeriod-fastPeriod {
			macdLine[i] = 0
		} else {
			macdLine[i] = fastEMA[i] - slowEMA[i-(slowPeriod-fastPeriod)]
		}
	}

	// Calculate signal line (EMA of MACD line)
	signalLine := calculateEMAFromValues(macdLine[slowPeriod-fastPeriod:], signalPeriod)

	// Calculate histogram (MACD line - signal line)
	for i := 0; i < len(signalLine); i++ {
		macdValue := macdLine[i+slowPeriod-fastPeriod]
		signalValue := signalLine[i]
		histogramValue := macdValue - signalValue

		result.Values = append(result.Values, map[string]interface{}{
			"timestamp": data[i+slowPeriod+signalPeriod-1].Timestamp,
			"macd":      macdValue,
			"signal":    signalValue,
			"histogram": histogramValue,
		})
	}

	return result, nil
}

// GetDefaultParameters gets the default parameters for the MACD indicator
func (i *MACDIndicator) GetDefaultParameters() map[string]interface{} {
	return map[string]interface{}{
		"fastPeriod":   float64(12),
		"slowPeriod":   float64(26),
		"signalPeriod": float64(9),
		"price":        "close",
	}
}

// Validate validates the parameters for the MACD indicator
func (i *MACDIndicator) Validate(params map[string]interface{}) error {
	// Check required parameters
	if _, ok := params["fastPeriod"]; !ok {
		params["fastPeriod"] = i.GetDefaultParameters()["fastPeriod"]
	}
	if _, ok := params["slowPeriod"]; !ok {
		params["slowPeriod"] = i.GetDefaultParameters()["slowPeriod"]
	}
	if _, ok := params["signalPeriod"]; !ok {
		params["signalPeriod"] = i.GetDefaultParameters()["signalPeriod"]
	}
	if _, ok := params["price"]; !ok {
		params["price"] = i.GetDefaultParameters()["price"]
	}

	// Validate periods
	fastPeriod, ok := params["fastPeriod"].(float64)
	if !ok {
		return errors.New("fastPeriod must be a number")
	}
	if fastPeriod < 1 {
		return errors.New("fastPeriod must be greater than 0")
	}

	slowPeriod, ok := params["slowPeriod"].(float64)
	if !ok {
		return errors.New("slowPeriod must be a number")
	}
	if slowPeriod < 1 {
		return errors.New("slowPeriod must be greater than 0")
	}

	signalPeriod, ok := params["signalPeriod"].(float64)
	if !ok {
		return errors.New("signalPeriod must be a number")
	}
	if signalPeriod < 1 {
		ret
(Content truncated due to size limit. Use line ranges to read in chunks)