/*******
* @Author:qingmeng
* @Description:
* @File:check
* @Date2022/2/16
 */

package tool

import (
	"golang.org/x/crypto/bcrypt"
	"log"
	"regexp"
)

//密码加盐
func AddSalt(pwd string) (pwdHash string, err error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return
	}
	pwdHash = string(hash)
	return
}

//CheckPassword 校验密码
func CheckPassword(pwdHash, pwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(pwdHash), []byte(pwd))
	return err == nil
}

// CheckIfSensitive 检查敏感词和sql关键符号，防止sql注入，存在返回true
func CheckIfSensitive(s string) bool {
	return defaultTrie.CheckWords(s)
}

// CheckPasswordLever 检查密码强度(0.不合规范 1~4为含多少种类别)
func CheckPasswordLever(ps string) int {
	if len(ps) < 6 {
		return 0
	}

	//初始密码强度
	n := 4
	num := `[0-9]{1}`
	a_z := `[a-z]{1}`
	A_Z := `[A-Z]{1}`
	symbol := `[!@#~$%^&*()+|_]{1}`
	space := `[ '-]`

	//不允许存在空格或-或'，防sql注入
	ok, err := regexp.MatchString(space, ps)
	if err != nil {
		log.Println("match space err:", err)
		return 0
	}
	if ok {
		return 0
	}
	ok, err = regexp.MatchString(num, ps)
	if err != nil {
		log.Println("match num err:", err)
		return 0
	}
	if !ok {
		n--
	}

	ok, err = regexp.MatchString(a_z, ps)
	if err != nil {
		log.Println("match a_z err:", err)
		return 0
	}
	if !ok {
		n--
	}

	ok, err = regexp.MatchString(A_Z, ps)
	if err != nil {
		log.Println("match A_Z err:", err)
		return 0
	}
	if !ok {
		n--
	}

	ok, err = regexp.MatchString(symbol, ps)
	if err != nil {
		log.Println("match symbol err:", err)
		return 0
	}
	if !ok {
		n--
	}

	return n
}
