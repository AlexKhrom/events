<template>
  <v-footer
      class="footer"
      color="primary lighten-1"
      padless
      style="border-radius: 4px"
  >
    <v-row
        justify="center"
        no-gutters
        style="border-radius: 5px"
    >
      <v-btn style="display: flex;margin-right: 50px" @click="$router.push('/add')" color="white"
             text
             rounded
             class="my-2">Add Event
      </v-btn>

      <v-btn @click="$router.push('/events')"
             color="white"
             text
             rounded
             class="my-2">Show all
      </v-btn>
    </v-row>
  </v-footer>
</template>

<script>
export default {
  name: "Footer",
  methods:{
    calcTop(){
      console.log(document.documentElement.clientHeight)
      return document.documentElement.clientHeight-50
    }
  },
  mounted() {
    this.calcTop()
  }
}
</script>

<style scoped>
.footer {
  position: sticky;
  margin-top: calc(40px);
  bottom: 20px;
}
</style>