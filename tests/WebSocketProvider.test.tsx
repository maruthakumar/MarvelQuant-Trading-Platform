import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import '@testing-library/jest-dom';
import { WebSocketProvider, WebSocketContext, useWebSocket } from '../implementation/websocket/WebSocketProvider';

// Mock WebSocketClient component
jest.mock('../implementation/websocket/WebSocketClient', () => {
  return {
    __esModule: true,
    default: jest.fn(({ onMessage, onConnectionStateChange }) => {
      // Store the callbacks for later use in tests
      (global as any).mockWebSocketCallbacks = {
        onMessage,
        onConnectionStateChange
      };
      return null;
    }),
    ConnectionState: {
      CONNECTING: 'connecting',
      OPEN: 'open',
      CLOSING: 'closing',
      CLOSED: 'closed',
      ERROR: 'error'
    },
    MessageType: {
      AUTHENTICATION: 'authentication',
      SIGNAL: 'signal',
      ACKNOWLEDGMENT: 'acknowledgment',
      STATUS_UPDATE: 'status_update',
      ERROR: 'error'
    }
  };
});

// Mock WebSocketStatus component
jest.mock('../implementation/websocket/WebSocketStatus', () => {
  return {
    __esModule: true,
    default: jest.fn(() => <div data-testid="websocket-status">WebSocket Status</div>)
  };
});

// Test component that uses the WebSocket context
const TestComponent = () => {
  const { connectionState, lastMessage, sendSignal } = useWebSocket();
  
  const handleClick = () => {
    sendSignal({
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
    });
  };
  
  return (
    <div>
      <div data-testid="connection-state">{connectionState}</div>
      <div data-testid="last-message">
        {lastMessage ? JSON.stringify(lastMessage) : 'No message'}
      </div>
      <button data-testid="send-signal-button" onClick={handleClick}>
        Send Signal
      </button>
    </div>
  );
};

describe('WebSocketProvider', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });
  
  test('should provide WebSocket context to children', () => {
    render(
      <WebSocketProvider wsUrl="ws://localhost:8080">
        <TestComponent />
      </WebSocketProvider>
    );
    
    // Check if the connection state is initially CLOSED
    expect(screen.getByTestId('connection-state')).toHaveTextContent('closed');
    
    // Check if the last message is initially null
    expect(screen.getByTestId('last-message')).toHaveTextContent('No message');
    
    // Check if the WebSocketStatus component is rendered
    expect(screen.getByTestId('websocket-status')).toBeInTheDocument();
  });
  
  test('should update connection state when WebSocketClient calls onConnectionStateChange', async () => {
    render(
      <WebSocketProvider wsUrl="ws://localhost:8080">
        <TestComponent />
      </WebSocketProvider>
    );
    
    // Simulate WebSocketClient calling onConnectionStateChange with CONNECTING
    (global as any).mockWebSocketCallbacks.onConnectionStateChange('connecting');
    
    // Check if the connection state is updated
    expect(screen.getByTestId('connection-state')).toHaveTextContent('connecting');
    
    // Simulate WebSocketClient calling onConnectionStateChange with OPEN
    (global as any).mockWebSocketCallbacks.onConnectionStateChange('open');
    
    // Check if the connection state is updated
    expect(screen.getByTestId('connection-state')).toHaveTextContent('open');
  });
  
  test('should update last message when WebSocketClient calls onMessage', async () => {
    render(
      <WebSocketProvider wsUrl="ws://localhost:8080">
        <TestComponent />
      </WebSocketProvider>
    );
    
    // Create test message
    const testMessage = {
      type: 'signal',
      payload: { test: 'data' },
      timestamp: Date.now()
    };
    
    // Simulate WebSocketClient calling onMessage with test message
    (global as any).mockWebSocketCallbacks.onMessage(testMessage);
    
    // Check if the last message is updated
    expect(screen.getByTestId('last-message')).toHaveTextContent(JSON.stringify(testMessage));
  });
  
  test('should throw error when useWebSocket is used outside WebSocketProvider', () => {
    // Suppress console.error for this test
    const originalError = console.error;
    console.error = jest.fn();
    
    // Expect render to throw an error
    expect(() => {
      render(<TestComponent />);
    }).toThrow('useWebSocket must be used within a WebSocketProvider');
    
    // Restore console.error
    console.error = originalError;
  });
});
