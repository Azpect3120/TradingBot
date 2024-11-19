## Steps

### Get a list of stocks from the [[API]]
- HTTP api? Or does Go have bindings?
- There doesn't appear to be bindings

### Loop through the list of stocks
- Filter for stocks that are within the criteria

### Get stock data from the last 50-200 (periods)
- Convert this data into an array of prices

### Use price data to generate indicator results
- Squeeze value
- MA relation (above below, at)
- Recent highs and lows
- Momentum *(If I can find the calculation JHF uses)*

### Store stocks that fit into a certain criteria
- Good squeeze
- MA relation
- Price (> 5)
- Volume (> 2M)
- Market cap (>= 10B)
