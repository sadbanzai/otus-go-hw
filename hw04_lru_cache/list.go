package hw04lrucache

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

type list struct {
	len   int
	front *ListItem
	back  *ListItem
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
	var item *ListItem
	if l.front == nil {
		item = &ListItem{Value: v, Next: nil, Prev: nil}
		l.front = item
		l.back = item
	} else {
		item = &ListItem{Value: v, Next: l.front, Prev: nil}
		l.front.Prev = item
		l.front = item
	}
	l.len++
	return l.front
}

func (l *list) PushBack(v interface{}) *ListItem {
	var item *ListItem
	if l.front == nil {
		item = &ListItem{Value: v, Next: nil, Prev: nil}
		l.front = item
		l.back = item
	} else {
		item = &ListItem{Value: v, Next: nil, Prev: l.back}
		l.back.Next = item
		l.back = item
	}
	l.len++
	return l.back
}

func (l *list) Remove(i *ListItem) {
	if l.front == nil {
		panic(i)
	}
	if i.Prev == nil {
		l.front = i.Next
	} else {
		i.Prev.Next = i.Next
	}
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

func NewList() List {
	return new(list)
}
