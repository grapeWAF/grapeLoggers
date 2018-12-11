package rules

import (
	"fmt"
	"reflect"
)

// 比较2个对象合并成一个新对象
// 大量反射效率不高 ，不能用于大量处理
func CmpAndCopy(src, dist interface{}) (interface{}, error) {
	srcType := reflect.Indirect(reflect.ValueOf(src))
	distType := reflect.Indirect(reflect.ValueOf(dist))

	if srcType.Kind() != distType.Kind() {
		return nil, fmt.Errorf("异常错误，类型不同无法COPY")
	}

	newData := reflect.Indirect(reflect.New(srcType.Type()))
	// 循环比较并COPY到新数值
	t := srcType.Type()
	for i := 0; i < distType.NumField(); i++ {
		// 比较
		srcFV := srcType.Field(i)
		srcTV := t.Field(i)
		newFV := newData.FieldByName(srcTV.Name)
		distFV := distType.FieldByName(srcTV.Name) // 查找相同的数值
		if distFV.IsValid() == false {
			return nil, fmt.Errorf("异常错误，两者结构不同，无法进行copy操作")
		}

		if !srcFV.CanSet() {
			continue
		}

		// 开始比较
		newFV.Set(srcFV)
		if reflect.DeepEqual(srcFV.Interface(), distFV.Interface()) == false {
			newFV.Set(distFV)
		}
	}

	return newData.Addr().Interface(), nil
}

func RemoveDuplicatesAndEmpty(a []string) (ret []string) {
	a_len := len(a)
	for i := 0; i < a_len; i++ {
		if (i > 0 && a[i-1] == a[i]) || len(a[i]) == 0 {
			continue
		}
		ret = append(ret, a[i])
	}
	return
}
