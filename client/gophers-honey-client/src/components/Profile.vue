<template>
  <div>
    <b-row>
      <b-col>
        <b-input disabled></b-input>
      </b-col>
    </b-row>
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
      dismissSecs: 3,
      alert: "",
      variant: "",
      user: {
        firstName: "",
        lastName: "",
        email: "",
        role: "",
        notifications_enabled: false,
        username: ""
      }
    }
  },
  created() {
    axios.defaults.headers.common['Authorization'] = 'Bearer '+ this.$cookies.get("token")
    this.getUser()
  },
  methods: {
    getUser: function() {
      let token = this.$cookies.get("token")
      let tokenPayload = JSON.parse(atob(token.split(".")[1]))
      axios({
        url: process.env.VUE_APP_API_ROOT+"/users?user="+tokenPayload.username,
        method: 'GET',
      }).then(function (response){
        if (response.data.error == null) {
          this.user = response.data
          window.console.log(this.user)
        }
      }.bind(this))
    }
  }
}
</script>

<style scoped>

</style>