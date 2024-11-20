# Azpect's Trading Bot


## How is the Rating Calculated

Squeeze is generated using JHF's SqueezePro indicator:

`explain the squeeze formula here`

1. Squeeze of the last 24 hours (3 days) is calculated
2. If the squeeze at the last hours generated is still active:
    1. Generate another 24 hours (3 days) of squeeze
    2. Test active squeeze again and repeat until the squeeze is not active
    3. Save the number of hours the squeeze was active for and which values 
       are found. eg. how many tight, how many normal...
3. Longer squeeze will increase the rating. Tighter squeeze will increase the 
   rating. If the squeeze was tight but loosened, the rating will not be as 
   high as if the squeeze was tight and stayed tight until recent time.

