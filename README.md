```golang
假设 users = [A, B, C, D]，index=1（要删除 B）：

users[:index] 是 [A]
users[index+1:] 是 [C, D]
users[index+1:]... 会展开为 C, D
最终 append([A], C, D) 的结果是 [A, C, D]，实现了删除 B 的效果。
```