# API说明文档

## 约定说明

### 无特殊说明，以下文档的“响应内容”均为成功时返回; 错误返回内容格式如下

```json
{
    "code": -1,             //错误码, 具体约定查看状态码列表
    "msg": "错误提示"
}
```

### 签名哈希算法:

php例子

```php
$timestamp = time();
$key = '密钥';
$signature = hash_hmac('md5', $appid . $timestamp, $key);
```

## 服务API

### 获取access token

GET /token?appid=xxxxx&timestamp=xxxxx&signature=xxxxx

响应内容

```json
{
    "code": 0,
    "token": "xxxxxx"
}
```

### 刷新access token

GET /refresh-token?appid=xxxxx&timestamp=xxxxx&signature=xxxxx

响应内容

```json
{
    "code": 0
}
```

### 获取jsapi ticket

GET /ticket?appid=xxxxx&timestamp=xxxxx&signature=xxxxx

响应内容

```json
{
    "code": 0,
    "ticket": "xxxx"
}
```

### 刷新jsapi ticket

GET /refresh-ticket?appid=xxxxx&timestamp=xxxxx&signature=xxxxx

响应内容

```json
{
    "code": 0
}
```