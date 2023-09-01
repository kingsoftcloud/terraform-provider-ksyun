package ksyun

import (
	"encoding/json"
	"reflect"
	"sort"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type LenFunc func([]interface{}) int
type SwapFunc func([]interface{}, int, int)
type LessFunc func([]interface{}, int, int) bool

type AnyList struct {
	list []interface{}

	lenFunc  LenFunc
	swapFunc SwapFunc
	lessFunc LessFunc
}

func (l AnyList) Len() int {
	if l.lenFunc == nil {
		return len(l.list)
	}
	return l.lenFunc(l.list)
}

func (l AnyList) Swap(i, j int) {
	if l.swapFunc == nil {
		l.list[i], l.list[j] = l.list[j], l.list[i]
		return
	}
	l.swapFunc(l.list, i, j)
}

func (l AnyList) Less(i, j int) bool {
	return l.lessFunc(l.list, i, j)
}

func (l AnyList) SetLessFunc(lessFunc LessFunc) {
	l.lessFunc = lessFunc
}

func (l AnyList) SetLenFunc(lenFunc LenFunc) {
	l.lenFunc = lenFunc
}

func (l AnyList) SetSwapFunc(swapFunc SwapFunc) {
	l.swapFunc = swapFunc
}

func (l AnyList) SetList(any []interface{}) {
	l.list = any
}

func (l AnyList) IsCompleted() bool {
	if l.lessFunc == nil {
		return false
	}
	return true
}

func (l AnyList) GetList() []interface{} {
	return l.list
}

func NewAnyList(lessFunc LessFunc) AnyList {
	return AnyList{
		lessFunc: lessFunc,
	}
}

func albRuleSetLessFunc(list []interface{}, i, j int) bool {
	idx1 := list[i]

	val1If, idx1Ok := idx1.(map[string]interface{})
	if !idx1Ok {
		return false
	}
	if reflect.DeepEqual(val1If["alb_rule_type"], "domain") {
		return true
	}
	return false
}

func albRuleSetStateFunc(lessFunc LessFunc) schema.SchemaStateFunc {
	return func(i interface{}) string {
		switch i.(type) {
		case string:
			return i.(string)
		case int:
			return strconv.Itoa(i.(int))
		case []interface{}:
			anyList := NewAnyList(lessFunc)
			anyList.SetList(i.([]interface{}))

			sort.Sort(anyList)

			valStr, err := json.Marshal(anyList.GetList())
			if err != nil {
				panic(err)
			}
			return string(valStr)
		}

		return ""
	}
}
