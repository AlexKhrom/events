<template>
  <div>
    <v-form
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
          label="email"
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
          required
          solo
      ></v-text-field>


      <v-btn
          style="height: 40px;width: 100px"
          :disabled="!valid"
          color="primary"
          class="my-2"
          @click="validate"
      >
        Login
      </v-btn>

      <v-btn
          style="height: 40px;margin-left: 10px"
          class="my-2"
          @click="$router.push('/signUp');email='';login='';password=''"
      >
        Sign up
      </v-btn>

    </v-form>


    <v-form v-if="page==='signUpIsOk'">
      <div class="text-h2" style="color: limegreen">login is successful</div>
    </v-form>
    <v-form v-if="page==='wrong'">
      <div class="text-h2" style="color: red">Wrong login or password</div>
    </v-form>
    <!--    <v-btn-->
    <!--        style="height: 40px;"-->
    <!--        color="primary"-->
    <!--        @click="signUp"-->
    <!--    >   send email-->
    <!--    </v-btn>-->
  </div>
</template>

<script>

export default {
  // components: {SignUp},
  data: () => ({
    page: 'login',
    valid: true,
    show: false,
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
    validate() {
      this.$refs.form.validate()
      this.loginReq({email: this.$store.state.user.email, password: this.$store.state.user.password})
    },
    reset() {
      this.$refs.form.reset()
    },
    resetValidation() {
      this.$refs.form.resetValidation()
    },
    marginTop() {
      const devices = new RegExp('Android|webOS|iPhone|iPad|iPod|BlackBerry|BB|PlayBook|IEMobile|Windows Phone|Kindle|Silk|Opera Mini', "i");
      if (devices.test(navigator.userAgent)) {
        return '100px'
      } else {
        return '40px'
      }
    },
    async loginReq(user) {
      console.log('make request :', user)
      let response = await fetch('/api/login', {
        method: 'POST',
        body: JSON.stringify(user)
      });


      if (response.ok) { // если HTTP-статус в диапазоне 200-299
        // получаем тело ответа (см. про этот метод ниже)
        console.log("response = ", response)
        this.page = 'signUpIsOk'
        setTimeout(this.changePage, 1500)
        console.log("okkkk")
        console.log("resp = ", response.status)
      } else if (response.status > 299 && response.status < 500) {
        console.log("resp = ", response.status)
        this.page = 'wrong'
      } else if (response.status > 500) {
        console.log("some wrong on backend")
      }
    },
    changePage() {
      window.location.replace("/")
    },


    getCookie(name) {
      const value = `; ${document.cookie}`;
      const parts = value.split(`; ${name}=`);
      if (parts.length === 2) return parts.pop().split(';').shift();
    }
  },
}
</script>

<style scoped>
.v-text-field input {
  font-size: 30px;
}

</style>