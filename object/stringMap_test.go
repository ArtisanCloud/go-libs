package object

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_MergeStringMap(t *testing.T) {

	base := &StringMap{
		"key1": "123",
		"key2": "xxx",
		"key3": "",
		"key4": "101112",
	}

	toMap := &StringMap{
		"key1": "",
		"key2": "456",
		"key3": "789",
	}

	toMap = MergeStringMap(toMap, base)

	assert.EqualValues(t, &StringMap{
		"key1": "123",
		"key2": "456",
		"key3": "789",
		"key4": "101112",
	}, toMap)

}

func Test_ReplaceStringMapRecursive(t *testing.T) {

	base := &StringMap{
		"key1": "123",
		"key2": "456",
		"key3": "789",
		"key4": "nil",
	}

	base2 := &StringMap{
		"key1": "456",
		"key2": "base456",
		"key3": "",
		"key4": "nil",
		"key5": "&StringMap{}",
	}

	toMap := &StringMap{
		"key2": "",
		"key3": "123",
		"key4": "",
		"key5": "nil",
		"key6": "&StringMap{}",
	}

	toMap = ReplaceStringMapRecursive(toMap, base, base2)

	assert.EqualValues(t, &StringMap{
		"key1": "456",
		"key2": "base456",
		"key3": "",
		"key4": "nil",
		"key5": "&StringMap{}",
		"key6": "&StringMap{}",
	}, toMap)

}
