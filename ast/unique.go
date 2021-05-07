package ast

import (
	"encoding/binary"
	"fmt"
	"reflect"
	"strconv"

	"github.com/sirkon/protoast/internal/errors"
	"github.com/spaolacci/murmur3"
)

// Unique интерфейс для уникализации инстансов в AST привязывающий к данному объекту некое уникальное для контекстного
// множества значение
type Unique interface {
	setUniqueKey(ctx UniqueContext)
	getUniqueKey() string
}

// UniqueContext контекст для создания уникальных значений
type UniqueContext map[string]struct{}

var _ Unique = &unique{}

type unique struct {
	value string
}

func (k *unique) setUniqueKey(ctx UniqueContext) {
	value := murmur3.New32WithSeed(uint32(len(ctx)))
	for {
		v := strconv.FormatUint(uint64(value.Sum32()), 16)
		if _, ok := ctx[v]; !ok {
			ctx[v] = struct{}{}
			k.value = v
			break
		}

		_ = binary.Write(value, binary.LittleEndian, v)
	}
}

func (k *unique) getUniqueKey() string {
	return k.value
}

// SetUnique устанавливает уникальное в рамках UniqueContext значение для данного Unique
func SetUnique(k Unique, ctx UniqueContext) {
	if len(k.getUniqueKey()) == 0 {
		k.setUniqueKey(ctx)
	}
}

// GetUnique получает ключ для данного k
func GetUnique(k Unique) string {
	return k.getUniqueKey()
}

// GetFieldKey получает ключ для поля данного k
func GetFieldKey(k Unique, fieldAddr interface{}) string {
	if reflect.TypeOf(fieldAddr).Kind() != reflect.Ptr {
		panic(errors.Newf("second parameter must be a pointer to one of k object field, got %T instead", fieldAddr))
	}

	val := reflect.ValueOf(k)
	for val.Type().Kind() != reflect.Ptr {
		panic(errors.Newf("invalid incoming object: must be a pointer to struct, got %T", k))
	}
	val = val.Elem()
	if val.Type().Kind() != reflect.Struct {
		panic(errors.Newf("invalid incoming object: must be a pointer to struct, got %T", k))
	}
	rawFieldAddr := reflect.ValueOf(fieldAddr).Pointer()
	var name string
	for i := 0; i < val.NumField(); i++ {
		realFieldAddr := val.Field(i).Addr().Pointer()
		if realFieldAddr == rawFieldAddr {
			name = val.Type().Field(i).Name
			break
		}
		if realFieldAddr > rawFieldAddr {
			break
		}
	}
	if len(name) == 0 {
		panic(errors.Newf("given pointer does not address any field in the given object"))
	}
	return fmt.Sprintf("%s::%s", GetUnique(k), name)
}
