package mcpclient_test

import (
	"mcp/client/mcpclient"
	"testing"
)

func TestMcpClient(t *testing.T) {
	client := mcpclient.NewMCPClient("http://localhost:8000")

	err := client.CreateSession("test-session", []mcpclient.ChatMessage{
		{
			Role:    "user",
			Content: "你好,这是第一次对话, 请记住以下内容: 测试消息111",
		},
		{
			Role:    "assistant",
			Content: "好的,请问你有什么疑问",
		},
	})

	if err != nil {
		t.Fatal(err)
	}

	result, err := client.Chat("test-session", "刚才让你记住的内容是啥")
	if err != nil {
		t.Fatal(err)
		return
	}

	print(result)

}

func TestStartBiliBili(t *testing.T) {
	client := mcpclient.NewMCPClient("http://localhost:8000")

	err := client.CreateSession("test-session", []mcpclient.ChatMessage{})

	if err != nil {
		t.Fatal(err)
	}

	result, err := client.Chat("test-session", "打开哔哩哔哩")
	if err != nil {
		t.Fatal(err)
		return
	}

	t.Log(result)
}
