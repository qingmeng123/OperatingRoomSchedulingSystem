/*******
* @Author:qingmeng
* @Description:
* @File:obj_to_map
* @Date2022/1/31
 */

package tool

import "reflect"

// 结构体转换为map
func ObjToMap(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}
