# golang async

## 简介

简单快速实现异步并发任务,避免繁琐代码
并将任务结果集中返回
类似python asyncio.gather的语法

## 应用场景

这在多个耗时长的网络请求（如：调用API接口）时非常有用。其可将顺序执行变为并行计算，可极大提高程序的执行效率。也能更好的发挥出多核CPU的优势。

## 使用

```go
go get github.com/jianmaikj/async
```

## demo

```go
package main

import (
	"github.com/jianmaikj/async"
	"runtime"
	"fmt"
	"time"
)

func request(params ...interface{}) (res interface{}) {
	//sql request...
	return
}

func main() {
	// 建议程序开启多核支持
	runtime.GOMAXPROCS(runtime.NumCPU())
	// 耗时操作

	// 开启并发操作
	t := time.Now()
	task1 := &async.Task{
		Name:    "1",
		Handler: request1,
		Params:  []interface{}{1},
	}
	task2 := &async.Task{
		Name:    "2",
		Handler: request2,
		Params:  []interface{}{2},
	}
	res := async.Gather(task1, task2)
	//结果是一个map[string][]reflect.Value返回值,可以根据Name取结果,如果没有设定Name则默认为任务顺序数字的字符串:"1","2",...,例如res["1"]即第一个任务结果集
	fmt.Println(res["1"], "time::", time.Now().Sub(t))

	// 或
	tasks := async.NewTasks()
	for i := 0; i <= 2; i++ {
		tasks = append(tasks, &async.Task{
			Handler: r1,
			Params:  []interface{}{1},
		})
	}
	res := async.Gather(tasks...)
	fmt.Println(res["1"])
}


```




