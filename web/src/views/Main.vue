<script setup>
import {onMounted, reactive, ref} from 'vue'
import {Document, OfficeBuilding, User,} from '@element-plus/icons-vue'

import {RouterView} from 'vue-router';
import router from "@/router/index.js";

import {useCounterStore} from "@/stores/counter";
import {ElMessage} from "element-plus";
import {get} from "@/net/index.js";

let infoDialogVisible = ref(false)
const store = useCounterStore()
let tabIndex = 1
const formRef = ref()
const editableTabsValue = ref('1')
const editableTabs = ref([
  {
    title: '首页',
    name: '1',
    path: '/index'
  }
])


//修改密码JS代码 start
const validateNewPwd = (rule, value, callback) => {
  if (value === '') {
    callback(new Error('请输入用户密码'))
  } else if (!/^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)[A-Za-z\d@$!%*?&]{8,}$/.test(value)) {
    callback(new Error('长度8-16，同时包含大小写字母和数字'))
  } else {
    callback()
  }
}

const modifyRules = {
  newPwd: [
    {required: true, validator: validateNewPwd, trigger: ['blur', 'change']},
    {min: 8, max: 16, message: '密码的长度必须在8-16个字符之间', trigger: ['blur', 'change']}
  ],
  oldPwd: [
    {required: true, trigger: ['blur', 'change']}
  ]
}

//定义修改密码动态显示变量
const modifyPwdVisible = ref(false)


// 定义一个等待函数，等待指定的毫秒数后 resolve
const wait = (milliseconds) => {
  return new Promise((resolve) => {
    setTimeout(() => {
      resolve();
    }, milliseconds);
  });
};

//提交修改的信息
function submitModify() {
  formRef.value.validate((isValid) => {
    if (isValid) {
      post('/api/user/password', {
        username: store.auth.user.username,
        old_password: modifyPwdInfo.oldPwd,
        new_password: modifyPwdInfo.newPwd

      }, (message) => {
        modifyPwdVisible.value = false
        modifyPwdInfo.newPwd = ''
        modifyPwdInfo.oldPwd = ''
        ElMessage.success(message)
        wait(1000).then(() => {
          router.push('/login').then(() => {
            store.auth.user = null
            router.go(0)
          });
        });
      }, (message) => {
        ElMessage.error(message)
      })
    } else {
      ElMessage.warning('请按要求填写表单内容！')
    }
  })
}

const modifyPwdInfo = reactive({
  newPwd: '',
  oldPwd: ''
})

//修改密码JS代码 END

//单击人员管理选项卡触发
const onUserManager = () => {
  router.push("/index/peopleManager")
  addTab("人员管理", "/index/peopleManager")
}

//单击职位管理选项卡触发
const onPostManager = () => {
  router.push("/index/postManager")
  addTab("职位管理", "/index/postManager")
}

//单击职位管理选项卡触发
const onOperatingRoomManager = () => {
  router.push("/index/operatingRoomManager")
  addTab("手术室管理", "/index/operatingRoomManager")
}

//单击排班管理选项卡触发
const onSchedulingManager = () => {
  router.push("/index/schedulingManager")
  addTab("手术排程", "/index/schedulingManager")
}

//单击排班查询选项卡触发
const onDutyQuery = () => {
  router.push("/index/dutyQuery")
  addTab("排程查询", "/index/dutyQuery")
}

//选中选项卡触发事件
const selectTab = (activeName) => {
  const index = parseInt(activeName)
  const path = editableTabs.value[index - 1].path
  router.push(path)
}

//添加选项卡
const addTab = (targetName, path) => {
  const newTabName = `${++tabIndex}`
  editableTabs.value.push({
    title: targetName,
    name: newTabName,
    path: path,
    close: 'closable'
  })
  console.log(editableTabs.value)
  console.log(editableTabsValue.value)
  editableTabsValue.value = newTabName
  console.log(editableTabsValue.value)

}

//删除选项卡
const removeTab = (targetName) => {
  const tabs = editableTabs.value
  let activeName = editableTabsValue.value
  if (activeName === targetName) {
    tabs.forEach((tab, index) => {
      if (tab.name === targetName) {
        const nextTab = tabs[index + 1] || tabs[index - 1]
        if (nextTab) {
          activeName = nextTab.name
        }
      }
    })
  }

  editableTabsValue.value = activeName
  editableTabs.value = tabs.filter((tab) => tab.name !== targetName)
}

//下拉菜单项触发事件
const handleCommand = (command) => {
  if (command === 'logout') {
    router.push('/')
    store.auth.user = null
  }
  if (command === 'modifyPassword') {
    modifyPwdVisible.value = true
  }
}

/*时间js开始*/
const time = ref('')
const dateTime = () => {
  let now = new Date()
  let year = now.getFullYear(); //获取完整年份
  let month = now.getMonth() + 1; //获取当前月份（0-11，0代表1月）
  let day = now.getDate(); //获取当前日（1-31）
  let hour = now.getHours();//获取当前小时数
  let minute = now.getMinutes();  //获取当前分钟数
  let second = now.getSeconds(); //获取当前秒数
  if (month < 10) {
    month = '0' + month
  }
  if (day < 10) {
    day = '0' + day
  }
  if (hour < 10) {
    hour = '0' + hour
  }
  if (minute < 10) {
    minute = '0' + minute
  }
  if (second < 10) {
    second = '0' + second
  }
  time.value = year + '-' + month + '-' + day + ' ' + hour + ':' + minute + ':' + second
}
setInterval(() => {
  dateTime()
}, 1000)
/*时间js结束*/
let jobs = ref([]);
let info = ref()
onMounted(() => {
      get('/api/user/info', (message) => {
        info = message
      })

      get('/api/post', (message) => {
        jobs = message
      }, (message) => {
        ElMessage.error(message)
      })

    }
)

const getPostName = (postId) => {
  const post = jobs.find(post => post.id === postId)
  return post ? post.name : '暂无职务' // 如果找不到对应的职务名称，则显示'未知'
}

</script>

<template>
  <div class="common-layout">
    <el-container>
      <el-header class="header">
        <div style="cursor: pointer; font-weight: bold;font-size: 20px; margin-bottom: 30px;"
             @click="router.push('/index')">手术室排程系统
        </div>
        <div style="color:white;margin-bottom: 30px; font-weight: bold;font-size: 20px;margin-left: 600px">
          欢迎您！{{ store.auth.user.name }}&nbsp;&nbsp;当前时间：{{ time }}
        </div>
        <el-dropdown style="width: 500px;position: absolute;right:1%;" @command="handleCommand">
          <div class="avatar" style="margin-bottom: 30px;margin-left: 300px">
            <el-avatar src="src/assets/avatar.png"/>
          </div>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item @click="infoDialogVisible=true">查看信息</el-dropdown-item>
            </el-dropdown-menu>
            <el-dropdown-menu>
              <el-dropdown-item command="modifyPassword">修改密码</el-dropdown-item>
            </el-dropdown-menu>
            <el-dropdown-menu>
              <el-dropdown-item command="logout">退出登录</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </el-header>
      <el-container style="height: 100%;width: 100%;position: fixed;left: 0%;top: 7%;background-color: #F7F7F7;">
        <el-aside class="aside">
          <el-menu default-active="/index" class="el-menu-vertical-demo"
                   style="height: 100%;background-color: #F7F7F7;" router>
            <el-menu-item index="1" @click="onUserManager()" v-if="store.auth.user.group_id == 1">
              <template #title>
                <el-icon>
                  <User/>
                </el-icon>
                <span>人员管理</span>
              </template>
            </el-menu-item>
            <el-menu-item index="2" @click="onPostManager()" v-if="store.auth.user.group_id == 1">
              <template #title>
                <el-icon>
                  <User/>
                </el-icon>
                <span>职位管理</span>
              </template>
            </el-menu-item>
            <el-menu-item index="3" @click="onOperatingRoomManager()" v-if="store.auth.user.group_id == 1">
              <template #title>
                <el-icon>
                  <OfficeBuilding/>
                </el-icon>
                <span>手术室管理</span>
              </template>
            </el-menu-item>
            <el-menu-item index="4" @click="onSchedulingManager()" v-if="store.auth.user.group_id == 1">
              <template #title>
                <el-icon>
                  <User/>
                </el-icon>
                <span>手术排程</span>
              </template>
            </el-menu-item>
            <el-menu-item index="5" route="/index/published" @click="onDutyQuery()">
              <template #title>
                <el-icon>
                  <Document/>
                </el-icon>
                <span>排程查询</span>
              </template>
            </el-menu-item>
          </el-menu>
        </el-aside>
        <el-main class="main">
          <el-tabs v-model="editableTabsValue" type="card" class="demo-tabs" @tab-remove="removeTab"
                   @tab-change="selectTab" style="position: relative;top: -3%;">
            <el-tab-pane v-for="item in editableTabs" :key="item.name" :label="item.title" :name="item.name"
                         :closable="item.close">
              <RouterView/>
            </el-tab-pane>
          </el-tabs>
        </el-main>
      </el-container>
    </el-container>
  </div>

  <!--  修改密码代码  -->
  <el-dialog v-model="modifyPwdVisible" title="密码修改" width="50%" align-center center>
    <div style="display: flex;justify-content: center;">
      <el-form :model="modifyPwdInfo" label-width="140px" :rules="modifyRules" ref="formRef">
        <el-form-item label="请输入旧密码验证" prop="oldPwd">
          <el-input v-model="modifyPwdInfo.oldPwd" type="password" show-password/>
        </el-form-item>
        <el-form-item label="请输入新密码" prop="newPwd">
          <el-input v-model="modifyPwdInfo.newPwd" type="password" show-password/>
        </el-form-item>
      </el-form>
    </div>
    <template #footer>
            <span class="dialog-footer">
                <el-button @click="modifyPwdVisible = false">取消</el-button>
                <el-button type="primary" @click="submitModify">
                    保存
                </el-button>
            </span>
    </template>
  </el-dialog>

  <el-dialog v-model="infoDialogVisible" title="个人信息" width="30%" align-center center>
    <div style="display: flex;justify-content: center;">
      <el-card style="max-width: 300px;font-size: medium">
        <p>{{ '姓名：' + info.name }}</p>
        <p>{{ '账号：' + info.username }}</p>
        <p>{{ '职位：' + getPostName(info.post_id) }}</p>
        <p>{{ '电话：' + info.phone }}</p>
        <p>性别：{{ info.gender ? '女' : '男' }}</p>
        <p>{{ '工龄：' + info.work_year }}</p>
        <p>权限：{{ info.group_id == 0 ? '访客' : '管理员' }}</p>
      </el-card>
    </div>
    <template #footer>
            <span class="dialog-footer">
                <el-button @click="infoDialogVisible = false">关闭</el-button>
            </span>
    </template>
  </el-dialog>
</template>

<style scoped>
.header {
  position: fixed;
  display: inline-flex;
  align-items: center;
  left: 0%;
  top: 0%;
  width: 100%;
  height: 100px;
  right: 20px;
  background-color: #539BD0;
  color: #fff;
}

.aside {
  /* margin-top: 50px; */
  width: 200px;
  height: 100%;
  background-color: aliceblue;
}

.main {
  width: 80%;
  height: 100%;
}

.avatar:not(:last-child) {
  border-right: 1px solid var(--el-border-color);
}
</style>