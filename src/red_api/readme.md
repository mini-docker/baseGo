### 红包api项目

端口号:8866

## 请求加密说明
- api使用rsa加密，加密方法为整个body加密
- header内需要传递的参数：



   
   |字段|类型|说明|
   |:---|:----:|:------:|
   |lineId|string|线路id|
   |agencyId|string|代理id|
   |sign|string|签名|
   |Content-Type|string|链接方式|
   
     
- 关于Content-Type:填写text/plain

- 关于sign签名:整个body使用rsa加密后拼接盐字符串的md5值，公式为Md5(base64(请求参数rsa加密后密文) + salt)  

   
- RSA公共密钥


-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDfw1/P15GQzGGYvNwVmXIGGxea
8Pb2wJcF7ZW7tmFdLSjOItn9kvUsbQgS5yxx+f2sAv1ocxbPTsFdRc6yUTJdeQol
DOkEzNP0B8XKm+Lxy4giwwR5LJQTANkqe4w/d9u129bRhTu/SUzSUIr65zZ/s6TU
GQD6QzKY1Y8xS+FoQQIDAQAB
-----END PUBLIC KEY----- 
 
 - salt
 
 9ad0749d493bd15c12118ed43519e920
 
 
请求示例：

POST http://localhost:8866/api/login


R17OuR1f1vf4r0hXRq86GuuRQvGHy3aZLert2BbLWEGCZDfGS6/PBdmrg64VZBegeoryjxRKXXr8XIoiuPkMf+fEvZRi8S0ayjDO9R8EfkVgSMUoNRahQbLuSvzNDgvUNRDyn1pKhORcoKVaqUm18viqWTQGX3BxiwIeQNUovy0=


### 登陆

- 接口地址: `/api/login`
- Method: POST
- 负责人: blazk

- url参数

|字段|类型|必填项|说明|
|:---|:----:|----:|:------:|
|account|string|true|账号|
|password|string|true|密码|
|ip|string|true|登陆ip|


- 返回值



data内数据

|字段|类型|说明|
|:---|:----:|:------:|
|sessionId|string|登陆成功后的sessionId|
|gameUrl|string|游戏登陆的url|



其他字段的数组元素结构为:

|字段|类型|说明|
|:---|:----:|:------:|
|code|int|状态码(0成功，5001账号不存在，5011密码错误,6000账号被停用,6001登陆失败)|
|message|string|状态信息|


返回示例:

```
{
    "code": 0,
    "message": "success",
    "data": {
        "sessionId": "11f2930d-5302-4fb5-87aa-7b9c9827c1ee_12_0",
        "gameUrl": "http://redfront.pkbeta.com"
    }
}

```

### 注册并登陆

- 接口地址: `/api/register`
- Method: POST
- 负责人: blazk

- url参数

|字段|类型|必填项|说明|
|:---|:----:|----:|:------:|
|account|string|true|账号|
|password|string|true|密码|
|ip|string|true|登陆ip|


- 返回值




data内数据

|字段|类型|说明|
|:---|:----:|:------:|
|sessionId|string|登陆成功后的sessionId|
|gameUrl|string|游戏登陆的url|



其他字段的数组元素结构为:

|字段|类型|说明|
|:---|:----:|:------:|
|code|int|状态码(0成功，5000账号已经存在，5012添加失败,6001登陆失败)|
|message|string|状态信息|


返回示例:

```
{
    "code": 0,
    "message": "success",
    "data": {
        "sessionId": "11f2930d-5302-4fb5-87aa-7b9c9827c1ee_12_0",
        "gameUrl": "http://redfront.pkbeta.com"
    }
}

```

### 额度转入

- 接口地址: `/api/transferredIn`
- Method: POST
- 负责人: blazk

- url参数

|字段|类型|必填项|说明|
|:---|:----:|----:|:------:|
|account|string|true|账号|
|amount|double|true|额度|


- 返回值


|字段|类型|说明|
|:---|:----:|:------:|
|data|string|订单号|



其他字段的数组元素结构为:

|字段|类型|说明|
|:---|:----:|:------:|
|code|int|状态码(200成功，6000账号被停用,5001账号不存在，1009修改失败)|
|message|string|状态信息|


返回示例:

```
{
    "code": 0,
    "message": "success",
    "data": 200
}

```

### 额度转出

- 接口地址: `/api/transferredOut`
- Method: POST
- 负责人: blazk

- url参数

|字段|类型|必填项|说明|
|:---|:----:|----:|:------:|
|account|string|true|账号|
|amount|double|true|额度|


- 返回值


|字段|类型|说明|
|:---|:----:|:------:|
|data|double|剩余额度|


其他字段的数组元素结构为:

|字段|类型|说明|
|:---|:----:|:------:|
|code|int|状态码(200成功，6000账号被停用,5001账号不存在，1009修改失败,6003额度不足)|
|message|string|状态信息|


返回示例:

```
{
    "code": 0,
    "message": "success",
    "data": 15549.62
}

```

### 用户查询

- 接口地址: `/api/userinfo`
- Method: POST
- 负责人: blazk

- url参数

|字段|类型|必填项|说明|
|:---|:----:|----:|:------:|
|account|string|true|账号|


- 返回值

data内数据

|字段|类型|说明|
|:---|:----:|:------:|
|id|int|id|
|lineId|string|lindId|
|agencyId|string|agencyId|
|account|string|账号|
|balance|double|余额|
|createTime|int|创建时间|
|editTime|int|修改时间|
|capital|double|红包押金|
|availableBalance|double|可用金额|
|lastLoginIp|string|上次登陆ip|
|lastLoginTime|int|上次登陆时间|


其他字段的数组元素结构为:

|字段|类型|说明|
|:---|:----:|:------:|
|code|int|状态码(200成功，5001账号不存在)|
|message|string|状态信息|


返回示例:

```
{
    "code": 0,
    "message": "success",
    "data": {
        "id": 12,
        "lineId": "",
        "agencyId": "",
        "account": "aaa_1_js108",
        "balance": 16313.53,
        "createTime": 1576838350,
        "editTime": 0,
        "capital": 3500,
        "available_balance": 12813.53,
        "lastLoginIp": "127.0.0.1",
        "lastLoginTime": 1577410874
    }
}

```

### 注单采集

- 接口地址: `/api/collect`
- Method: POST
- 负责人: blazk


- 返回值

data内数据

|字段|类型|说明|
|:---|:----:|:------:|
|id|int|id|
|lineId|string|lindId|
|agencyId|string|agencyId|
|settlementInfo|string|结算信息(json)|
|collectStatus|int|采集状态 1未采集 2 已采集|
|createTime|int|创建时间|

settlementInfo内数据

|字段|类型|说明|
|:---|:----:|:------:|
|id|int|id|
|lineId|string|lindId|
|agencyId|string|agencyId|
|account|string|账号|
|redSender|string|发包者|
|createTime|int|创建时间|
|gameType|int|游戏类型 1牛牛 2扫雷|
|gamePlay|int|游戏玩法|
|roomId|int|群id|
|roomName|string|群名称|
|orderNo|string|注单号|
|redId|int|红包ID|
|redMoney|double|发包金额|
|redNum|int|红包个数|
|receiveMoney|double|领取金额|
|royalty|double|抽水比例|
|royaltyMoney|double|抽水金额|
|odds|double|赔率|
|thunderNum|int|雷值|
|adminNum|int|庄家牛数|
|memberNum|int|玩家牛数|
|money|double|输赢金额|
|realMoney|double|实际输赢金额|
|gameTime|int|游戏时间|
|receiveTime|int|红包领取时间|
|redStartTime|int|红包开始时间|
|status|int|状态 0未结算，1赢，2输，3无效|

其他字段的数组元素结构为:

|字段|类型|说明|
|:---|:----:|:------:|
|code|int|状态码(200成功，5001账号不存在)|
|message|string|状态信息|


返回示例:

```
扫雷
{
    "code": 0,
    "message": "success",
    "data": [
        {
            "id": 304,
            "lineId": "aaa",
            "agencyId": "a",
            "userId": 12,
            "account": "aaa_1_js108",
            "redSender": "aaa_1_js108",
            "gameType": 1,
            "gamePlay": 1,
            "roomId": 1,
            "roomName": "1231",
            "orderNo": "slhb20191227102808463599",
            "redId": 298,
            "redMoney": 2,
            "redNum": 2,
            "receiveMoney": 1.11,
            "royalty": 10,
            "royaltyMoney": 0,
            "money": -4,
            "realMoney": 0,
            "gameTime": 1,
            "receiveTime": 1577417288,
            "redStartTime": 1577417288,
            "status": 3,
            "extra": "{\"thunderNum\":2,\"odds\":2}",
            "isRobot": 2,
            "isFreeDeath": 2
        }
    ]
}

牛牛
{
    "code": 0,
    "message": "success",
    "data": [
        {
            "id": 190,
            "lineId": "aaa",
            "agencyId": "a",
            "userId": 1,
            "account": "aaa_1_js029",
            "redSender": "aaa_1_js108",
            "gameType": 2,
            "gamePlay": 1,
            "roomId": 4,
            "roomName": "1321",
            "orderNo": "nnhb20191224185034230824",
            "redId": 183,
            "redMoney": 20,
            "redNum": 2,
            "receiveMoney": 4.5,
            "royalty": 10,
            "royaltyMoney": 0.45,
            "money": 4.5,
            "realMoney": 4.05,
            "gameTime": 2,
            "receiveTime": 1577188234,
            "redStartTime": 1577188225,
            "status": 1,
            "extra": "{\"adminNum\":2,\"memberNum\":4}",
            "isRobot": 2,
            "isFreeDeath": 2
        }
    ]
}

普通红包
{
    "code": 0,
    "message": "success",
    "data": [
        {
            "id": 192,
            "lineId": "aaa",
            "agencyId": "a",
            "userId": 1,
            "account": "aaa_1_js029",
            "redSender": "aaa_1_js108",
            "gameType": 0,
            "gamePlay": 0,
            "roomId": 4,
            "roomName": "1321",
            "orderNo": "pthb20191224185034230824",
            "redId": 182,
            "redMoney": 20,
            "redNum": 2,
            "receiveMoney": 4.5,
            "royalty": 0,
            "royaltyMoney": 0.00,
            "money": 4.5,
            "realMoney": 4.5,
            "gameTime": 2,
            "receiveTime": 1577188234,
            "redStartTime": 1577188225,
            "status": 1,
            "extra": "",
            "isRobot": 2,
            "isFreeDeath": 2
        }
    ]
}
```

### 注单补采集

- 接口地址: `/api/collect`
- Method: POST
- 负责人: blazk

- url参数

|字段|类型|必填项|说明|
|:---|:----:|----:|:------:|
|start|int|true|开始时间|
|end|int|true|结束时间|


- 返回值

data内数据

|字段|类型|说明|
|:---|:----:|:------:|
|id|int|id|
|lineId|string|lindId|
|agencyId|string|agencyId|
|settlementInfo|string|结算信息(json)|
|collectStatus|int|采集状态 1未采集 2 已采集|
|createTime|int|创建时间|

settlementInfo内数据

|字段|类型|说明|
|:---|:----:|:------:|
|id|int|id|
|lineId|string|lindId|
|agencyId|string|agencyId|
|account|string|账号|
|redSender|string|发包者|
|createTime|int|创建时间|
|gameType|int|游戏类型 1牛牛 2扫雷|
|gamePlay|int|游戏玩法|
|roomId|int|群id|
|roomName|string|群名称|
|orderNo|string|注单号|
|redId|int|红包ID|
|redMoney|double|发包金额|
|redNum|int|红包个数|
|receiveMoney|double|领取金额|
|royalty|double|抽水比例|
|royaltyMoney|double|抽水金额|
|odds|double|赔率|
|thunderNum|int|雷值|
|adminNum|int|庄家牛数|
|memberNum|int|玩家牛数|
|money|double|输赢金额|
|realMoney|double|实际输赢金额|
|gameTime|int|游戏时间|
|receiveTime|int|红包领取时间|
|redStartTime|int|红包开始时间|
|status|int|状态 0未结算，1赢，2输，3无效|

其他字段的数组元素结构为:

|字段|类型|说明|
|:---|:----:|:------:|
|code|int|状态码(200成功，5001账号不存在)|
|message|string|状态信息|


返回示例:

```
{
    "code": 0,
    "message": "success",
    "data": [
        {
            "id": 0,
            "lineId": "aaa",
            "agencyId": "1",
            "settlementInfo": "{\"id\":305,\"account\":\"aaa_1_js029\",\"redSender\":\"aaa_1_js108\",\"gameType\":1,\"gamePlay\":1,\"roomId\":1,\"roomName\":\"1231\",\"orderNo\":\"slhb20191227102815458149\",\"redId\":298,\"redMoney\":2,\"redNum\":2,\"receiveMoney\":0.89,\"royalty\":10,\"royaltyMoney\":0,\"odds\":0,\"thunderNum\":0,\"adminNum\":3,\"memberNum\":7,\"memberMine\":0,\"money\":4,\"realMoney\":3.6,\"gameTime\":1,\"receiveTime\":1577417295,\"redStartTime\":1577417288,\"status\":1}",
            "collectStatus": 1,
            "createTime": 1577417426
        },
        {
            "id": 0,
            "lineId": "aaa",
            "agencyId": "1",
            "settlementInfo": "{\"id\":304,\"account\":\"aaa_1_js108\",\"redSender\":\"aaa_1_js108\",\"gameType\":1,\"gamePlay\":1,\"roomId\":1,\"roomName\":\"1231\",\"orderNo\":\"slhb20191227102808463599\",\"redId\":298,\"redMoney\":2,\"redNum\":2,\"receiveMoney\":1.11,\"royalty\":10,\"royaltyMoney\":0,\"odds\":0,\"thunderNum\":0,\"adminNum\":3,\"memberNum\":3,\"memberMine\":0,\"money\":-4,\"realMoney\":0,\"gameTime\":1,\"receiveTime\":1577417288,\"redStartTime\":1577417288,\"status\":3}",
            "collectStatus": 1,
            "createTime": 1577417426
        }
    ]
}
```
