<script setup lang="ts">
import {computed, onMounted, reactive, ref, watch, watchEffect} from "vue";
import {ElMessage} from "element-plus";
import {useRoute} from "vue-router";
import {Predict} from "../../../../wailsjs/go/main/App";
import blackChess from '../../../assets/images/black-chess.png'
import whiteChess from '../../../assets/images/white-chess.png'

const route = useRoute()
const first = Number(route.query.first)
const back = Number(route.query.back)

onMounted(() => {
  initBoard()
})
let ctx: CanvasRenderingContext2D
const initBoard = () => {
  // 绘制棋盘
  ctx = document.getElementById('canvas').getContext('2d')
  // 总宽高为650*650，一共14个格子，每个格子40px，剩余90。边缘45px
  // 绘制背景
  ctx.fillStyle = 'rgb(227, 189, 10)'
  ctx.fillRect(0, 0, 650, 650)
  ctx.fillStyle = 'black'
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
const isAI = computed(() => {
  return (cur.value === 1 && first !== 0) || (cur.value === -1 && back !== 0)
})

let oldX = -1
let oldY = -1
watch(cur, async () => {
  // 判断是否为AI下棋
  let res
  if (cur.value === 1 && first !== 0) {
    res = await Predict(boards.value, cur.value, first)
  } else if (cur.value === -1 && back !== 0) {
    res = await Predict(boards.value, cur.value, back)
  } else {
    return
  }
  console.log(res)
  if (res[0]===-1 || res[1]===-1 || boards.value[res[0]][res[1]]!==0) {
    ElMessage.error('AI下棋失败')
    return
  }
  down(res[0], res[1])
}, {
  immediate: true
})

// 根据鼠标位置，判断可以落子
const curMouse = reactive([-1, -1])
const canDown = computed(() => {
  if (cur.value === 0 || isAI.value) {
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

// 鼠标落子
const mouseDown = (e) => {
  if (cur.value === 0 || isAI.value) {
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
  down(x, y)
}
// 下棋
const down = (x, y) => {
  boards.value[x][y] = cur.value
  // 绘制棋子
  let img = new Image()
  if (boards.value[x][y] == 1) {
    img.src = blackChess
  } else if (boards.value[x][y] == -1) {
    img.src = whiteChess
  }
  img.onload = () => {
    ctx.drawImage(img, x*40+45-15, y*40+45-15, 30, 30)
    // 绘制红色的点，显示当前下的棋子
    ctx.fillStyle = 'red'
    ctx.beginPath()
    ctx.arc(45+40*x, 45+40*y, 2, 0, 2*Math.PI)
    ctx.fill()
    // 擦除先前绘制的红点
    if (oldX !== -1 && oldY !== -1) {
      let oldImg = new Image()
      if (boards.value[oldX][oldY] == 1) {
        oldImg.src = blackChess
      } else if (boards.value[oldX][oldY] == -1) {
        oldImg.src = whiteChess
      }
      oldImg.onload = () => {
        ctx.drawImage(oldImg, oldX*40+45-15, oldY*40+45-15, 30, 30)
        oldX = x
        oldY = y
      }
    } else {
      oldX = x
      oldY = y
    }
  }
  cur.value *= -1
  if (gameOver(x, y)) {
    cur.value = 0
  }
}

</script>

<template>
  <div class="board">
    <canvas id="canvas" width="650" height="650" @click="mouseDown" :style="{
      cursor: canDown?'pointer':'default'
    }" @mousemove="curMouse[0]=$event.offsetX;curMouse[1]=$event.offsetY"/>
  </div>
</template>

<style scoped lang="scss">
.board {

}
</style>