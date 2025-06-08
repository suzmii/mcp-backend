在最开始，请你记住，思考大于行动，在正式开始你的工作前，请先思考整个项目中不清晰的地方，然后提问，在完全明确你需要做什么之后，你再开始动手。
你是一名富有设计理念的前端开发者，接下来你将会使用 vue3 构建一个精美且完善的 web 应用，包括登录，注册，聊天,和 session 列表页面。

接下来我会描述各个页面的功能和页面跳转逻辑，以及页面细节
首先是登录注册页面
用户只有 username 和 password，注册也只需要填写这两个信息，登录时也是.
登录可以拿到 access_token 和 refresh_token。在访问部分接口时，你需要把 access_token 放在请求头中。
access_token 过期之后，访问这部分接口会返回 401，客户端（也就是你要写的这个 web 应用，后面简称为客户端）应该请求 refresh_token 接口，拿到新的 access_token ，然后重试请求。应该有类似的框架，你可以自由选择技术栈。
如果 refresh_token 过期，你将无法通过 refresh_token 获取新的 access_token，这时你需要重新登录。

然后是 session 列表
session 列表页面需要做分页，从后面的 api 中你能得到分页的参数。
每个列表项包括 sessionID 和 session Hint，你需要把每个列表项做成一个卡片。

最后是聊天页面，这个页面请你仿照现在的 AI 聊天页面，暂时先不支持流式传输。
以下是对话基本逻辑：
拿到用户的输入
if 是新对话
`call mcp` chat 总结对话主题
`call后端` create session，传入对话主题，拿到 SesionId
`call mcp` new session，传入 sessionId，拿到 reply
else
`call mcp` new session，传入 sessionId 和 messagelist(从 session 详情拿到)，拿到 reply
fi

`call后端` append message，传入 user message 和 reply

---

关于 api，请你先把两个 Api 的 base url 留空，并在注释中给出明确提示。

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
// 以下是mcp接口，这些都不需要token
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
