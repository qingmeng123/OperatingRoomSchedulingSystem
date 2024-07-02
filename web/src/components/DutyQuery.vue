<template>
  <div>
    <div class="toolbox">
        <el-date-picker
            v-model="searchInfo.date"
            type="date"
            placeholder="选择日期"
            format="YYYY/MM/DD"
            value-format="YYYY-MM-DD"
        />
      <el-input style="width: 200px;" v-model="searchInfo.username" placeholder="请输入名字" :prefix-icon="Search" />
      <el-button style="" type="success" :icon="Search" @click="query">搜索</el-button>
    </div>
    <div class="gstc-wrapper"  ref="gstc"></div>
  </div>
</template>

<script>
import GSTC from "gantt-schedule-timeline-calendar/dist/gstc.wasm.esm.min.js";
import {Plugin as TimelinePointer} from "gantt-schedule-timeline-calendar/dist/plugins/timeline-pointer.esm.min.js";
import {Plugin as Selection} from "gantt-schedule-timeline-calendar/dist/plugins/selection.esm.min.js";
import {Plugin as ItemResizing} from "gantt-schedule-timeline-calendar/dist/plugins/item-resizing.esm.min.js";
import {Plugin as ItemMovement} from "gantt-schedule-timeline-calendar/dist/plugins/item-movement.esm.min.js";
import {Plugin as Bookmarks} from "gantt-schedule-timeline-calendar/dist/plugins/time-bookmarks.esm.min.js";

import "gantt-schedule-timeline-calendar/dist/style.css";
import {reactive, ref} from "vue";
import {get, post} from "@/net/index.js";
import {ElMessage} from "element-plus";
import {Search} from "@element-plus/icons-vue";

let gstc, state;

const colors = ['#E74C3C', '#DA3C78', '#7E349D', '#0077C0', '#07ABA0', '#0EAC51', '#F1892D'];
function getRandomColor(state) {
  if (state===0){
    return '#E74C3C'
  }
  return '#07ABA0'
}

const columnsFromDB = [
  {
    id: 'id',
    label: 'label',
    data: ({ row }) => Number(GSTC.api.sourceID(row.id)), // show original id

    sortable: 'label', // sort by id converted to number
    width: 80,
    header: {
      content: '序号',
    },
  },
  {
    id: 'label',
    /*data: 'label',*/
    data({ row, vido }) {
      return vido.html`<div @click=${() => onClick(row)}>${row.label}</div>`;
    },
    sortable:'label',
    width: 150,
    header: {
      content: '医护人员',
    },
  },
];

const hours = [
  {
    zoomTo: 100, // we want to display this format for all zoom levels until 100
    period: 'hour',
    periodIncrement: 10,
    format({ timeStart }) {
      return timeStart.format('DD MMMM YYYY');
      },
  },
];

const minutes = [
  {
    zoomTo: 100,
    period: 'minute',
    periodIncrement: 15,
    main: true,
    format({ timeStart, vido }) {
      return vido.html`<div style="text-align:center;">${timeStart.format('HH:mm')}</div>`; // full list of formats: https://day.js.org/docs/en/display/format
    },
  },
];

const today = new Date(); // 获取当前日期
const formattedDate = today.toISOString().split('T')[0]; // 将日期格式化为 'YYYY-MM-DD'

const data=reactive({
  rooms:[],
  jobs:[],
  items:[],
  rows:[],
  name:'',
  phone:'',
  job:''
})
const searchInfo=reactive({
  id:'',
  username:'',
  date:''
})

let infoDialogVisible = ref(false)

function onClick(row) {
  get('/api/user?id='+Number(GSTC.api.sourceID(row.id)),(message)=>{
    data.name=message.name
    data.job=getPostName(message.post_id)
    data.phone=message.phone

    alert("姓名: " + data.name+ "\n职位: " + data.job + "\n电话: " + data.phone);
  },(message)=>{
    alert(message)
  })

}

function getUserinfo(id){
  get('/api/user?id='+id,(message)=>{
    data.name=message.name
    data.phone=message.phone
  },(message)=>{
    ElMessage.error(message);
  })
}

function convertDateFormat(dateString) {
  const date = new Date(dateString); // 解析日期字符串为 Date 对象
  const year = date.getFullYear(); // 获取年份
  const month = ('0' + (date.getMonth() + 1)).slice(-2); // 获取月份并补零
  const day = ('0' + date.getDate()).slice(-2); // 获取日期并补零
  const hours = ('0' + date.getHours()).slice(-2); // 获取小时并补零
  const minutes = ('0' + date.getMinutes()).slice(-2); // 获取分钟并补零
  const seconds = ('0' + date.getSeconds()).slice(-2); // 获取秒数并补零

  // 返回所需格式的日期字符串
  return `${year}-${month}-${day}T${hours}:${minutes}:${seconds}`;
}

function getSurgeries() {
  post('/api/admin/surgery/list', {
    date:searchInfo.date,
    pageNum:-1
  }, (message) => {

    // 将消息数组中的每个元素转换为适合 item 对象的格式，并推送到 items 数组中
    if (message.pageList && message.pageList.length > 0) {
      data.items=[]
      data.rows=[]
      message.pageList.forEach(surgery => {
        surgery.users.forEach(user => {
          data.rows.push({
            id: user.id,
            label: user.name,
          })
          data.items.push({
            id: surgery.id + ' ' + user.id,
            label: getRoomName(surgery.room_id) + ' ' + surgery.name,
            rowId: user.id,
            style: { background: getRandomColor(surgery.state)},
            time: {
              start: GSTC.api.date(convertDateFormat(surgery.start_time)).valueOf(),
              end: GSTC.api.date(convertDateFormat(surgery.end_time)).valueOf(),

            }
          });
        })
        console.log("color:",surgery.state)
      });

    }
  }, (message) => {
    ElMessage.error(message);
  });

}

function queryByUsername(){
  post('/api/admin/surgery/list', {
    date:searchInfo.date,
    username:searchInfo.username,
    pageNum:-1
  },
      (message)=>{
        if (message.pageList && message.pageList.length > 0) {
          message.pageList.forEach(surgery => {
              data.rows[0]={
                id: surgery.users[0].id,
                label: surgery.users[0].name,
              }
              data.items.push({
                id: surgery.id + ' ' + surgery.users[0].id,
                label: getRoomName(surgery.room_id) + ' ' + surgery.name,
                rowId: surgery.users[0].id,
                style: { background: getRandomColor(surgery.state)},
                time: {
                  start: GSTC.api.date(convertDateFormat(surgery.start_time)).valueOf(),
                  end: GSTC.api.date(convertDateFormat(surgery.end_time)).valueOf(), // 假设消息中有结束时间的属性
                },
            })

          });
        }
      },(message) => {
        ElMessage.error(message);
      }
  )
}

function query(){
  if (searchInfo.username!==''){
    queryByUsername()
  }else{
    getSurgeries()
  }


  const config=newConfig()

  state.update("config.list.rows",config.list.rows)

  state.update("config.chart.items",config.chart.items)
  state.update("config.chart.time.from",config.chart.time.from)
  state.update("config.chart.time.to",config.chart.time.to)

}

// 时间戳：1637244864707
/* 时间戳转换为时间 */
function timestampToTime(timestamp) {
  timestamp = timestamp ? timestamp : null;
  let date = new Date(timestamp);//时间戳为10位需*1000，时间戳为13位的话不需乘1000
  let Y = date.getFullYear() + '-';
  let M = (date.getMonth() + 1 < 10 ? '0' + (date.getMonth() + 1) : date.getMonth() + 1) + '-';
  let D = (date.getDate() < 10 ? '0' + date.getDate() : date.getDate()) + ' ';
  let h = (date.getHours() < 10 ? '0' + date.getHours() : date.getHours()) + ':';
  let m = (date.getMinutes() < 10 ? '0' + date.getMinutes() : date.getMinutes()) + ':';
  let s = date.getSeconds() < 10 ? '0' + date.getSeconds() : date.getSeconds();
  return Y + M + D + h + m + s;
}


function clickAction(element, data) {
  function onClick(event) {
    // data variable will be updated in update method below so it will be always actual
    console.log("staye:",data.item)
    alert(`${data.item.label}\n${timestampToTime(data.item.time.start)}-${timestampToTime(data.item.time.end)}`);
  }

  element.addEventListener("click", onClick);

  return {
    update(element, newData) {
      data = newData; // data from parent scope updated
    },

    destroy(element, data) {
      element.removeEventListener("click", onClick);
    },
  };
}

function newConfig(){
  return {
    licenseKey:
        "====BEGIN LICENSE KEY====\nXOfH/lnVASM6et4Co473t9jPIvhmQ/l0X3Ewog30VudX6GVkOB0n3oDx42NtADJ8HjYrhfXKSNu5EMRb5KzCLvMt/pu7xugjbvpyI1glE7Ha6E5VZwRpb4AC8T1KBF67FKAgaI7YFeOtPFROSCKrW5la38jbE5fo+q2N6wAfEti8la2ie6/7U2V+SdJPqkm/mLY/JBHdvDHoUduwe4zgqBUYLTNUgX6aKdlhpZPuHfj2SMeB/tcTJfH48rN1mgGkNkAT9ovROwI7ReLrdlHrHmJ1UwZZnAfxAC3ftIjgTEHsd/f+JrjW6t+kL6Ef1tT1eQ2DPFLJlhluTD91AsZMUg==||U2FsdGVkX1/SWWqU9YmxtM0T6Nm5mClKwqTaoF9wgZd9rNw2xs4hnY8Ilv8DZtFyNt92xym3eB6WA605N5llLm0D68EQtU9ci1rTEDopZ1ODzcqtTVSoFEloNPFSfW6LTIC9+2LSVBeeHXoLEQiLYHWihHu10Xll3KsH9iBObDACDm1PT7IV4uWvNpNeuKJc\npY3C5SG+3sHRX1aeMnHlKLhaIsOdw2IexjvMqocVpfRpX4wnsabNA0VJ3k95zUPS3vTtSegeDhwbl6j+/FZcGk9i+gAy6LuetlKuARjPYn2LH5Be3Ah+ggSBPlxf3JW9rtWNdUoFByHTcFlhzlU9HnpnBUrgcVMhCQ7SAjN9h2NMGmCr10Rn4OE0WtelNqYVig7KmENaPvFT+k2I0cYZ4KWwxxsQNKbjEAxJxrzK4HkaczCvyQbzj4Ppxx/0q+Cns44OeyWcwYD/vSaJm4Kptwpr+L4y5BoSO/WeqhSUQQ85nvOhtE0pSH/ZXYo3pqjPdQRfNm6NFeBl2lwTmZUEuw==\n====END LICENSE KEY====",
    plugins: [TimelinePointer(), Selection(), ItemResizing(), ItemMovement(), Bookmarks()],
    list: {
      columns: {
        data: GSTC.api.fromArray(columnsFromDB),
      },
      rows: GSTC.api.fromArray(data.rows),
    },
    chart: {
      items: GSTC.api.fromArray(data.items),
      calendarLevels: [hours, minutes],
      time: {
        zoom: 14,
        period: 'minute',

        to: GSTC.api.date(searchInfo.date + 'T23:00:00').endOf('hour').valueOf(),
      },
    },
    actions: {
      "chart-timeline-items-row-item": [clickAction],
    },
  }
}

const getRoomName = (roomId) => {
  const room = data.rooms.find(room => room.id === roomId)
  return room ? room.name : '暂无手术室名' // 如果找不到对应的职务名称，则显示'未知'
}

const getPostName = (postId) => {
  const post = data.jobs.find(post => post.id === postId)
  return post ? post.name : '暂无职务' // 如果找不到对应的职务名称，则显示'未知'
}

export default {
  name: "GSTC",
  computed: {
    Search() {
      return Search
    }
  },
  mounted() {
    get('/api/post',(message)=>{
      data.jobs=message
    },(message)=>{
      ElMessage.error(message)
    })
    get('/api/admin/operatingRoom',(message)=>{
      data.rooms=message
    },(message)=>{
      ElMessage.error(message)
    })
    searchInfo.date=formattedDate
    getSurgeries()
    state = GSTC.api.stateFromConfig(newConfig());
    gstc = GSTC({
      element: this.$refs.gstc,
      state,
    });
  },

  beforeUnmount() {
    if (gstc) gstc.destroy();
  },

  methods: {
    query,
    getRoomName
  },

  data(){
    return {
      infoDialogVisible:infoDialogVisible,
      searchInfo: searchInfo,
      formattedDate:formattedDate
    };
  }
};

</script>
<style scoped>
.gstc-component {
  margin: 0;
  padding: 0;
}
.toolbox {
  padding: 10px;
}
</style>