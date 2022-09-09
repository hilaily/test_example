package main

import (
	"testing"

	. "github.com/agiledragon/gomonkey/v2"
	"github.com/stretchr/testify/assert"
)

/********************************************

	介绍 gomonkey 包的使用方法

********************************************/

/*
	给参数打桩
*/

var num = 10

// TestApplyGlobalVar 给全局参数打桩
func TestApplyGlobalVar(t *testing.T) {
	assert.Equal(t, num, 10)

	patches := ApplyGlobalVar(&num, 150)
	// Reset 用于恢复 patch
	defer patches.Reset()

	assert.Equal(t, num, 150)
}

/*
	给函数打桩
*/

func Echo(name string, age int) (string, int) {
	return name, age
}

// 模拟一个包级别的函数
func TestApplyFunc(t *testing.T) {
	patches := ApplyFunc(Echo, func(string, int) (string, int) {
		return "ace", 11
	})
	defer patches.Reset()

	name, age := Echo("bob", 0)

	assert.Equal(t, name, "ace")
	assert.Equal(t, age, 11)
}

// 模拟一个函数的返回值
func TestApplyFuncReturn(t *testing.T) {
	patches := ApplyFuncReturn(Echo, "ace", 11)
	defer patches.Reset()

	name, age := Echo("bob", 0)

	assert.Equal(t, name, "ace")
	assert.Equal(t, age, 11)
}

// 批量模拟函数的一组返回值
func TestApplyFuncSeq(t *testing.T) {
	outputs := []OutputCell{
		// 每一组数据是一次的返回值，Times 表示这个值返回几次，默认是 1 次
		// 下面这段表示前两次调用返回 Clark，第三次返回 Dave
		{Values: []interface{}{"Clark", 8}, Times: 2},
		{Values: []interface{}{"Dave", 9}, Times: 1},
	}

	patches := ApplyFuncSeq(Echo, outputs)
	defer patches.Reset()

	name, age := Echo("bob", 0)
	assert.Equal(t, name, "Clark")
	assert.Equal(t, age, 8)

	name, age = Echo("bob", 0)
	assert.Equal(t, name, "Clark")
	assert.Equal(t, age, 8)

	name, age = Echo("bob", 0)
	assert.Equal(t, name, "Dave")
	assert.Equal(t, age, 9)

}

// 模拟一个函数变量
func TestApplyFuncVar(t *testing.T) {
	// 假如我们不是直接使用的一个函数，而是使用的一个函数类型的变量，如下面的 echoFunc
	// ApplyFunc 无法直接mock echoFunc，当然直接使用 ApplyFunc 是 mock Echo 方法，echoEcho 也会被修改
	// 另外一种方式就是使用 ApplyFuncVar 去 mock echoFunc 这个函数变量。
	echoFunc := Echo

	patches := ApplyFuncVar(&echoFunc, func(string, int) (string, int) {
		return "ace", 11
	})
	defer patches.Reset()

	name, age := echoFunc("bob", 0)

	assert.Equal(t, name, "ace")
	assert.Equal(t, age, 11)
}

/*
	给方法打桩
*/

type AStruct struct{}

func (a *AStruct) Echo(name string, age int) (string, int) {
	return name, age
}

func (a *AStruct) echo(name string, age int) (string, int) {
	return name, age
}

// 模拟一个结构体的方法
func TestApplyMethodFunc(t *testing.T) {
	myStruct := &AStruct{}

	var a *AStruct
	patches := ApplyMethodFunc(a, "Echo", func(string, int) (string, int) {
		return "ace", 11
	})
	defer patches.Reset()

	name, age := myStruct.Echo("bob", 0)

	assert.Equal(t, name, "ace")
	assert.Equal(t, age, 11)
}

// 模拟一个结构体的导出方法的返回值
func TestApplyMethodReturn(t *testing.T) {
	myStruct := &AStruct{}

	var a *AStruct
	//  第一个参数是结构体，如果是指针方法就传指针，如果是值方法就传值
	patches := ApplyMethodReturn(a, "Echo", "ace", 11)
	defer patches.Reset()

	name, age := myStruct.Echo("bob", 0)

	assert.Equal(t, name, "ace")
	assert.Equal(t, age, 11)
}

// 模拟结构体私有方法
func TestApplyMethodPrivateFunc(t *testing.T) {
	myStruct := &AStruct{}

	var a *AStruct
	//  第一个参数是结构体，如果是指针方法就传指针，如果是值方法就传值
	patches := ApplyPrivateMethod(a, "echo", func(string, int) (string, int) {
		return "ace", 11
	})
	defer patches.Reset()

	name, age := myStruct.echo("bob", 0)

	assert.Equal(t, name, "ace")
	assert.Equal(t, age, 11)
}
