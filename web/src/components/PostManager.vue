<script setup>
import { Search } from '@element-plus/icons-vue'
import {ref, reactive, onMounted} from 'vue'
import {get, post,put} from "@/net";
import {ElMessage, ElMessageBox} from "element-plus";
const multipleTableRef = ref()
const dialogVisible = ref(false)
const addDialogVisible = ref(false)
const formRef = ref()

const rules = {
  name: [
    { required: true, message:'请输入职位名', trigger: ['blur', 'change'] },
    { min: 1, max: 10, message: '职位名的长度必须在1-10个字符之间', trigger: ['blur', 'change'] },
  ],
}


const Posts=reactive({
  id:'',
  name:'',
  number:''
})

const clearPosts = () => {
  Posts.id=''
  Posts.name=''
}

const data=reactive({
  tabalData:[]
})


/*页面打开自动查询职位信息*/
onMounted(() => {
  get('/api/post', (message) => {
    data.tableData=message
  }, (message) => {
    ElMessage.error(message)
  })
})

function modify(row){
  Posts.id=row.id
  Posts.name=row.name
  dialogVisible.value = true
}


function submitAddInfo(){
  formRef.value.validate((isValid) => {
    if(isValid) {
      put('/api/admin/post',{
        name:Posts.name
      },(message) => {
        addDialogVisible.value = false
        clearPosts()
        ElMessage.success(message)
            get('/api/post', (message) => {
              data.tableData=message
            }, (message) => {
              ElMessage.error(message)
            })
          },
          (message)=>{
          ElMessage.error(message)
          }
      )
    }else {
      ElMessage.warning('请按要求填写表单内容！')
    }
  })
}

function submitModify(){
  dialogVisible.value = false
  formRef.value.validate((isValid) => {
    if(isValid) {
      post('/api/admin/post',{
            id:Posts.id,
            name:Posts.name
          },(message) => {
            clearPosts()
            ElMessage.success(message)
            get('/api/post', (message) => {
              data.tableData=message
            }, (message) => {
              ElMessage.error(message)
            })
          },
          (message)=>{
            ElMessage.error(message)
          }
      )
    }else {
      ElMessage.warning('请按要求填写表单内容！')
    }
  })
}

</script>


<template>
<!--
  <div style="display: flex;margin-right: 50px; justify-content: flex-end;">
-->
  <el-row>
      <el-button  style="margin-right: 10px;" type="success" @click="addDialogVisible = true">新增职位</el-button>
  </el-row>
<!--  </div>-->

  <el-table
      ref="multipleTableRef"
      :data="data.tableData"
      style="width: 50%"
      border

  >

    <el-table-column type="selection" width="55" />
    <el-table-column label="序号" width="100" >
      <template #default="{ row, $index }">
        <span>{{ $index + 1 }}</span>
      </template>
    </el-table-column>
    <el-table-column prop="name" label="职位名" width="200" sortable/>
    <el-table-column prop="number" label="人员数" width="200" sortable/>
    <el-table-column label="操作" align="center" #default="scope">
      <el-space wrap>
      <div>
        <el-button size="small" @click="modify(scope.row)">修改</el-button>
      </div>
      </el-space>
    </el-table-column>
  </el-table>

<!--  &lt;!&ndash;  分页代码  &ndash;&gt;
  <div style="width: 100%;margin-top: 20px;display: flex;justify-content: center;" v-if="pagingOn">
    <el-pagination background layout="prev, pager, next" :page-count="pages" :current-page="pageNum"
                   @current-change="nextPage"/>
  </div>-->

  <el-dialog v-model="addDialogVisible" title="新增职位信息" width="30%" align-center center>
    <div style="display: flex;justify-content: center;">
      <el-form :model="Posts" label-width="120px" :rules="rules" ref="formRef">
        <el-form-item label="职位名" prop="name">
          <el-input v-model="Posts.name" type="text"/>
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

  <el-dialog v-model="dialogVisible" title="修改职位信息" width="30%" align-center center>
    <div style="display: flex;justify-content: center;">
      <el-form :model="Posts" label-width="120px" :rules="rules" ref="formRef">
        <el-form-item label="职位名" prop="name">
          <el-input v-model="Posts.name" type="text"/>
        </el-form-item>
      </el-form>
    </div>
    <template #footer>
            <span class="dialog-footer">
                <el-button @click="addDialogVisible = false">取消</el-button>
                <el-button type="primary" @click="submitModify">
                    保存
                </el-button>
            </span>
    </template>
  </el-dialog>

</template>
<style scoped></style>