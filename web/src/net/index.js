import { ElMessage } from "element-plus";
import axios from "axios";

const defaultError = () => ElMessage.error("出错了，请联系管理员！")
const defaultFailure = (message) => ElMessage.warning(message)

function post(url, data, success, failure, error){
    axios.post(url, data, {
        headers: {
            'Content-type': 'application/x-www-form-urlencoded'
        },
        withCredentials: true
    }).then(({data}) => {
        if(data.status){
            success(data.data)
        }else{
            failure(data.data)
        }
    }).catch(error)
}

function put(url, data, success, failure, error){
    axios.put(url, data, {
        headers: {
            'Content-type': 'application/x-www-form-urlencoded'
        },
        withCredentials: true
    }).then(({data}) => {
        if(data.status){
            success(data.data)
        }else{
            failure(data.data)
        }
    }).catch(error)
}


function get(url, success, failure, error){
    axios.get(url, {
        withCredentials: true
    }).then(({data}) => {
        if(data.status){
            success(data.data)
        }else{
            failure(data.data)
        }
    }).catch(error)
}

function jsonPost(url, data, success, failure, error){
    axios.request({
        url:url,
        method:"post",
        data: data,
        headers:{
            'content-type':"application/json"
        }
    }).then(({data}) => {
        if(data.status){
            success(data.data)
        }else{
            failure(data.data)
        }
    }).catch(error)
}

export {get, post,put, jsonPost}