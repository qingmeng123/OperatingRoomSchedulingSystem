<script setup>
import { Search } from '@element-plus/icons-vue'
import {ref, reactive, onMounted} from 'vue'
import {get, post} from "@/net";
import {ElMessage, ElMessageBox} from "element-plus";


//表单校验规则
// 表单校验rules中校验规则的名字必须和绑定在form表单上的数据体中的待校验的数据项名字一致
const formRef = ref()

const validateTel = (rule, value, callback) => {
  if (value === '') {
    callback(new Error('请输入电话号码'))
  } else if(!/^\d{11}$/.test(value)){
    callback(new Error('请输入11位数字电话号码'))
  } else {
    callback()
  }
}

const validatePwd = (rule, value, callback) => {
  if(!/^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)[A-Za-z\d@$!%*?&]{8,}$/.test(value)){
    callback(new Error('长度8-16，同时包含大小写字母和数字'))
  } else {
    callback()
  }
}

const validatePassword = (rule, value, callback) => {
  if (value === '') {
    callback(new Error('请再次输入密码'))
  } else if (value !== User.pwd) {
    callback(new Error("两次输入的密码不一致"))
  } else {
    callback()
  }
}

const rules = {
  userName: [
    { required: true, message:'请输入姓名', trigger: ['blur', 'change'] },
    { min: 2, max: 10, message: '姓名的长度必须在2-10个字符之间', trigger: ['blur', 'change'] },
  ],
  loginName: [
    { required: true,message:'请输入用户名', trigger: ['blur', 'change'] },
    { min: 2, max: 15, message: '登录名的长度必须在2-15个字符之间', trigger: ['blur', 'change'] }
  ],
  pwd: [
    { required: true, validator: validatePwd, trigger: ['blur', 'change'] },
    { min: 8, max: 16, message: '密码的长度必须在8-16个字符之间', trigger: ['blur', 'change'] }
  ],
  tel: [
    { required: true, validator: validateTel, trigger: 'blur'},
  ],
  rePassword: [
    { required: true, validator: validatePassword, trigger: ['blur', 'change'] },
  ],
  pid:[
    {required:true,message:'请选择用户职位',trigger: ['blur', 'change']},
    {min:0,max:10,trigger: ['blur', 'change']}
  ]
}
const modifyRules = {
  userName: [
    { required: false,  trigger: ['blur', 'change'] },
    { min: 2, max: 15, message: '用户名的长度必须在2-15个字符之间', trigger: ['blur', 'change'] },
  ],
  tel: [
    { required: false,  trigger: 'blur'},
    {min:11,max:11,message:'电话长度11位'}
  ],
  pwd: [
    { required: false, trigger: ['blur', 'change'] },
    { min: 8, max: 16, message: '密码的长度必须在8-16个字符之间', trigger: ['blur', 'change'] }
  ],
}
//表单校验规则js结束

/*分页参数定义*/
const pagingOn = ref(true)
const pageNum = ref(1)
const pages = ref(0)

// 搜索参数列表
const searchInfo = reactive({
  searchUserName: '',
  searchTel: ''
})

//修改新增窗口显示隐藏控制变量
const dialogVisible = ref(false)
const addDialogVisible = ref(false)

// 模糊查询方法
function search(){
  post('/api/admin/query/userList',{
    name: searchInfo.searchUserName,
    phone: searchInfo.searchTel,
    pageNum: pageNum.value
  },(message) => {
    const map = new Map(Object.entries(message))
    tableData.value = map.get('pageList')
    pages.value = map.get('totalPage')
  },(message) => {
    ElMessage.error(message)
  })
}

// 修改方法
function modify(row){
  User.id = row.id
  User.userName = row.name
  User.loginName = row.username
  User.tel = row.tel
  User.job = row.job
  User.work_year=row.work_year
  User.post_id=row.post_id
  User.pwd = row.pwd
  User.gender=row.gender
  dialogVisible.value = true
}

//提交修改的信息
function submitModify(){
  dialogVisible.value = false
  formRef.value.validate((isValid) => {
    if(isValid) {
      post('/api/admin/modify', {
        id:User.id,
        username:User.loginName,
        password: User.pwd,
        name: User.userName,
        job: User.job,
        phone: User.tel,
        work_year:User.work_year,
        post_id:User.post_id,
        gender:User.gender,
        group_id:User.group_id
      }, (message) => {
        dialogVisible.value = false
        clearUser()
        ElMessage.success(message)
        get('/api/admin/userList', (message) => {
          const map = new Map(Object.entries(message))
          tableData.value = map.get('pageList')
          pages.value = map.get('totalPage')
          pageNum.value = map.get('pageNum')
        }, (message) => {
          ElMessage.error(message)
        })
      }, (message) => {
        ElMessage.error(message)
      })
    }else {
      ElMessage.warning('请按要求填写表单内容！')
    }
  })
}

/*新增用户信息*/
function submitAddInfo(){
  formRef.value.validate((isValid) => {
    if(isValid) {
      post('/api/user/register', {
        username: User.loginName,
        password: User.pwd,
        name: User.userName,
        post_id:User.post_id,
        work_year:User.work_year,
        phone: User.tel,
        gender:User.gender
      }, (message) => {
        addDialogVisible.value = false
        clearUser()
        ElMessage.success(message)
        get('/api/admin/userList', (message) => {
          const map = new Map(Object.entries(message))
          tableData.value = map.get('pageList')
          pages.value = map.get('totalPage')
          pageNum.value = map.get('pageNum')
        }, (message) => {
          ElMessage.error(message)
        })
      }, (message) => {
        ElMessage.error(message)
      })
    } else {
      ElMessage.warning('请按要求填写表单内容！')
    }
  })
}


/*表格js*/
//定义表格需要变量
const tableData = ref([])
const User = reactive({
  id: '',
  loginName: '',
  userName: '',
  tel: '',
  job: '',
  post_id:'',
  work_year:'',
  pwd:'',
  rePassword: '',
  group_id:0,
  gender:''
})

let jobs=ref([]);
const multipleTableRef = ref()
const multipleSelection = ref([])

//清空User数据体方法
const clearUser = () => {
  User.id = ''
  User.userName = ''
  User.loginName = ''
  User.tel = ''
  User.job = ''
  User.post_id=''
  User.work_year=''
  User.pwd = ''
  User.rePassword = ''
  User.group_id=0
  User.gender=''
}

/*页面打开自动查询用户信息*/
onMounted(() => {
  get('/api/post',(message)=>{
    jobs=message
    },(message)=>{
    ElMessage.error(message)
  })
  get('/api/admin/userList', (message) => {
    const map = new Map(Object.entries(message))
    tableData.value = map.get('pageList')
    pages.value = map.get('totalPage')
    pageNum.value = map.get('pageNum')
  }, (message) => {
    ElMessage.error(message)
  })
})

/*下一页方法实现*/
function nextPage(val){
  pageNum.value = val
  post('/api/admin/query/userList',{
    pageNum:pageNum.value,
    name:searchInfo.searchUserName,
    phone:searchInfo.searchTel
  },
      (message) => {
    const map = new Map(Object.entries(message))
    tableData.value = map.get('pageList')
    pages.value = map.get('totalPage')
  }, (message) => {
    ElMessage.error(message)
  })
}

/*删除方法*/

//封装删除方法公共代码
const delUser = (row) => {
  ElMessageBox.confirm(
      '确定要删除选中数据吗？',
      '警告',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
  )
      .then(() => {
        post('/api/admin/delUser', {
         id:row.id
        }, (message) => {
          get('/api/admin/userList', (message) => {
            const map = new Map(Object.entries(message))
            tableData.value = map.get('pageList')
            pages.value = map.get('totalPage')
            pageNum.value = map.get('pageNum')
          }, (message) => {
            ElMessage.error(message)
          })
          ElMessage({
            type: 'success',
            message: '删除成功',
          })
        }, (message) => {
          ElMessage.error(message)
        })

      })
      .catch(() => {
        ElMessage({
          type: 'info',
          message: '取消删除',
        })
      })
}


/*多选实现方法，
会将选中数据放入multipleSelection*/
const handleSelectionChange = (val) => {
  multipleSelection.value = val
}

const getPostName = (postId) => {
  const post = jobs.find(post => post.id === postId)
  return post ? post.name : '暂无职务' // 如果找不到对应的职务名称，则显示'未知'
}

</script>

<template>
  <el-space wrap>
    <el-input style="width: 200px;margin-bottom: 10px;" v-model="searchInfo.searchUserName" placeholder="请输入姓名" :prefix-icon="Search" />
    <el-input style="width: 200px;margin-bottom: 10px;" v-model="searchInfo.searchTel" placeholder="请输入电话号码" :prefix-icon="Search" />
    <el-button style="margin-bottom: 10px;" type="success" :icon="Search" @click="search">搜索</el-button>
    <el-button style="margin-bottom: 10px;" type="success" @click="addDialogVisible = true">新增人员</el-button>
  </el-space>

  <el-table
      ref="multipleTableRef"
      :data="tableData"
      style="width: 100%"
      border
      @selection-change="handleSelectionChange"
  >
    <el-table-column type="selection" width="55" />
    <el-table-column label="序号" width="100" >
    <template #default="{ row, $index }">
      <span>{{ $index + 1 }}</span>
    </template>
    </el-table-column>
    <el-table-column property="username" label="用户名" width="200" />
    <el-table-column property="name" label="姓名" width="200" prop="name" sortable/>
    <el-table-column property="phone" label="手机号码" width="300" prop="phone" sortable/>
    <el-table-column label="职务" width="200" prop="post_id" sortable>
      <template #default="{ row }">
        {{ getPostName(row.post_id) }}
      </template>
    </el-table-column>
    <el-table-column property="gender" label="性别" width="150" sortable>
      <template #default="{ row }">
        {{ row.gender?'女':'男' }}
      </template>
    </el-table-column>
    <el-table-column property="work_year" label="工龄" width="150" sortable/>
    <el-table-column label="操作" align="center" #default="scope">
      <el-space wrap>
        <div>
          <el-button size="small" @click="modify(scope.row)">修改</el-button>
          <el-button size="small" type="danger" @click="delUser(scope.row)">删除</el-button>
        </div>
      </el-space>
    </el-table-column>
  </el-table>

  <!--  分页代码  -->
  <div style="width: 100%;margin-top: 20px;display: flex;justify-content: center;" v-if="pagingOn">
    <el-pagination background layout="prev, pager, next" :page-count="pages" :current-page="pageNum"
                   @current-change="nextPage"/>
  </div>

  <!--  新增人员代码  -->
  <el-dialog v-model="addDialogVisible" title="新增人员信息" width="50%" align-center center>
    <div style="display: flex;justify-content: center;">
      <el-form :model="User" label-width="120px" :rules="rules" ref="formRef">
        <el-form-item label="姓名" prop="userName">
          <el-input v-model="User.userName" type="text"/>
        </el-form-item>
        <el-form-item label="用户名" prop="loginName">
          <el-input v-model="User.loginName" style="width: 300px;" />
        </el-form-item>
        <el-form-item label="用户电话号码" prop="tel">
          <el-input v-model="User.tel" style="width: 300px;" />
        </el-form-item>
        <el-form-item label="用户职务">
          <el-select v-model="User.post_id" placeholder="请选择用户职务" style="width: 300px;" prop="pid">
            <el-option
                v-for="item in jobs"
                :key="item.id"
                :label="item.name"
                :value="item.id">
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="用户性别">
          <el-select v-model="User.gender" placeholder="请选择用户性别" style="width: 300px;">
            <el-option label="男" :value="false"></el-option>
            <el-option label="女" :value="true"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="工作年限" prop="tel">
          <el-input v-model="User.work_year" style="width: 300px;" />
        </el-form-item>
        <el-form-item label="用户密码" prop="pwd">
          <el-input v-model="User.pwd" type="password" show-password/>
        </el-form-item>
        <el-form-item label="重复密码" prop="rePassword">
          <el-input v-model="User.rePassword" type="password" show-password/>
        </el-form-item>
      </el-form>
    </div>
    <template #footer>
            <span class="dialog-footer">
                <el-button @click="addDialogVisible = false">取消</el-button>
                <el-button type="primary" @click="submitAddInfo">
                    保存
                </el-button>
            </span>
    </template>
  </el-dialog>

  <!--  修改人员代码  -->
  <el-dialog v-model="dialogVisible" title="用户信息修改" width="50%" align-center center>
    <div style="display: flex;justify-content: center;">
      <el-form :model="User" label-width="120px" :rules="modifyRules" ref="formRef">
        <el-form-item label="姓名" prop="userName">
          <el-input v-model="User.userName" type="text"/>
        </el-form-item>
        <el-form-item label="用户电话号码" prop="tel">
          <el-input v-model="User.tel" style="width: 300px;" />
        </el-form-item>
        <el-form-item label="用户职务">
          <el-select v-model="User.post_id" placeholder="请选择用户职务" style="width: 300px;">
            <el-option
                v-for="item in jobs"
                :key="item.id"
                :label="item.name"
                :value="item.id">
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="用户性别">
          <el-select v-model="User.gender" placeholder="请选择用户性别" style="width: 300px;">
            <el-option label="男" :value="false"></el-option>
            <el-option label="女" :value="true"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="工作年龄" prop="work_year">
          <el-input v-model="User.work_year" style="width: 300px;" />
        </el-form-item>
        <el-form-item label="权限" >
          <el-select v-model="User.group_id" placeholder="请选择用户权限" style="width: 300px;">
            <el-option label="普通权限" :value="0"></el-option>
            <el-option label="管理员权限" :value="1"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="用户密码" prop="pwd">
          <el-input v-model="User.pwd" style="width: 300px;" />
        </el-form-item>
      </el-form>
    </div>
    <template #footer>
            <span class="dialog-footer">
                <el-button @click="dialogVisible = false">取消</el-button>
                <el-button type="primary" @click="submitModify">
                    保存
                </el-button>
            </span>
    </template>
  </el-dialog>

</template>
<style scoped></style>