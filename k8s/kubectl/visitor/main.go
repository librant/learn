package main

import (
	"log"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	var visitor Visitor
	var visitors []Visitor

	visitor = Visitor1{}
	visitors = append(visitors, visitor)

	visitor = Visitor2{
		visitor: VisitorList(visitors),
	}

	visitor = Visitor3{visitor}
	visitor.Visit(func() error {
		log.Printf("In visitFunc")
		return nil
	})
}
