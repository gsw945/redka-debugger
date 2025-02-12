package redkacore

import (
	"fmt"

	"github.com/nalgeon/redka"
)

func RedkaGet(db *redka.DB, key string) {
	k, err := db.Key().Get(key)
	if err != nil {
		panic(err)
	}
	intType := int(k.Type)
	typeName := RedkaTypeName(intType)
	fmt.Printf("key(ID=%d, Key=%s, Type=%v, Version=%d, ETime=%v, MTime=%v)\n", k.ID, k.Key, typeName, k.Version, k.ETime, k.MTime)
	switch intType {
	case 0:
		fmt.Println("not supported")
	case 1:
		valueString, err := db.Str().Get(k.Key)
		if err != nil {
			panic(err)
		}
		fmt.Printf("value: %s\n", valueString)
	case 2:
		valueList, err := db.List().Range(k.Key, 0, -1)
		if err != nil {
			panic(err)
		}
		for i, v := range valueList {
			fmt.Printf("value[%d]: %v\n", i, v)
		}
	case 3:
		valueSet, err := db.Set().Items(k.Key)
		if err != nil {
			panic(err)
		}
		for i, v := range valueSet {
			fmt.Printf("value[%d]: %v\n", i, v)
		}
	case 4:
		valueHash, err := db.Hash().Items(k.Key)
		if err != nil {
			panic(err)
		}
		for k, v := range valueHash {
			fmt.Printf("value[%s]: %v\n", k, v)
		}
	case 5:
		valueZSet, err := db.ZSet().Range(k.Key, 0, -1)
		if err != nil {
			panic(err)
		}
		for i, v := range valueZSet {
			fmt.Printf("value[%d]: %v\n", i, v)
		}
	default:
		fmt.Println("not supported")
	}
}
