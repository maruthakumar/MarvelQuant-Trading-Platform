# Known Issues and Limitations

This document provides a comprehensive list of known issues and limitations in the Trading Platform v9.7.7.

## Known Issues

### Backend Issues

1. **Order Execution Service Memory Leak**
   - **Issue**: Under high load conditions (>500 orders/second), the order execution service may experience a small memory leak.
   - **Workaround**: Restart the service daily during off-hours. This issue is scheduled to be fixed in v9.8.0.
   - **Impact**: Minimal impact during normal operation; may require attention during extended high-volume trading sessions.

2. **Market Data Service Connection Drops**
   - **Issue**: Connections to the market data service may occasionally drop when the upstream provider experiences latency spikes.
   - **Workaround**: The service automatically reconnects, but clients should implement retry logic with exponential backoff.
   - **Impact**: Temporary data gaps possible during reconnection periods.

3. **Authentication Token Refresh**
   - **Issue**: In rare cases, token refresh may fail if the original token was issued during a leap second.
   - **Workaround**: Users will need to log in again. This edge case occurs extremely rarely.
   - **Impact**: Minimal; affects only users who logged in during a leap second (occurs once every few years).

4. **Database Query Performance with Large Datasets**
   - **Issue**: Queries involving historical data spanning multiple years may experience performance degradation.
   - **Workaround**: Use time-partitioned queries or implement data aggregation for long time periods.
   - **Impact**: Affects reporting and analysis functions with very large historical datasets.

5. **WebSocket Connection Limit**
   - **Issue**: The WebSocket server has a default limit of 5,000 concurrent connections.
   - **Workaround**: Implement connection pooling on client side or deploy multiple WebSocket servers with load balancing.
   - **Impact**: May affect deployments with very large user bases.

### Frontend Issues

1. **Chart Rendering in Safari**
   - **Issue**: Some advanced chart features may render incorrectly in Safari browsers.
   - **Workaround**: Use Chrome or Firefox for optimal experience. Safari compatibility improvements are planned for v9.8.0.
   - **Impact**: Affects users who prefer Safari as their primary browser.

2. **Large Order Book Display**
   - **Issue**: Order books with more than 1000 levels may cause performance degradation in the UI.
   - **Workaround**: Limit order book display to 500 levels. Performance optimization is scheduled for v9.8.0.
   - **Impact**: Affects traders who need to view very deep order books.

3. **Strategy Backtesting UI**
   - **Issue**: The strategy backtesting UI may become unresponsive with extremely large datasets (>1 year of minute data).
   - **Workaround**: Limit backtesting to shorter time periods or use lower resolution data for long-term tests.
   - **Impact**: Affects users performing extensive backtesting on long time periods.

4. **React Component Re-rendering**
   - **Issue**: Some complex dashboard components may experience unnecessary re-rendering under specific conditions.
   - **Workaround**: Implement React.memo and useMemo more extensively in affected components.
   - **Impact**: Minor performance impact on dashboard pages with many components.

5. **Mobile Responsiveness on Complex Screens**
   - **Issue**: Some advanced trading screens do not adapt well to very small mobile screens (<375px width).
   - **Workaround**: Use tablet or larger mobile devices for complex trading operations.
   - **Impact**: Affects users attempting to use advanced features on small mobile devices.

### C++ Execution Engine Issues

1. **Memory Usage Optimization**
   - **Issue**: The C++ execution engine may use more memory than necessary when processing certain complex order types.
   - **Workaround**: Configure memory limits in the service configuration. Optimization is planned for v9.8.0.
   - **Impact**: Higher memory usage than optimal, but within acceptable limits for most deployments.

2. **Thread Contention**
   - **Issue**: Under extreme load (>1000 orders/second), thread contention may occur in the execution engine.
   - **Workaround**: Adjust thread pool configuration based on available CPU cores. Further optimization is planned for v9.8.0.
   - **Impact**: May cause slight latency increases during peak load periods.

3. **Compiler Compatibility**
   - **Issue**: The C++ execution engine requires GCC 9.3.0 or later; earlier versions may cause compilation errors.
   - **Workaround**: Ensure development and deployment environments use compatible compiler versions.
   - **Impact**: Affects development and deployment on systems with older compiler versions.

4. **Debug Mode Performance**
   - **Issue**: Running the C++ execution engine in debug mode causes significant performance degradation.
   - **Workaround**: Use release builds for performance testing and production deployment.
   - **Impact**: Affects development workflow when debugging performance-critical code.

### Integration Issues

1. **Multi-broker Order Routing**
   - **Issue**: Order routing logic may not always select the optimal broker under certain market conditions.
   - **Workaround**: Manually specify preferred broker for critical orders.
   - **Impact**: May result in slightly higher execution costs in some scenarios.

2. **XTS API Rate Limiting**
   - **Issue**: XTS API enforces strict rate limits that may be reached during high-volume trading.
   - **Workaround**: Implement request batching and throttling in high-volume scenarios.
   - **Impact**: May delay execution of some orders during extremely active trading periods.

3. **Zerodha Integration Timezone Handling**
   - **Issue**: Zerodha integration may incorrectly handle orders placed exactly at market open/close in certain timezones.
   - **Workaround**: Avoid placing orders within 1 minute of market open/close when using Zerodha.
   - **Impact**: Affects users in non-standard timezones trading at market open/close boundaries.

4. **Interactive Brokers TWS Connection**
   - **Issue**: Connection to Interactive Brokers TWS may be lost if TWS automatically restarts for its daily maintenance.
   - **Workaround**: The system automatically attempts to reconnect, but users should be aware of the daily TWS maintenance window.
   - **Impact**: Potential brief interruption in trading capabilities during TWS maintenance.

### SIM Environment Issues

1. **SIM Environment Data Synchronization**
   - **Issue**: Market data in the SIM environment may have a slight delay (up to 5 seconds) compared to the LIVE environment.
   - **Workaround**: Account for potential data delay when evaluating SIM trading strategies.
   - **Impact**: May affect accuracy of high-frequency trading simulations.

2. **SIM User Permission Management**
   - **Issue**: Changes to SIM user permissions require cache refresh to take effect immediately.
   - **Workaround**: Manually trigger cache refresh after permission changes or wait for automatic refresh (every 5 minutes).
   - **Impact**: Brief delay in permission changes taking effect.

## Limitations

### Performance Limitations

1. **Maximum Order Rate**
   - The system is designed to handle up to 1000 orders per second in the current configuration.
   - Higher throughput requires additional hardware resources and configuration adjustments.
   - Benchmark testing shows performance degradation beyond 1500 orders per second on reference hardware.

2. **Market Data Update Frequency**
   - Real-time market data updates are limited to 10 updates per second per symbol in the standard configuration.
   - Higher update frequencies require additional licensing and configuration.
   - WebSocket connections may experience increased latency with more than 100 simultaneous symbol subscriptions.

3. **Concurrent Users**
   - The system is designed to support up to 1000 concurrent users in the standard configuration.
   - Higher user counts require additional hardware resources and configuration adjustments.
   - UI performance may degrade with more than 500 simultaneous active users on reference hardware.

4. **Database Scaling**
   - The current database schema is optimized for up to 10 million orders per day.
   - Higher volumes require database sharding or additional optimization.
   - Historical data queries spanning more than 5 years may require data aggregation for acceptable performance.

5. **WebSocket Connection Limits**
   - Each server instance supports up to 5,000 concurrent WebSocket connections.
   - Deployments requiring more connections should implement a load-balanced WebSocket cluster.
   - Each client connection should limit subscriptions to 100 symbols for optimal performance.

### Functional Limitations

1. **Supported Order Types**
   - Currently supports: Market, Limit, Stop, Stop Limit, Trailing Stop, and OCO orders.
   - Complex conditional orders beyond these types are not supported in this version.
   - Custom order types require implementation at both backend and broker integration layers.

2. **Supported Asset Classes**
   - Currently supports: Equities, Options, Futures, and Forex.
   - Crypto assets and bonds are planned for future releases.
   - Support for derivatives is limited to exchange-traded options and futures.

3. **Backtesting Limitations**
   - Backtesting is limited to 1 year of historical data at minute resolution.
   - Tick-level backtesting is supported but limited to 1 week of data due to performance considerations.
   - Backtesting does not account for market impact of large orders.

4. **Strategy Optimization**
   - Strategy optimization is limited to 5 parameters with up to 10 values each.
   - More complex optimization requires custom implementation.
   - Grid search is the only supported optimization method; genetic algorithms planned for future release.

5. **Risk Management**
   - Position risk calculations are updated at 1-minute intervals.
   - Real-time risk calculations are limited to delta, gamma, theta, and vega for options.
   - Portfolio-level VaR calculations are based on end-of-day positions only.

6. **Reporting**
   - Standard reports are limited to daily, weekly, and monthly periods.
   - Custom reporting periods require manual data export and processing.
   - Report generation for periods longer than 1 year may experience performance issues.

### Integration Limitations

1. **Broker Connectivity**
   - Currently supports integration with XTS, Interactive Brokers, and Zerodha.
   - Additional broker integrations require custom development.
   - Each broker integration has specific limitations based on the broker's API capabilities.

2. **External Data Sources**
   - Currently supports integration with XTS, Polygon.io, and Alpha Vantage.
   - Additional data source integrations require custom development.
   - Historical data availability varies by source and asset class.

3. **Authentication Methods**
   - Currently supports username/password, JWT, and OAuth2 authentication.
   - Biometric and hardware token authentication are planned for future releases.
   - Single sign-on integration is limited to OAuth2-compatible providers.

4. **API Rate Limits**
   - Public API endpoints are rate-limited to 100 requests per minute per user.
   - WebSocket connections are limited to 100 symbol subscriptions per connection.
   - Bulk operations are limited to 1000 items per request.

5. **Mobile Platform Support**
   - Full functionality is available on iOS 14+ and Android 10+.
   - Older mobile OS versions may experience limited functionality or UI issues.
   - Some advanced features are optimized for tablet or desktop use.

### SIM Environment Limitations

1. **Simulation Accuracy**
   - SIM environment uses simplified market impact models.
   - Order execution in SIM environment assumes ideal conditions and may not reflect real-world slippage.
   - Market data in SIM environment may have slight delays compared to LIVE environment.

2. **Virtual Balance Management**
   - Virtual account balances are updated at end-of-day, not in real-time.
   - Margin calculations in SIM environment use simplified models.
   - SIM environment does not simulate broker-specific margin requirements.

3. **Strategy Testing**
   - SIM environment supports up to 50 concurrent strategies per user.
   - Strategy scheduling is limited to 1-minute minimum intervals.
   - Complex multi-asset strategies may experience performance limitations.

## Planned Resolutions

Most of the issues and limitations mentioned above are scheduled to be addressed in the upcoming releases:

- **v9.8.0 (May 2025)**: 
  - Performance optimizations for C++ execution engine
  - Safari compatibility fixes
  - Improved order book rendering
  - Enhanced WebSocket connection handling
  - Memory usage optimizations

- **v9.9.0 (June 2025)**: 
  - Additional broker integrations
  - Enhanced authentication methods
  - Improved backtesting capabilities
  - Mobile responsiveness improvements
  - Database query optimization

- **v10.0.0 (July 2025)**: 
  - Support for crypto assets and bonds
  - Advanced strategy optimization using genetic algorithms
  - Real-time risk calculations
  - Enhanced simulation accuracy
  - Comprehensive reporting system

## Reporting New Issues

For any issues not listed here, please report them through the issue tracking system with the following information:
1. Detailed description of the issue
2. Steps to reproduce
3. Expected behavior
4. Actual behavior
5. Environment details (browser, OS, etc.)
6. Screenshots or logs (if applicable)

## Requesting Feature Enhancements

To request new features or enhancements to address current limitations:
1. Check if the feature is already planned in upcoming releases
2. Submit a feature request through the issue tracking system
3. Include detailed use cases and business justification
4. Specify desired implementation timeline

## Document Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0 | April 4, 2025 | Trading Platform Team | Initial document creation |
| 1.1 | April 4, 2025 | Trading Platform Team | Added SIM environment issues and limitations |
| 1.2 | April 4, 2025 | Trading Platform Team | Expanded performance limitations section |
| 1.3 | April 4, 2025 | Trading Platform Team | Added reporting and feature request guidelines |
