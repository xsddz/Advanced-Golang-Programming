package kh

/*
 * 有效的括号字符串（ref: https://leetcode-cn.com/problems/valid-parenthesis-string/）
 * 给定一个只包含三种字符的字符串：（ ，） 和 *，写一个函数来检验这个字符串是否为有效字符串。有效字符串具有如下规则：
 * - 任何左括号 ( 必须有相应的右括号 )。
 * - 任何右括号 ) 必须有相应的左括号 ( 。
 * - 左括号 ( 必须在对应的右括号之前 )。
 * - * 可以被视为单个右括号 ) ，或单个左括号 ( ，或一个空字符串。
 * - 一个空字符串也被视为有效字符串。
 */

import (
	"fmt"
)

func init() {
	fmt.Println("init:kh")
}

type Stack struct {
	data []int
	len  int
}

func (s *Stack) Push(c int) {
	s.data = append([]int{c}, s.data...)
	s.len++
}

func (s *Stack) Pop() int {
	c := s.data[0]
	s.data = s.data[1:]
	s.len--
	return c
}

func J(s string) bool {
	leftStack := Stack{}
	starStack := Stack{}

	for i, c := range s {
		if c == '(' {
			// push
			leftStack.Push(i)
			continue
		}

		if c == '*' {
			// push
			starStack.Push(i)
			continue
		}

		if c == ')' {
			// check stack or pop
			if leftStack.len == 0 && starStack.len == 0 {
				return false
			} else if leftStack.len > 0 {
				leftStack.Pop()
			} else if starStack.len > 0 {
				starStack.Pop()
			}
		}
	}

	for leftStack.len > 0 && starStack.len > 0 {
		if leftStack.Pop() > starStack.Pop() {
			return false
		}
	}

	return leftStack.len == 0
}

func Run() {
	ss := []string{
		"((*)",
		"()",
		"(*)",
		"(*))",
		"*(",
	}
	for _, s := range ss {
		fmt.Println(s, "\t:", J(s))
	}
}
