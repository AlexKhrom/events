<template>
  <div>
    <v-form v-if="page==='signUp'"
            class="form"
            :style="'margin-top:'+marginTop()"
            ref="form"
            v-model="valid"
            lazy-validation
    >
      <v-text-field
          size="30px"
          v-model="$store.state.user.email"
          :rules="emailRules"
          @input="errowMes=''"
          label="E-mail"
          required
          solo
      ></v-text-field>

      <v-text-field
          size="30px"
          v-model="$store.state.user.login"
          :rules="loginRules"
          @input="errowMes=''"
          label="Login"
          required
          solo
      ></v-text-field>

      <v-text-field
          v-model="$store.state.user.password"
          :rules="passwordRules"
          label="Password"
          :append-icon="show ? 'mdi-eye' : 'mdi-eye-off'"
          :type="show ? 'text' : 'password'"
          @click:append="show = !show"
          solo
      ></v-text-field>


      <v-btn
          style="height: 40px;"
          class="my-2"
          color="primary"
          @click="signUp"
          :disabled="!valid"
      >
        Sign up
      </v-btn>
    </v-form>

    <Code @changePage="changePage" v-if="page==='code'"/>

    <v-form v-if="errowMes==='errorLogin'">
      <div class="text-h6" style="color: limegreen">This login is already exist</div>
    </v-form>
    <v-form v-if="errowMes==='errorEmail'">
      <div class="text-h6" style="color: red">This Email is already exist</div>
    </v-form>
  </div>
</template>

<script>
import Code from "@/components/login/Code";


export default {
  name: "signUp",
  components: {Code},
  data: () => ({
    valid: true,
    show: false,
    page: 'signUp',
    errowMes:"",
    login: '',
    email: '',
    code: '',
    loginRules: [
      v => !!v || 'Login is required',
      v => (v && v.length <= 6) || 'Login must be less than 4 characters',
    ],
    password: '',
    passwordRules: [
      v => !!v || 'Password is required',
      v => (v && v.length <= 10) || 'Password must be less than 10 characters',
    ],
    codeRules: [
      v => !!v || 'code is required',
      v => (v && v.length === 4) || 'code must be 4 characters',
    ],
    emailRules: [
      v => !!v || 'E-mail is required',
      v => /.+@.+\..+/.test(v) || 'E-mail must be valid',
    ],

  }),
  methods: {
    changePage(page) {
      console.log('chenge page = ', page)
      this.page = page
    },
    marginTop() {
      const devices = new RegExp('Android|webOS|iPhone|iPad|iPod|BlackBerry|BB|PlayBook|IEMobile|Windows Phone|Kindle|Silk|Opera Mini', "i");
      if (devices.test(navigator.userAgent)) {
        return '100px'
      } else {
        return '40px'
      }
    },
    signUp() {
      this.$refs.form.validate()
      this.signUpReq({
        email: this.$store.state.user.email,
        login: this.$store.state.user.login,
        password: this.$store.state.user.password
      })
      this.page = 'code'

    },
    async signUpReq(user) {
      console.log('make request :', user)
      let response = await fetch('/api/signUp', {
        method: 'POST',
        body: JSON.stringify(user)
      });

      if (response.ok) { // если HTTP-статус в диапазоне 200-299
        // получаем тело ответа (см. про этот метод ниже)
        console.log("response = ", response)
        console.log("okkkk")
      } else {
        if (response.status === 425) {
          this.errowMes = 'errorEmail'
        }else if(response.status === 426){
          this.errowMes = 'errorLogin'
        }
        console.log("not ok")
      }
    },
  }
}
</script>

<style scoped>

</style>