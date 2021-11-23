<template>
  <div style="height: 100%">
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
    <div style="overflow-y: auto; height: 100%">
      <h5 class="text-center">SMTP setup</h5>
      <b-form @submit.prevent="updateSmtp" class="container" style="height: fit-content">
        <b-form-row>
          <b-col class="input">
            <b-form-group
                id="input-group-1"
                label="SMTP host"
                label-for="input-1"
            >
              <b-form-input
                  id="input-1"
                  v-model="smtpUpdate.smtp_host"
                  type="text"
                  :placeholder="conf.smtp_server.smtp_host"
                  @input.native="checkSmtpForm()"
              >
              </b-form-input>
            </b-form-group>
          </b-col>
          <b-col class="input">
            <b-form-group
                id="input-group-2"
                label="SMTP port"
                label-for="input-2"
            >
              <b-form-input
                  id="input-2"
                  v-model="smtpUpdate.smtp_port"
                  type="number"
                  :placeholder="toString(conf.smtp_server.smtp_port)"
                  @input.native="checkSmtpForm()"
              >
              </b-form-input>
            </b-form-group>
          </b-col>
        </b-form-row>
        <b-form-row>
          <b-col class="input">
            <b-form-group
                id="input-group-3"
                label="Username"
                label-for="input-3"
            >
              <b-form-input
                  id="input-3"
                  v-model="smtpUpdate.username"
                  type="text"
                  placeholder="Username"
                  @input.native="checkSmtpForm()"
              >
              </b-form-input>
            </b-form-group>
          </b-col>
          <b-col class="input">
            <b-form-group
                id="input-group-4"
                label="Password"
                label-for="input-4"
            >
              <b-form-input
                  id="input-4"
                  v-model="smtpUpdate.password"
                  type="password"
                  placeholder="Password"
                  @input.native="checkSmtpForm()"
              >
              </b-form-input>
            </b-form-group>
          </b-col>
        </b-form-row>
        <div style="margin: auto; width: fit-content">
          <b-button type="submit" :disabled="!smtpFormValid">Save SMTP settings</b-button>
          <b-button v-on:click="testEmail" style="background-color: cadetblue; margin-left: 15px">Send a test email
          </b-button>
        </div>
      </b-form>
    </div>
  </div>
</template>

<script>
import axios from "axios";

export default {
  name: "Service",
  data() {
    return {
      dismissCountDown: 0,
      dismissSecs: 10,
      alert: "",
      variant: "",
      smtpFormValid: false,
      conf: {
        smtp_server: {
          smtp_host: "",
          smtp_port: 0,
          username: "",
          password: "",
        },
        whitelist: []
      },
      smtpUpdate: {
        smtp_host: "",
        smtp_port: 0,
        username: "",
        password: "",
      },
      whitelistUpdate: [],
    }
  },
  created() {
    this.getConf()
  },
  methods: {
    showAlert: function (variant) {
      this.variant = variant
      this.dismissCountDown = this.dismissSecs
    },
    countDownChanged: function (dismissCountDown) {
      this.dismissCountDown = dismissCountDown
    },
    removeItem: function (array, key, value) {
      const index = array.findIndex(obj => obj[key] === value)
      return index >= 0 ? [
        ...array.slice(0, index),
        ...array.slice(index + 1)
      ] : array;
    },
    checkSmtpForm: function () {
      if (parseInt(this.smtpUpdate.smtp_port) !== this.conf.smtp_server.smtp_port ||
          this.smtpUpdate.smtp_host.length > 0 || (this.smtpUpdate.username.length > 0 && this.smtpUpdate.password.length > 0)) {
        // If username is password is filled, check if the other has also been filled.
        if ((this.smtpUpdate.username.length > 0 && this.smtpUpdate.password.length === 0) ||
            (this.smtpUpdate.username.length === 0 && this.smtpUpdate.password.length > 0)) {
          this.smtpFormValid = false
          return
        }
        this.smtpFormValid = true
      } else {
        this.smtpFormValid = false
      }
    },
    updateSmtp: function () {
      if (this.smtpUpdate.smtp_host.length === 0) {
        this.smtpUpdate.smtp_host = this.conf.smtp_server.smtp_host
      }
      this.smtpUpdate.smtp_port = parseInt(this.smtpUpdate.smtp_port)
      window.console.log("updated smtp conf", this.smtpUpdate)
      axios({
        url: process.env.VUE_APP_API_ROOT + "/config",
        method: "PATCH",
        data: this.smtpUpdate
      }).then(function (response) {
        if (response.status === 200) {
          if (response.data.error === "") {
            this.alert = "Succesfully updated smtp data"
            this.showAlert("success")
            this.smtpUpdate = {
              smtp_host: "",
              smtp_port: this.smtpUpdate.smtp_port,
              username: "",
              password: "",
            }
            this.smtpFormValid = false
            this.getConf()
          } else {
            this.alert = "Error updating smtp data: " + response.data.error
            this.showAlert("danger")
          }
        } else {
          this.alert = "Server sent response code " + response.status
          this.showAlert("danger")
        }
      }.bind(this))
    },
    testEmail: function () {
      window.console.log("Sending testemail")
      axios({
        url: process.env.VUE_APP_API_ROOT + "/config/testEmail",
        method: "GET"
      }).then(function (response) {
        if (response.status === 200) {
          if (response.data.error === "") {
            this.alert = "Succesfully sent email"
            this.showAlert("success")
          } else {
            this.alert = "Error sending test email: " + response.data.error
            this.showAlert("danger")
          }
        } else {
          this.alert = "Server sent response code: " + response.status
          this.showAlert("danger")
        }
      }.bind(this))
    },
    getConf: function () {
      axios({
        url: process.env.VUE_APP_API_ROOT + "/config",
        method: "GET",
      }).then(function (response) {
        if (response.status === 200) {
          window.console.log("Config from server: ", response.data)
          this.conf = response.data
          this.smtpUpdate.smtp_port = this.conf.smtp_server.smtp_port
        } else {
          window.console.log("Response code ", response.status)
        }
      }.bind(this))
    }
  }
}
</script>

<style scoped>

</style>