// Package factory provides factory methods for creating broker clients
package factory

import (
	"errors"
	
	"github.com/trading-platform/backend/internal/broker/common"
	"github.com/trading-platform/backend/internal/broker/xts/client"
	"github.com/trading-platform/backend/internal/broker/xts/pro"
	"github.com/trading-platform/backend/internal/broker/zerodha"
)

// NewBrokerClient creates a new broker client based on the provided configuration
func NewBrokerClient(config *common.BrokerConfig) (common.BrokerClient, error) {
	switch config.BrokerType {
	case common.BrokerTypeXTSPro:
		if config.XTSPro == nil {
			return nil, errors.New("XTS Pro configuration is required")
		}
		return pro.NewXTSProClient(config.XTSPro)
	case common.BrokerTypeXTSClient:
		if config.XTSClient == nil {
			return nil, errors.New("XTS Client configuration is required")
		}
		return client.NewXTSClientImpl(config.XTSClient)
	case common.BrokerTypeZerodha:
		if config.Zerodha == nil {
			return nil, errors.New("Zerodha configuration is required")
		}
		return zerodha.NewZerodhaAdapter(config.Zerodha)
	default:
		return nil, errors.New("unsupported broker type")
	}
}
