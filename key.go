package goose

import (
	"fmt";
	"reflect";
)

type Tag interface{}

type key struct {
	typeLiteral reflect.Type
	tag Tag
}

type Key interface {
}

func CreateKeyForType(typeLiteral reflect.Type) Key {
	return key {
		typeLiteral,
		nil,
	}
}

func CreateKeyForTaggedType(typeLiteral reflect.Type, tag Tag) Key {
	return key {
		typeLiteral,
		tag,
	}
}

func (this key) String() {
	if this.tag == nil {
		fmt.Sprintf("%v", this.typeLiteral)
	} else {
		fmt.Sprintf("%v(%v)", this.typeLiteral, this.tag)
	}
}
