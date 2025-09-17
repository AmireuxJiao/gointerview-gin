```golang
假设 users = [A, B, C, D]，index=1（要删除 B）：

users[:index] 是 [A]
users[index+1:] 是 [C, D]
users[index+1:]... 会展开为 C, D
最终 append([A], C, D) 的结果是 [A, C, D]，实现了删除 B 的效果。
```

# Basic Routing 测试框架详解

## TestGetAllUsers 简单GET请求

一个简单的Get请求测试模板，先完整的建立整个服务端，并且构建路由表，然后，通过**httptest.NewRequest**创建一个请求，之后采用**httptest.NewRecorder**来记录路由的返回，之后**r.ServeHTTP(w, req)**运行整个请求，并且获取服务端输出，之后就是验证，返回和预期是否相同。

```go
func TestGetAllUsers(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)

	// Check if users data is returned
	data, ok := response.Data.([]interface{})
	assert.True(t, ok)
	assert.Equal(t, 3, len(data))
}
```



## 虚拟服务端

`setupRouter()`构建了一个和实际运行过程中完全相同的服务端和路由表，但是并没有运行。这个虚拟服务端接收`ServeHTTP`发出的请求，按照流程回复对应的服务端返回。复现基础。相当于**模拟器**

## 请求器(Request)

`httptest.NewRequest` 是 Go 标准库提供的工具，用于创建一个 “虚拟的 HTTP 请求”。它的作用相当于：

- 浏览器输入 URL 发起请求
- 或者 curl 命令发送请求

这里需要指定请求方法（GET/POST）、路径 + 参数、请求体（POST 时需要），完全模拟真实客户端会发送的内容。

## 记录器(Recorder)

`httptest.NewRecorder` 是一个 “虚拟的响应接收器”，相当于：

- 浏览器的 “开发者工具”（记录服务器返回的内容）
- 或者 curl 命令的输出窗口

它会保存服务端返回的所有信息：状态码、响应头、响应体等，方便后续验证。

## 触发(ServeHTTP)

`ServeHTTP`是一个触发器，用于在内存中构建请求并且发送到**虚拟服务端**

1. 接收虚拟请求`req`
2. 按照`setupRouter`定义的路由规则，找到对应的处理函数（比如`/hello`的处理逻辑）
3. 执行处理函数，生成响应
4. 把响应结果写入到`w`（记录器）中

这个过程和 “真实客户端访问真实服务器” 的内部逻辑完全一致，只是在内存中完成，不需要网络传输，效率极高。





