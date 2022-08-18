<!--evets = []event{-->
<!--  dt: date of creat event-->
<!--  title: name of event-->
<!--  timeStart: date of start doing smth-->
<!--  timeEnd: date of finish doing smth-->
<!--}-->
<template>
  <div>
    <template v-if="isChangeEvent">
      <div class="picker" @click="isChangeEvent=false">
        <div class="picker__content" @click.stop>
          <slot>
            <v-row justify="space-around">
              <v-col>
                <v-time-picker
                    format="24hr"
                    v-model="timePicker"
                ></v-time-picker>
              </v-col>
              <v-col>
                <v-date-picker
                    v-model="datePicker"
                    color="primary lighten-1"
                ></v-date-picker>
              </v-col>
              <v-col>
                <v-btn class="mb-4" width="290px"
                       @click="saveEventChanges"
                       color="white">save
                </v-btn>
              </v-col>
            </v-row>
          </slot>
        </div>
      </div>
    </template>

    <template v-if="isChangeEventWithDuration">
      <div class="picker" @click="isChangeEventWithDuration=false">
        <div class="picker__content" @click.stop>
          <slot>
            <div>
              <h1>Change your event:</h1>
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
                </v-col>
                <v-col>
                  <v-date-picker
                      v-model="datePicker"
                      color="primary lighten-1"
                  ></v-date-picker>
                </v-col>
                <v-col>
                  <v-btn class="mb-2" width="290px"
                         @click="saveEventWithDurationChanges"
                         height="50px"
                         color="white">save
                  </v-btn>
                </v-col>
              </v-row>
            </div>
          </slot>
        </div>
      </div>
    </template>


    <div v-for="(event,index) in this.$store.state.events"
         :key="event.id"
    >
      <div v-if="isNewDay(index)" style="margin-top: 20px;color: #555555">{{ new Date(event.dt) | moment }}</div>
      <div
          style="display: flex;justify-content: space-between;margin-top: 20px;padding:5px;">
        <div>
          <div :style="event.deleted?'opacity: 0.5':''"
               style="font-size: 18px;text-transform: capitalize;display: flex;align-items: center;">{{ event.title }}
          </div>
          <span v-if="!event.deleted" style="color: #555;font-size: 11px;display: block">{{
              getHourAndMinutes(event)
            }}</span>
        </div>
        <div>
          <span v-if="event.deleted" style="color: #555;font-size: 12px;padding-right: 10px">{{
              event.waitSeconds
            }}с</span>
          <span v-if="event.timeStart!==0&&event.timeStart!==undefined&&!event.deleted" style="color: #555;font-size: 14px;">
            {{ printEventDuration(event) }}
          </span>
          <v-btn v-if="!event.deleted" @click="changeEvent(event.id)" icon>✏️</v-btn>
          <v-btn v-if="event.deleted" @click="camebackEvent(event.id)" text small>вернуть ↩︎</v-btn>
          <v-btn v-else @click="deleteEvent(event.id)" icon>🗑</v-btn>
        </div>
      </div>
      <div :style="{border: '2px solid '+borderColorForEvent(index),width: calcTimePercent(index) +'%'}"
      ></div>
    </div>
    <Footer >

    </Footer>
  </div>
</template>

<script>
import moment from "moment";
import Footer from "@/components/Footer";


export default {
  name: "listEvents",
  components: {Footer},
  filters: {
    moment(date) {
      return moment(date).format("D MMMM");
    },
  },
  data() {
    return {
      isChangeEvent: false,
      isChangeEventWithDuration: false,
      startEventTime: null,
      endEventTime: null,
      eventToChangeId: -1,
      timePicker: "",
      percentOfEventTime: 0,
      datePicker: (new Date(Date.now() - (new Date()).getTimezoneOffset() * 60000)).toISOString().substr(0, 10),

    }
  },
  computed: {
    eventsArray() {
      // console.log(this.$store.state.events.slice().sort((a, b) =>
      //     a.dt < b.dt ? 1 : -1))
      return this.$store.state.events.slice().sort((a, b) =>
          a.dt < b.dt ? 1 : -1)
    }
  },
  methods: {
    deleteEvent(id) {
      let index = this.findIndexEvents(id)
      let event = this.$store.state.events[index]

      // console.log("index = ", index)
      // console.log({event})

      let waitSeconds = 3
      event.waitSeconds = waitSeconds
      let intervalId = setInterval(() => {
        if (!event.deleted) {
          clearInterval(intervalId)
          localStorage.setItem('events', JSON.stringify(this.$store.state.events))
          this.$forceUpdate()
          return
        }
        if (--waitSeconds < 0) {
          clearInterval(intervalId);
          let index = this.findIndexEvents(id)
          this.deleteEventFromBD(id)
          this.$store.state.events.splice(index, 1)
        } else {
          event.waitSeconds = waitSeconds
        }
        this.$forceUpdate()
      }, 1000);

      // this.$store.state.events.splice(index, 1)

      event.deleted = Date.now();
      this.$forceUpdate();
      localStorage.setItem('events', JSON.stringify(this.$store.state.events))
    },
    camebackEvent(id) {
      let index = this.findIndexEvents(id)
      this.$store.state.events[index].deleted = 0
      this.$forceUpdate();
    },
    changeEvent(id) {
      let index = this.findIndexEvents(id)
      let date = new Date(this.$store.state.events[index].dt)
      this.selectedEventIndex = index
      if (this.$store.state.events[index].timeEnd !== 0) {
        let wasDateStart = new Date(this.$store.state.events[index].timeStart)
        let wasDateEnd = new Date(this.$store.state.events[index].timeEnd)

        this.startEventTime = this.makeHoursAndMinutesStr(wasDateStart)
        this.endEventTime = this.makeHoursAndMinutesStr(wasDateEnd)
        this.datePicker = (new Date(date - (new Date()).getTimezoneOffset() * 60000)).toISOString().substr(0, 10)

        console.log(this.startEventTime)
        console.log(this.endEventTime)
        this.isChangeEventWithDuration = true

      } else {

        this.timePicker = this.makeHoursAndMinutesStr(date)
        this.datePicker = (new Date(date - (new Date()).getTimezoneOffset() * 60000)).toISOString().substr(0, 10)
        this.isChangeEvent = true
        this.$forceUpdate();
      }
    },
    saveEventWithDurationChanges() {
      let startDate = this.makeDate(this.datePicker, this.startEventTime)
      // let startDate = new Date(this.datePicker + "T" + this.startEventTime)
      let endDate = this.makeDate(this.datePicker, this.endEventTime)
      // let endDate = new Date(this.datePicker + "T" + this.endEventTime)

      this.$store.state.events[this.selectedEventIndex].timeStart = startDate.getTime() + this.$store.state.events[this.selectedEventIndex].dt % 1000
      this.$store.state.events[this.selectedEventIndex].timeEnd = endDate.getTime() + this.$store.state.events[this.selectedEventIndex].dt % 1000
      this.$store.state.events[this.selectedEventIndex].dt = this.$store.state.events[this.selectedEventIndex].timeStart

      this.changeEventFromBD(this.$store.state.events[this.selectedEventIndex])
      localStorage.setItem('events', JSON.stringify(this.$store.state.events))
      this.isChangeEventWithDuration = false
      this.sortEvents()
    },
    saveEventChanges() {
      let date = this.makeDate(this.datePicker, this.timePicker)
      // let date = new Date(this.datePicker + "T" + this.timePicker)

      this.$store.state.events[this.selectedEventIndex].dt = date.getTime() + this.$store.state.events[this.selectedEventIndex].dt % 1000

      this.changeEventFromBD(this.$store.state.events[this.selectedEventIndex])
      localStorage.setItem('events', JSON.stringify(this.$store.state.events))
      this.isChangeEvent = false
      this.sortEvents()
    },
    calcTimePercent(index) {
      if (this.$store.state.events[index].timeEnd === 0) {
        return 100
      }
      if ((this.$store.state.events[index].timeEnd - this.$store.state.events[index].timeStart) <= 0)
        return 0

      let timeForDay = 0
      for (let i = 0; i < this.eventsArray.length; i++) {
        if ((new Date(this.$store.state.events[i].dt)).getDate() === (new Date(this.$store.state.events[index].dt)).getDate()
            && this.$store.state.events[i].timeStart !== 0) {
          timeForDay += this.$store.state.events[i].timeEnd - this.$store.state.events[i].timeStart
        }
      }
      return this.percentOfEventTime = (this.$store.state.events[index].timeEnd - this.$store.state.events[index].timeStart) / timeForDay * 100
    },

    async getEvents() {
      let response = await fetch('/api/getEvents', {
        method: 'GET',
      });

      if (response.ok) { // если HTTP-статус в диапазоне 200-299
        // получаем тело ответа (см. про этот метод ниже)
        let events = await response.json()
        this.$store.state.events = events
        console.log(this.$store.state.events)
        this.sortEvents()
        localStorage.setItem('events', JSON.stringify(this.$store.state.events))

        console.log("okkkk")
      } else if (response.status > 299 && response.status < 500) {
        console.log("err")
      } else if (response.status > 500) {
        console.log("some wrong on backend")
      }
    },

    async changeEventFromBD(event) {
      let response = await fetch('/api/changeEvent', {
        method: 'PUT',
        body: JSON.stringify(event),
      });
      if (response.ok) { // если HTTP-статус в диапазоне 200-299
        // получаем тело ответа (см. про этот метод ниже)
        console.log(this.$store.state.events)
        this.sortEvents()
        localStorage.setItem('events', JSON.stringify(this.$store.state.events))

        console.log("okkkk")
      } else if (response.status > 299 && response.status < 500) {
        console.log("err")

      } else if (response.status > 500) {
        console.log("some wrong on backend")
      }
    },
    async deleteEventFromBD(id) {
      let response = await fetch('/api/deleteEvent', {
        method: 'DELETE',
        body: JSON.stringify({id: id}),
      });
      if (response.ok) { // если HTTP-статус в диапазоне 200-299
        // получаем тело ответа (см. про этот метод ниже)
        console.log(this.$store.state.events)
        this.sortEvents()
        localStorage.setItem('events', JSON.stringify(this.$store.state.events))

        console.log("okkkk")
      } else if (response.status > 299 && response.status < 500) {
        console.log("err")

      } else if (response.status > 500) {
        console.log("some wrong on backend")
      }
    },


    borderColorForEvent(index) {
      // console.log(this.$store.state.events[index].timeStart === undefined ? "blue" : "green")
      return this.$store.state.events[index].timeStart === 0 ? "none" : "green"//rgb(47,101,186)-blue
    },
    sortEvents() {
      this.$store.state.events = this.$store.state.events.slice().sort((a, b) =>
          a.dt < b.dt ? 1 : -1)
    },

    isNewDay(index) {
      if (index === 0) {
        return true
      }
      console.log(index, (new Date(this.eventsArray[index - 1].dt)).getDate(), (new Date(this.eventsArray[index].dt)).getDate())
      return (new Date(this.eventsArray[index - 1].dt)).getDate() !== (new Date(this.eventsArray[index].dt)).getDate();
    },

    printEventDuration(event) {
      let hours = Math.floor((event.timeEnd - event.timeStart) / (3600000))
      let minutes = Math.floor((event.timeEnd - event.timeStart) % (3600000) / 60000)
      minutes = minutes.toString().length === 1 ? '0' + minutes : minutes
      return hours + ":" + minutes
    },

    getHourAndMinutes(event) {
      // let hour = moment(event.dt).hour()
      // let minutes = moment(event.dt).minutes().toString().length === 1 ? "0" + moment(event.dt).minutes() : moment(event.dt).minutes()
      if (event.timeEnd === 0 || event.timeEnd === undefined) {
        return this.makeHoursAndMinutesStr(new Date(event.dt))
      } else {
        return this.makeHoursAndMinutesStr(new Date(event.dt)) + " - " + this.makeHoursAndMinutesStr(new Date(event.timeEnd))

      }
    },
    getEndOfEventStr(event) {
      return this.makeHoursAndMinutesStr(new Date(event.timeEnd))
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
        seconds = seconds.length === 1 ? '0' + seconds : seconds
      }
      seconds = seconds === undefined ? "" : ":" + seconds
      milliseconds = milliseconds === undefined ? "" : "." + milliseconds
      console.log("make Date = ", date + time + seconds + milliseconds)

      return new Date(date + time + seconds + milliseconds)
    },
    findIndexEvents(id) {
      for (let i = 0; i < this.$store.state.events.length; i++) {
        if (this.$store.state.events[i].id === id) {
          return i
        }
      }
      return -1
    },

  },
  beforeMount() {
    this.getEvents()
  },
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
  padding: 0;

  min-height: 100%;
  justify-content: space-between;
  z-index: 200;
}
</style>