# 问

EmployeeController中的GetProfile接口,下面的代码中:

```go
userID, _ := c.Get("userID")
role, _ := c.Get("userRole")
```

userID和role是从哪里获取的？

# 答

在 Gin 框架中，`userID` 和 `role` 是通过 **中间件** 从 JWT 令牌中解析并注入到请求上下文（`gin.Context`）中的。以下是它们的来源和流程说明：

---

### **1. 数据来源：JWT 令牌**
当用户登录成功后，后端会生成一个 **JWT 令牌**，其中包含以下关键信息（在 `utils/jwt.go` 中定义）：
```go
type Claims struct {
    UserID uint   `json:"user_id"`  // 用户ID
    Role   string `json:"role"`     // 用户角色（如 admin/employee）
    // ...
}
```

---

### **2. 中间件解析令牌并注入上下文**
在 `middleware/auth.go` 的鉴权中间件中，解析 JWT 令牌后，将用户信息存入 Gin 的上下文：
```go
// middleware/auth.go
func JWTAuth() gin.HandlerFunc {
    return func(c *gin.Context) {
        // ...（解析令牌逻辑）

        // 从令牌中提取 Claims
        claims, err := utils.ParseJWT(tokenString)
        
        // 将用户ID和角色注入上下文
        c.Set("userID", claims.UserID)
        c.Set("userRole", claims.Role)

        c.Next() // 继续处理后续请求
    }
}
```

---

### **3. 控制器从上下文中读取数据**
在 `controllers/user.go` 的 `GetProfile` 方法中，通过 `c.Get("key")` 从上下文获取数据：
```go
func GetProfile(c *gin.Context) {
    // 从上下文获取 userID 和 role
    userID, _ := c.Get("userID")     // 值来自中间件设置的 claims.UserID
    role, _ := c.Get("userRole")     // 值来自中间件设置的 claims.Role

    // ...后续逻辑
}
```

---

### **关键流程总结**
1. **用户登录** → 后端生成 JWT 令牌（含 `user_id` 和 `role`）。
2. **后续请求** → 前端在 `Authorization` 头中携带令牌。
3. **中间件拦截** → 解析令牌并将 `user_id` 和 `role` 存入请求上下文。
4. **控制器读取** → 从上下文直接获取用户信息，无需重复解析令牌。

---

### **验证示例**
#### **1. 中间件设置的值**
假设用户令牌解析后的 `claims` 为：
```json
{
  "user_id": 123,
  "role": "employee"
}
```
中间件会通过 `c.Set("userID", 123)` 和 `c.Set("userRole", "employee")` 存入上下文。

#### **2. 控制器获取的值**
在 `GetProfile` 中：
```go
userID, _ := c.Get("userID")   // 123（类型为 uint）
role, _ := c.Get("userRole")   // "employee"（类型为 string）
```

---

### **注意事项**
1. **数据类型转换**  
   `c.Get()` 返回的是 `interface{}` 类型，需确保后续使用时转换正确类型：
   ```go
   userID := c.Get("userID").(uint)
   role := c.Get("userRole").(string)
   ```

2. **中间件顺序**  
   确保 `JWTAuth()` 中间件在需要鉴权的路由之前执行（如 `authGroup.Use(middleware.JWTAuth())`）。

3. **上下文生命周期**  
   上下文数据仅在当前请求内有效，不同请求之间的数据隔离。

---

通过以上机制，你可以安全地在控制器中获取用户身份信息，实现基于角色的权限控制！