mcp client server 简称 mcp

对话基本逻辑：
拿到用户的输入
if 是新对话
`call mcp` chat 总结对话主题
`call后端` create session，传入对话主题，拿到 SesionId
`call mcp` new session，传入 sessionId，拿到 reply
else
`call mcp` new session，传入 sessionId 和 messagelist(从 session 详情拿到)，拿到 reply
fi

`call后端` append message，传入 user message 和 reply

```json
// 以下是后端API
// 其中以mcp开头的需要在Headers中加入Authorization: Bearer xxx，xxx为access token
// 如果某接口Unauthenticated了，则需要调用token刷新接口
// 如果刷新之后还是Unauthenticated，需要弹出提示有内部错误（debug用），为达到此目的你需要设置token的过期时间，暂时写死是3分钟（实际的过期时间更长）
// 如果调用token刷新接口失败则需要重新登陆
{
  "apis": [
    {
      // 创建session
      "path": "/mcp/session",
      "method": "POST",
      "request": {
        "hint": "string"
      },
      "response": {
        "session_uuid": "string"
      }
    },
    {
      // 获取session列表
      "path": "/mcp/sessions",
      "method": "GET",
      "request": null,
      "response": {
        "sessions": [
          {
            "session_uuid": "string",
            "session_hint": "string"
          }
        ]
      }
    },
    {
      // 获取session详情（包括对话记录）
      "path": "/mcp/sessions/:id",
      "method": "GET",
      "request": null,
      "response": {
        "session": {
          "session_uuid": "string",
          "session_hint": "string"
        },
        "messages": [
          {
            "role": "string",
            "content": "string"
          }
        ]
      }
    },
    {
      // 记录对话
      "path": "/mcp/message",
      "method": "POST",
      "request": {
        "session_uuid": "string",
        "messages": [
          {
            "role": "string",
            "content": "string"
          }
        ]
      },
      "response": null
    },
    {
      // 登录
      "path": "/user/login",
      "method": "POST",
      "request": {
        "username": "string",
        "password": "string"
      },
      "response": {
        "access_token": "string",
        "refresh_token": "string"
      }
    },
    {
      // 注册
      "path": "/user/register",
      "method": "POST",
      "request": {
        "username": "string",
        "password": "string"
      },
      "response": null
    },
    {
      // 刷新token
      "path": "/user/token/refresh",
      "method": "POST",
      "request": {
        "refresh_token": "string"
      },
      "response": {
        "access_token": "string"
      }
    }
  ],
  "definitions": {
    "Message": {
      "role": "string",
      "content": "string"
    },
    "SessionInfo": {
      "session_uuid": "string",
      "session_hint": "string"
    },
    "Claims": {
      "user_id": "number",
      "permissions": "Permission[]"
    }
  }
}
```

```json
// 以下是mcp接口
[
  {
    // 创建session
    "path": "/sessions",
    "method": "POST",
    "request": {
      "session_id": "string",
      "context": "List[Dict<string, string>]"
    },
    "response": {
      "success": "boolean"
    }
  },
  {
    // 有session的对话
    "path": "/chat",
    "method": "POST",
    "request": {
      "session_id": "string",
      "content": "string"
    },
    "response": {
      "response": "string"
    },
    "errors": {
      "404": "会话不存在",
      "500": "服务器内部错误"
    }
  },
  {
    // 无sesion对话，用于生成对话主题
    "path": "/nosession-chat",
    "method": "POST",
    "request": {
      "content": "string"
    },
    "response": {
      "response": "string"
    },
    "errors": {
      "500": "服务器内部错误"
    }
  }
]
```
