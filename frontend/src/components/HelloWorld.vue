<script setup>
import {reactive, ref} from 'vue'
import {Greet} from '../../wailsjs/go/main/App'
import { ElLoading, ElMessageBox, ElMessage } from 'element-plus'

const data = reactive({
  name: "",
  result: [],
  cententLength: '',
  status: '',
  code: ''})
let show = ref(false)

function greet() {
  data.name = data.name.trim()
  console.log(data.name)
  const loading = ElLoading.service({
    lock: true,
    text: '搜索中...',
    background: 'rgba(0, 0, 0, 0.7)',
  })
  Greet(data.name).then(res => {
    console.log(JSON.parse(res))
    res = JSON.parse(res)
    loading.close()
    show.value = false
    if(res.code == 0){
      data.result = res.data
      data.cententLength = res.centent_length
      data.status = res.status
      data.code = res.raw_code
    }else{
      data.result = []
      data.cententLength = ''
      data.status = ''
      data.code = ''
      ElMessage.error(res.msg);
    }

    
  })
}

function copyContent(){
  let text = data.result.join('\n')
  navigator.clipboard
    .writeText(text)
    .then(() => {
      ElMessage.success("复制成功");
    })
    .catch((err) => {
      ElMessage.error("复制失败，请重试！");
    });
}

</script>

<template>
  <div class="page">
    <div class="header"></div>
    <div class="inner">
      <img class="bg" src="../assets/images/bg.png"/>
    </div>
    <div class="search-mode">
      <div class="page-title">
        <img class="icon" src="../assets/images/search.png"/>
        二级子域名获取
      </div>
      <div id="input" class="input-box">
        <input id="name" placeholder="请输入正确的URL,以http或https开头" v-model="data.name"  @keydown.enter.exact.prevent="greet" autocomplete="off" class="input" type="text"/>
        <button class="btn" @click="greet">获取</button>
      </div>
      
      <div id="result" class="result">
        <p class="title">
          获取结果
          
        </p>

        <el-descriptions v-if="data.status">
          <el-descriptions-item label="状态码：">{{data.status}}</el-descriptions-item>
          <el-descriptions-item label="大小：">{{data.cententLength}}KB</el-descriptions-item>
          <el-descriptions-item label="源码："><el-tag size="small" @click="show = true">查看源码</el-tag></el-descriptions-item>
        </el-descriptions>
        
        <div class="name-mode">
          <p class="name-title" v-if="data.status">以下是获取到的二级子域名列表
            <button class="copy" @click="copyContent" >复制全部</button>
          </p>
          <span class="name" v-for="(item, index) in data.result" :key="index">{{item}}
            <text v-if="index == 0">（当前访问域名）</text>
          </span>
          <p class="null" v-if="!data.result.length">
            <img class="icon" src="../assets/images/null.png"/>
            暂无数据
          </p>
        </div>
      </div>
    </div>
    
    <el-dialog v-model="show" title="源码详情" width="80vh" align-center>
      <p class="code">{{data.code}}</p>
    </el-dialog>
    
  </div>
  <!-- <p class="message">操作成功</p> -->
</template>

<style>
.name-title{
  font-size: 14px;
  color: #666;
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.el-overlay-dialog{
  overflow: hidden;
}
.el-descriptions__cell{
  text-align: center;
}
.code{
  height: 70vh;
  overflow-y: scroll;
  line-height: 30px;
  white-space: pre-line;
  overflow-x: hidden;
}
.el-tag{
  cursor: pointer;
}
.null{
  font-size: 16px;
}
.null .icon{
  width: 230px;
  height: auto;
  display: block;
  margin: 50px auto;
}
.page{
  position: relative;
}
.page-title{
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 30px 0 20px;
  font-size: 42px;
  color: #fff;
  font-weight: bold;
}
.page-title .icon{
  width: 64px;
  height: 64px;
  margin-right: 8px;
}
.header{
  width: 100%;
  height: 230px;
  background: linear-gradient(54deg,#68a9ff 16%, #026bf7 97%);
}
.inner{
  width: 100%;
  height: calc(100vh - 230px);
}
.bg{
  width: 100%;
  height: 100%;
  display: block;
  object-fit: cover;
}
.search-mode{
  width: 100%;
  position: absolute;
  top: 0;
  left: 0;
}
.search-mode ::-webkit-scrollbar {
  width: 5px !important;
  height: calc(96vh - 350px);
}
.search-mode ::-webkit-scrollbar-thumb {
  width: 4px !important;
  border-radius: 2px !important;
  background-clip: content-box !important;
  background: rgba(0, 0, 0, 0.2);
}
.message{
  position: fixed;
  top: 20px;
  left: 50%;
  transform: translateX(-50%);
  width: 180px;
  height: 60px;
}
.result {
  width: 60vw;
  height: calc(100vh - 255px);
  overflow-y: scroll;
  margin: 20px auto;
  background: #fff;
  border-radius: 10px;
  padding: 0 20px 20px;
  box-shadow: 0px 3px 6px 0px rgba(0,0,0,0.15); 
}
.copy{
  width: 100px;
}
.title{
  height: 50px;
  font-size: 18px;
  color: #333;
  display: flex;
  align-items: center;
  justify-content: flex-start;
  margin: 0 auto 20px;
}
.name-mode{
  margin-top: 20px;
}
.name{
  display: block;
  font-size: 14px;
  line-height: 30px;
  padding: 0;
  margin: 0 20px 0 0;
}
.copy{
  width: 70px;
  height: 30px;
  line-height: 30px;
  border-radius: 10px;
  border: none;
  background: #CDE3FF;
  border: none;
  color: #0E79FF;
  font-size: 12px;
  cursor: pointer;
  margin-left: 10px;
}
.input-box{
  width: 60vw;
  height: 50px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 20px;
  border: 2px solid #0e79ff;
  border-radius: 25px;
  box-shadow: 0px 3px 6px 0px rgba(0,0,0,0.15); 
  background: #fff;
}


.input-box .btn {
  width: 80px;
  height: 30px;
  line-height: 30px;
  padding: 0;
  border-radius: 0 20px 20px 0;
  border: none;
  cursor: pointer;
  background: #fff;
  border-left: 1px solid #0E79FF;
}



.input-box .input {
  width: calc(100% - 80px);
  border: none;
  border-radius: 3px;
  outline: none;
  height: 50px;
  line-height: 50px;
  padding: 0 20px;
  box-sizing: border-box;
  background-color: #fff;
  -webkit-font-smoothing: antialiased;
  font-size: 16px;
  border-radius: 25px 0 0 25px;
}
</style>
