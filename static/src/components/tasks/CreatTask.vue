<template>
  <div>
    <h1>Creat your task:</h1>
    <template>
      <v-row>
        <v-col
            cols="12"
            sm="6"
        >
          <v-date-picker
              v-model="datesTask"
              range
              color="black lighten-1"
          ></v-date-picker>
        </v-col>
        <v-col
            cols="12"
            sm="6"
        >
          <v-text-field
              v-model="dateRangeText"
              label="Date range"
              style="background-color: white;padding: 20px 20px 0 20px;border-radius: 3px"
              prepend-icon="mdi-calendar"
              readonly
          ></v-text-field>
        </v-col>
      </v-row>
    </template>
    <v-text-field
        style="margin-top: 20px"
        label="😳 to do ..."
        solo
        v-model="taskTitle"
    ></v-text-field>
    <v-container grid-list-md text-xs-center>
      <v-layout row wrap>
        <v-row>
          <template v-if="typeOfTaskSelect==='time'">
            <v-flex xs4 style="margin-top: -10px">
              <v-text-field
                  label="hours"
                  solo
                  :rules="codeRules"
                  v-model="taskHours"
              ></v-text-field>
            </v-flex>
            <v-flex xs4 style="margin-top: -10px">
              <v-text-field
                  label="minutes"
                  solo
                  v-model="taskMinutes"
              ></v-text-field>
            </v-flex>
          </template>
          <template style="margin-top: -10px" v-if="typeOfTaskSelect==='my type'">
            <v-flex xs4 style="margin-top: -10px">
              <v-text-field
                  label="count"
                  solo
                  :rules="codeRules"
                  v-model="taskCount"
              ></v-text-field>
            </v-flex>
            <v-flex xs4 style="margin-top: -10px">
              <v-text-field
                  label="type"
                  solo
                  :rules="codeRules"
                  v-model="taskMyType"
              ></v-text-field>
            </v-flex>
          </template>
          <v-flex xs4>
            <v-select
                style="margin-top: -10px"
                v-model="typeOfTaskSelect"
                :items="items"
                label="time"
                solo
            ></v-select>
          </v-flex>
        </v-row>
      </v-layout>
    </v-container>
    <v-btn @click="creatTask" color="white" style="margin: -10px 0 50px;width: 100%;height: 40px">creat</v-btn>

  </div>
</template>

<script>


export default {
  name: "creatTask",
  data() {
    return {
      isNewTask: false,
      datesTask: ['2019-09-10', '2019-09-20'],
      taskTitle: "",
      items: ['time', 'my type'],
      typeOfTaskSelect: "time",
      taskHours: "",
      taskMinutes: "",
      taskCount: "",
      taskMyType: "",

      codeRules: [
        v => !!v || 'code is required',
        v => (v && v.length <= 20) || 'code must be 4 characters',
      ],
    }
  },
  computed: {
    dateRangeText() {
      return this.datesTask.join(' ~ ')
    },
  },
  methods: {
    creatTask() {
      let num = 0;
      let type = "";
      if (this.typeOfTaskSelect === "time") {
        if (isNaN(parseInt(this.taskHours))) {
          alert("hours should be integer")
          return
        }
        if (isNaN(parseInt(this.taskMinutes))) {
          alert("minutes should be integer")
          return;
        }
        num = (parseInt(this.taskHours) * 60 + parseInt(this.taskMinutes)) * 60
        type = "time"
      } else {
        if (isNaN(parseInt(this.taskCount))) {
          alert("count should be integer")
          return;
        }
        num = parseInt(this.taskCount)
        type = this.taskMyType
      }

      if (this.makeDate(this.datesTask[0]).getTime() / 1000 > this.makeDate(this.datesTask[1]).getTime() / 1000) {
        let tmp = this.datesTask[0]
        this.datesTask[0] = this.datesTask[1]
        this.datesTask[1] = tmp
      }


      let newTask = {
        timeStart: this.makeDate(this.datesTask[0], this.makeHoursAndMinutesStr(new Date())).getTime() / 1000,//seconds
        timeEnd: this.makeDate(this.datesTask[1], this.makeHoursAndMinutesStr(new Date())).getTime() / 1000,//seconds
        title: this.taskTitle,
        target: num,//if type==time - seconds
        num: 0,
        type: type,
      }
      this.$store.state.tasks.push(newTask)
      localStorage.setItem('tasks', JSON.stringify(this.$store.state.tasks))
      this.$store.state.makeReq('/api/newTask','POST',JSON.stringify(newTask))
      console.log(newTask)
      this.$emit("isNewTaskChange", false)
    },
    makeHoursAndMinutesStr(date) {
      let hours = date.getHours()
      let minutes = date.getMinutes()
      hours = hours.toString().length === 1 ? '0' + hours : hours
      minutes = minutes.toString().length === 1 ? '0' + minutes : minutes
      return hours + ':' + minutes
    },
    makeDate(date, time, seconds, milliseconds) {
      date = date === undefined ? "" : date
      time = time === undefined ? "" : "T" + time
      if (seconds !== undefined) {
        seconds = seconds.toString().length === 1 ? '0' + seconds : seconds
      }
      seconds = seconds === undefined ? "" : ":" + seconds
      milliseconds = milliseconds === undefined ? "" : "." + milliseconds
      console.log("make Date = ", date + time + seconds + milliseconds)

      return new Date(date + time + seconds + milliseconds)
    },
    getCurrentWeek() {
      // console.log(parseInt(new Date().get()))

      let day = (new Date()).getDay() - 1
      if (day === -1) {
        day = 6
      }

      let startDate = new Date((new Date()).getTime() - day * 1000 * 60 * 60 * 24)
      let endDate = new Date(startDate.getTime() + 1000 * 60 * 60 * 24 * 6)

      let startMonth = startDate.getUTCMonth() + 1
      let endMOnth = endDate.getUTCMonth() + 1

      startMonth.toString().length === 1 ? startMonth = "0" + startMonth : startMonth = "" + startMonth
      endMOnth.toString().length === 1 ? endMOnth = "0" + endMOnth : startMonth = "" + endMOnth

      let startDay = startDate.getUTCDate()
      let endDay = endDate.getUTCDate()

      startDay.toString().length === 1 ? startDay = "0" + startDay : startDay = "" + startDay
      endDay.toString().length === 1 ? endDay = "0" + endDay : endDay = "" + endDay

      let startDateStr = startDate.getFullYear() + "-" + startMonth + "-" + startDay
      let endDateStr = endDate.getFullYear() + "-" + endMOnth + "-" + endDay
      console.log(startDateStr)
      console.log(endDateStr)

      console.log(new Date(startDateStr))
      console.log(new Date(endDateStr))
      this.datesTask = [startDateStr, endDateStr]
    },
  },

  beforeMount() {
    this.getCurrentWeek()
  },
}
</script>

<style scoped>

</style>

