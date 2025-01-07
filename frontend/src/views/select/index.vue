<script setup lang="ts">
import {reactive} from "vue";
import router from "../../router";

// 先手后手的选择
const handSelect = reactive({
  first: 0,
  back: 0,
  model: 1
})
// 每种选择对应的值
const selects = ['人类', '简单AI', '中等AI', '困难AI']

// 开始游戏
const start = () => {
  router.push({
    path: '/game',
    query: handSelect
  })
}
// 重置
const reset = () => {
  Object.assign(handSelect, {
    first: 0,
    back: 0
  })
}

</script>

<template>
  <div class="select">
    <div class="title">
      <span>模式选择</span>
    </div>
    <div class="middle">
      <el-form class="form" label-width="80px">
        <el-form-item label="先手">
          <el-radio-group v-model="handSelect.first">
            <el-radio v-for="(item, i) in selects" :label="item"
                      :value="i"
            ></el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="后手">
          <el-radio-group v-model="handSelect.back">
            <el-radio v-for="(item, i) in selects" :label="item"
                      :value="i"
            ></el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="AI模型">
          <el-radio-group v-model="handSelect.model">
            <el-radio label="传统模型" :value="1"></el-radio>
            <el-radio label="卷积模型" :value="2" disabled></el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="start">开始</el-button>
          <el-button @click="reset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<style scoped lang="scss">
.select {
  .title {
    color: #409eff;
    margin-top: 20px;
    font-size: 42px;
    font-weight: 500;
    text-align: center;
  }
  .middle {
    display: flex;
    justify-content: center;
    margin-top: 200px;
    .form {
      width: 450px;
    }
  }
}
</style>