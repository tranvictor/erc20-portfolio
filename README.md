# erc20-portfolio

Monitors transaction done by an address to build ERC20 portfolio report for it.

## How it works

1. Using etherscan to get the list of transaction done by the address (trade txs, transfer tx)
2. Store them so it will not have to requery from etherscan for old data
3. Analyze though those txs' receipt in order to see which token trade with
