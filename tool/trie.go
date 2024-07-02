/*******
* @Author:qingmeng
* @Description:
* @File:trie
* @Date2021/12/11
 */

package tool

var defaultSensitiveWords = []string{"你妈", "傻逼", "反中国", "反人类", " ", "\"", "'", "=", "--", "%", "*", "?", "^"}
var defaultTrie = NewTrie()

//初始化敏感词字典树
func init() {
	defaultTrie.Insert(defaultSensitiveWords)
}

type Trie struct {
	root *TrieNode
}

func NewTrie() *Trie {
	return &Trie{
		root: newTrieNode(),
	}
}

//字典树结点
type TrieNode struct {
	children map[interface{}]*TrieNode
	end      bool //标记该节点是否为尾结点
}

func newTrieNode() *TrieNode {
	return &TrieNode{
		children: make(map[interface{}]*TrieNode),
		end:      true,
	}
}

//添加结点
func (trie *Trie) Insert(words []string) {
	if len(words) == 0 {
		return
	}
	for i := 0; i < len(words); i++ {
		node := trie.root
		word := []rune(words[i])
		for _, v := range word {
			//不存在该字段时创建新结点
			if _, ok := node.children[v]; !ok {
				node.end = false
				node.children[v] = newTrieNode()
			}
			node = node.children[v]
		}
	}
	trie.root.end = false
}

//检查是否存在敏感词，存在返回true
func (trie *Trie) CheckWords(word string) bool {
	node := trie.root
	for _, v := range []rune(word) {
		//匹配上字符
		if _, ok := node.children[v]; ok {
			//fmt.Println(string(v))
			//当一个敏感词匹配到末尾时，匹配成功
			if node.children[v].end {
				return true
			}
			node = node.children[v]
		} else { //从根节点开始重新匹配
			node = trie.root
		}
	}
	return false
}
