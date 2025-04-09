import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import '@testing-library/jest-dom';
import TradingViewSignalProcessor from '../implementation/tradingview/TradingViewSignalProcessor';
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

describe('TradingViewSignalProcessor', () => {
  const mockOnSignalReceived = jest.fn();
  const mockOnSignalProcessed = jest.fn();
  
  beforeEach(() => {
    jest.clearAllMocks();
  });
  
  test('should process signal when received from WebSocket', async () => {
    const { useWebSocket } = require('../implementation/websocket/WebSocketProvider');
    
    // Render component
    render(
      <WebSocketProvider wsUrl="ws://localhost:8080">
        <TradingViewSignalProcessor
          onSignalReceived={mockOnSignalReceived}
          onSignalProcessed={mockOnSignalProcessed}
        />
      </WebSocketProvider>
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
    
    // Simulate receiving signal from WebSocket
    useWebSocket.mockReturnValue({
      connectionState: 'open',
      lastMessage: {
        type: 'signal',
        payload: testSignal,
        timestamp: Date.now()
      },
      sendSignal: jest.fn(() => true),
      reconnect: jest.fn()
    });
    
    // Re-render to trigger useEffect with new lastMessage
    render(
      <WebSocketProvider wsUrl="ws://localhost:8080">
        <TradingViewSignalProcessor
          onSignalReceived={mockOnSignalReceived}
          onSignalProcessed={mockOnSignalProcessed}
        />
      </WebSocketProvider>
    );
    
    // Check if onSignalReceived was called with the test signal
    expect(mockOnSignalReceived).toHaveBeenCalledWith(expect.objectContaining({
      multileg: true,
      type: 'ENTRY',
      strategy: expect.objectContaining({
        type: 'IronCondor',
        tag: 'DEFAULT'
      })
    }));
    
    // Wait for signal processing to complete
    await waitFor(() => {
      expect(mockOnSignalProcessed).toHaveBeenCalled();
    });
    
    // Check if onSignalProcessed was called with the test signal and success=true
    expect(mockOnSignalProcessed).toHaveBeenCalledWith(
      expect.objectContaining({
        multileg: true,
        type: 'ENTRY'
      }),
      true
    );
  });
  
  test('should parse TradingView alert text format', async () => {
    const { useWebSocket } = require('../implementation/websocket/WebSocketProvider');
    
    // Render component
    render(
      <WebSocketProvider wsUrl="ws://localhost:8080">
        <TradingViewSignalProcessor
          onSignalReceived={mockOnSignalReceived}
          onSignalProcessed={mockOnSignalProcessed}
        />
      </WebSocketProvider>
    );
    
    // Create test alert text
    const testAlertText = `MULTILEG: YES
TYPE: ENTRY
OPT: IronCondor
STAG: DEFAULT
LOTS: 2
PRODUCT: MIS
SYMBOL: NIFTY`;
    
    // Simulate receiving alert text from WebSocket
    useWebSocket.mockReturnValue({
      connectionState: 'open',
      lastMessage: {
        type: 'signal',
        payload: testAlertText,
        timestamp: Date.now()
      },
      sendSignal: jest.fn(() => true),
      reconnect: jest.fn()
    });
    
    // Re-render to trigger useEffect with new lastMessage
    render(
      <WebSocketProvider wsUrl="ws://localhost:8080">
        <TradingViewSignalProcessor
          onSignalReceived={mockOnSignalReceived}
          onSignalProcessed={mockOnSignalProcessed}
        />
      </WebSocketProvider>
    );
    
    // Check if onSignalReceived was called with the parsed signal
    expect(mockOnSignalReceived).toHaveBeenCalledWith(expect.objectContaining({
      multileg: true,
      type: 'ENTRY',
      strategy: expect.objectContaining({
        type: 'IronCondor',
        tag: 'DEFAULT'
      }),
      instrument: expect.objectContaining({
        symbol: 'NIFTY',
        lots: 2,
        product: 'MIS'
      })
    }));
    
    // Wait for signal processing to complete
    await waitFor(() => {
      expect(mockOnSignalProcessed).toHaveBeenCalled();
    });
  });
  
  test('should handle invalid signal format gracefully', async () => {
    const { useWebSocket } = require('../implementation/websocket/WebSocketProvider');
    
    // Render component
    render(
      <WebSocketProvider wsUrl="ws://localhost:8080">
        <TradingViewSignalProcessor
          onSignalReceived={mockOnSignalReceived}
          onSignalProcessed={mockOnSignalProcessed}
        />
      </WebSocketProvider>
    );
    
    // Create invalid signal
    const invalidSignal = {
      invalid: 'data'
    };
    
    // Simulate receiving invalid signal from WebSocket
    useWebSocket.mockReturnValue({
      connectionState: 'open',
      lastMessage: {
        type: 'signal',
        payload: invalidSignal,
        timestamp: Date.now()
      },
      sendSignal: jest.fn(() => true),
      reconnect: jest.fn()
    });
    
    // Re-render to trigger useEffect with new lastMessage
    render(
      <WebSocketProvider wsUrl="ws://localhost:8080">
        <TradingViewSignalProcessor
          onSignalReceived={mockOnSignalReceived}
          onSignalProcessed={mockOnSignalProcessed}
        />
      </WebSocketProvider>
    );
    
    // Check if onSignalReceived was called with a fallback signal
    expect(mockOnSignalReceived).toHaveBeenCalledWith(expect.objectContaining({
      multileg: true,
      type: expect.any(String),
      strategy: expect.objectContaining({
        type: expect.any(String),
        tag: expect.any(String)
      }),
      instrument: expect.objectContaining({
        symbol: expect.any(String),
        lots: expect.any(Number),
        product: expect.any(String)
      })
    }));
    
    // Wait for signal processing to complete
    await waitFor(() => {
      expect(mockOnSignalProcessed).toHaveBeenCalled();
    });
  });
});
