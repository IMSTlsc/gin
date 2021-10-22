package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func main() {
	timeCalc := func() gin.HandlerFunc {
		return func(ctx *gin.Context) {
			if ctx.Query("a") == "" {
				ctx.Abort()                            //abort中断调用链路
				ctx.JSON(http.StatusBadRequest, gin.H{ //返回错误信息
					"message": "a 参数错误",
				})
				return
			}
			start := time.Now()
			fmt.Println("next 之前")

			ctx.Next() //在中间件内部调用ctx.Next()方法即调用链路中下一级，这里会调用具体的接口doSometing
			//如果不显式调用ctx.Next()，则需要等该中间件完成后再执行下一链路

			fmt.Println("next 之后") //接口调用完成后会继续执行该中间件
			cost := time.Since(start)
			fmt.Printf("花费时间%d", cost.Microseconds())
		}
	}

	e := gin.Default() //创建路由
	e.GET("/time", timeCalc(), func(ctx *gin.Context) {
		fmt.Println("进入接口函数")
		ctx.JSON(http.StatusOK, gin.H{
			"a": ctx.Query("a"),
		})
	}, doSomething(new(gin.Context)))

	e.Run(":8080")

	e.GET("/send", func(ctx *gin.Context) {
		type Param struct {
			A string `form:"a" binding:"required"` //解析a的数据，必须绑定
			B int    `form:"b" binding:"required"`
		}
		param := new(Param)
		if err := ctx.ShouldBind(param); err != nil {// 将数据绑定到param
			ctx.JSON(400, gin.H{
				"err": err.Error(),
			})
			return
		}
		ctx.JSON(200, gin.H{
			"Content-Type": ctx.ContentType(),
			"a":            param.A,
			"b":            param.B,
		})
	})

}
func doSomething(ctx *gin.Context) gin.HandlerFunc {
	time.Sleep(1 * time.Second)
	return func(context *gin.Context) {
	}
}

//获取post请求的请求体，请求localhost:8080/send -d "b=2&c=3"
func getPara(e *gin.Engine) {
	e.POST("/send", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"name": ctx.Param("name"),     //获取路径中:之后的name
			"qu":   ctx.Query("age"),      //获取路径中？之后的age
			"js":   ctx.PostForm("a"),     //获取路径中？之后json类型的数据
			"arr":  ctx.QueryArray("arr"), //获取路径中？之后数组类型的数据
		})
	})
}

//绑定路由规则(http://localhost:8080/hello)，执行该路径下绑定的函数，函数可以多个
func hello(r *gin.Engine) {
	r.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusAccepted, "hello")
	}, checkout())

	r.Run(":8080") //绑定监听端口
}
func checkout() gin.HandlerFunc {
	return func(c *gin.Context) { //匿名函数作为返回值，形参的实际值由外边的函数确定
		c.String(200, "ajahahhaha")
	}
}

//中间件调用顺序：全局中间件 -> group中间件 -> 单个接口中间件 -> 接口方法
func mid() {
	e := gin.New()
	e.Use(gin.Logger()) //全局中间件

	e.GET("/benchmark", BenchMarkMiddleware()) //单个接口中间件

	auth := e.Group("/auth") //一组的中间件
	auth.Use(AuthGroupMiddleware())
	{
		auth.POST("/login", LoginMiddleware()) //一组中间中各个接口可以再定义自己的中间件
		auth.POST("/test", TestMiddleware())   //
	}
	e.Run(":8080")
}

func BenchMarkMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
	}
}
func AuthGroupMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
	}
}
func LoginMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
	}
}
func TestMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
	}
}
