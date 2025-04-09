import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import '@testing-library/jest-dom';
import WebSocketClient, { ConnectionState, MessageType } from '../implementation/websocket/WebSocketClient';

// Mock WebSocket
class MockWebSocket {
  onopen: any = null;
  onclose: any = null;
  onmessage: any = null;
  onerror: any = null;
  readyState = 0;
  
  constructor(public url: string) {
    setTimeout(() => {
      this.readyState = 1;
      if (this.onopen) this.onopen();
    }, 100);
  }
  
  send(data: string) {
    // Mock send functionality
    return true;
  }
  
  close() {
    setTimeout(() => {
      this.readyState = 3;
      if (this.onclose) this.onclose();
    }, 100);
  }
  
  // Helper to simulate incoming messages
  mockReceiveMessage(data: any) {
    if (this.onmessage) {
      this.onmessage({ data: JSON.stringify(data) });
    }
  }
  
  // Helper to simulate errors
  mockError(error: any) {
    if (this.onerror) {
      this.onerror(error);
    }
  }
}

// Mock window.WebSocket
global.WebSocket = MockWebSocket as any;

describe('WebSocketClient', () => {
  const mockOnMessage = jest.fn();
  const mockOnConnectionStateChange = jest.fn();
  
  beforeEach(() => {
    jest.clearAllMocks();
  });
  
  test('should connect to WebSocket server on mount', async () => {
    const { rerender } = render(
      <WebSocketClient
        url="ws://localhost:8080"
        onMessage={mockOnMessage}
        onConnectionStateChange={mockOnConnectionStateChange}
      />
    );
    
    // Should call onConnectionStateChange with CONNECTING state
    expect(mockOnConnectionStateChange).toHaveBeenCalledWith(ConnectionState.CONNECTING);
    
    // Wait for connection to open
    await waitFor(() => {
      expect(mockOnConnectionStateChange).toHaveBeenCalledWith(ConnectionState.OPEN);
    });
  });
  
  test('should send authentication message when authToken is provided', async () => {
    const sendSpy = jest.spyOn(MockWebSocket.prototype, 'send');
    
    render(
      <WebSocketClient
        url="ws://localhost:8080"
        authToken="test-token"
        onMessage={mockOnMessage}
        onConnectionStateChange={mockOnConnectionStateChange}
      />
    );
    
    // Wait for connection to open and auth message to be sent
    await waitFor(() => {
      expect(sendSpy).toHaveBeenCalled();
    });
    
    // Verify auth message format
    const sentMessage = JSON.parse(sendSpy.mock.calls[0][0]);
    expect(sentMessage.type).toBe(MessageType.AUTHENTICATION);
    expect(sentMessage.payload.token).toBe("test-token");
  });
  
  test('should handle incoming messages', async () => {
    render(
      <WebSocketClient
        url="ws://localhost:8080"
        onMessage={mockOnMessage}
        onConnectionStateChange={mockOnConnectionStateChange}
      />
    );
    
    // Wait for connection to open
    await waitFor(() => {
      expect(mockOnConnectionStateChange).toHaveBeenCalledWith(ConnectionState.OPEN);
    });
    
    // Get instance of MockWebSocket
    const mockWs = (global.WebSocket as any).mock.instances[0];
    
    // Simulate receiving a message
    const testMessage = {
      type: MessageType.SIGNAL,
      payload: { test: 'data' },
      timestamp: Date.now()
    };
    
    mockWs.mockReceiveMessage(testMessage);
    
    // Verify onMessage was called with the correct message
    expect(mockOnMessage).toHaveBeenCalledWith(testMessage);
  });
  
  test('should handle connection close', async () => {
    render(
      <WebSocketClient
        url="ws://localhost:8080"
        onMessage={mockOnMessage}
        onConnectionStateChange={mockOnConnectionStateChange}
      />
    );
    
    // Wait for connection to open
    await waitFor(() => {
      expect(mockOnConnectionStateChange).toHaveBeenCalledWith(ConnectionState.OPEN);
    });
    
    // Get instance of MockWebSocket
    const mockWs = (global.WebSocket as any).mock.instances[0];
    
    // Simulate connection close
    mockWs.onclose();
    
    // Verify onConnectionStateChange was called with CLOSED state
    expect(mockOnConnectionStateChange).toHaveBeenCalledWith(ConnectionState.CLOSED);
  });
  
  test('should handle connection errors', async () => {
    render(
      <WebSocketClient
        url="ws://localhost:8080"
        onMessage={mockOnMessage}
        onConnectionStateChange={mockOnConnectionStateChange}
      />
    );
    
    // Wait for connection to open
    await waitFor(() => {
      expect(mockOnConnectionStateChange).toHaveBeenCalledWith(ConnectionState.OPEN);
    });
    
    // Get instance of MockWebSocket
    const mockWs = (global.WebSocket as any).mock.instances[0];
    
    // Simulate connection error
    mockWs.onerror(new Error('Test error'));
    
    // Verify onConnectionStateChange was called with ERROR state
    expect(mockOnConnectionStateChange).toHaveBeenCalledWith(ConnectionState.ERROR);
  });
});
