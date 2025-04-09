// PlaceOrder places a new order with the XTS Client API
// This implementation includes the apiOrderSource parameter which is specific to XTS Client
func (c *XTSClientImpl) PlaceOrder(order *common.Order) (*common.OrderResponse, error) {
	if c.token == "" {
		return nil, errors.New("not logged in")
	}
	
	if order == nil {
		return nil, errors.New("order is required")
	}
	
	url := fmt.Sprintf("%s/interactive/orders", c.baseURL)
	
	// Prepare the request parameters
	params := url.Values{}
	params.Set("exchangeSegment", order.ExchangeSegment)
	params.Set("exchangeInstrumentID", order.ExchangeInstrumentID)
	params.Set("productType", order.ProductType)
	params.Set("orderType", order.OrderType)
	params.Set("orderSide", order.OrderSide)
	params.Set("timeInForce", order.TimeInForce)
	params.Set("disclosedQuantity", fmt.Sprintf("%d", order.DisclosedQuantity))
	params.Set("orderQuantity", fmt.Sprintf("%d", order.OrderQuantity))
	params.Set("limitPrice", fmt.Sprintf("%f", order.LimitPrice))
	params.Set("stopPrice", fmt.Sprintf("%f", order.StopPrice))
	params.Set("orderUniqueIdentifier", order.OrderUniqueIdentifier)
	
	// Add the apiOrderSource parameter which is specific to XTS Client
	if order.APIOrderSource != "" {
		params.Set("apiOrderSource", order.APIOrderSource)
	} else {
		// Default value if not provided
		params.Set("apiOrderSource", "WEBAPI")
	}
	
	// Add clientID parameter if provided and not an investor client
	if order.ClientID != "" && !c.isInvestor {
		params.Set("clientID", order.ClientID)
	}
	
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	req.URL.RawQuery = params.Encode()
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", c.token)
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	
	var response struct {
		Type        string `json:"type"`
		Code        int    `json:"code"`
		Description string `json:"description"`
		Result      struct {
			AppOrderID       string `json:"AppOrderID"`
			OrderGeneratedID string `json:"OrderGeneratedID"`
			OrderStatus      string `json:"OrderStatus"`
			OrderRejected    bool   `json:"OrderRejected"`
			RejectReason     string `json:"RejectReason"`
		} `json:"result"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	
	if response.Type != "success" {
		return nil, fmt.Errorf("place order failed: %s", response.Description)
	}
	
	// Convert the response to the common OrderResponse model
	orderResponse := &common.OrderResponse{
		OrderID:         response.Result.AppOrderID,
		ExchangeOrderID: response.Result.OrderGeneratedID,
		Status:          response.Result.OrderStatus,
	}
	
	if response.Result.OrderRejected {
		orderResponse.Status = "REJECTED"
		orderResponse.RejectionReason = response.Result.RejectReason
	}
	
	return orderResponse, nil
}
