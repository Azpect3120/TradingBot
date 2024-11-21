# Azpect's Trading Bot


## How is the Rating Calculated

Squeeze is generated using JHF's SqueezePro indicator:

`explain the squeeze formula here`


Rating is scored out of 100 points.
Current Squeeze: 15/100
Past squeeze: 60/100
50H MA: 15/100
9H MA: 10/100

Total: 100/100

WIP: 
- Swings and relative locations to swings
- Support and resistance levels

1. The current hourly squeeze. (15 possible points)

    1. Ultra tight squeeze:    +15  points
    2. Tight squeeze:          +10  points
    3. Normal squeeze:         +5   points
    4. Wide squeeze:           +2.5 points
    5. No squeeze:             -100 points *This is an automatic fail, this will never appear on my screener*

2. The past hourly squeeze. (60 possible points)

    - Long/short positions:
        1. Ultra tight squeeze (per last 14):    +5    point(s) (max 25)
        2. Tight squeeze (per last 14):          +3    point(s) (max 15)
        3. Normal squeeze (per last 14):         +1    point(s) (max 5)
        4. Wide squeeze (per last 14):           +0.5  point(s) (max 2.5)
        5. No squeeze (per last 14):             -1    point(s) (max -14)
        6. Squeeze continuation >= wide(per 7):  +0.75 point(s) (max 7.5) *cannot be just wide squeezes, if they're all wide, then the current 7 is 0 and the calculation is halted*
        7. Squeeze is increasing (last 7):       +5    points 
        8. Squeeze is decreasing (last 7):       -5    points
        9. Squeeze is constant (last 7):         +2.5  points

3. The 50 hour moving average. (15 possible points)

    - For long positions:
        1. Price is above 50H MA (most recent):  +5   points
        2. Price is below 50H MA (most recent):  -2.5 points
        3. Price has crossed up 50H MA (past 7): +2.5 points
        4. Price has crossed dn 50H MA (past 7): -2.5 points
        5. Price has been above 50H MA (per 7):  +0.5 point(s) (max 2.5)
        6. 50H MA is below lows of (past 21):    +2.5 points
        7. 50H MA is increasing (per 7):         +0.5 point(s) (max 2.5)

    - For short positions:
        1. Price is above 50H MA (most recent):  -2.5 points
        2. Price is below 50H MA (most recent):  +5   points
        3. Price has crossed up 50H MA (past 7): -2.5 points
        4. Price has crossed dn 50H MA (past 7): +2.5 points
        5. Price has been bellow 50H MA (per 7): +0.5 point(s) (max 2.5)
        6. 50H MA is above highs of (past 21):   +2.5 points
        7. 50H MA is decreasing (per 7):         +0.5 point(s) (max 2.5)

4. The 9 hour moving average. (10 possible points)

    - For long positions:
        1. Price is above 9H MA (most recent):  +3   points
        2. Price is below 9H MA (most recent):  -2   points
        3. Price has crossed up 9H MA (past 7): +2   points
        4. Price has crossed dn 9H MA (past 7): -2   points
        5. 9H MA is btw high and low (past 7):  +2   points
        6. Price has been above 9H MA (per 7):  +0.5 point(s) (max 1.5)
        7. 9H MA is increasing (per 7):         +0.5 point(s) (max 1.5)

    - For short positions:
        1. Price is above 9H MA (most recent):  -2   points
        2. Price is below 9H MA (most recent):  +3   points
        3. Price has crossed up 9H MA (past 7): -2   points
        4. Price has crossed dn 9H MA (past 7): +2   points
        5. 9H MA is btw high and low (past 7):  +2   points
        6. Price has been below 9H MA (per 7):  +0.5 point(s) (max 1.5)
        7. 9H MA is decreasing (per 7):         +0.5 point(s) (max 1.5)
