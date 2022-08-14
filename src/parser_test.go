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
	assert.Equal(t, p4.(*FunctionNode).Name, "add")
	assert.Equal(t, len(p4.(*FunctionNode).Args), 3)
	p5, _ := Parse(`=Concat("Hello", ",", "World!")`)
	assert.Equal(t, p5.(*FunctionNode).Name, "concat")
}
