/*******
* @Author:qingmeng
* @Description:
* @File:cache
* @Date:2024/5/9
 */

package cache

import "github.com/go-redis/redis/v8"

var RedisClient *redis.Client

func InitCache() {
	ctx := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	RedisClient = ctx
}
