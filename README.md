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

## TestUpdateUser带有JSON内容的测试

测试思路相同，构建请求，触发请求，验证返回内容，PUT为更新内容，则需要同样验证是否更新完成，需要验证数据库。

在 HTTP 请求 / 响应中，**Header（头部）** 是一组键值对（`Key: Value`），用于传递  『元数据』 —— 也就是描述请求 / 响应的额外信息（比如数据格式、认证信息、缓存规则等）。它不直接包含业务数据（业务数据在请求体 / 响应体中），但能让客户端和服务器 『互相理解』 对方的意图和数据格式。

```go
func TestUpdateUser_Success(t *testing.T) {
	router := setupRouter()

	updatedUser := User{
		Name:  "John Updated",
		Email: "john.updated@example.com",
		Age:   31,
	}

	jsonData, _ := json.Marshal(updatedUser)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/users/1", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)

	// Check updated user data
	userData, ok := response.Data.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "John Updated", userData["name"])
}
```

### Header核心作用

- 快递盒里的物品是「请求体 / 响应体」（实际数据）
- 快递单上的信息（收件人、电话、是否易碎、保价金额）就是「Header」—— 描述物品的属性和处理要求

HTTP 中的 Header 同理：

- 客户端通过请求 Header 告诉服务器：「我发的是 JSON 数据」，「我需要中文」，「我已经登录了（带 Token）」
- 服务器通过响应 Header 告诉客户端：「我返回的是 HTML」，「数据缓存 1 小时」，「这个资源不存在（404）」

### Content-Type

1. **`application/json`**

   - 含义：请求体 / 响应体是 JSON 格式的数据

   - 场景：API 接口通信（前后端分离项目中最常用）

   - 示例请求体：

     ```json
     {"name":"张三","age":20}
     ```

   - 注意：发送 JSON 时，必须在请求 Header 中指定 `Content-Type: application/json`，否则服务器可能解析失败。

2. **`application/x-www-form-urlencoded`**

   - 含义：请求体是「表单键值对」格式（类似 URL 查询参数）

   - 场景：传统 HTML 表单提交（`method="post"`时默认格式）

   - 示例请求体：

     ```plaintext
     username=zhangsan&password=123&hobby=reading
     ```

   - 特点：数据是明文键值对，适合简单文本数据，不适合二进制文件。

3. **`multipart/form-data`**

   - 含义：请求体是「多部分数据」，可同时包含文本和二进制文件

   - 场景：文件上传（如上传图片、文档），或同时提交文本和文件

   - 示例结构（简化）：

     ```plaintext
     --分隔符
     Content-Disposition: form-data; name="username"
     
     zhangsan
     --分隔符
     Content-Disposition: form-data; name="avatar"; filename="head.jpg"
     Content-Type: image/jpeg
     
     [二进制图片数据]
     --分隔符--
     ```

   - 特点：会自动生成分隔符区分不同部分数据，支持二进制，是文件上传的标准格式。

4. **`text/html`**

   - 含义：响应体是 HTML 网页内容
   - 场景：浏览器请求网页（服务器返回 HTML 时默认用这个）
   - 示例：服务器返回`<html><body>Hello</body></html>`时，响应 Header 会带 `Content-Type: text/html; charset=utf-8`（charset 指定字符编码）。

5. **`text/plain`**

   - 含义：纯文本格式（无特殊结构）
   - 场景：简单的字符串数据（如返回一段日志、提示文本）
   - 示例：响应体是`"操作成功"`，Header 带 `Content-Type: text/plain`。

























