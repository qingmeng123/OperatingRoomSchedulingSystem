<script setup>
import {ref} from "vue";

const time = ref('')

const dateTime = () => {
    let now = new Date()
    let year = now.getFullYear(); //获取完整年份
    let month = now.getMonth() + 1; //获取当前月份（0-11，0代表1月）
    let day = now.getDate(); //获取当前日（1-31）
    let hour = now.getHours();//获取当前小时数
    let minute = now.getMinutes();  //获取当前分钟数
    let second = now.getSeconds(); //获取当前秒数
    if (month < 10){
        month = '0' + month
    }
    if (day < 10){
        day = '0' + day
    }
    if (hour < 10){
        hour = '0' + hour
    }
    if (minute < 10){
        minute = '0' + minute
    }
    if (second < 10){
        second = '0' + second
    }
    time.value = year + '-' + month + '-' + day + ' ' + hour + ':' + minute + ':' + second
}
setInterval(() => {
    dateTime()
},1000)

</script>

<template>
    <div class="background-image">

            <div style=" font-weight: bold;font-size: 80px; color: white;">
                欢迎使用手术室排程系统</div>

    </div>
</template>


<style scoped>
.background-image {
  min-height: 85vh;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
  background-image: url("@/assets/bg.png");
  background-size: cover;
}

</style>