package hw04lrucache

import "fmt"

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

func (li *ListItem) String() string {
	var nextValue interface{}
	if next := li.Next; next != nil {
		nextValue = next.Value
	}

	var prevValue interface{}
	if prev := li.Prev; prev != nil {
		prevValue = prev.Value
	}

	return fmt.Sprintf("Value is %v, Next value is %v, Prev value is %v", li.Value, nextValue, prevValue)
}

type list struct {
	len   int
	front *ListItem
	back  *ListItem
}

func (l *list) String() string {
	values := make([]int, 0)
	for i := l.Front(); i != nil; i = i.Next {
		values = append(values, i.Value.(int))
	}

	return fmt.Sprintf("%v", values)
}

func NewList() List {
	return new(list)
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.front
}

func (l *list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	currFrontItem := l.Front()
	newFrontItem := &ListItem{
		Value: v,
		Next:  currFrontItem,
		Prev:  nil,
	}
	if currFrontItem != nil {
		currFrontItem.Prev = newFrontItem
	}

	l.front = newFrontItem
	if l.Len() == 0 {
		l.back = newFrontItem
	}

	l.len++

	return newFrontItem
}

func (l *list) PushBack(v interface{}) *ListItem {
	currBackItem := l.Back()
	newBackItem := &ListItem{
		Value: v,
		Next:  nil,
		Prev:  currBackItem,
	}
	if currBackItem != nil {
		currBackItem.Next = newBackItem
	}
	l.back = newBackItem

	if l.Len() == 0 {
		l.front = newBackItem
	}

	l.len++

	return newBackItem
}

func (l *list) Remove(i *ListItem) {
	// if element is front
	if i.Prev == nil {
		l.front = i.Next
	} else {
		i.Prev.Next = i.Next
	}

	// if element is back
	if i.Next == nil {
		l.back = i.Prev
	} else {
		i.Next.Prev = i.Prev
	}
	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	l.Remove(i)
	l.PushFront(i.Value)
}
