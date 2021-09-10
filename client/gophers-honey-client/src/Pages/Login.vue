<template>
  <div class="container">
    <div class="sub-container">
      <div class="text-center login-container">
        <img src="../assets/gopher_honey.png" alt="Gophers-Honey" width="170" height="170" class="mb-4">
        <b-row class="input">
          <b-col md="4" style="margin: auto">
            <b-input
              v-model="userinfo.username"
              placeholder="Username"
              @keyup.enter="login"
            >
            </b-input>
          </b-col>
        </b-row>
        <b-row class="input">
          <b-col md="4" style="margin: auto">
            <b-input
                v-model="userinfo.password"
                placeholder="Password"
                @keyup.enter="login"
            >
            </b-input>
          </b-col>
        </b-row>
        <b-row>
          <b-col md="4" class="login-alert">
            <b-alert
                :show="dismissCountDown"
                dismissible
                :variant="variant"
                fade
                @dismissed="dismissCountDown=0"
                @dismiss-count-down="countDownChanged"
            >
              {{ alert }}
            </b-alert>
          </b-col>
        </b-row>
        <b-row>
          <b-col md="4" class="submit">
            <b-button
                :disabled="(userinfo.username.length == 0 || userinfo.password.length == 0)"
                @click="login"
                type="submit"
                class="submit-button"
            >Login</b-button>
          </b-col>
        </b-row>
      </div>
    </div>
  </div>
</template>

<script>

import axios from 'axios';

export default{
  name: "Login",
  data: function() {
    return {
      userinfo: {username: "", password: ""},
      dismissCountDown: 0,
      dismissSecs: 3,
      alert: "",
      variant: "",
    }
  },
  mounted() {

  },
  methods: {
    countDownChanged: function (dismissCountDown) {
      this.dismissCountDown = dismissCountDown
    },
    showAlert: function (variant) {
      this.variant = variant
      this.dismissCountDown = this.dismissSecs
    },
    login: function() {
      let that = this
      if (!this.userinfo.username.length || !this.userinfo.password.length) {
        this.alert = "Username or password is empty"
        this.showAlert("danger")
        return
      }
      let userinfoJson = JSON.stringify(this.userinfo)

      axios.post(
        process.env.VUE_APP_API_ROOT+"/users/login", userinfoJson
      ).then(response => {
        if (response.status == 200) {
          that.$cookies.set("session",response.data,"30min","/")
        }
      })
    }
  }
}
</script>

<style>
p {
  text-align: left;
  margin-bottom: 0;
}
.input{
  margin-top: 5px;
  margin-bottom: 15px;
}
.login-container{
  position: relative;
  margin: auto;
  top: 50%;
  left: 50%;
  transform: translate(-50%,-50%);
  align-items: center;
}
.submit {
  margin: auto;
}
.login-alert{
  margin: auto;
}
.submit-button{
  width: 100%;
}
.container{
  height: 100vh
}
.sub-container{
  height: 100vh
}
body{
  height: 100vh
}
</style>