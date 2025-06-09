1. 注册
   调用 user 服务，对密码哈希后储存
2. 登录
   调用 user 服务，user 服务检查密码，然后 user 服务调用 auth 服务拿 token，返回
3. 刷新 token
   调用 auth 服务验证 refresh token，然后签发 access token
