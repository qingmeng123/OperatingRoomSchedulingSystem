<template>

  <!--  新增人员代码  -->
  <el-dialog v-model="addDialogVisible" title="新增人员信息" width="50%" align-center center>
    <div style="display: flex;justify-content: center;">
      <el-form :model="User" label-width="120px" :rules="registerRules" ref="formRef">
        <el-form-item label="姓名" prop="userName">
          <el-input v-model="User.userName" type="text"/>
        </el-form-item>
        <el-form-item label="账号" prop="loginName">
          <el-input v-model="User.loginName" style="width: 300px;" />
        </el-form-item>
        <el-form-item label="用户电话号码" prop="tel">
          <el-input v-model="User.tel" style="width: 300px;" />
        </el-form-item>
        <el-form-item label="用户职务" >
          <el-select v-model="User.post_id" placeholder="请选择用户职务" style="width: 300px;" prop="pid">
            <el-option
                v-for="item in jobs"
                :key="item.id"
                :label="item.name"
                :value="item.id">
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="工作年限">
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

  <div>
    <div class="login-container">
      <div style="width: 350px" class="login-box">
        <div style="font-weight: bold; font-size: 24px; text-align: center; margin-bottom: 30px">登录</div>
        <el-form :model="form" ref="formRef" :rules="rules">
          <el-form-item prop="username">
            <el-input prefix-icon="User" v-model="form.username" placeholder="请输入用户名"/>
          </el-form-item>

          <el-form-item prop="password">
            <el-input type="password" prefix-icon="Lock" v-model="form.password" placeholder="请输入密码" show-password/>
          </el-form-item>

          <el-form-item>
            <el-button type="primary" style="width: 100%" @click="login">登录</el-button>
          </el-form-item>
        </el-form>
        <el-form-item>
          <el-button type="primary" style="width: 100%" @click="addDialogVisible = true">注册</el-button>
        </el-form-item>
      </div>

    </div>
  </div>
</template>

<script setup>
import {get,post} from '../net'
import {onMounted, reactive, ref} from "vue";
import {ElMessage} from "element-plus";
import {useCounterStore} from "@/stores/counter.js";
import router from "@/router/index.js";

const form = reactive({
  username: '',
  password: '',
  code: ''
})
let jobs=ref([]);

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
  if (value === '') {
    callback(new Error('请输入用户密码'))
  } else if(!/^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)[A-Za-z\d@$!%*?&]{8,}$/.test(value)){
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

const rules = reactive({
  username: [
    {required: true, message: '请输入用户名', trigger: 'blur'},
    {min: 3, max: 15, message: '长度必须在3到15', trigger: 'blur'},
  ],
  password: [
    {required: true, message: '请输入密码', trigger: 'blur'},
    {min: 3, max: 15, message: '长度必须在3到15', trigger: 'blur'},
  ],
})



const registerRules = {
  userName: [
    { required: true, message: '请输入姓名',  trigger: ['blur', 'change'] },
    { min: 2, max: 10, message: '姓名的长度必须在2-10个字符之间', trigger: ['blur', 'change'] },
  ],
  loginName: [
    { required: true,message: '请输入用户名',trigger: ['blur', 'change'] },
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



const formRef = ref()
const store = useCounterStore()
const login = () => {
  formRef.value.validate((valid) => {
    console.log(valid)
    if (valid) {
      post('/api/user/login', {
            username: form.username,
            password: form.password
          }, (message) => {
            ElMessage.success('登录成功')
            window.localStorage.setItem("token", message.token)
            get('/api/user/info',(message)=>{
              store.auth.user=message
              router.push('/index')

            },(message)=>{
              store.auth.user=null
            })
          },
          (message) => {
            ElMessage.error(message)
          }
      )
    }
  })
}

/*新增用户信息*/
const User = reactive({
  id: '',
  loginName: '',
  userName: '',
  tel: '',
  job: '',
  post_id:'',
  work_year:'',
  pwd:'',
  rePassword: ''
})
const addDialogVisible = ref(false)
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
}
onMounted(() => {
  get('/api/post',(message)=>{
    jobs=message
  },(message)=>{
    ElMessage.error(message)
  })
})

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
      }, (message) => {
        addDialogVisible.value = false
        clearUser()
        ElMessage.success("注册成功")
        router.go(0);
      }, (message) => {
        ElMessage.error(message)
      })
    } else {
      ElMessage.warning('请按要求填写表单内容！')
    }
  })
}

</script>

<style scoped>
.login-container {
  min-height: 100vh;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
  background-image: url("@/assets/bg.png");
  background-size: cover;
}

.login-box {
  background-color: rgba(255, 255, 255, .8);
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
  padding: 30px;
  border-radius: 5px;
}
</style>