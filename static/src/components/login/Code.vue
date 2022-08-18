<template>
  <div>
    <v-form
        class="form"
        :style="'margin-top:'+marginTop()"
        ref="form"
        v-model="valid"
        lazy-validation
    >
      <h2 style="margin-bottom: 20px">Check your email for code</h2>
      <v-row style="margin-bottom: 0">
        <v-btn
            @click="changePage('signUp')"
            color="grey"
            plain
        >
          change email
        </v-btn>
      </v-row>
      <v-text-field
          style="width: 150px"
          v-model="code"
          :rules="codeRules"
          v-mask="'####'"
          label="X-X-X-X"
          @input="checkCode"
          solo
      ></v-text-field>
      <v-btn
          style="height: 40px;"
          color="primary"
          @click="signUp"
      >
        reset code
      </v-btn>
    </v-form>
    <v-form v-if="page==='signUpIsOk'">
      <div class="text-h2" style="color: limegreen">sign up is successful</div>
    </v-form>
    <v-form v-if="page==='wrong'">
      <div class="text-h2" style="color: red">bad code</div>
    </v-form>
  </div>
</template>

<script>
export default {
  name: "code",
  data: () => ({
    wasReq: false,
    valid: true,
    show: false,
    page: '',
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
    checkCode() {
      if (this.code.length === 4 && !this.wasReq) {
        this.checkCodeReq({token: this.getCookie('codeId'), code: this.code})
        this.wasReq = true
        console.log('code = 4')
        // this.page = 'signUpIsOk'
        // setTimeout(this.changePage('signUpIsOk'), 3000)
      } else {
        this.wasReq = false
      }
    },
    signUp() {
      this.$refs.form.validate()
      this.signUpReq({
        email: this.$store.state.user.email,
        login: this.$store.state.user.login,
        password: this.$store.state.user.password
      })
    },
    changePage() {
      console.log("replace!!")
      window.location.replace("/")
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
        console.log("not ok")
      }
    },
    async checkCodeReq(codes) {
      console.log('make request :', codes)
      let response = await fetch('/api/checkCode', {
        method: 'POST',
        body: JSON.stringify(codes)
      });

      if (response.ok) { // если HTTP-статус в диапазоне 200-299
        // получаем тело ответа (см. про этот метод ниже)
        console.log("response = ", response)
        console.log("okkkk")
        this.page = 'signUpIsOk'
        console.log("this.page = 'signUpIsOk'")
        setTimeout(this.changePage, 1500)
      } else {
        this.page = 'wrong'
        console.log("not ok")
      }
    },
    marginTop() {
      const devices = new RegExp('Android|webOS|iPhone|iPad|iPod|BlackBerry|BB|PlayBook|IEMobile|Windows Phone|Kindle|Silk|Opera Mini', "i");
      if (devices.test(navigator.userAgent)) {
        return '100px'
      } else {
        return '40px'
      }
    },
    getCookie(name) {
      const value = `; ${document.cookie}`;
      const parts = value.split(`; ${name}=`);
      if (parts.length === 2) return parts.pop().split(';').shift();
    }
  }
}
</script>

<style scoped>

</style>