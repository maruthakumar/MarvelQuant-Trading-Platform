import unittest
import json
import time
import threading
from unittest.mock import MagicMock, patch
from websocket import WebSocketApp

# Import the MultiLegClient class
import sys
sys.path.append('/home/ubuntu/workspace/implementation/python')
from MultiLegClient import MultiLegClient

class TestMultiLegClient(unittest.TestCase):
    """Test cases for the MultiLegClient class."""
    
    def setUp(self):
        """Set up test fixtures."""
        # Mock WebSocketApp
        self.mock_ws = MagicMock(spec=WebSocketApp)
        
        # Patch the websocket.WebSocketApp to return our mock
        self.patcher = patch('websocket.WebSocketApp', return_value=self.mock_ws)
        self.mock_websocket_app = self.patcher.start()
        
        # Create client with test configuration
        self.client = MultiLegClient(
            ws_url="ws://test.example.com/ws",
            auth_token="test-token",
            auto_reconnect=False  # Disable auto-reconnect for testing
        )
        
        # Mock thread
        self.client.ws_thread = MagicMock()
    
    def tearDown(self):
        """Tear down test fixtures."""
        self.patcher.stop()
    
    def test_initialization(self):
        """Test client initialization."""
        self.assertEqual(self.client.ws_url, "ws://test.example.com/ws")
        self.assertEqual(self.client.auth_token, "test-token")
        self.assertFalse(self.client.auto_reconnect)
        self.assertFalse(self.client.is_connected)
        self.assertEqual(self.client.reconnect_attempts, 0)
        self.assertEqual(self.client.message_queue, [])
        
        # Verify WebSocketApp was created with correct URL
        self.mock_websocket_app.assert_called_once_with(
            "ws://test.example.com/ws",
            on_open=self.client._on_open,
            on_message=self.client._on_message,
            on_error=self.client._on_error,
            on_close=self.client._on_close
        )
        
        # Verify WebSocketApp.run_forever was started in a thread
        self.mock_ws.run_forever.assert_not_called()  # Not called directly
        self.client.ws_thread.start.assert_called_once()  # Thread was started
    
    def test_connect_disconnect(self):
        """Test connect and disconnect methods."""
        # Test successful connection
        result = self.client.connect()
        self.assertTrue(result)
        
        # Test successful disconnection
        result = self.client.disconnect()
        self.assertTrue(result)
        self.mock_ws.close.assert_called_once()
        self.assertFalse(self.client.is_connected)
    
    def test_on_open(self):
        """Test on_open handler."""
        # Register mock callback
        on_connect_callback = MagicMock()
        self.client.on_connect(on_connect_callback)
        
        # Mock _send_authentication
        self.client._send_authentication = MagicMock()
        
        # Call on_open
        self.client._on_open(self.mock_ws)
        
        # Verify state changes
        self.assertTrue(self.client.is_connected)
        self.assertEqual(self.client.reconnect_attempts, 0)
        
        # Verify authentication was sent
        self.client._send_authentication.assert_called_once()
        
        # Verify callback was called
        on_connect_callback.assert_called_once()
    
    def test_on_message(self):
        """Test on_message handler."""
        # Register mock callbacks
        on_signal_callback = MagicMock()
        on_status_update_callback = MagicMock()
        on_error_callback = MagicMock()
        
        self.client.on_signal(on_signal_callback)
        self.client.on_status_update(on_status_update_callback)
        self.client.on_error(on_error_callback)
        
        # Test signal message
        signal_data = {
            'type': 'signal',
            'payload': {'test': 'data'}
        }
        self.client._on_message(self.mock_ws, json.dumps(signal_data))
        on_signal_callback.assert_called_once_with({'test': 'data'})
        
        # Test status update message
        status_data = {
            'type': 'status_update',
            'payload': {'status': 'ok'}
        }
        self.client._on_message(self.mock_ws, json.dumps(status_data))
        on_status_update_callback.assert_called_once_with({'status': 'ok'})
        
        # Test error message
        error_data = {
            'type': 'error',
            'payload': 'Test error'
        }
        self.client._on_message(self.mock_ws, json.dumps(error_data))
        on_error_callback.assert_called_once_with('Test error')
        
        # Test invalid JSON
        self.client._on_message(self.mock_ws, "invalid json")
        # No additional callbacks should be called
        on_signal_callback.assert_called_once()
        on_status_update_callback.assert_called_once()
        on_error_callback.assert_called_once()
    
    def test_on_error(self):
        """Test on_error handler."""
        # Register mock callback
        on_error_callback = MagicMock()
        self.client.on_error(on_error_callback)
        
        # Call on_error
        self.client._on_error(self.mock_ws, Exception("Test error"))
        
        # Verify callback was called
        on_error_callback.assert_called_once_with("WebSocket error: Test error")
    
    def test_on_close(self):
        """Test on_close handler."""
        # Register mock callback
        on_disconnect_callback = MagicMock()
        self.client.on_disconnect(on_disconnect_callback)
        
        # Mock reconnect
        self.client.reconnect = MagicMock()
        
        # Call on_close
        self.client._on_close(self.mock_ws, 1000, "Normal closure")
        
        # Verify state changes
        self.assertFalse(self.client.is_connected)
        
        # Verify callback was called
        on_disconnect_callback.assert_called_once_with(1000, "Normal closure")
        
        # Verify reconnect was not called (auto_reconnect is False)
        self.client.reconnect.assert_not_called()
        
        # Test with auto_reconnect enabled
        self.client.auto_reconnect = True
        self.client._on_close(self.mock_ws, 1001, "Going away")
        self.client.reconnect.assert_called_once()
    
    def test_send_message(self):
        """Test send_message method."""
        # Test sending when connected
        self.client.is_connected = True
        message = {'type': 'test', 'payload': 'data'}
        
        result = self.client._send_message(message)
        self.assertTrue(result)
        self.mock_ws.send.assert_called_once_with(json.dumps(message))
        
        # Test sending when disconnected
        self.client.is_connected = False
        self.mock_ws.send.reset_mock()
        
        result = self.client._send_message(message)
        self.assertFalse(result)
        self.mock_ws.send.assert_not_called()
        self.assertEqual(self.client.message_queue, [message])
    
    def test_process_message_queue(self):
        """Test process_message_queue method."""
        # Add messages to queue
        messages = [
            {'type': 'test1', 'payload': 'data1'},
            {'type': 'test2', 'payload': 'data2'}
        ]
        self.client.message_queue = messages.copy()
        
        # Mock _send_message
        self.client._send_message = MagicMock(return_value=True)
        
        # Process queue
        self.client._process_message_queue()
        
        # Verify messages were sent
        self.assertEqual(self.client._send_message.call_count, 2)
        self.client._send_message.assert_any_call(messages[0])
        self.client._send_message.assert_any_call(messages[1])
        
        # Verify queue was cleared
        self.assertEqual(self.client.message_queue, [])
    
    def test_send_signal(self):
        """Test send_signal method."""
        # Valid signal
        valid_signal = {
            'multileg': True,
            'type': 'ENTRY',
            'strategy': {
                'type': 'IronCondor',
                'tag': 'DEFAULT'
            },
            'instrument': {
                'symbol': 'NIFTY',
                'lots': 1,
                'product': 'MIS'
            }
        }
        
        # Mock _send_message
        self.client._send_message = MagicMock(return_value=True)
        
        # Send valid signal
        result = self.client.send_signal(valid_signal)
        self.assertTrue(result)
        self.client._send_message.assert_called_once()
        
        # Invalid signal (missing required field)
        invalid_signal = {
            'type': 'ENTRY',
            'strategy': {
                'type': 'IronCondor'
            }
        }
        
        self.client._send_message.reset_mock()
        result = self.client.send_signal(invalid_signal)
        self.assertFalse(result)
        self.client._send_message.assert_not_called()
    
    def test_validate_signal(self):
        """Test _validate_signal method."""
        # Valid signal
        valid_signal = {
            'multileg': True,
            'type': 'ENTRY',
            'strategy': {
                'type': 'IronCondor',
                'tag': 'DEFAULT'
            },
            'instrument': {
                'symbol': 'NIFTY',
                'lots': 1,
                'product': 'MIS'
            }
        }
        
        result = self.client._validate_signal(valid_signal)
        self.assertTrue(result)
        
        # Invalid signals
        # Missing multileg
        invalid_signal1 = {
            'type': 'ENTRY',
            'strategy': {'type': 'IronCondor'},
            'instrument': {'symbol': 'NIFTY'}
        }
        self.assertFalse(self.client._validate_signal(invalid_signal1))
        
        # Invalid strategy (not a dict)
        invalid_signal2 = {
            'multileg': True,
            'type': 'ENTRY',
            'strategy': 'IronCondor',
            'instrument': {'symbol': 'NIFTY'}
        }
        self.assertFalse(self.client._validate_signal(invalid_signal2))
        
        # Invalid instrument (missing symbol)
        invalid_signal3 = {
            'multileg': True,
            'type': 'ENTRY',
            'strategy': {'type': 'IronCondor'},
            'instrument': {'lots': 1}
        }
        self.assertFalse(self.client._validate_signal(invalid_signal3))
    
    def test_create_entry_signal(self):
        """Test create_entry_signal method."""
        # Create entry signal
        signal = self.client.create_entry_signal(
            strategy_type='IronCondor',
            strategy_tag='DEFAULT',
            symbol='NIFTY',
            lots=2,
            product='NRML',
            parameters={'width': 100}
        )
        
        # Verify signal format
        self.assertEqual(signal['multileg'], True)
        self.assertEqual(signal['type'], 'ENTRY')
        self.assertEqual(signal['strategy']['type'], 'IronCondor')
        self.assertEqual(signal['strategy']['tag'], 'DEFAULT')
        self.assertEqual(signal['instrument']['symbol'], 'NIFTY')
        self.assertEqual(signal['instrument']['lots'], 2)
        self.assertEqual(signal['instrument']['product'], 'NRML')
        self.assertEqual(signal['parameters']['width'], 100)
    
    def test_create_exit_signal(self):
        """Test create_exit_signal method."""
        # Create exit signal
        signal = self.client.create_exit_signal(
            strategy_type='IronCondor',
            strategy_tag='DEFAULT',
            symbol='NIFTY',
            lots=2,
            product='NRML'
        )
        
        # Verify signal format
        self.assertEqual(signal['multileg'], True)
        self.assertEqual(signal['type'], 'EXIT')
        self.assertEqual(signal['strategy']['type'], 'IronCondor')
        self.assertEqual(signal['strategy']['tag'], 'DEFAULT')
        self.assertEqual(signal['instrument']['symbol'], 'NIFTY')
        self.assertEqual(signal['instrument']['lots'], 2)
        self.assertEqual(signal['instrument']['product'], 'NRML')
        self.assertNotIn('parameters', signal)

if __name__ == '__main__':
    unittest.main()
