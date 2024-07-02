<template>
  <el-row>
    <el-input style="width: 200px;margin-bottom: 10px;" v-model="searchInfo.searchName" placeholder="请输入手术室名" :prefix-icon="Search" />
    <el-input style="width: 200px;margin-bottom: 10px;" v-model="searchInfo.searchUserName" placeholder="请输入医生名称" :prefix-icon="Search" />
    <el-form-item  style="width: 200px;margin-bottom: 10px;" >
      <el-select v-model="searchInfo.searchCategoryId" placeholder="请选择科室" prop="pid">
        <el-option
            v-for="item in  data.categorys"
            :key="item.id"
            :label="item.name"
            :value="item.id">
        </el-option>
      </el-select>
    </el-form-item>
    <el-form-item  style="width: 200px;margin-bottom: 10px;" >
      <el-select v-model="searchInfo.searchRoomId" placeholder="请选择手术房间" prop="pid">
        <el-option
            v-for="item in  data.rooms"
            :key="item.id"
            :label="item.name"
            :value="item.id">
        </el-option>
      </el-select>
    </el-form-item>
    <el-button style="margin-bottom: 10px;" type="success" :icon="Search" @click="search">搜索</el-button>
    <el-button  style="margin-right: 10px;" type="success" @click="addDialogVisible = true">预约手术</el-button>
  </el-row>

  <el-table
      ref="multipleTableRef"
      :data="tableData"
      style="width: 100%"
      border
      :default-sort="{ prop: 'date', order: 'ascending' }"
  >
    <el-table-column type="selection" width="55" />
    <el-table-column label="序号" width="100" >
      <template #default="{ row, $index }">
        <span>{{ $index + 1 }}</span>
      </template>
    </el-table-column>
    <el-table-column property="name" label="手术名" sortable />
    <el-table-column label="科室" sortable>
      <template #default="{ row }">
        {{ getCategoryName(row.category_id) }}
      </template>
    </el-table-column>
    <el-table-column label="手术室" width="100" sortable>
      <template #default="{ row }">
        {{ getRoomName(row.room_id) }}
      </template>
    </el-table-column>
    <el-table-column label="人员安排" width="200" >
      <template #default="{ row }">
        <el-scrollbar height="50px">
          <p v-for="user in row.users" :key="user">
            {{ user.name+' '+getPostName(user.post_id) }}</p>
        </el-scrollbar>
      </template>
    </el-table-column>
    <el-table-column  prop="date" label="手术时间" sortable width="400" >
      <template #default="{ row }">
       {{parseTimeString(row.start_time) }}至{{parseTimeString(row.end_time)}}
      </template>
    </el-table-column>
    <el-table-column prop="state" label="进度" sortable width="100" >
    <template #default="{ row }">
      {{ getStateName(row.state) }}
    </template>
    </el-table-column>
    <el-table-column label="操作" align="center" #default="scope">
      <el-space wrap>
        <div>
          <el-button size="small" type="success" @click="finish(scope.row)">完成手术</el-button>
          <el-button size="small" type="danger" @click="delSurgery(scope.row)">取消预约</el-button>
        </div>
      </el-space>
    </el-table-column>
  </el-table>
  <div style="width: 100%;margin-top: 20px;display: flex;justify-content: center;" v-if="pagingOn">
    <el-pagination background layout="prev, pager, next" :page-count="pages" :current-page="pageNum"
                   @current-change="nextPage"/>
  </div>


<!--  预约手术代码-->
  <el-dialog v-model="addDialogVisible" title="预约手术信息"  align-center center>
    <div style="display: flex;justify-content: center;">
      <el-form :inline="true" :model="Scheduling" ref="formRef" class="form-inline" label-width="120px" style="max-width: 800px">
        <el-form-item label="手术名">
          <el-input v-model="Scheduling.name"  placeholder="请输入手术名"/>
        </el-form-item>
        <el-form-item label="科室" >
          <el-select v-model="Scheduling.category_id" placeholder="请选择科室"  >
            <el-option
                v-for="item in data.categorys"
                :key="item.id"
                :label="item.name"
                :value="item.id">
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="手术时长"  >
          <el-input v-model="Scheduling.duration" placeholder="请输入时长:(分)"/>
        </el-form-item>
        <el-form-item  label=" " label-width="130">
          <el-button type="primary" @click="queryTime()" >
            获取推荐时间
          </el-button>
        </el-form-item>

        <el-form-item label="选择时间" >
          <el-time-select
              v-model="Scheduling.start_time"
              style="width: 150px"
              :max-time="Scheduling.end_time"
              placeholder="开始时间"
              :start="data.recommendTime"
              step="00:15"
              end="23:00"
              format="HH:mm"
          />
          <el-time-select
              v-model="Scheduling.end_time"
              style="width: 150px"
              :min-time="Scheduling.start_time"
              placeholder="结束时间"
              :start="Scheduling.start_time"
              step="00:15"
              end="23:00"
              format="HH:mm"
          />
        </el-form-item>
        <el-form-item >
          <el-button type="primary" @click="subTime()" class="alignButton" >
            查询手术室和医护人员
          </el-button>
        </el-form-item>
        <el-form-item label="手术室">
          <el-select v-model="Scheduling.room_id" placeholder="请选择手术室" style="flex: 1;"  >
            <el-option
                v-for="item in data.avalableRooms"
                :key="item.id"
                :label="item.name"
                :value="item.id">
            </el-option>
          </el-select>

        </el-form-item>
        <el-form-item label="医护人员">

          <el-checkbox-group v-model="Scheduling.users" style="display: flex; flex-wrap: wrap;" >
            <el-checkbox v-for="user in data.users" :key="user" :label="user" :value="user" style="flex: 0 0 20%;">
              {{ user.name }} {{getPostName(user.post_id)}}
            </el-checkbox>
          </el-checkbox-group>
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

</template>

<script setup>
import {onMounted, reactive, ref} from "vue";
import {get, jsonPost, post, put} from "@/net/index.js";
import {ElMessage, ElMessageBox} from "element-plus";
import {Search} from "@element-plus/icons-vue";
const formRef = ref()

const addDialogVisible = ref(false)

const Scheduling=reactive({
  id:0,
  name:'',
  category_id:'',
  start_time:'',
  end_time:'',
  room_id:'',
  duration:'',
  users:[]
})

const clearScheduling=()=>{
  Scheduling.id=0
  Scheduling.name=''
  Scheduling.category_id=''
  Scheduling.start_time=''
  Scheduling.end_time=''
  Scheduling.room_id=''
  Scheduling.duration=''
  Scheduling.users=[]
}

const tableData=ref([])
const data = reactive({
  recommendTime:'',
  categorys:[],
  jobs:[],
  rooms:[],
  users:[],
  avalableRooms:[]
})

/*分页参数定义*/
const pagingOn = ref(true)
const pageNum = ref(1)
const pages = ref(0)

// 搜索参数列表
const searchInfo = reactive({
  searchName: '',
  searchCategoryId: '',
  searchRoomId:'',
  searchUserName:''
})

const now = new Date();
let currentMinute = now.getMinutes();

// 将分钟变为 15 的倍数
currentMinute = Math.ceil(currentMinute / 15) * 15;

/*页面打开自动查询用户信息*/
onMounted(() => {
  data.recommendTime=`${now.getHours()}:${currentMinute}`
  get('/api/admin/category',(message)=>{
    data.categorys=message
  },(message)=>{
    ElMessage.error(message)
  })

  get('/api/admin/operatingRoom',(message)=>{
    data.rooms=message
  },(message)=>{
    ElMessage.error(message)
  })

  get('/api/post',(message)=>{
    data.jobs=message
  },(message)=>{
    ElMessage.error(message)
  })

  post('/api/admin/surgery/list',{
  },(message)=>{
        const map = new Map(Object.entries(message))
        tableData.value= map.get('pageList')
        pages.value = map.get('totalPage')
        pageNum.value = map.get('pageNum')
  },(message)=>{
    ElMessage.error(message)
      }
  )
})

/*下一页方法实现*/
function nextPage(val){
  pageNum.value = val
  post('/api/admin/surgery/list',{
        name: searchInfo.searchName,
        category_id: searchInfo.searchCategoryId,
        room_id:searchInfo.searchRoomId,
        username:searchInfo.searchUserName,
        pageNum:pageNum.value,
      },
      (message) => {
        const map = new Map(Object.entries(message))
        tableData.value = map.get('pageList')
        pages.value = map.get('totalPage')
      }, (message) => {
        ElMessage.error(message)
      })
}


function submitAddInfo(){
  jsonPost('/api/admin/surgery/reserve',
    Scheduling,(message)=>{
    addDialogVisible.value=false
    clearScheduling()
        data.recommendTime=`${now.getHours()}:${currentMinute}`
    ElMessage.success(message)
        post('/api/admin/surgery/list',{
            },(message)=>{
              const map = new Map(Object.entries(message))
              tableData.value= map.get('pageList')
              pages.value = map.get('totalPage')
              pageNum.value = map.get('pageNum')
            },(message)=>{
              ElMessage.error(message)
            }
        )
  },(message)=>{
    ElMessage.error(message)
  })

}

const getPostName = (postId) => {
  const post = data.jobs.find(post => post.id === postId)
  return post ? post.name : '暂无职务' // 如果找不到对应的职务名称，则显示'未知'
}

const getCategoryName = (categoryId) => {
  const category = data.categorys.find(category => category.id === categoryId)
  return category ? category.name : '暂无类别名' // 如果找不到对应的职务名称，则显示'未知'
}

const getRoomName = (roomId) => {
  const room = data.rooms.find(room => room.id === roomId)
  return room ? room.name : '暂无手术室名' // 如果找不到对应的职务名称，则显示'未知'
}

//获取进度名
const getStateName = (code) => {
  if (code===1){
    return "已完成"
  }else{
    return "未完成"
  }
}


function queryTime(){
  post('/api/admin/surgery/recommend',{
    category_id:Scheduling.category_id,
    duration:Scheduling.duration
  },(message)=>{
    data.recommendTime= parseTimeHMString(message)
  },(message)=>{
    ElMessage.error(message)
  })
}

function parseTimeString(timeString) {
  var parsedDate = new Date(timeString);
  var year = parsedDate.getFullYear();
  var month = String(parsedDate.getMonth() + 1).padStart(2, '0');
  var day = String(parsedDate.getDate()).padStart(2, '0');
  var hours = String(parsedDate.getHours()).padStart(2, '0');
  var minutes = String(parsedDate.getMinutes()).padStart(2, '0');
  return year + "-" + month + "-" + day + " " + hours + ":" + minutes;
}

function parseTimeHMString(timeString) {
  var parsedDate = new Date(timeString);
  var hours = String(parsedDate.getHours()).padStart(2, '0');
  var minutes = String(parsedDate.getMinutes()).padStart(2, '0');
  minutes=Math.ceil(minutes / 15) * 15
  return  hours + ":" + minutes;
}

function subTime(){
  post('/api/admin/query/operatingRoomList',{
    category:Scheduling.category_id,
    start_time:Scheduling.start_time,
    end_time:Scheduling.end_time
  },(message)=>{
    if (message.pageList===null){
      ElMessage.error("时段忙碌，请重新选择手术时间")
    }
    data.avalableRooms = message.pageList

  },(message)=>{
    ElMessage.error(message)
  })

  post('/api/admin/surgery/users',{
    start_time:Scheduling.start_time,
    end_time:Scheduling.end_time
  },(message)=>{
    console.log("surgery users:",message)
    data.users = message
    data.users.sort((a, b) => {
      return a.post_id - b.post_id;
    });
  },(message)=>{
    ElMessage.error(message)
  })
}

const finish = (row) => {
  ElMessageBox.confirm(
      '手术已经完成了吗？',
      '提示',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'info',
      }
  )
      .then(() => {
        post('/api/admin/surgery/finish', {
          id:row.id,
          state:1
        }, (message) => {
         row.state=1
          ElMessage({
            type: 'success',
            message: '完成成功',
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

//封装删除方法公共代码
const delSurgery = (row) => {
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
        post('/api/admin/surgery/del', {
          id:row.id
        }, (message) => {
          post('/api/admin/surgery/list', {},(message) => {
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

// 模糊查询方法
function search(){
  console.log("info:",searchInfo)
  post('/api/admin/surgery/list',{
    name: searchInfo.searchName,
    category_id: searchInfo.searchCategoryId,
    room_id:searchInfo.searchRoomId,
    username:searchInfo.searchUserName,
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

<style>
.formLayout {
  display: flex;
  flex-direction: column;
}
.alignButton {
  margin-bottom: 10px;
  margin-left: 50px;
}

 .form-inline .el-input {
   --el-input-width: 220px;
 }

.form-inline .el-select {
  --el-select-width: 220px;
}

</style>