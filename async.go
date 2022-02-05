package async

// 异步执行类，提供异步执行的功能，可快速方便的开启异步执行
// 通过NewAsync() 来创建一个新的异步操作对象
// 通过调用 Add 函数来向异步任务列表中添加新的任务
// 通过调用 Run 函数来获取一个接收返回的channel，当返回结果时将会返回一个map[string][]interface{}
// 的结果集，包括每个异步函数所返回的所有的结果
import (
	"log"
	"reflect"
	"strconv"
)

// 异步执行所需要的数据
type fun struct {
	Name    string
	Handler reflect.Value
	Params  []reflect.Value
}

// Task 异步并发任务,第一个参数为该异步请求的任务名,第二个参数任务函数,第三个为函数的参数列表
type Task struct {
	Name    string
	Handler interface{}
	Params  []interface{}
}

//funRes 函数结果
type funRes struct {
	Name   string
	Result []reflect.Value
}

// Async 异步执行对象
type Async struct {
	tasks []fun
}

// NewTasks 创建一个新的异步任务空集合
func NewTasks() []*Task {
	return make([]*Task, 0)
}

// New 创建一个新的异步执行对象
func New() *Async {
	return &Async{tasks: make([]fun, 0)}
}

// Add 添加异步执行任务
// name 任务名，结果返回时也将放在任务名中
// handler 任务执行函数，将需要被执行的函数导入到程序中
// params 任务执行函数所需要的参数
func (a *Async) Add(f *Task) bool {
	handler := f.Handler
	params := f.Params
	handlerValue := reflect.ValueOf(handler)
	if handlerValue.Kind() == reflect.Func {

		paramNum := len(params)
		paramsValue := make([]reflect.Value, 0)

		if paramNum > 0 {
			for _, v := range params {
				paramsValue = append(paramsValue, reflect.ValueOf(v))
			}
		}
		a.tasks = append(a.tasks, fun{
			Name:    f.Name,
			Handler: handlerValue,
			Params:  paramsValue,
		})

		return true
	}

	return false
}

// Run 任务执行函数，成功时将返回一个用于接受结果的channel
// 在所有异步任务都运行完成时，结果channel将会返回一个map[string][]interface{}的结果。
func (a *Async) Run() (chan map[string][]reflect.Value, bool) {
	tasks := a.tasks
	count := len(tasks)
	if count < 1 {
		return nil, false
	}
	result := make(chan map[string][]reflect.Value)
	chans := make(chan *funRes, count)

	go func(result chan map[string][]reflect.Value, chans chan *funRes) {
		rs := make(map[string][]reflect.Value)
		defer func(rs map[string][]reflect.Value) {
			result <- rs
		}(rs)
		for {
			if count < 1 {
				break
			}

			select {
			case res := <-chans:
				count--
				rs[res.Name] = res.Result
			}
		}
	}(result, chans)

	for _, v := range tasks {
		go func(chans chan *funRes, fun fun) {
			result := make([]reflect.Value, 0)
			defer func(name string, chans chan *funRes) {
				chans <- &funRes{Name: name, Result: result}
			}(fun.Name, chans)

			values := fun.Handler.Call(fun.Params)

			if valuesNum := len(values); valuesNum > 0 {
				result = values
				return
			}
		}(chans, v)
	}

	return result, true
}

// Clean 清空任务队列.
func (a *Async) Clean() {
	a.tasks = make([]fun, 0)
}

// Gather 快速启动并发任务
func Gather(funs ...*Task) map[string][]reflect.Value {

	tasks := New()
	defer tasks.Clean()

	for i, v := range funs {
		name := v.Name
		if name == "" {
			name = strconv.Itoa(i)
		}
		tasks.Add(&Task{
			Name:    name,
			Handler: v.Handler,
			Params:  v.Params,
		})
	}
	//执行
	if chans, ok := tasks.Run(); ok {
		// 将数据从通道中取回
		res := <-chans
		// 这里判断下是否所有的异步请求都已经执行成功
		if len(res) == 2 {
			return res
		} else {
			log.Println("async not execution all task")
		}
	}
	// 重置任务
	tasks.Clean()
	return nil
}
