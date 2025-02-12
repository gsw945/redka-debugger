package redkacore

import (
	"fmt"

	"github.com/nalgeon/redka"
)

func RedkaTypeName(t int) string {
	switch t {
	case 0:
		return "Any(0)"
	case 1:
		return "String(1)"
	case 2:
		return "List(2)"
	case 3:
		return "Set(3)"
	case 4:
		return "Hash(4)"
	case 5:
		return "ZSet(5)"
	default:
		return "<Unknown>"
	}
}

func RedkaKeys(db *redka.DB, pattern string) {
	keys, err := db.Key().Keys(pattern)
	if err != nil {
		panic(err)
	}
	for _, key := range keys {
		fmt.Printf("key(ID=%d, Key=%s, Type=%v, Version=%d, ETime=%v, MTime=%v)\n", key.ID, key.Key, RedkaTypeName(int(key.Type)), key.Version, key.ETime, key.MTime)
	}
}
