<script setup lang="ts">
import Board from './board/index.vue'
import router from "../../router";
import {computed, ref} from "vue";
import {useRoute} from "vue-router";

const route = useRoute()
const first = Number(route.query.first)
const back = Number(route.query.back)
const selects = ['人类', '简单AI', '中等AI', '困难AI']

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
const isAI = computed(() => {
  return (cur.value === 1 && first !== 0) || (cur.value === -1 && back !== 0)
})

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
      <span v-if="cur===0" :style="{color: 'red'}">游戏结束</span>
      <span v-else-if="isAI">AI正在思考</span>
      <span v-else>请下棋</span>
    </div>
    <div class="middle">
      <div class="left">
        <Board v-model="boards" v-model:cur="cur"></Board>
      </div>
      <div class="right">
        <div class="top">
          <div v-if="cur !== 0">当前：{{cur==1?'黑棋':'白棋'}} ({{selects[cur==1?first:back]}})</div>
        </div>
        <div class="button">
          <el-button @click="reset">重新游戏</el-button>
          <el-button @click="goSelect">返回</el-button>
        </div>
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
    .left {
      margin-left: 20px;
    }
    .right {
      margin-left: 50px;
      display: flex;
      flex-direction: column;
      justify-content: space-between;
      .top {
        div {
          margin-bottom: 20px;
        }
      }
    }
  }
}
</style>