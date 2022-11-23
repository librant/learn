package main

import (
	"log"
)

type VisitorFunc func() error

type Visitor interface {
	Visit(VisitorFunc) error
}

type VisitorList []Visitor

// Visit 实现 visitor 接口
func (l VisitorList) Visit(fn VisitorFunc) error {
	for i := range l {
		if err := l[i].Visit(func() error {
			log.Printf("In VisitorList before fn")
			if err := fn(); err != nil {
				log.Printf("VisitorList fn failed: %v", err)
			}
			log.Printf("In VisitorList after fn")
			return nil
		}); err != nil {
			return err
		}
	}
	return nil
}

// Visitor1 结构 1
type Visitor1 struct {}

// Visit 实现
func (v Visitor1) Visit(fn VisitorFunc) error {
	log.Printf("In Visitor1 before fn")
	if err := fn(); err != nil {
		log.Printf("Visitor1 fn failed: %v", err)
	}
	log.Printf("In Visitor1 after fn")
	return nil
}

type Visitor2 struct {
	visitor Visitor
}

func (v Visitor2) Visit(fn VisitorFunc) error {
	if err := v.visitor.Visit(func() error {
		log.Printf("In Visitor2 before fn")
		if err := fn(); err != nil {
			log.Printf("Visitor2 fn failed: %v", err)
		}
		log.Printf("In Visitor2 after fn")
		return nil
	}); err != nil {
		log.Printf("Visitor2 failed: %v", err)
		return err
	}
	return nil
}

type Visitor3 struct {
	visitor Visitor
}

func (v Visitor3) Visit(fn VisitorFunc) error {
	if err := v.visitor.Visit(func() error {
		log.Printf("In Visitor3 before fn")
		if err := fn(); err != nil {
			log.Printf("Visitor3 fn failed: %v", err)
		}
		log.Printf("In Visitor3 after fn")
		return nil
	}); err != nil {
		log.Printf("Visitor3 failed: %v", err)
		return err
	}
	return nil
}

