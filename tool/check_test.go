/*******
* @Author:qingmeng
* @Description:
* @File:check_Test
* @Date:2022/7/29
 */

package tool

import (
	"fmt"
	"strconv"
	"testing"
)

func TestCheckPasswordLever(t *testing.T) {
	ps := "ah 'aaaaa"
	lever := CheckPasswordLever(ps)
	fmt.Println(lever)
}

var users = []struct {
	name string
	out  bool
}{
	{"你妈", true},
	{"你好", false},
	{"aa=", true},
	{"aa-", false},
	{"aa--", true},
	{"aa%", true},
	{"aa'", true},
	{"\"", true},
	{"username=?", true},
}

func TestCheckIfSensitive(t *testing.T) {
	for i, user := range users {
		t.Run("test"+strconv.Itoa(i), func(t *testing.T) {
			get := CheckIfSensitive(user.name)
			if get != user.out {
				t.Errorf("got %v,want %v", get, user.out)
			}
		})
	}

}
