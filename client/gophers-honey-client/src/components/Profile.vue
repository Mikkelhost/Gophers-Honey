<template>
  <div>
    <b-form @submit.prevent="submitUser">
      <b-row>
        <b-col>
          <label type="input" for="first_name">First name</label>
          <b-form-input id="first_name" disabled :placeholder="user.first_name"></b-form-input>
        </b-col>
        <b-col>
          <label type="input" for="last_name">Last name</label>
          <b-form-input id="last_name" disabled :placeholder="user.last_name"></b-form-input>
        </b-col>
        <b-col>
          <label type="input" for="role">Role</label>
          <b-form-input id="role" disabled :placeholder="user.role"></b-form-input>
        </b-col>
      </b-row>
      <b-row>
        <b-col md="4">
          <label type="input" for="username">Username</label>
          <b-form-input id="username" disabled :placeholder="user.username"></b-form-input>
        </b-col>
        <b-col md="6">
          <label type="input" for="email">Email</label>
          <b-form-input id="email" type="email" @input.native="checkUserForm()" :placeholder="user.email" v-model="form.email"></b-form-input>
        </b-col>
        <b-col md="2">
          <label type="checkbox" for="notifications">Email alerts</label>
          <b-form-checkbox id="notifications" @change="checkUserForm()" v-model="form.notifications_enabled" switch></b-form-checkbox>
        </b-col>
      </b-row>
      <b-row>
        <b-col>
          <label type="password" for="currpassword">Current password*</label>
          <b-form-input id="currpassword" type="password" @input.native="checkUserForm()" placeholder="Current password" v-model="form.curr_password" required></b-form-input>
        </b-col>
        <b-col>
          <label type="password" for="password">Password</label>
          <b-form-input id="password" type="password" @input.native="checkUserForm()" placeholder="Password" v-model="form.password"></b-form-input>
        </b-col>
        <b-col>
          <label type="password" for="confirmpw">Confirm Password</label>
          <b-form-input id="confirmpw" type="password" @input.native="checkUserForm()" placeholder="Confirm Password" v-model="form.confirmPw"></b-form-input>
        </b-col>
      </b-row>
      <div style="margin: auto; width: fit-content">
        <b-button type="submit" :disabled="!profileFormValid">Save profile settings</b-button>
      </div>
      <b-col md="12">
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
    </b-form>
  </div>
</template>

<script>
import axios from "axios";

export default {
  name: "Profile",
  data: function(){
    return{
      loading: false,
      dismissCountDown: 0,
      dismissSecs: 5,
      alert: "",
      variant: "",
      profileFormValid: false,
      form: {
        email: "",
        notifications_enabled: false,
        username: "",
        password: "",
        curr_password: "",
        confirmPw: "",
      },
      user: {
        first_name: "",
        last_name: "",
        email: "",
        role: "",
        notifications_enabled: false,
        username: ""
      },
      claims: null,
    }
  },
  created() {
    axios.defaults.headers.common['Authorization'] = 'Bearer '+ this.$cookies.get("token")
    let token = this.$cookies.get("token")
    this.claims = JSON.parse(atob(token.split(".")[1]))
    this.getUser()
  },
  methods: {
    showAlert: function (variant) {
      this.variant = variant
      this.dismissCountDown = this.dismissSecs
    },
    countDownChanged: function (dismissCountDown) {
      this.dismissCountDown = dismissCountDown
    },
    checkUserForm: function(){
      if ( (this.form.email.length === 0 && this.form.password.length === 0 && this.form.confirmPw.length === 0 && this.form.notifications_enabled === this.user.notifications_enabled) || this.form.curr_password.length === 0){
        this.profileFormValid = false
      } else {
        this.profileFormValid = true
      }
    },
    getUser: function() {
      axios({
        url: process.env.VUE_APP_API_ROOT+"/users?user="+this.claims.username,
        method: 'GET',
      }).then(function (response){
        if (response.data.error == null) {
          this.user = response.data
          this.form.notifications_enabled = this.user.notifications_enabled
          window.console.log(this.user)
        }
      }.bind(this))
    },
    submitUser: function() {
      window.console.log("Submitting profile settings")
      if (this.form.password.length > 0 && this.form.password !== this.form.confirmPw){
        this.alert = "Passwords has to match"
        this.showAlert("danger")
        return
      }
      axios({
        url: process.env.VUE_APP_API_ROOT+"/users",
        method: "PUT",
        data: this.form
      }).then(function(response){
        window.console.log(response.data)
        if (response.data.error == "") {
          this.user.notifications_enabled = this.form.notifications_enabled
          if (this.form.email.length > 0) {
            this.user.email = this.form.email
          }
          this.form.password = ""
          this.form.confirmPw = ""
          this.form.email = ""
          this.checkUserForm()
          this.alert = "Succesfully updated user profile"
          this.showAlert("success")
        } else {
          this.alert = response.data.error
          this.showAlert("danger")
        }
      }.bind(this))
    }
  }
}
</script>

<style scoped>
.row {
  margin: 5px 0 10px 0;
}
</style>