package src

import (
	"github.com/alecthomas/assert/v2"
	"testing"
)

func TestParse(t *testing.T) {
	p1, _ := Parse(`5`)
	assert.Equal(t, p1.(*LiteralNode).Value.(int), 5)
	p2, _ := Parse(`="Hello, World!"`)
	assert.Equal(t, p2.(*LiteralNode).Value.(string), "Hello, World!")
	p3, _ := Parse(`=5`)
	assert.Equal(t, p3.(*LiteralNode).Value.(int), 5)
	p4, _ := Parse(`=Add(1, 2 ,3)`)
	assert.Equal(t, p4.(*FunctionNode).Name, "ADD")
	assert.Equal(t, len(p4.(*FunctionNode).Args), 3)
	p5, _ := Parse(`=Concat("Hello", ",", "World!")`)
	assert.Equal(t, p5.(*FunctionNode).Name, "CONCAT")
	p6, _ := Parse(`=A1`)
	assert.Equal(t, p6.(*ReferenceNode).Row, 0)
	assert.Equal(t, p6.(*ReferenceNode).Col, 0)
	p7, _ := Parse(`=Concat(A1, AA45)`)
	assert.Equal(t, p7.(*FunctionNode).Name, "CONCAT")
	assert.Equal(t, len(p7.(*FunctionNode).Args), 2)
	assert.Equal(t, p7.(*FunctionNode).Args[0].(*ReferenceNode).Row, 0)
	assert.Equal(t, p7.(*FunctionNode).Args[0].(*ReferenceNode).Col, 0)
	assert.Equal(t, p7.(*FunctionNode).Args[1].(*ReferenceNode).Row, 44)
	assert.Equal(t, p7.(*FunctionNode).Args[1].(*ReferenceNode).Col, 26)
}
