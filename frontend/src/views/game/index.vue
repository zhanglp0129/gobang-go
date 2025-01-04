<script setup lang="ts">
import Board from './board/index.vue'
import router from "../../router";
import {ref} from "vue";

// 当前棋盘，-1表示白棋，1表示黑棋，0表示空位
const boards = ref([])
for (let i=0;i<15;i++) {
  let temp = []
  for (let j=0;j<15;j++) {
    temp.push(0)
  }
  boards.value.push(temp)
}
// 当前应该下棋的一方 1黑棋 -1白棋 0禁止下棋
const cur = ref(1)

// 返回到选择界面
const goSelect = () => {
  router.push('/select')
}
// 重新游戏
const reset = () => {
  location.reload()
}
</script>

<template>
  <div class="game">
    <div class="title">
      <span>正在游戏</span>
    </div>
    <div class="middle">
      <div class="left">
        <Board v-model="boards" v-model:cur="cur"></Board>
      </div>
      <div class="right">
        <el-button @click="reset">重新游戏</el-button>
        <el-button @click="goSelect">返回</el-button>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
.game {
  .title {
    color: #409eff;
    margin-top: 20px;
    font-size: 42px;
    font-weight: 500;
    text-align: center;
    user-select: none;
  }
  .middle {
    display: flex;
    margin-top: 20px;
    justify-content: space-between;
    .left {
      margin-left: 20px;
    }
    .right {
      margin-right: 20px;
    }
  }
}
</style>