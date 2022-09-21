<template>
  <div>
    <div v-for="(task) in this.$store.state.tasks"
         :key="task.id"
    >
      <div @click="isTaskInfo=true"
           style="display: flex;justify-content: space-between;margin-top: 20px;padding:5px;">
        <div>
          <div :style="task.deleted?'opacity: 0.5':''"
               style="font-size: 18px;text-transform: capitalize;display: flex;align-items: center;">{{ task.title }}
          </div>
          <span :style="task.deleted?'opacity: 0.5':''"
                style="color: #555;font-size: 15px;display: block;padding-top: 10px">{{
              timeLeft(new Date(task.timeEnd * 1000))
            }}
          </span>

        </div>
        <div @click.stop>
          <span v-if="task.deleted" style="color: #555;font-size: 12px;padding-right: 10px">{{
              task.waitSeconds
            }}с</span>
          <!--          <v-btn v-if="!event.deleted" @click="changeEvent(event.id)" icon>✏️</v-btn>-->
          <v-btn v-if="!task.deleted" @click="openChangeTask(task.id)" icon>✏️</v-btn>
          <v-btn v-if="task.deleted" @click="camebackTask(task.id)" text small>вернуть ↩︎</v-btn>
          <v-btn v-else @click="deleteTask(task.id)" icon>🗑</v-btn>
          <span v-if="!task.deleted"
                style="color: #555;font-size: 13px;display: block;text-align: right;padding-right: 5px">
              {{ calcProgress(task.id) }}

          </span>
        </div>
      </div>
      <div :style="{border: '2px solid green',opacity:'0.3'}"
      ></div>
      <div
          :style="{position:'relative',top:'-4px',border: '2px solid green',width:''+calcPercent(task)+'%',display:calcPercent(task)===0?'none':''}"
      ></div>
    </div>

    <template v-if="isAddNum">
      <div class="picker" @click="isAddNum=false">
        <div class="picker__content" @click.stop>
          <slot>
            <div style="background-color: white;padding: 20px;border-radius: 3px; width: 320px">
              <v-container grid-list-md text-xs-center>
                <v-layout row wrap>
                  <v-row>
                    <div style="display: flex;justify-content: space-between;padding:5px;color: #777777;width: 100%">
                      <div>
                        <div class="text-h5" style="max-width: 200px;word-wrap: break-word;">
                          {{ $store.state.tasks[findIndexTasks(chosenTaskId)].title }}
                        </div>

                      </div>
                      <div style="white-space: nowrap;">
                        {{ calcProgress(chosenTaskId) }}
                      </div>
                    </div>
                    <template v-if="$store.state.tasks[findIndexTasks(chosenTaskId)].type==='time'">
                      <v-flex xs6>
                        <v-text-field
                            label="hours"
                            outlined
                            v-model="taskHours"
                        ></v-text-field>
                      </v-flex>
                      <v-flex xs6>
                        <v-text-field
                            label="minutes"
                            outlined
                            v-model="taskMinutes"
                        ></v-text-field>
                      </v-flex>
                    </template>
                    <template v-else>
                      <v-flex xs12>
                        <v-text-field
                            :label="$store.state.tasks[findIndexTasks(chosenTaskId)].type"
                            outlined
                            v-model="taskNum"
                        ></v-text-field>
                      </v-flex>
                    </template>
                    <v-btn class="mb-2" width="290px"
                           @click="changeTask"
                           height="50px"
                           color="white">add
                    </v-btn>
                  </v-row>
                </v-layout>
              </v-container>
            </div>
          </slot>
        </div>
      </div>
    </template>


    <template v-if="isTaskInfo">
      <div class="picker" @click="isTaskInfo=false">
        <div class="picker__content" @click.stop>
          <slot>
            <h1>Hello World</h1>
          </slot>
        </div>
      </div>
    </template>


    <Footer>

    </Footer>
  </div>
</template>

<script>
import Footer from "@/components/Footer";
import moment from "moment";

export default {
  name: "ListTasks",
  components: {Footer},
  data() {
    return {
      taskHours: 0,
      taskMinutes: 0,
      taskNum: 0,
      isAddNum: false,
      isTaskInfo: false,
      chosenTaskId: -1,
      taskPieces: [],
    }
  },
  methods: {
    async getTasks() {
      console.log('hi there too')
      let response = await this.$store.state.makeReq('api/getTasks', 'GET')

      let tasks = await response.json()
      console.log('tasls = ', tasks)
      if (tasks === null) {
        tasks = []
      }

      this.$store.state.tasks = tasks
      console.log('resp = ', this.$store.state.tasks)
      localStorage.setItem('tasks', JSON.stringify(this.$store.state.tasks))
      console.log(this.$store.state.tasks)
    },

    openChangeTask(taskId) {
      this.isAddNum = true
      this.chosenTaskId = taskId
      this.getTaskPieces()
      console.log('hi')
    },
    async getTaskPieces(){
      let response = await this.$store.state.makeReq("api/getTaskPieces", "POST", JSON.stringify({id:this.chosenTaskId}))
      let tasksPieces = await response.json()
      console.log('tasls = ', tasksPieces)
      if (tasksPieces === null) {
        tasksPieces = []
      }

      this.taskPieces = tasksPieces
      console.log("pieces = ", tasksPieces)

    },
    async changeTask() {
      let index = this.findIndexTasks(this.chosenTaskId)
      let task = this.$store.state.tasks[index]
      let num = 0
      let deltaNum = 0
      console.log({task})
      if (parseInt(this.taskNum) > 0 || parseInt(this.taskHours) >= 0 && parseInt(this.taskMinutes) >= 0) {
        if (task.type === 'time') {
          console.log("hi there")
          console.log('calc = ', parseInt(this.$store.state.tasks[index].num) + parseInt((parseInt(this.taskHours) * 60 + parseInt(this.taskMinutes)) * 60))
          num = parseInt(this.$store.state.tasks[index].num) + parseInt((parseInt(this.taskHours) * 60 + parseInt(this.taskMinutes)) * 60)
          deltaNum = parseInt((parseInt(this.taskHours) * 60 + parseInt(this.taskMinutes)) * 60)
        } else {
          num = parseInt(this.$store.state.tasks[index].num) + parseInt(this.taskNum)
          deltaNum = parseInt(this.taskNum)
        }
        this.$store.state.tasks[index].num = num
      }
      localStorage.setItem('tasks', JSON.stringify(this.$store.state.tasks))
      console.log(task)
      let newTask = {
        id: task.id,
        title: task.title,
        timeStart: task.timeStart,
        timeEnd: task.timeEnd,
        type: task.type,
        target: task.target,
        num: num,
        deltaNum: deltaNum,
      }
      let response = await this.$store.state.makeReq('api/changeTask', 'PUT', JSON.stringify(newTask))
      console.log(response)
      this.isAddNum = false
      this.taskHours = 0
      this.taskMinutes = 0
      this.taskNum = 0
    },
    deleteTask(id) {
      let index = this.findIndexTasks(id)
      let task = this.$store.state.tasks[index]

      // console.log("index = ", index)
      // console.log({event})

      let waitSeconds = 3
      task.waitSeconds = waitSeconds
      let intervalId = setInterval(() => {
        if (!task.deleted) {
          clearInterval(intervalId)
          localStorage.setItem('tasks', JSON.stringify(this.$store.state.tasks))
          this.$forceUpdate()
          return
        }
        if (--waitSeconds < 0) {
          clearInterval(intervalId);
          let index = this.findIndexTasks(id)
          this.$store.state.makeReq("api/deleteTask", "DELETE", JSON.stringify({id: id}))
          this.$store.state.tasks.splice(index, 1)
        } else {
          task.waitSeconds = waitSeconds
        }
        this.$forceUpdate()
      }, 1000);

      // this.$store.state.events.splice(index, 1)

      task.deleted = Date.now();
      this.$forceUpdate();
      localStorage.setItem('events', JSON.stringify(this.$store.state.events))
    },
    camebackTask(id) {
      let index = this.findIndexTasks(id)
      this.$store.state.tasks[index].deleted = 0
      this.$forceUpdate();
    },
    timeLeft(date) {
      // console.log(date)
      // let arr = moment(date).fromNow().split(' ')
      // console.log(arr)
      // arr.splice(0,1)
      // console.log(arr)
      // return "осталось "+arr.join(" ")
      return moment(date).fromNow()
    },
    calcProgress(taskId) {
      let task = this.$store.state.tasks[this.findIndexTasks(taskId)]
      if (task.type === 'time') {
        return this.makeMinutesAndHoursStr(parseInt(parseInt(task.num) / 3600), parseInt(task.num) % 3600 / 60) + " / " + this.makeMinutesAndHoursStr(parseInt(parseInt(task.target) / 3600), parseInt(task.target) % 3600 / 60)
      } else {
        return task.num + ' / ' + task.target
      }
    },
    calcPercent(task) {
      let num = 0
      if (task.num !== undefined) {
        num = task.num
      }
      // console.log('percentege = ', num/task.target*100)
      return num / task.target * 100 > 100 ? 100 : num / task.target * 100
    },
    makeMinutesAndHoursStr(hours, minutes) {
      hours = hours.toString().length === 1 ? '0' + hours : hours
      minutes = minutes.toString().length === 1 ? '0' + minutes : minutes
      return hours + ':' + minutes
    },
    findIndexTasks(id) {
      for (let i = 0; i < this.$store.state.tasks.length; i++) {
        if (this.$store.state.tasks[i].id === id) {
          return i
        }
      }
      return -1
    },
  },
  beforeMount() {
    console.log('hi there')
    this.getTasks()
  }
}
</script>

<style scoped>
.picker {
  color: white;
  top: 0;
  bottom: 0;
  right: 0;
  left: 0;
  background-color: rgba(0, 0, 0, 0.5);
  position: fixed;
  z-index: 100;
  overflow: auto;
  display: flex;
  padding-top: 80px;
}

.picker__content {
  /*overflow-y: scroll;*/
  display: block;
  width: 290px;
  margin: auto;


  min-height: 100%;
  justify-content: space-between;
  z-index: 200;
}
</style>