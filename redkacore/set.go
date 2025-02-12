package redkacore

import (
	"fmt"

	"github.com/nalgeon/redka"
)

func RedkaSet(db *redka.DB, key string, value string, field string) {
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
		err := db.Str().Set(k.Key, value)
		if err != nil {
			panic(err)
		}
	case 2:
		fmt.Println("not supported")
	case 3:
		fmt.Println("not supported")
	case 4:
		_, err := db.Hash().Set(k.Key, field, value)
		if err != nil {
			panic(err)
		}
	case 5:
		fmt.Println("not supported")
	default:
		fmt.Println("not supported")
	}
}
