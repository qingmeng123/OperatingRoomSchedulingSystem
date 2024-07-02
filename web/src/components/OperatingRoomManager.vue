<script setup>
import { Search } from '@element-plus/icons-vue'
import {ref, reactive, onMounted} from 'vue'
import {get, post,put} from "@/net";
import {ElMessage, ElMessageBox} from "element-plus";
const multipleTableRef = ref()
const dialogVisible = ref(false)
const dialogCateVisible = ref(false)
const addDialogVisible = ref(false)
const addCategoryDialogVisible = ref(false)
const formRef = ref()
const tableData = ref([])
const data=reactive({
  category:[]
})
const rules = {
  name: [
    { required: true, message:'请输入手术室名', trigger: ['blur', 'change'] },
    { min: 1, max: 10, message: '职位名的长度必须在1-10个字符之间', trigger: ['blur', 'change'] },
  ],
  category: [
    { required: true, message:'请输入手术室类别', trigger: ['blur', 'change'] },
    { min: 1, max: 10, message: '职位名的长度必须在1-10个字符之间', trigger: ['blur', 'change'] },
  ],
}

/*分页参数定义*/
const pagingOn = ref(true)
const pageNum = ref(1)
const pages = ref(0)

const Cate=reactive({
  id:'',
  name:''
})

const clearCate = () => {
  Cate.id=''
  Cate.name=''
}

const Room=reactive({
  id:'',
  name:'',
  category:''
})

const clearRoom = () => {
  Room.id=''
  Room.name=''
  Room.category=''
}


/*页面打开自动查询职位信息*/
onMounted(() => {
  get('/api/admin/category',(message)=>{
    data.category=message
  },(message)=>{
    ElMessage.error(message)
  })
  post('/api/admin/query/operatingRoomList', {},(message) => {
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
  console.log(pageNum.value)
  post('/api/admin/query/operatingRoomList', {
    pageNum:pageNum.value,
    name:searchInfo.searchName,
    category:searchInfo.searchCategory

  },(message) => {
    const map = new Map(Object.entries(message))
    tableData.value = map.get('pageList')
    pages.value = map.get('totalPage')
  }, (message) => {
    ElMessage.error(message)
  })
}

function modifyCate(row){
  Cate.id=row.id
  Cate.name=row.name
  dialogCateVisible.value = true
}

function modify(row){
  Room.id=row.id
  Room.name=row.name
  Room.category=row.category
  dialogVisible.value = true
}

const getCategoryName = (categoryId) => {
  const cate =  data.category.find(cate => cate.id === categoryId)
  return cate ? cate.name : '暂无该类型手术' // 如果找不到对应的职务名称，则显示'未知'
}

function submitAddInfo(){
  formRef.value.validate((isValid) => {
    if(isValid) {
      put('/api/admin/operatingRoom',{
        name:Room.name,
        category:Room.category
      },(message) => {
        addDialogVisible.value = false
        clearRoom()
        ElMessage.success(message)
            post('/api/admin/query/operatingRoomList', {},(message) => {
              const map = new Map(Object.entries(message))
              tableData.value = map.get('pageList')
              pages.value = map.get('totalPage')
              pageNum.value = map.get('pageNum')
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

function submitAddCateInfo(){
  formRef.value.validate((isValid) => {
    if(isValid) {
      put('/api/admin/category',{
            name:Cate.name,
          },(message) => {
            addCategoryDialogVisible.value = false
            clearCate()
            ElMessage.success(message)
            get('/api/admin/category',(message) => {
              data.category=message
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
      post('/api/admin/operatingRoom',{
            id:Room.id,
            name:Room.name,
            category:Room.category
          },(message) => {
            clearRoom()
            ElMessage.success(message)
            post('/api/admin/query/operatingRoomList', {},(message) => {
              const map = new Map(Object.entries(message))
              tableData.value = map.get('pageList')
              pages.value = map.get('totalPage')
              pageNum.value = map.get('pageNum')
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

function submitCateModify(){
  dialogCateVisible.value = false
  formRef.value.validate((isValid) => {
    if(isValid) {
      post('/api/admin/category',{
            id:Cate.id,
            name:Cate.name,
          },(message) => {
            clearCate()
            ElMessage.success(message)
            get('/api/admin/category',(message) => {
              data.category=message
              console.log("cate:",category)
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

//封装删除方法公共代码
const delRoom = (row) => {
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
        post('/api/admin/delOperatingRoom', {
          id:row.id
        }, (message) => {
          post('/api/admin/query/operatingRoomList', {},(message) => {
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

const delCate = (row) => {
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
        post('/api/admin/delCategory', {
          id:row.id
        }, (message) => {
          get('/api/admin/category',(message) => {
            data.category=message
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
// 搜索参数列表
const searchInfo = reactive({
  searchName: '',
  searchCategory: ''
})

// 模糊查询方法
function search(){
  post('/api/admin/query/operatingRoomList',{
    name: searchInfo.searchName,
    category: searchInfo.searchCategory,
    pageNum: pageNum.value
  },(message) => {
    const map = new Map(Object.entries(message))
    tableData.value = map.get('pageList')
    pages.value = map.get('totalPage')
  },(message) => {
    ElMessage.error(message)
  })
}

</script>


<template>
<!--
  <div style="display: flex;margin-right: 50px; justify-content: flex-end;">
-->
  <el-row>
    <el-input style="width: 200px;margin-bottom: 10px;" v-model="searchInfo.searchName" placeholder="请输入手术室名" :prefix-icon="Search" />
    <el-form-item  style="width: 200px;margin-bottom: 10px;" >
      <el-select v-model="searchInfo.searchCategory" placeholder="请选择手术类别" prop="pid">
        <el-option
            v-for="item in  data.category"
            :key="item.id"
            :label="item.name"
            :value="item.id">
        </el-option>
      </el-select>
    </el-form-item>
    <el-button style="margin-bottom: 10px;" type="success" :icon="Search" @click="search">搜索</el-button>
    <el-button  style="margin-right: 10px;" type="success" @click="addDialogVisible = true">新增手术室房间</el-button>
    <el-button  style="margin-right: 10px;" type="success" @click="addCategoryDialogVisible = true">新增科室类别</el-button>
  </el-row>
<!--  </div>-->

  <el-container>
  <el-table
      ref="multipleTableRef"
      :data="tableData"
      style="width: 50%"
      border
  >

    <el-table-column type="selection" width="55" />
    <el-table-column label="序号" width="100" >
      <template #default="{ row, $index }">
        <span>{{ $index + 1 }}</span>
      </template>
    </el-table-column>
    <el-table-column prop="name" label="手术室名" width="200" sortable/>
    <el-table-column label="科室类别" width="250" sortable>
      <template #default="{ row }">
        {{ getCategoryName(row.category) }}
      </template>
    </el-table-column>
    <el-table-column label="操作" align="center" #default="scope">
      <el-space wrap>
      <div>
        <el-button size="small" @click="modify(scope.row)">修改</el-button>
        <el-button size="small" type="danger" @click="delRoom(scope.row)">删除</el-button>
      </div>
      </el-space>
    </el-table-column>
  </el-table>


<!--    类别-->
  <el-table
      ref="multipleTableRef"
      :data="data.category"
      style="width: 35%;margin-left: 100px"
      border
      marg
  >

    <el-table-column type="selection" width="55" />
    <el-table-column label="序号" width="100" >
      <template #default="{ row, $index }">
        <span>{{ $index + 1 }}</span>
      </template>
    </el-table-column>
    <el-table-column prop="name" label="科室类别" width="150" sortable/>
    <el-table-column label="操作" align="center" #default="scope">
      <el-space wrap >
        <div >
          <el-button size="small" @click="modifyCate(scope.row)">修改</el-button>
          <el-button size="small" type="danger" @click="delCate(scope.row)">删除</el-button>
        </div>
      </el-space>
    </el-table-column>
  </el-table>
    </el-container>


  <!--  分页代码  -->
  <div style="width: 50%;margin-top: 20px;display: flex;justify-content: center;" v-if="pagingOn">
    <el-pagination background layout="prev, pager, next" :page-count="pages" :current-page="pageNum"
                   @current-change="nextPage"/>
  </div>

  <el-dialog v-model="addCategoryDialogVisible" title="新增科室类型信息" width="30%" align-center center>
    <div style="display: flex;justify-content: center;">
      <el-form :model="Cate" label-width="120px" :rules="rules.category" ref="formRef">
        <el-form-item label="科室类别" prop="name">
          <el-input v-model="Cate.name" type="text"/>
        </el-form-item>
      </el-form>
    </div>
    <template #footer>
            <span class="dialog-footer">
                <el-button @click="addCategoryDialogVisible = false">取消</el-button>
                <el-button type="primary" @click="submitAddCateInfo">
                    保存
                </el-button>
            </span>
    </template>
  </el-dialog>

  <el-dialog v-model="addDialogVisible" title="新增手术室信息" width="30%" align-center center>
    <div style="display: flex;justify-content: center;">
      <el-form :model="Room" label-width="120px" :rules="rules" ref="formRef">
        <el-form-item label="手术室名" prop="name">
          <el-input v-model="Room.name" type="text"/>
        </el-form-item>

        <el-form-item label="科室类别">
          <el-select v-model="Room.category" placeholder="请选择科室类别" style="width: 300px;" prop="pid">
            <el-option
                v-for="item in  data.category"
                :key="item.id"
                :label="item.name"
                :value="item.id">
            </el-option>
          </el-select>
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



  <el-dialog v-model="dialogVisible" title="修改手术室信息" width="30%" align-center center>
    <div style="display: flex;justify-content: center;">
      <el-form :model="Room" label-width="120px" :rules="rules" ref="formRef">
        <el-form-item label="手术室名" prop="name">
          <el-input v-model="Room.name" type="text"/>
        </el-form-item>

        <el-form-item label="科室类别" prop="name">
          <el-input v-model="Room.category" type="text"/>
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

  <el-dialog v-model="dialogCateVisible" title="修改科室类别信息" width="30%" align-center center>
    <div style="display: flex;justify-content: center;">
      <el-form :model="Room" label-width="120px" :rules="rules.category" ref="formRef">
        <el-form-item label="类别" prop="name">
          <el-input v-model="Cate.name" type="text"/>
        </el-form-item>

      </el-form>
    </div>
    <template #footer>
            <span class="dialog-footer">
                <el-button @click="addDialogVisible = false">取消</el-button>
                <el-button type="primary" @click="submitCateModify">
                    保存
                </el-button>
            </span>
    </template>
  </el-dialog>

</template>
<style scoped></style>