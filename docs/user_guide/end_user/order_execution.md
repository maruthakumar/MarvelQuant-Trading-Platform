# Order Execution

## Introduction

This guide covers all aspects of order execution in the Trading Platform, including different order types, advanced order features, and order management. Understanding these concepts is essential for effective trading and implementing your investment strategies.

## Order Types

The Trading Platform supports a variety of order types to accommodate different trading strategies and market conditions.

### Market Order

A market order is an instruction to buy or sell a security immediately at the best available current price.

**When to use:**
- When immediate execution is more important than the exact price
- In highly liquid markets where price slippage is minimal
- When you want to ensure your order is filled

**How to place a market order:**
1. Navigate to the Order Execution screen
2. Select the instrument you want to trade
3. Choose "Market" from the order type dropdown
4. Enter the quantity you wish to buy or sell
5. Select Buy or Sell
6. Review the order details
7. Click "Submit Order"

**Example:**
If you want to buy 100 shares of AAPL immediately at the current market price:
- Instrument: AAPL
- Order Type: Market
- Quantity: 100
- Action: Buy

**Considerations:**
- Market orders are executed immediately during market hours
- The execution price may differ from the last traded price due to market movement
- In volatile markets, price slippage can be significant
- Market orders are guaranteed to execute if the market is open and the security is trading

### Limit Order

A limit order is an instruction to buy or sell a security at a specified price or better.

**When to use:**
- When you want to specify the maximum price you're willing to pay or minimum price you're willing to accept
- When you're not in a hurry to execute the trade
- When you want to control the execution price

**How to place a limit order:**
1. Navigate to the Order Execution screen
2. Select the instrument you want to trade
3. Choose "Limit" from the order type dropdown
4. Enter the quantity you wish to buy or sell
5. Enter your limit price
6. Select Buy or Sell
7. Choose the time-in-force (how long the order remains active)
8. Review the order details
9. Click "Submit Order"

**Example:**
If you want to buy 100 shares of AAPL, but only if the price drops to $150 or lower:
- Instrument: AAPL
- Order Type: Limit
- Quantity: 100
- Limit Price: $150.00
- Action: Buy
- Time-in-force: Day

**Considerations:**
- Limit orders are not guaranteed to execute
- The order will only execute if the market price reaches your limit price
- Partial fills are possible if there's not enough volume at your limit price
- Limit orders provide price protection but may result in missed opportunities if the market moves away from your limit price

### Stop Order

A stop order (also called a stop-loss order) is an instruction to buy or sell a security when its price reaches a specified stop price. Once triggered, it becomes a market order.

**When to use:**
- To limit potential losses on existing positions
- To protect profits on existing positions
- To enter a position when the price breaks above or below a certain level

**How to place a stop order:**
1. Navigate to the Order Execution screen
2. Select the instrument you want to trade
3. Choose "Stop" from the order type dropdown
4. Enter the quantity you wish to buy or sell
5. Enter your stop price
6. Select Buy or Sell
7. Choose the time-in-force
8. Review the order details
9. Click "Submit Order"

**Example:**
If you own 100 shares of AAPL and want to sell if the price drops to $140 or lower:
- Instrument: AAPL
- Order Type: Stop
- Quantity: 100
- Stop Price: $140.00
- Action: Sell
- Time-in-force: Good Till Canceled (GTC)

**Considerations:**
- Once triggered, a stop order becomes a market order and is subject to market conditions
- The execution price may be different from the stop price, especially in volatile markets
- Stop orders are not guaranteed to limit losses to the specified price
- Stop orders are not displayed in the order book until triggered

### Stop-Limit Order

A stop-limit order combines features of both stop and limit orders. When the stop price is reached, it activates a limit order rather than a market order.

**When to use:**
- When you want to control both the trigger price and the execution price
- When you want to avoid excessive slippage in volatile markets
- When you need precise control over your entry or exit price

**How to place a stop-limit order:**
1. Navigate to the Order Execution screen
2. Select the instrument you want to trade
3. Choose "Stop-Limit" from the order type dropdown
4. Enter the quantity you wish to buy or sell
5. Enter your stop price (trigger price)
6. Enter your limit price (maximum/minimum execution price)
7. Select Buy or Sell
8. Choose the time-in-force
9. Review the order details
10. Click "Submit Order"

**Example:**
If you own 100 shares of AAPL and want to sell if the price drops to $140, but not below $138:
- Instrument: AAPL
- Order Type: Stop-Limit
- Quantity: 100
- Stop Price: $140.00
- Limit Price: $138.00
- Action: Sell
- Time-in-force: Good Till Canceled (GTC)

**Considerations:**
- Stop-limit orders provide more control but are not guaranteed to execute
- If the market price gaps beyond your limit price, your order may not execute
- The limit portion of the order works like any other limit order
- This order type requires setting both a stop price and a limit price

## Advanced Order Features

### Bracket Orders

A bracket order is a set of three orders that help you enter a position and then manage it with a take-profit and a stop-loss order.

**Components:**
1. **Primary Order**: The initial order to enter a position (market or limit)
2. **Take-Profit Order**: A limit order to exit with a profit
3. **Stop-Loss Order**: A stop or stop-limit order to limit potential losses

**How to place a bracket order:**
1. Navigate to the Order Execution screen
2. Select the instrument you want to trade
3. Choose "Bracket" from the order type dropdown
4. Configure your primary order (type, quantity, price if applicable)
5. Set your take-profit level (price or percentage)
6. Set your stop-loss level (price or percentage)
7. Review all three components of the order
8. Click "Submit Order"

**Example:**
You want to buy 100 shares of AAPL at market price, take profit at $165, and limit losses at $145:
- Primary Order: Buy 100 AAPL at Market
- Take-Profit: Sell 100 AAPL at Limit $165.00
- Stop-Loss: Sell 100 AAPL at Stop $145.00

**Considerations:**
- When either the take-profit or stop-loss order executes, the other is automatically canceled
- All three components are submitted simultaneously
- You can modify the take-profit and stop-loss levels while the position is open
- Bracket orders help enforce disciplined trading by setting exit points in advance

### One-Cancels-Other (OCO) Orders

An OCO order is a pair of orders with the condition that if one order executes, the other is automatically canceled.

**When to use:**
- When you want to set both upside and downside exit strategies
- When you're looking to enter a position at one of two possible price points
- When you need to implement "either/or" trading decisions

**How to place an OCO order:**
1. Navigate to the Order Execution screen
2. Select the instrument you want to trade
3. Choose "OCO" from the order type dropdown
4. Configure the first order (type, quantity, price)
5. Configure the second order (type, quantity, price)
6. Review both components of the order
7. Click "Submit Order"

**Example:**
You own 100 shares of AAPL and want to either sell if it reaches $170 (take profit) or if it drops to $145 (stop loss):
- First Order: Sell 100 AAPL at Limit $170.00
- Second Order: Sell 100 AAPL at Stop $145.00

**Considerations:**
- Only one of the two orders will execute; the other will be canceled
- Both orders must be for the same instrument and quantity
- OCO orders can combine different order types (limit, stop, stop-limit)
- This order type helps manage risk in uncertain market conditions

### Trailing Stop Orders

A trailing stop order is a stop order that automatically adjusts the stop price as the market price moves in your favor.

**When to use:**
- To protect profits while allowing for further gains
- When you want your stop level to adjust automatically with market movement
- For trending markets where you want to "let profits run"

**How to place a trailing stop order:**
1. Navigate to the Order Execution screen
2. Select the instrument you want to trade
3. Choose "Trailing Stop" from the order type dropdown
4. Enter the quantity you wish to buy or sell
5. Set the trailing amount (fixed amount or percentage)
6. Select Buy or Sell
7. Choose the time-in-force
8. Review the order details
9. Click "Submit Order"

**Example:**
You own 100 shares of AAPL currently trading at $160, and want to set a trailing stop 5% below the market price:
- Instrument: AAPL
- Order Type: Trailing Stop
- Quantity: 100
- Trailing Amount: 5%
- Action: Sell
- Time-in-force: Good Till Canceled (GTC)

Initially, your stop price would be at $152 (5% below $160). If AAPL rises to $170, your stop price would automatically adjust to $161.50 (5% below $170).

**Considerations:**
- The stop price adjusts only when the market moves in your favor
- The stop price never moves in the opposite direction
- Once triggered, a trailing stop becomes a market order
- Trailing stops can be set as a fixed amount or a percentage

### Iceberg Orders

An iceberg order is a large order that is divided into smaller, visible portions to minimize market impact.

**When to use:**
- When trading large quantities that might move the market
- When you want to hide the true size of your order
- When you need to execute a large position gradually

**How to place an iceberg order:**
1. Navigate to the Order Execution screen
2. Select the instrument you want to trade
3. Choose "Iceberg" from the order type dropdown
4. Enter the total quantity you wish to buy or sell
5. Set the visible quantity (the portion visible to the market)
6. Enter your limit price
7. Select Buy or Sell
8. Choose the time-in-force
9. Review the order details
10. Click "Submit Order"

**Example:**
You want to buy 10,000 shares of AAPL at $155, but only want to show 500 shares at a time:
- Instrument: AAPL
- Order Type: Iceberg
- Total Quantity: 10,000
- Visible Quantity: 500
- Limit Price: $155.00
- Action: Buy
- Time-in-force: Day

**Considerations:**
- Only the visible portion appears in the order book
- As the visible portion executes, another portion automatically becomes visible
- Iceberg orders are always limit orders
- Some exchanges charge additional fees for iceberg orders

## Order Management

### Viewing Open Orders

To view your current open orders:

1. Navigate to the "Order Book" tab in the Order Execution screen
2. All open orders are displayed with the following information:
   - Order ID
   - Instrument
   - Order Type
   - Quantity
   - Price (if applicable)
   - Status
   - Submission Time
   - Expiration (for orders with time-in-force)

You can filter and sort the order book by various criteria to quickly find specific orders.

### Modifying Orders

To modify an existing open order:

1. Locate the order in the Order Book
2. Click the "Modify" button or right-click and select "Modify"
3. A modification form will appear with the current order details
4. Change the parameters you wish to modify (price, quantity, etc.)
5. Click "Submit Changes"

**Modifiable parameters:**
- Price (for limit, stop, and stop-limit orders)
- Quantity (can be decreased, but increasing may require canceling and replacing the order)
- Time-in-force
- Take-profit and stop-loss levels (for bracket orders)
- Trailing amount (for trailing stop orders)

**Considerations:**
- Some modifications may change your order priority in the exchange's order book
- Not all parameters can be modified for all order types
- Orders that are in the process of executing cannot be modified
- Some exchanges treat significant modifications as new orders

### Cancelling Orders

To cancel an existing open order:

1. Locate the order in the Order Book
2. Click the "Cancel" button or right-click and select "Cancel"
3. Confirm the cancellation when prompted

You can also cancel multiple orders at once:
1. Select multiple orders using Ctrl+Click (Cmd+Click on Mac)
2. Right-click and select "Cancel Selected Orders"
3. Confirm the cancellation when prompted

To cancel all open orders for a specific instrument:
1. Navigate to the instrument's detail page
2. Click "Cancel All Orders" in the Orders section
3. Confirm the cancellation when prompted

**Considerations:**
- Cancellation is not guaranteed if the order is in the process of executing
- There may be a brief delay between your cancellation request and confirmation
- The platform will notify you if a cancellation request is rejected
- Canceled orders appear in your order history with a "Canceled" status

### Order Status Tracking

The platform provides real-time updates on the status of your orders:

- **Created**: Order has been created in the platform but not yet sent to the exchange
- **Sent**: Order has been transmitted to the exchange
- **Acknowledged**: Exchange has confirmed receipt of the order
- **Partially Filled**: Some but not all of the order quantity has been executed
- **Filled**: Order has been completely executed
- **Canceled**: Order has been canceled before execution
- **Rejected**: Order was not accepted by the exchange or broker
- **Expired**: Order reached its time-in-force limit without executing

You can set up notifications for order status changes in your account preferences.

### Order History

To view your historical orders:

1. Navigate to the "Order History" tab in the Order Execution screen
2. Set the date range for the orders you want to view
3. Apply any filters (instrument, order type, status)
4. View detailed information about past orders

The order history provides comprehensive information:
- Order details (instrument, type, quantity, price)
- Execution details (execution time, execution price)
- Fill details for partially filled orders
- Reason for rejection or cancellation
- Associated fees and commissions

You can export your order history to CSV or Excel format for record-keeping or analysis.

## Best Practices for Order Execution

### Pre-Trade Checklist

Before submitting any order, consider the following:

1. **Verify the instrument**: Double-check the symbol to ensure you're trading the correct security
2. **Confirm the order type**: Select the appropriate order type for your trading strategy
3. **Check quantity and price**: Verify that the quantity and price are as intended
4. **Review your account**: Ensure you have sufficient funds or margin for the trade
5. **Consider market conditions**: Be aware of current volatility and liquidity
6. **Check for corporate actions**: Be aware of any upcoming dividends, splits, or other events
7. **Review risk parameters**: Ensure the trade aligns with your risk management strategy

### Risk Management Strategies

Implement these risk management techniques in your order execution:

1. **Position Sizing**: Limit each position to a small percentage of your portfolio
2. **Stop-Loss Orders**: Always use stop-loss orders to limit potential losses
3. **Take-Profit Targets**: Set realistic profit targets based on technical or fundamental analysis
4. **Diversification**: Spread your trades across different instruments and sectors
5. **Correlation Awareness**: Be mindful of correlated positions that could amplify losses
6. **Volatility Adjustment**: Adjust position sizes based on the instrument's volatility
7. **Scenario Planning**: Consider what could go wrong and have contingency plans

### Common Mistakes to Avoid

Be aware of these common order execution pitfalls:

1. **Fat Finger Errors**: Accidentally entering the wrong price or quantity
2. **Chasing the Market**: Repeatedly modifying orders to chase a moving price
3. **Overtrading**: Excessive trading that increases costs and can lead to emotional decisions
4. **Ignoring Liquidity**: Not considering the impact of your order on thinly traded markets
5. **Misusing Order Types**: Using inappropriate order types for the current market conditions
6. **Neglecting Time Zones**: Forgetting about market hours in different time zones
7. **Emotional Trading**: Making impulsive decisions based on fear or greed

## Troubleshooting

### Common Order Issues

**Order Rejected**
- **Possible Causes**: Insufficient funds, invalid price, market closed, trading halted
- **Solution**: Check the rejection reason in the order details and address the specific issue

**Order Stuck in "Sent" Status**
- **Possible Causes**: Communication issue with exchange, system processing delay
- **Solution**: Wait a few minutes, then contact support if the status doesn't update

**Unexpected Partial Fill**
- **Possible Causes**: Insufficient liquidity at your price level
- **Solution**: Wait for additional fills or modify the remaining quantity

**Price Slippage on Market Orders**
- **Possible Causes**: Fast-moving markets, low liquidity
- **Solution**: Consider using limit orders for more price control

### Getting Help with Orders

If you encounter issues with order execution:

1. **Check the Knowledge Base**: Many common issues are addressed in our help center
2. **Contact Trading Support**: Our trading desk can assist with order-related issues
   - Email: trading@tradingplatform.example.com
   - Phone: +1-800-TRADING ext. 2
   - Live Chat: Available during market hours

3. **Submit a Support Ticket**: For complex issues, create a detailed support ticket
   - Include order IDs and screenshots
   - Describe the issue and steps you've already taken
   - Note any error messages you received

## Next Steps

Now that you understand order execution, explore these related guides:

- [Portfolio Management](./portfolio_management.md) - Learn how to track and analyze your positions
- [Strategy Management](./strategy_management.md) - Discover how to create and implement trading strategies
- [Risk Management](./risk_management.md) - Master techniques for protecting your capital
- [Advanced Order Strategies](./advanced_order_strategies.md) - Explore sophisticated order combinations
