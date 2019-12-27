package SmallComponents

import (
	"fmt"
)

/**
栈：限制插入和删除只能在一个位置上进行的表，该位置是表的末端，叫做栈的顶(top)
对栈的基本操作有push(进栈)和pop(出栈)。
基本算法：
进栈(push):
1.若top>=n时，作出错误处理(进栈前先检查栈是否已满，满则溢出，不满则进入2)
2.置top = top + 1(栈指针加1,指向进栈地址)
3.s(top) = x ,结束(x为新进栈的元素)
出栈(pop):
1.若top <=0，则给出下溢信息，作出错处理(出栈前先检查栈是否为空，空则下溢，不空走2)
2.x = s(top),出栈后的元素赋值给x
3.top = top -1 ，栈指针减1,指向栈顶
*/

// 定义常量栈的初始大小
const initSize int = 2

// 栈结构体
type Stack struct {
	// 容量
	size int
	// 栈顶
	top int
	// 用slice作容器，定义为interface{}
	data []interface{}
}

// 创建并初始化栈，返回strck
func CreateStack() Stack {
	s := Stack{}
	s.size = initSize
	s.top = -1
	s.data = make([]interface{}, initSize)
	return s
}

// 判断栈是否为空
func (s *Stack) IsEmpty() bool {
	return s.top == -1
}

// 判断栈是否已满
func (s *Stack) IsFull() bool {
	return s.top == s.size-1
}

// 入栈
func (s *Stack) Push(data interface{}) bool {
	// 首先判断栈是否已满
	if s.IsFull() {
		fmt.Println("stack is full, push failed")
		return false
	}
	// 栈顶指针+1
	s.top++
	// 把当前的元素放在栈顶的位置
	s.data[s.top] = data
	return true
}

// pop,返回栈顶元素
func (s *Stack) Pop() interface{} {
	// 判断是否是空栈
	if s.IsEmpty() {
		fmt.Println("stack is empty , pop error")
		return nil
	}
	// 把栈顶的元素赋值给临时变量tmp
	tmp := s.data[s.top]
	// 栈顶指针-1
	s.top--
	return tmp
}

// 栈的元素的长度
func (s *Stack) GetLength() int {
	length := s.top + 1
	return length
}

// 清空栈
func (s *Stack) Clear() {
	s.top = -1
}

// 遍历栈
func (s *Stack) Traverse() {
	// 是否为空栈
	if s.IsEmpty() {
		fmt.Println("stack is empty")
	}
	for i := 0; i <= s.top; i++ {
		fmt.Println(s.data[i], " ")
	}
}

// 打印当前栈的信息
func (s *Stack) PrintInfo() {
	fmt.Println("栈容量：", s.size)
	fmt.Println("栈顶指针：", s.top)
	fmt.Println("栈元素长度:", s.GetLength())
	fmt.Println("栈是否为空:", s.IsEmpty())
	fmt.Println("栈是否已满:", s.IsFull())
	fmt.Println("========遍历======")
	s.Traverse()
	fmt.Println("=========遍历结束========")
}
