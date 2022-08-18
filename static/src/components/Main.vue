<template>
  <v-app style="padding: 20px;">
    <div class="content">
      <AddEvent>

      </AddEvent>
    </div>
  </v-app>
</template>

<script>
// import {Bar} from 'vue-chartjs/legacy'


import AddEvent from "@/components/AddEvent";

export default {

  name: 'Main',
  components: { AddEvent},
  data() {
    return {
      page: 'add',
    }
  },
  methods: {},
}
</script>

<style scoped>
body {
  min-height: 80%;
}

.content{
  flex: 1 0 auto;
  margin-bottom: 20px;
}


</style>
