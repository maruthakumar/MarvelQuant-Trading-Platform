import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import '@testing-library/jest-dom';
import TradingViewIntegration from '../implementation/tradingview/TradingViewIntegration';

// Mock WebSocketProvider component
jest.mock('../implementation/websocket/WebSocketProvider', () => {
  const originalModule = jest.requireActual('../implementation/websocket/WebSocketProvider');
  return {
    ...originalModule,
    WebSocketProvider: jest.fn(({ children }) => (
      <div data-testid="websocket-provider">{children}</div>
    ))
  };
});

// Mock TradingViewChart component
jest.mock('../implementation/tradingview/TradingViewChart', () => {
  return {
    __esModule: true,
    default: jest.fn(({ symbol, interval }) => (
      <div data-testid="tradingview-chart">
        <div data-testid="chart-symbol">{symbol}</div>
        <div data-testid="chart-interval">{interval}</div>
      </div>
    ))
  };
});

// Mock TradingViewSignalProcessor component
jest.mock('../implementation/tradingview/TradingViewSignalProcessor', () => {
  return {
    __esModule: true,
    default: jest.fn(({ onSignalReceived, onSignalProcessed }) => {
      // Store callbacks for testing
      (global as any).mockSignalCallbacks = {
        onSignalReceived,
        onSignalProcessed
      };
      return <div data-testid="tradingview-signal-processor"></div>;
    })
  };
});

describe('TradingViewIntegration', () => {
  const mockOnSignalReceived = jest.fn();
  const mockOnSignalProcessed = jest.fn();
  
  beforeEach(() => {
    jest.clearAllMocks();
  });
  
  test('should render all child components', () => {
    render(
      <TradingViewIntegration
        wsUrl="ws://localhost:8080"
        defaultSymbol="NIFTY"
        defaultInterval="1D"
        onSignalReceived={mockOnSignalReceived}
        onSignalProcessed={mockOnSignalProcessed}
      />
    );
    
    // Check if WebSocketProvider is rendered
    expect(screen.getByTestId('websocket-provider')).toBeInTheDocument();
    
    // Check if TradingViewChart is rendered with correct props
    expect(screen.getByTestId('tradingview-chart')).toBeInTheDocument();
    expect(screen.getByTestId('chart-symbol')).toHaveTextContent('NIFTY');
    expect(screen.getByTestId('chart-interval')).toHaveTextContent('1D');
    
    // Check if TradingViewSignalProcessor is rendered
    expect(screen.getByTestId('tradingview-signal-processor')).toBeInTheDocument();
    
    // Check if signal history section is rendered
    expect(screen.getByText('Signal History')).toBeInTheDocument();
    expect(screen.getByText('No signals received yet')).toBeInTheDocument();
  });
  
  test('should update symbol and interval when changed', () => {
    render(
      <TradingViewIntegration
        wsUrl="ws://localhost:8080"
        defaultSymbol="NIFTY"
        defaultInterval="1D"
      />
    );
    
    // Find symbol input and update button
    const symbolInput = screen.getByLabelText('Symbol:');
    const updateButton = screen.getByText('Update');
    
    // Change symbol
    fireEvent.change(symbolInput, { target: { value: 'BANKNIFTY' } });
    fireEvent.click(updateButton);
    
    // Check if TradingViewChart is updated with new symbol
    expect(screen.getByTestId('chart-symbol')).toHaveTextContent('BANKNIFTY');
    
    // Find interval select
    const intervalSelect = screen.getByLabelText('Interval:');
    
    // Change interval
    fireEvent.change(intervalSelect, { target: { value: '5' } });
    
    // Check if TradingViewChart is updated with new interval
    expect(screen.getByTestId('chart-interval')).toHaveTextContent('5');
  });
  
  test('should handle signal received and processed', async () => {
    render(
      <TradingViewIntegration
        wsUrl="ws://localhost:8080"
        defaultSymbol="NIFTY"
        defaultInterval="1D"
        onSignalReceived={mockOnSignalReceived}
        onSignalProcessed={mockOnSignalProcessed}
      />
    );
    
    // Create test signal
    const testSignal = {
      multileg: true,
      type: 'ENTRY',
      strategy: {
        type: 'IronCondor',
        tag: 'DEFAULT'
      },
      instrument: {
        symbol: 'NIFTY',
        lots: 1,
        product: 'MIS'
      }
    };
    
    // Simulate signal received
    (global as any).mockSignalCallbacks.onSignalReceived(testSignal);
    
    // Check if onSignalReceived was called
    expect(mockOnSignalReceived).toHaveBeenCalledWith(testSignal);
    
    // Check if signal history is updated
    expect(screen.queryByText('No signals received yet')).not.toBeInTheDocument();
    expect(screen.getByText('ENTRY')).toBeInTheDocument();
    expect(screen.getByText('IronCondor')).toBeInTheDocument();
    expect(screen.getByText('NIFTY')).toBeInTheDocument();
    expect(screen.getByText('1')).toBeInTheDocument();
    expect(screen.getByText('Pending')).toBeInTheDocument();
    
    // Simulate signal processed
    (global as any).mockSignalCallbacks.onSignalProcessed(testSignal, true);
    
    // Check if onSignalProcessed was called
    expect(mockOnSignalProcessed).toHaveBeenCalledWith(testSignal, true);
    
    // Check if signal history is updated
    expect(screen.getByText('Processed')).toBeInTheDocument();
  });
});
