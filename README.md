# SimpleCryptoPractice

## 實踐功能
### 部署一ERC20合約於Goerli測試網
### 透過Golang操作Smart contract
  ```
  go run main.go server
  ```
  * 查詢總發行量 TotalSupply
  * 合約擁有者 Owner
  * Token名稱 Name
  * Token代號 Symbol
  * 查訊地址持有Token數量 BalanceOf
  * 造幣 Mint
  * 授權轉移數量 Approve
  * 轉移授權的Token TransferFrom
  * 轉移Token Transfer
  
### 串接MetaMask以取得主要操作地址
  * 因操作地址為使用MetaMask建立的地址，需另外取出地址密鑰簽署以完成交易
  
### 使用之前端為MetaMask官方提供
  ```
   yarn serve
  ```
