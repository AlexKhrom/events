<template>
  <div style="height: 100%;">
<!--    <div style="display:inline-block;margin: auto;width:100%">-->
<!--      <v-container>-->
<!--        <v-layout row wrap align-center>-->
    <div style="display: grid;grid-template-columns: repeat(2, 1fr);gap: 20px">
      <div v-for="(title,index) in $store.state.titles" :key="title.title">
        <v-btn @touchstart="deleteTitleDown('touch')" @mousedown="deleteTitleDown('click')"
               @touchend="deleteTitleUp(title,index,'touch')" @mouseup="deleteTitleUp(title,index,'click')"
               style="height: 100px;width: 100%">
          {{ title.title }}
        </v-btn>
      </div>
    </div>
<!--        </v-layout>-->
<!--      </v-container>-->
<!--    </div>-->
    <div style="margin-top:20px;display: flex;justify-content: center;text-align: center">
      <v-text-field
          style="margin:20px 20px 20px 0"
          label="🥴 новое событие"
          solo
          v-model="newEvent"
          @cancel="createTitle"
      ></v-text-field>

      <v-checkbox
          style="padding-top: 15px;min-width: 130px"

          v-model="isEventWithDuration"
          label="with duration"
      ></v-checkbox>
    </div>

    <v-btn @click="createTitle" color="primary" style="margin-top: -30px;width: 100%;height: 50px">add</v-btn>
    <v-btn @click="this.openNewTask" color="primary" style="margin-top:20px;width: 100%;height: 50px">new task
    </v-btn>

    <template v-if="isNewTask">
      <div class="picker" @click="isNewTask=false;">
        <div class="picker__content" @click.stop>
          <slot>
            <CreatTask @isNewTaskChange="changeIsNewTask">

            </CreatTask>
          </slot>
        </div>
      </div>
    </template>

    <template v-if="isPickEventDuration">
      <div class="picker" @click="isPickEventDuration=false">
        <div class="picker__content" @click.stop>
          <slot>
            <div>
              <h1>Plan your event:</h1>
              <v-row
                  justify="space-around"
                  align="center"
              >
                <v-col style="width: 350px; flex: 0 1 auto;">
                  <h2>Start:</h2>
                  <v-time-picker
                      format="24hr"
                      v-model="startEventTime"
                      :max="endEventTime"
                  ></v-time-picker>
                </v-col>
                <v-col style="width: 350px; flex: 0 1 auto;">
                  <h2>End:</h2>
                  <v-time-picker
                      format="24hr"
                      v-model="endEventTime"
                      :min="startEventTime"
                  ></v-time-picker>
                  <v-btn class="mb-2" width="290px"
                         @click="addEventWithDuration"
                         style="margin-top: 30px;"
                         height="50px"
                         color="white">save
                  </v-btn>
                </v-col>
<!--                <div style="display:block;height: 150px;width: 1px"></div>-->
              </v-row>
            </div>
          </slot>
        </div>
      </div>
    </template>


    <template>
      <v-row justify="center">
        <v-dialog
            v-model="dialog"
            persistent
            max-width="290"
        >
          <v-card>
            <v-card-title class="text-h5">
              Delete {{ titleToDelete.title }} ?
            </v-card-title>
            <!--            <v-card-text>Let Google help apps determine location. This means sending anonymous location data to Google,-->
            <!--              even when no apps are running.-->
            <!--            </v-card-text>-->
            <v-card-actions>
              <v-spacer></v-spacer>
              <v-btn
                  color="darken-1"
                  text
                  @click="dialog = false;"
              >
                No
              </v-btn>
              <v-btn
                  color="green darken-1"
                  text
                  @click="dialog = false;deleteTitleWithIndex(titleIndex)"
              >
                Yes
              </v-btn>
            </v-card-actions>
          </v-card>
        </v-dialog>
      </v-row>
    </template>

    <Footer >

    </Footer>
  </div>

</template>

<script>
import CreatTask from "@/components/tasks/CreatTask";
import Footer from "@/components/Footer";

export default {
  components: {Footer, CreatTask},

  name: "addEvent",
  data() {
    return {
      timePicker: "",
      selectedEventIndex: -1,
      newEventToAdd: {},
      newEvent: "",
      timerDeleteTitle: false,
      timerId: 0,
      dialog: false,
      titleToDelete: {},
      titleIndex: 0,
      titleToDeleteID: -1,
      isNewTask:false,

      startEventTime: null,
      endEventTime: null,
      isPickEventDuration: false,
      isEventWithDuration: false,
      datePicker: (new Date(Date.now() - (new Date()).getTimezoneOffset() * 60000)).toISOString().substr(0, 10),

    }
  },
  computed: {
    dateRangeText() {
      return this.datesTask.join(' ~ ')
    },
  },
  methods: {
    sortEvents() {
      this.$store.state.events.sort((a, b) =>
          a.dt < b.dt ? 1 : -1)
    },
    async makeRequest(event) {
      console.log('make request :', event)
      let response = await fetch('/api/newEvent', {
        method: 'POST',
        body: JSON.stringify(event)
      });

      if (response.ok) { // если HTTP-статус в диапазоне 200-299
        // получаем тело ответа (см. про этот метод ниже)
        console.log("okkkk")
      } else {
        console.log("not ok")
      }
    },
    addEvent(newTitle, index, param) {
      console.log("hi addevent. title -", newTitle.title, "\nclick on : ", param)

      if (newTitle.withDuration) {
        this.isPickEventDuration = true

        this.startEventTime = this.makeHoursAndMinutesStr(new Date())
        this.endEventTime = this.makeHoursAndMinutesStr(new Date())
        this.newEventToAdd = {
          title: newTitle.title,
          dt: Date.now(),
          withDuration: newTitle.withDuration,
        }
      } else {
        let event = {
          title: newTitle.title,
          dt: Date.now(),
          withDuration: newTitle.withDuration,
        }
        this.$store.state.events.push(event)

        this.makeRequest(event)

        this.sortEvents()
        localStorage.setItem('events', JSON.stringify(this.$store.state.events))

        this.page = 'list';
      }

      this.$store.state.titles.splice(index, 1)
      this.$store.state.titles.unshift(newTitle)
      localStorage.setItem('titles', JSON.stringify(this.$store.state.titles))
    },
    addEventWithDuration() {
      let date = new Date()
      date.setHours(0, 0, 0, 0)

      let nowMonth = date.getMonth().toString().length < 2 ? '0' + (date.getMonth() + 1) : (date.getMonth() + 1)
      let day = date.getDate().toString().length < 2 ? '0' + (date.getDate()) : (date.getDate())
      // console.log(date.getFullYear() + "-" + nowMonth + "-" + date.getDate() + "T" + this.startEventTime)
      // console.log((new Date(date.getFullYear() + "-" + nowMonth + "-" + date.getDate() + "T" + this.startEventTime)).getTime())
      console.log("make start date:", this.makeDate(date.getFullYear() + "-" + nowMonth + "-" + day, this.startEventTime, (new Date).getSeconds(), (new Date).getMilliseconds()).getTime(),)

      // console.log(new Date(date.getFullYear() + "-" + nowMonth + "-" + date.getDate() + " " + this.startEventTime))
      let event = {
        title: this.newEventToAdd.title,
        dt: this.makeDate(date.getFullYear() + "-" + nowMonth + "-" + day, this.startEventTime, (new Date).getSeconds(), (new Date).getMilliseconds()).getTime(),
        timeStart: this.makeDate(date.getFullYear() + "-" + nowMonth + "-" + day, this.startEventTime).getTime(),
        timeEnd: this.makeDate(date.getFullYear() + "-" + nowMonth + "-" + day, this.endEventTime).getTime(),
        withDuration: this.newEventToAdd.withDuration,
      }
      console.log("event = ", event)
      this.$store.state.events.push(event)

      this.makeRequest(event)
      console.log(this.$store.state.events)
      this.startEventTime = null
      this.endEventTime = null

      this.sortEvents()
      localStorage.setItem('events', JSON.stringify(this.$store.state.events))
      this.isPickEventDuration = false
      this.page = 'list';
    },
    findIndexEvents(id) {
      for (let i = 0; i < this.$store.state.events.length; i++) {
        if (this.$store.state.events[i].dt === id) {
          return i
        }
      }
      return -1
    },
    deleteTitleDown(event) {
      if (event === 'click' && ('ontouchstart' in window || navigator.msMaxTouchPoints)) {
        return
      }
      this.timerId = setTimeout(this.handlerFunc, 500)
    },
    handlerFunc() {
      this.timerDeleteTitle = true
      console.log('hi')
    },
    deleteTitleUp(title, index, event) {
      if (event === 'click' && ('ontouchstart' in window || navigator.msMaxTouchPoints)) {
        return
      }
      if (this.timerDeleteTitle) {
        console.log(title, index)
        this.titleIndex = index
        this.titleToDeleteID = title.id
        this.titleToDelete = title
        this.dialog = true
      } else {
        this.addEvent(title, index, event)
      }
      this.timerDeleteTitle = false
      clearTimeout(this.timerId)
    },
    deleteTitleWithIndex(index) {
      this.$store.state.titles.splice(index, 1)
      this.makeDeleteTitleResp(this.titleToDeleteID)
      this.$forceUpdate();
      localStorage.setItem('titles', JSON.stringify(this.$store.state.titles))
      console.log('delete');
    },
    createTitle() {
      console.log(this.makeDate('2019-09-10'))
      if (this.newEvent !== "") {
        // alert("newEvent")
        let newTitle = {
          title: this.newEvent,
          withDuration: this.isEventWithDuration
        }
        this.makeAddTitleResp(newTitle)
        this.$store.state.titles.push(newTitle)
        console.log(this.$store.state.titles)
        this.newEvent = ""
        this.isEventWithDuration = false
        localStorage.setItem('titles', JSON.stringify(this.$store.state.titles))
        this.$forceUpdate();
      }
    },

    async makeGetTitlesResp() {
      let response = await fetch('/api/getTitles', {
        method: 'GET',
      });

      if (response.ok) { // если HTTP-статус в диапазоне 200-299
        // получаем тело ответа (см. про этот метод ниже)
        let titles = await response.json()
        if (titles === null) {
          titles = []
        }
        console.log(titles)
        this.$store.state.titles = titles
        console.log(this.$store.state.titles)
        localStorage.setItem('titles', JSON.stringify(this.$store.state.titles))

        console.log("okkkk")
      } else if (response.status > 299 && response.status < 500) {
        console.log("err")
      } else if (response.status > 500) {
        console.log("some wrong on backend")
      }
    },
    async makeAddTitleResp(title) {
      let response = await fetch('/api/newTitle', {
        method: 'POST',
        body: JSON.stringify(title)
      });

      if (response.ok) { // если HTTP-статус в диапазоне 200-299
        // получаем тело ответа (см. про этот метод ниже)
        console.log("okkkk")
      } else {
        console.log("not ok")
      }
    },
    async makeDeleteTitleResp(id) {
      let response = await fetch('/api/deleteTitle', {
        method: 'DELETE',
        body: JSON.stringify({id: id}),
      });
      if (response.ok) { // если HTTP-статус в диапазоне 200-299
        // получаем тело ответа (см. про этот метод ниже)
        console.log(this.$store.state.titles)
        localStorage.setItem('titles', JSON.stringify(this.$store.state.titles))

        console.log("okkkk")
      } else if (response.status > 299 && response.status < 500) {
        console.log("err")

      } else if (response.status > 500) {
        console.log("some wrong on backend")
      }
    },
    changeIsNewTask(bool){
      this.isNewTask = bool
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
    openNewTask() {
      this.isNewTask = true
    },
  },
  beforeMount() {
    this.makeGetTitlesResp()
  },
}
</script>

<style scoped>
.titles {
  margin: -3em 0 0 -2em;

  /* Выравнивание по центру */
  text-align: center;
}

.title {
  display: inline-block;

  vertical-align: top;

  /* Убираем выравнивание по центру */
  text-align: left;

  margin: 3em 0 0 2em;
}

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