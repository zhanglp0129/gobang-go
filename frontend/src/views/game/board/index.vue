<script setup lang="ts">
import {computed, onMounted, reactive, ref} from "vue";
import {ElMessage} from "element-plus";

onMounted(() => {
  initBoard()
})
let ctx: CanvasRenderingContext2D
const initBoard = () => {
  // 绘制棋盘
  ctx = document.getElementById('canvas').getContext('2d')
  // 总宽高为650*650，一共14个格子，每个格子40px，剩余90。边缘45px
  // 画横线，A-O
  ctx.font = '20px Arial'
  for (let i=0;i<15;i++) {
    ctx.fillText(String.fromCharCode(65+i), 15, 45+40*i+10)
    ctx.moveTo(45, 45+40*i)
    ctx.lineTo(650-45, 45+40*i)
    ctx.stroke()
  }
  // 画竖线，1-15
  for (let i=0;i<15;i++) {
    ctx.fillText(i+1+'', 45+40*i-10, 30)
    ctx.moveTo(45+40*i, 45)
    ctx.lineTo(45+40*i, 650-45)
    ctx.stroke()
  }
  // 画圆点
  function drawDot(x: number, y: number) {
    ctx.beginPath()
    ctx.arc(45+40*(x-1), 45+40*(y-1), 4, 0, 2*Math.PI)
    ctx.fill()
  }
  drawDot(4, 4)
  drawDot(4, 12)
  drawDot(12, 4)
  drawDot(12, 12)
  drawDot(8, 8)
}

// 获取棋盘数据
const boards = defineModel()
const cur = defineModel('cur')

// 根据鼠标位置，判断可以落子
const curMouse = reactive([-1, -1])
const canDown = computed(() => {
  if (cur.value === 0) {
    return false
  }
  let x=-1, y=-1
  for (let i=0;i<15;i++) {
    if (curMouse[0] >= 45+40*i-15 && curMouse[0] <= 45+40*i+15) {
      x = i
    }
    if (curMouse[1] >= 45+40*i-15 && curMouse[1] <= 45+40*i+15) {
      y = i
    }
  }
  return x!=-1 && y!=-1 && boards.value[x][y]==0
})

// 游戏结束
const ways = [[0, 1], [1, 0], [1, 1], [1, -1]]
const gameOver = (x, y): boolean => {
  for (let way of ways) {
    let dx=way[0], dy=way[1]
    let cnt = [[x, y]]
    for (let i=1;true;i++) {
      let cx=x+i*dx, cy=y+i*dy
      if (cx<0 || cx>=15 || cy<0 || cy>=15 || boards.value[cx][cy]!==cur.value) {
        break
      }
      cnt.push([cx, cy])
    }
    dx *= -1
    dy *= -1
    for (let i=1;true;i++) {
      let cx=x+i*dx, cy=y+i*dy
      if (cx<0 || cx>=15 || cy<0 || cy>=15 || boards.value[cx][cy]!==cur.value) {
        break
      }
      cnt.push([cx, cy])
    }
    if (cnt.length >= 5) {
      // 给五子加上边框
      for (let c of cnt) {
        ctx.beginPath()
        ctx.strokeStyle = 'red'
        ctx.lineWidth = 3
        ctx.arc(c[0]*40+45, c[1]*40+45, 15, 0, 2*Math.PI)
        ctx.stroke()
      }
      return true
    }
  }
  return false
}

// 落子
const down = (e) => {
  if (cur.value === 0) {
    return
  }
  // 获取当前落子位置，每个棋子半径为15
  let x=-1, y=-1
  for (let i=0;i<15;i++) {
    if (e.offsetX >= 45+40*i-15 && e.offsetX <= 45+40*i+15) {
      x = i
    }
    if (e.offsetY >= 45+40*i-15 && e.offsetY <= 45+40*i+15) {
      y = i
    }
  }
  if (x==-1 || y==-1 || boards.value[x][y]!=0) {
    return
  }
  boards.value[x][y] = cur.value
  cur.value *= -1
  if (gameOver(x, y)) {
    ElMessage.success(`游戏结束，${cur.value === 1?'黑棋':'白棋'}胜利`)
    cur.value = 0
  }
}

</script>

<template>
  <div class="board">
    <canvas id="canvas" width="650" height="650" @click="down" :style="{
      cursor: canDown?'pointer':'default'
    }" @mousemove="curMouse[0]=$event.offsetX;curMouse[1]=$event.offsetY"/>
    <img src="../../../assets/images/black-chess.png" alt="" v-for="i in 225"
         :key="i" v-show="boards[Math.floor((i-1)/15)][(i-1)%15] === 1" :style="{
           left: 45+40*Math.floor((i-1)/15)-15+'px',
           top: 45+40*((i-1)%15)-15+'px'
         }">
    <img src="../../../assets/images/white-chess.png" alt="" v-for="i in 225"
         :key="i" v-show="boards[Math.floor((i-1)/15)][(i-1)%15] === -1" :style="{
           left: 45+40*Math.floor((i-1)/15)-15+'px',
           top: 45+40*((i-1)%15)-15+'px'
         }">
  </div>
</template>

<style scoped lang="scss">
.board {
  position: relative;
  #canvas {
    background-color: rgb(227, 189, 10);
  }
  img {
    position: absolute;
  }
}
</style>