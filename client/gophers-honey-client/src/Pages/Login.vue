<template>
  <div class="container">
    <div class="sub-container">
      <div class="text-center login-container">
        <img src="../assets/gopher_honey.png" alt="Gophers-Honey" width="170" height="170" class="mb-4">
        <form @submit.prevent="login">
          <b-row class="input">
            <b-col md="4" style="margin: auto">
              <b-input
                  v-model="userinfo.username"
                  type="text"
                  required
                  placeholder="Username"
              >
              </b-input>
            </b-col>
          </b-row>
          <b-row class="input">
            <b-col md="4" style="margin: auto">
              <b-input
                  v-model="userinfo.password"
                  type="password"
                  required
                  placeholder="Password"
              >
              </b-input>
            </b-col>
          </b-row>
          <template v-if="loading">
            <div class="lds-ellipsis"><div></div><div></div><div></div><div></div></div>
          </template>
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
                  :disabled="(userinfo.username.length === 0 || userinfo.password.length === 0)"
                  type="submit"
                  class="submit-button"
              >Login</b-button>
            </b-col>
          </b-row>
        </form>
      </div>
    </div>
  </div>
</template>

<script>

import axios from 'axios';
import { router } from '../router';

export default{
  name: "Login",
  data: function() {
    return {
      userinfo: {username: "", password: ""},
      dismissCountDown: 0,
      dismissSecs: 3,
      alert: "",
      variant: "",
      loading: false,
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
      this.loading = true
      this.dismissCountDown = 0
      axios.post(
        process.env.VUE_APP_API_ROOT+"/users/login", userinfoJson
      ).then(response => {
        if (response.status === 200) {
          that.loading = false
          if (response.data.error === "") {
            that.$cookies.set("token",response.data.token,"24h","/")
            router.push('/').catch(()=>{})
          } else if (response.data.error === "Incorrect username or password") {
            that.alert = response.data.error
            that.showAlert("danger")
          } else {
            that.alert = "Internal server error"
            that.showAlert("danger")
          }
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

.lds-ellipsis {
  display: inline-block;
  position: relative;
  width: 80px;
  height: 80px;
}
.lds-ellipsis div {
  position: absolute;
  top: 33px;
  width: 13px;
  height: 13px;
  border-radius: 50%;
  background: #7e7e7e;
  animation-timing-function: cubic-bezier(0, 1, 1, 0);
}
.lds-ellipsis div:nth-child(1) {
  left: 8px;
  animation: lds-ellipsis1 0.6s infinite;
}
.lds-ellipsis div:nth-child(2) {
  left: 8px;
  animation: lds-ellipsis2 0.6s infinite;
}
.lds-ellipsis div:nth-child(3) {
  left: 32px;
  animation: lds-ellipsis2 0.6s infinite;
}
.lds-ellipsis div:nth-child(4) {
  left: 56px;
  animation: lds-ellipsis3 0.6s infinite;
}
@keyframes lds-ellipsis1 {
  0% {
    transform: scale(0);
  }
  100% {
    transform: scale(1);
  }
}
@keyframes lds-ellipsis3 {
  0% {
    transform: scale(1);
  }
  100% {
    transform: scale(0);
  }
}
@keyframes lds-ellipsis2 {
  0% {
    transform: translate(0, 0);
  }
  100% {
    transform: translate(24px, 0);
  }
}
</style>