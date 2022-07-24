# Gil, 一个 Golang 实现的 Web 框架
---

### 功能亮点
1）GET、POST 方法用于设置 handler 处理的请求的方法

2）Next()、Abort()、Set()、Get() 方法用于管理上下文

3）支持中间件

4）支持分组路由

5）支持动态路由

6）支持 URL 参数查询

7）支持表单处理

</br>

### 简单例子
``` go
import "github.com/RGBli/gil"

func main() {
    engine := gil.New()
    engine.GET("/", func(c *gil.Context) {
        c.String("Hello gil")
    })
    engine.Run(":8080")
}
```

</br>

### 技术
1）Golang 1.15

2）使用前缀树实现动态路由

3）单元测试

</br>

### 参考文章
https://github.com/geektutu/7days-golang/tree/master/gee-web

https://blog.dianduidian.com/post/gin-%E4%B8%AD%E9%97%B4%E4%BB%B6next%E6%96%B9%E6%B3%95%E5%8E%9F%E7%90%86%E8%A7%A3%E6%9E%90/

</br>