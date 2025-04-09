import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import '@testing-library/jest-dom';
import TradingViewChart from '../implementation/tradingview/TradingViewChart';
import { WebSocketProvider } from '../implementation/websocket/WebSocketProvider';

// Mock useWebSocket hook
jest.mock('../implementation/websocket/WebSocketProvider', () => {
  const originalModule = jest.requireActual('../implementation/websocket/WebSocketProvider');
  return {
    ...originalModule,
    useWebSocket: jest.fn(() => ({
      connectionState: 'open',
      lastMessage: null,
      sendSignal: jest.fn(() => true),
      reconnect: jest.fn()
    }))
  };
});

// Mock TradingView widget
global.TradingView = {
  widget: jest.fn().mockImplementation(function(config) {
    this.container = config.container;
    this.symbol = config.symbol;
    this.interval = config.interval;
    this.theme = config.theme;
    this.autosize = config.autosize;
    this.height = config.height;
    this.width = config.width;
    this.onChartReady = (callback) => {
      this.chartReadyCallback = callback;
      return this;
    };
    this.chart = () => ({
      onIntervalChanged: () => ({
        subscribe: jest.fn((_, callback) => {
          this.intervalChangedCallback = callback;
          return this;
        })
      }),
      onSymbolChanged: () => ({
        subscribe: jest.fn((_, callback) => {
          this.symbolChangedCallback = callback;
          return this;
        })
      })
    });
    this.createButton = () => ({
      attr: jest.fn(() => ({
        text: jest.fn(() => ({
          on: jest.fn((event, callback) => {
            this.buttonCallback = callback;
            return this;
          })
        }))
      }))
    });
    this.remove = jest.fn();
    
    // Simulate chart ready
    setTimeout(() => {
      if (this.chartReadyCallback) {
        this.chartReadyCallback();
      }
    }, 100);
    
    return this;
  })
};

describe('TradingViewChart', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });
  
  test('should render TradingView chart container', () => {
    render(
      <WebSocketProvider wsUrl="ws://localhost:8080">
        <TradingViewChart symbol="NIFTY" />
      </WebSocketProvider>
    );
    
    // Check if the chart container is rendered
    expect(screen.getByClassName('tradingview-chart-container')).toBeInTheDocument();
    expect(screen.getByClassName('tradingview-chart')).toBeInTheDocument();
  });
  
  test('should initialize TradingView widget with correct props', async () => {
    render(
      <WebSocketProvider wsUrl="ws://localhost:8080">
        <TradingViewChart 
          symbol="NIFTY" 
          interval="1D" 
          theme="light" 
          autosize={true} 
          height={500} 
          width={800} 
          enableSignals={true} 
        />
      </WebSocketProvider>
    );
    
    // Wait for widget to be initialized
    await waitFor(() => {
      expect(global.TradingView.widget).toHaveBeenCalled();
    });
    
    // Check if widget was initialized with correct props
    const widgetInstance = (global.TradingView.widget as jest.Mock).mock.instances[0];
    expect(widgetInstance.symbol).toBe('NIFTY');
    expect(widgetInstance.interval).toBe('1D');
    expect(widgetInstance.theme).toBe('light');
    expect(widgetInstance.autosize).toBe(true);
  });
  
  test('should setup communication with TradingView chart when ready', async () => {
    const { useWebSocket } = require('../implementation/websocket/WebSocketProvider');
    const mockSendSignal = jest.fn(() => true);
    
    // Mock useWebSocket to return sendSignal function
    useWebSocket.mockReturnValue({
      connectionState: 'open',
      lastMessage: null,
      sendSignal: mockSendSignal,
      reconnect: jest.fn()
    });
    
    render(
      <WebSocketProvider wsUrl="ws://localhost:8080">
        <TradingViewChart 
          symbol="NIFTY" 
          interval="1D" 
          enableSignals={true} 
        />
      </WebSocketProvider>
    );
    
    // Wait for widget to be initialized and chart ready callback to be called
    await waitFor(() => {
      const widgetInstance = (global.TradingView.widget as jest.Mock).mock.instances[0];
      expect(widgetInstance.chartReadyCallback).toBeDefined();
    });
    
    // Get widget instance
    const widgetInstance = (global.TradingView.widget as jest.Mock).mock.instances[0];
    
    // Simulate button click
    if (widgetInstance.buttonCallback) {
      widgetInstance.buttonCallback();
    }
    
    // Check if sendSignal was called with correct signal
    expect(mockSendSignal).toHaveBeenCalledWith(expect.objectContaining({
      multileg: true,
      type: 'ENTRY',
      strategy: expect.objectContaining({
        type: 'IronCondor',
        tag: 'DEFAULT'
      }),
      instrument: expect.objectContaining({
        symbol: 'NIFTY',
        lots: 1,
        product: 'MIS'
      })
    }));
  });
  
  test('should cleanup widget on unmount', async () => {
    const { unmount } = render(
      <WebSocketProvider wsUrl="ws://localhost:8080">
        <TradingViewChart symbol="NIFTY" />
      </WebSocketProvider>
    );
    
    // Wait for widget to be initialized
    await waitFor(() => {
      expect(global.TradingView.widget).toHaveBeenCalled();
    });
    
    // Get widget instance
    const widgetInstance = (global.TradingView.widget as jest.Mock).mock.instances[0];
    
    // Unmount component
    unmount();
    
    // Check if widget.remove was called
    expect(widgetInstance.remove).toHaveBeenCalled();
  });
});
