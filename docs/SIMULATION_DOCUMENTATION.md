# Simulation System Documentation

## Overview

The Trading Platform Simulation System provides a comprehensive environment for backtesting, optimizing, and evaluating trading strategies using historical market data. This documentation covers all aspects of the simulation system, including its architecture, components, configuration, and usage.

## Table of Contents

1. [System Architecture](#system-architecture)
2. [Historical Data Replay](#historical-data-replay)
3. [Performance Analytics](#performance-analytics)
4. [Strategy Optimization](#strategy-optimization)
5. [Scenario Testing](#scenario-testing)
6. [Results Export](#results-export)
7. [Comparison Tools](#comparison-tools)
8. [User Interface](#user-interface)
9. [API Reference](#api-reference)
10. [Configuration](#configuration)
11. [Troubleshooting](#troubleshooting)

## System Architecture

The simulation system is designed with a modular architecture that allows for flexibility, extensibility, and maintainability. The system consists of the following main components:

- **Historical Data Replay**: Provides mechanisms for replaying historical market data at various speeds
- **Performance Analytics**: Calculates and visualizes performance metrics for trading strategies
- **Strategy Optimization**: Optimizes strategy parameters using various algorithms
- **Scenario Testing**: Tests strategies under different market scenarios
- **Results Export**: Exports simulation results in various formats
- **Comparison Tools**: Compares results from multiple simulations
- **User Interface**: Provides a web-based interface for interacting with the simulation system

These components interact with each other through well-defined interfaces, allowing for easy integration and extension.

### Component Diagram

```
+---------------------+     +---------------------+     +---------------------+
|                     |     |                     |     |                     |
| Historical Data     |---->| Performance         |---->| Results Export      |
| Replay              |     | Analytics           |     |                     |
|                     |     |                     |     |                     |
+---------------------+     +---------------------+     +---------------------+
         |                           |                           |
         |                           |                           |
         v                           v                           v
+---------------------+     +---------------------+     +---------------------+
|                     |     |                     |     |                     |
| Strategy            |<--->| Scenario Testing    |<--->| Comparison Tools    |
| Optimization        |     |                     |     |                     |
|                     |     |                     |     |                     |
+---------------------+     +---------------------+     +---------------------+
                                      |
                                      |
                                      v
                            +---------------------+
                            |                     |
                            | User Interface      |
                            |                     |
                            |                     |
                            +---------------------+
```

## Historical Data Replay

The Historical Data Replay component provides mechanisms for replaying historical market data at various speeds, allowing for realistic simulation of trading strategies.

### Features

- Support for multiple data sources (CSV, JSON, time-series databases)
- Variable replay speeds (real-time, accelerated, decelerated)
- Simulation of market events (gaps, halts, news events)
- Programmatic API and command-line interface

### Usage

#### Basic Usage

```python
from trading_platform.simulation.historical_replay import ReplayController

# Create replay controller
controller = ReplayController()

# Configure data source
controller.set_data_source('database', connection_params={
    'host': 'localhost',
    'port': 8086,
    'database': 'market_data'
})

# Configure replay parameters
controller.set_replay_parameters(
    symbols=['AAPL', 'MSFT', 'GOOGL'],
    start_date='2023-01-01',
    end_date='2023-12-31',
    speed=2.0  # 2x speed
)

# Register event handlers
controller.register_handler('price_update', on_price_update)
controller.register_handler('trade', on_trade)
controller.register_handler('market_event', on_market_event)

# Start replay
controller.start()

# Pause replay
controller.pause()

# Resume replay
controller.resume()

# Stop replay
controller.stop()
```

#### Command-Line Interface

```bash
# Start replay
python -m trading_platform.simulation.historical_replay.cli.replay_cli \
    --data-source database \
    --host localhost \
    --port 8086 \
    --database market_data \
    --symbols AAPL,MSFT,GOOGL \
    --start-date 2023-01-01 \
    --end-date 2023-12-31 \
    --speed 2.0
```

### Configuration

The Historical Data Replay component can be configured through the global configuration object or by passing configuration parameters directly to the ReplayController constructor.

```python
from trading_platform.simulation.config import Config

# Create configuration
config = Config({
    'historical_replay.data_source': 'database',
    'historical_replay.database.host': 'localhost',
    'historical_replay.database.port': 8086,
    'historical_replay.database.name': 'market_data',
    'historical_replay.default_speed': 1.0,
    'historical_replay.include_market_events': True
})

# Create replay controller with configuration
controller = ReplayController(config=config)
```

## Performance Analytics

The Performance Analytics component calculates and visualizes performance metrics for trading strategies, providing insights into strategy performance.

### Features

- Comprehensive performance metrics calculation
- Visualization of performance metrics
- Report generation in various formats
- Integration with other simulation components

### Performance Metrics

The Performance Analytics component calculates the following metrics:

#### Return Metrics
- Total Return
- Annualized Return
- Monthly Returns
- Daily Returns

#### Risk Metrics
- Volatility
- Drawdowns
- Value at Risk (VaR)
- Conditional Value at Risk (CVaR)

#### Risk-Adjusted Metrics
- Sharpe Ratio
- Sortino Ratio
- Calmar Ratio
- Information Ratio

#### Trade Metrics
- Win Rate
- Profit Factor
- Average Win/Loss
- Maximum Win/Loss

### Usage

```python
from trading_platform.simulation.performance_analytics import PerformanceMetrics, PerformanceVisualizer, PerformanceReporter

# Calculate performance metrics
metrics = PerformanceMetrics()
results = metrics.calculate_all_metrics(
    equity_curve=equity_curve,
    trades=trades,
    benchmark_returns=benchmark_returns
)

# Visualize performance
visualizer = PerformanceVisualizer()
fig_equity = visualizer.plot_equity_curve(equity_curve)
fig_drawdowns = visualizer.plot_drawdowns(equity_curve)
fig_returns = visualizer.plot_returns_distribution(equity_curve)
fig_monthly = visualizer.plot_monthly_returns_heatmap(equity_curve)

# Generate performance report
reporter = PerformanceReporter()
report_path = reporter.generate_report(
    equity_curve=equity_curve,
    trades=trades,
    benchmark_returns=benchmark_returns,
    format='html'
)
```

## Strategy Optimization

The Strategy Optimization component optimizes strategy parameters using various algorithms, helping to find the best parameter values for a given strategy.

### Features

- Multiple optimization methods (grid search, random search, genetic algorithm)
- Customizable objective functions
- Parallel optimization for improved performance
- Integration with other simulation components

### Optimization Methods

#### Grid Search

Grid search exhaustively searches through a specified parameter grid to find the best parameter combination.

```python
from trading_platform.simulation.strategy_optimization import ParameterOptimizer

# Create parameter optimizer
optimizer = ParameterOptimizer()

# Define parameter grid
param_grid = {
    'fast_ma': [5, 10, 15, 20],
    'slow_ma': [50, 100, 150, 200],
    'stop_loss': [0.01, 0.02, 0.03]
}

# Define objective function
def objective_function(params, data):
    # Implement strategy with given parameters
    # Return performance metric (e.g., Sharpe ratio)
    pass

# Run grid search
results = optimizer.grid_search(
    param_grid=param_grid,
    objective_function=objective_function,
    data=historical_data,
    maximize=True
)
```

#### Genetic Algorithm

The genetic algorithm uses evolutionary principles to find optimal parameter values.

```python
from trading_platform.simulation.strategy_optimization import GeneticOptimizer

# Create genetic optimizer
optimizer = GeneticOptimizer()

# Define parameter space
param_space = {
    'fast_ma': (5, 50),
    'slow_ma': (50, 200),
    'stop_loss': (0.01, 0.05)
}

# Define fitness function
def fitness_function(params, data):
    # Implement strategy with given parameters
    # Return performance metric (e.g., Sharpe ratio)
    pass

# Run genetic optimization
results = optimizer.optimize(
    param_space=param_space,
    fitness_function=fitness_function,
    data=historical_data,
    population_size=50,
    generations=20,
    maximize=True
)
```

## Scenario Testing

The Scenario Testing component tests strategies under different market scenarios, helping to evaluate strategy robustness.

### Features

- Multiple scenario types (market crash, volatility spike, trend reversal)
- Customizable scenario parameters
- Comprehensive scenario reporting
- Integration with other simulation components

### Scenario Types

#### Market Crash

Simulates a sharp market decline with optional recovery phases.

```python
from trading_platform.simulation.scenario_testing import ScenarioGenerator, ScenarioTester

# Create scenario generator
generator = ScenarioGenerator()

# Generate market crash scenario
crash_scenario = generator.generate_market_crash(
    base_data=historical_data,
    crash_magnitude=0.2,  # 20% crash
    crash_duration=5,     # 5 days
    recovery=True,
    recovery_duration=20  # 20 days
)

# Test strategy under scenario
tester = ScenarioTester()
results = tester.test_strategy(
    strategy=my_strategy,
    scenario_data=crash_scenario
)
```

#### Volatility Spike

Increases price volatility for a specified duration.

```python
# Generate volatility spike scenario
volatility_scenario = generator.generate_volatility_spike(
    base_data=historical_data,
    volatility_multiplier=3.0,  # 3x normal volatility
    duration=10                 # 10 days
)

# Test strategy under scenario
results = tester.test_strategy(
    strategy=my_strategy,
    scenario_data=volatility_scenario
)
```

#### Trend Reversal

Reverses existing market trends.

```python
# Generate trend reversal scenario
reversal_scenario = generator.generate_trend_reversal(
    base_data=historical_data,
    reversal_magnitude=1.5,  # 1.5x trend magnitude
    duration=30              # 30 days
)

# Test strategy under scenario
results = tester.test_strategy(
    strategy=my_strategy,
    scenario_data=reversal_scenario
)
```

## Results Export

The Results Export component exports simulation results in various formats, allowing for further analysis and reporting.

### Features

- Multiple export formats (CSV, JSON, Excel, XML, database)
- Comprehensive export types (equity curve, trades, performance metrics)
- Metadata support for improved organization
- Integration with other simulation components

### Export Formats

#### CSV

```python
from trading_platform.simulation.results_export import ResultsExporter

# Create results exporter
exporter = ResultsExporter()

# Export equity curve to CSV
csv_path = exporter.export_equity_curve(
    equity_curve=equity_curve,
    format='csv',
    output_path='equity_curve.csv'
)

# Export trades to CSV
csv_path = exporter.export_trades(
    trades=trades,
    format='csv',
    output_path='trades.csv'
)

# Export performance metrics to CSV
csv_path = exporter.export_performance_metrics(
    equity_curve=equity_curve,
    trades=trades,
    format='csv',
    output_path='metrics.csv'
)
```

#### JSON

```python
# Export equity curve to JSON
json_path = exporter.export_equity_curve(
    equity_curve=equity_curve,
    format='json',
    output_path='equity_curve.json'
)

# Export trades to JSON
json_path = exporter.export_trades(
    trades=trades,
    format='json',
    output_path='trades.json'
)

# Export performance metrics to JSON
json_path = exporter.export_performance_metrics(
    equity_curve=equity_curve,
    trades=trades,
    format='json',
    output_path='metrics.json'
)
```

#### Complete Results

```python
# Export complete results to ZIP archive
zip_path = exporter.export_complete_results(
    equity_curve=equity_curve,
    trades=trades,
    charts=True,
    output_path='simulation_results.zip'
)
```

## Comparison Tools

The Comparison Tools component compares results from multiple simulations, helping to identify the best strategies or parameters.

### Features

- Comprehensive comparison of performance metrics
- Visualization of comparative performance
- Strategy ensemble creation
- Integration with other simulation components

### Usage

```python
from trading_platform.simulation.comparison_tools import SimulationComparer, StrategyEnsemble

# Create simulation comparer
comparer = SimulationComparer()

# Compare equity curves
fig_equity = comparer.compare_equity_curves(
    equity_curves={
        'Strategy A': equity_curve_a,
        'Strategy B': equity_curve_b,
        'Strategy C': equity_curve_c
    },
    benchmark_equity=benchmark_equity
)

# Compare performance metrics
metrics_df = comparer.compare_metrics(
    equity_curves={
        'Strategy A': equity_curve_a,
        'Strategy B': equity_curve_b,
        'Strategy C': equity_curve_c
    },
    trades_dict={
        'Strategy A': trades_a,
        'Strategy B': trades_b,
        'Strategy C': trades_c
    },
    benchmark_equity=benchmark_equity
)

# Create comparison report
report_path = comparer.create_comparison_report(
    equity_curves={
        'Strategy A': equity_curve_a,
        'Strategy B': equity_curve_b,
        'Strategy C': equity_curve_c
    },
    trades_dict={
        'Strategy A': trades_a,
        'Strategy B': trades_b,
        'Strategy C': trades_c
    },
    benchmark_equity=benchmark_equity
)

# Create strategy ensemble
ensemble = StrategyEnsemble()

# Create equal-weight ensemble
ensemble_equity = ensemble.create_equal_weight_ensemble(
    equity_curves={
        'Strategy A': equity_curve_a,
        'Strategy B': equity_curve_b,
        'Strategy C': equity_curve_c
    }
)

# Create optimized-weight ensemble
ensemble_equity, weights = ensemble.create_optimized_weight_ensemble(
    equity_curves={
        'Strategy A': equity_curve_a,
        'Strategy B': equity_curve_b,
        'Strategy C': equity_curve_c
    },
    optimization_metric='sharpe_ratio'
)
```

## User Interface

The User Interface component provides a web-based interface for interacting with the simulation system, making it accessible to users without programming knowledge.

### Features

- Account management
- Backtesting configuration and execution
- Strategy optimization
- Scenario testing
- Results comparison
- Performance visualization

### Usage

```python
from trading_platform.simulation.ui import SimulationWebUI

# Create web UI
ui = SimulationWebUI()

# Start web UI
ui.start(host='0.0.0.0', port=5000, open_browser=True)
```

### Interface Sections

#### Account Management

The Account Management section allows users to create, view, edit, and delete simulation accounts.

#### Backtesting

The Backtesting section allows users to configure and run backtests on historical data.

#### Optimization

The Optimization section allows users to optimize strategy parameters using various algorithms.

#### Scenario Testing

The Scenario Testing section allows users to test strategies under different market scenarios.

#### Comparison

The Comparison section allows users to compare results from multiple simulations.

## API Reference

The simulation system provides a comprehensive API for programmatic access to all features.

### Historical Data Replay API

```python
from trading_platform.simulation.historical_replay import ReplayController, ReplayEngine, DataSource

# Create and configure replay controller
controller = ReplayController()
controller.set_data_source(...)
controller.set_replay_parameters(...)
controller.register_handler(...)
controller.start()
controller.pause()
controller.resume()
controller.stop()

# Access replay engine directly
engine = ReplayEngine()
engine.load_data(...)
engine.set_speed(...)
engine.start()
engine.pause()
engine.stop()

# Create custom data source
class MyDataSource(DataSource):
    def __init__(self, ...):
        super().__init__(...)
    
    def load_data(self, ...):
        # Implement data loading
        pass
    
    def get_data(self, ...):
        # Implement data retrieval
        pass
```

### Performance Analytics API

```python
from trading_platform.simulation.performance_analytics import PerformanceMetrics, PerformanceVisualizer, PerformanceReporter

# Calculate performance metrics
metrics = PerformanceMetrics()
results = metrics.calculate_all_metrics(...)
returns = metrics.calculate_returns(...)
volatility = metrics.calculate_volatility(...)
drawdowns = metrics.calculate_drawdowns(...)
sharpe = metrics.calculate_sharpe_ratio(...)

# Visualize performance
visualizer = PerformanceVisualizer()
fig_equity = visualizer.plot_equity_curve(...)
fig_drawdowns = visualizer.plot_drawdowns(...)
fig_returns = visualizer.plot_returns_distribution(...)
fig_monthly = visualizer.plot_monthly_returns_heatmap(...)

# Generate performance report
reporter = PerformanceReporter()
report_path = reporter.generate_report(...)
```

### Strategy Optimization API

```python
from trading_platform.simulation.strategy_optimization import ParameterOptimizer, GeneticOptimizer

# Use parameter optimizer
optimizer = ParameterOptimizer()
results = optimizer.grid_search(...)
results = optimizer.random_search(...)
results = optimizer.walk_forward_optimization(...)

# Use genetic optimizer
optimizer = GeneticOptimizer()
results = optimizer.optimize(...)
```

### Scenario Testing API

```python
from trading_platform.simulation.scenario_testing import ScenarioGenerator, ScenarioTester

# Generate scenarios
generator = ScenarioGenerator()
crash_scenario = generator.generate_market_crash(...)
volatility_scenario = generator.generate_volatility_spike(...)
reversal_scenario = generator.generate_trend_reversal(...)

# Test scenarios
tester = ScenarioTester()
results = tester.test_strategy(...)
results = tester.test_multiple_scenarios(...)
```

### Results Export API

```python
from trading_platform.simulation.results_export import ResultsExporter

# Export results
exporter = ResultsExporter()
path = exporter.export_equity_curve(...)
path = exporter.export_trades(...)
path = exporter.export_performance_metrics(...)
path = exporter.export_complete_results(...)
```

### Comparison Tools API

```python
from trading_platform.simulation.comparison_tools import SimulationComparer, StrategyEnsemble

# Compare simulations
comparer = SimulationComparer()
fig = comparer.compare_equity_curves(...)
df = comparer.compare_metrics(...)
path = comparer.create_comparison_report(...)

# Create ensembles
ensemble = StrategyEnsemble()
equity = ensemble.create_equal_weight_ensemble(...)
equity, weights = ensemble.create_optimized_weight_ensemble(...)
equity = ensemble.create_dynamic_weight_ensemble(...)
```

### User Interface API

```python
from trading_platform.simulation.ui import SimulationAccountManager, SimulationWebUI

# Manage accounts
manager = SimulationAccountManager()
account_id = manager.create_account(...)
account = manager.get_account(...)
account = manager.update_account(...)
success = manager.delete_account(...)
accounts = manager.list_accounts(...)

# Start web UI
ui = SimulationWebUI()
ui.start(...)
```

## Configuration

The simulation system can be configured through a global configuration object or by passing configuration parameters directly to component constructors.

### Global Configuration

```python
from trading_platform.simulation.config import Config

# Create configuration
config = Config({
    'historical_replay.data_source': 'database',
    'historical_replay.database.host': 'localhost',
    'historical_replay.database.port': 8086,
    'historical_replay.database.name': 'market_data',
    'historical_replay.default_speed': 1.0,
    'historical_replay.include_market_events': True,
    
    'performance_analytics.annualization_factor': 252,
    'performance_analytics.risk_free_rate': 0.02,
    
    'strategy_optimization.parallel_jobs': 4,
    'strategy_optimization.default_metric': 'sharpe_ratio',
    
    'scenario_testing.default_scenarios': ['market_crash', 'volatility_spike'],
    
    'results_export.default_format': 'csv',
    'results_export.output_directory': 'simulation_results',
    
    'comparison_tools.default_metrics': ['sharpe_ratio', 'max_drawdown', 'win_rate'],
    
    'ui.host': '0.0.0.0',
    'ui.port': 5000,
    'ui.static_directory': 'static',
    'ui.templates_directory': 'templates'
})

# Set as global configuration
from trading_platform.simulation import set_global_config
set_global_config(config)
```

### Component-Specific Configuration

```python
# Create component with specific configuration
from trading_platform.simulation.historical_replay import ReplayController

controller = ReplayController(config=Config({
    'data_source': 'database',
    'database.host': 'localhost',
    'database.port': 8086,
    'database.name': 'market_data',
    'default_speed': 1.0,
    'include_market_events': True
}))
```

## Troubleshooting

### Common Issues

#### Historical Data Replay

- **Issue**: Replay is too slow or too fast
  - **Solution**: Adjust the replay speed using `controller.set_speed()`

- **Issue**: Data is missing or incomplete
  - **Solution**: Check the data source configuration and ensure the data is available

#### Performance Analytics

- **Issue**: Performance metrics are incorrect
  - **Solution**: Ensure the equity curve and trades data are correctly formatted

- **Issue**: Visualizations are not displaying correctly
  - **Solution**: Check the matplotlib configuration and ensure the figure size is appropriate

#### Strategy Optimization

- **Issue**: Optimization is taking too long
  - **Solution**: Reduce the parameter space or use random search instead of grid search

- **Issue**: Optimization results are inconsistent
  - **Solution**: Use a fixed random seed for reproducibility

#### Scenario Testing

- **Issue**: Scenarios are not realistic
  - **Solution**: Adjust the scenario parameters to match historical market behavior

- **Issue**: Strategy performance is poor under all scenarios
  - **Solution**: Review the strategy logic and consider more robust approaches

#### Results Export

- **Issue**: Export fails with permission error
  - **Solution**: Ensure the output directory is writable

- **Issue**: Exported files are too large
  - **Solution**: Export only the necessary data or use compression

#### Comparison Tools

- **Issue**: Comparison charts are cluttered
  - **Solution**: Compare fewer strategies or use a subset of metrics

- **Issue**: Ensemble creation fails
  - **Solution**: Ensure all equity curves have the same index

#### User Interface

- **Issue**: Web UI is not accessible
  - **Solution**: Check the host and port configuration and ensure the server is running

- **Issue**: Forms are not submitting correctly
  - **Solution**: Check the form validation and ensure all required fields are filled

### Getting Help

If you encounter issues not covered in this documentation, please:

1. Check the API reference for detailed information on component usage
2. Review the example code in the `examples` directory
3. Search the issue tracker for similar problems
4. Contact support at support@tradingplatform.com

## Conclusion

The Trading Platform Simulation System provides a comprehensive environment for backtesting, optimizing, and evaluating trading strategies. By following this documentation, you should be able to effectively use all features of the simulation system to develop and test your trading strategies.
