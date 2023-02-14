// SPDX-License-Identifier: GPL-3.0
pragma solidity 0.8.17;

contract Token {
    // Token 名字 (幣的名字)11
    string public name = "Anonymous Mask Token";
    // 暱稱
    string public symbol = "AMT";
    // 使用的小數點位
    uint8 decimals = 5;
    address public owner;

    // 設定合約擁有者
    constructor() {
        owner = msg.sender;
    }

    // Event
    // 轉移事件，代幣轉移時都必須觸發，包含0值轉移，創建代幣時也須觸發事件，_from設置為0x000
    event  Transfer(address indexed _from, address indexed _to, uint256 _value);
    // 成功事件，任何成功的調用都要觸發
    event  Approval(address indexed _onwer, address indexed _spender,uint256 _value);


    //返回該帳戶此代幣的餘額，帳號餘額map like map[address]uin256 in golang
    mapping (address => uint256) public balanceOf;

    // A地址允許B地址提取的代幣數量
    // 返回_owner批准_spender的提取數量
    mapping (address => mapping(address => uint256)) public allowance;

    // 代幣總數，返回代幣的總供應量，此種寫法包含宣告變數及產生簡單的get func
    uint256 public totalSupply;
    /*
    效用同此
        uint256 totalSupply;
        function totalSupply() public view returns (uint256){
            return totalSupply;
        }
    */
    
    // 轉移代幣至_to地址，觸發Transfer event & throw error
    function transfer(address _to, uint256 _value) public returns (bool success){
        return transferFrom(msg.sender,_to,_value);
    }

    // 從_from轉移代幣至_to地址，觸發Transfer event & throw error
    function transferFrom(address _from,address _to, uint256 _value) public returns (bool success){
        // 檢查from有無足夠代幣，沒有的話throw error 
        require(balanceOf[_from] >= _value,"not enough token!");


        // 確認當前呼叫合約者是否為from
        if (msg.sender != _from) {
            // 此情況為第三者要轉移from地址的代幣到to
            // 確認from授予msg.sender的提領代幣數量有無大於此次轉移顆數
            require(allowance[_from][msg.sender] >= _value);
        }

        balanceOf[_from] -= _value;
        balanceOf[_to] += _value;

	    emit Transfer(_from,_to,_value);

        return true;
    }
    // 批准提領數，設置_spender可提取的代幣數量 
    function approve(address _spender, uint256 _value) public returns (bool success) {
        // @ msg.sender => 合約呼叫者
        // @ _spender => 提領者
        // msg.sender 允許 _spender 提領 value 顆 AMT
        allowance[msg.sender][_spender] = _value;

        // 觸發事件
	    emit Approval(msg.sender,_spender,_value);

        return true;
    }

    // 造幣
    function mint(address _to,uint value) public {
        
        balanceOf[msg.sender] += value;
        totalSupply += value;
        emit Transfer(address(0),_to,value);
    }
    
    // 摧毀合約
    function destroySmartContract(address payable _to) public {
        require(msg.sender == owner,"You are not the owner");
        selfdestruct(_to);
    }

}
