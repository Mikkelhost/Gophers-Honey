<template>
  <div style="height: 100%; overflow-y: hidden">
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
    <div style="overflow-y: auto; height: 100%; overflow-x: hidden;">
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
      <b-form @submit.prevent="addIp">
        <h5 class="text-center">Whitelist</h5>
        <p class="text-center description">Append or delete ips to be whitelisted. Whitelisted ips will not raise any
          alerts if they are making connections to the honeypots</p>
        <b-form-row>
          <b-col md="6" class="input" style="margin: auto;">
            <b-form-group
                id="input-group-5"
                label="Add IP"
                label-for="input-5"
            >
              <b-input-group>
                <b-form-input
                    id="input-5"
                    v-model="whitelistNewIp.ip_address"
                    type="text"
                    placeholder="Add new ip to whitelist"
                >
                </b-form-input>
                <b-input-group-append>
                  <b-button type="submit" style="background-color: olivedrab"
                            :disabled="!(whitelistNewIp.ip_address.length > 0)">
                    <b-icon-plus></b-icon-plus>
                  </b-button>
                </b-input-group-append>
              </b-input-group>
            </b-form-group>
          </b-col>
        </b-form-row>
      </b-form>
      <b-form>
        <b-form-row v-for="(ip, index) in conf.ip_whitelist" v-bind:key="index" style="margin: 10px 0 10px 0;">
          <b-col md="6" class="input" style="margin: auto;">
            <b-input-group>
              <b-form-input
                  type="text"
                  :placeholder="ip"
              >
              </b-form-input>
              <b-input-group-append>
                <b-button style="background-color: #b73333" v-on:click="deleteIp(ip, index)">
                  <b-icon-dash></b-icon-dash>
                </b-button>
              </b-input-group-append>
            </b-input-group>
          </b-col>
        </b-form-row>
      </b-form>
    </div>
  </div>
</template>

<script>
import axios from "axios";
import getEnv from '../utils/env'

export default {
  name: "Service",
  data() {
    return {
      apiRoot: getEnv('VUE_APP_API_ROOT'),
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
        ip_whitelist: []
      },
      smtpUpdate: {
        smtp_host: "",
        smtp_port: 0,
        username: "",
        password: "",
      },
      whitelistNewIp: {
        delete: false,
        ip_address: "",
      },
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
    removeIndex: function (array, index) {
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
        url: this.apiRoot + "/config",
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
        url: this.apiRoot + "/config/testEmail",
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
        url: this.apiRoot + "/config",
        method: "GET",
      }).then(function (response) {
        if (response.status === 200) {
          window.console.log("Config from server: ", response.data)
          this.conf = response.data
          if (response.data.ip_whitelist === null) {
            this.conf.ip_whitelist = []
          }
          this.smtpUpdate.smtp_port = this.conf.smtp_server.smtp_port
        } else {
          window.console.log("Response code ", response.status)
        }
      }.bind(this))
    },
    addIp: function () {
      window.console.log("Adding ip: ", this.whitelistNewIp.ip_address)
      this.whitelistNewIp.delete = false
      axios({
        url: this.apiRoot + "/config/whitelist",
        method: "PATCH",
        data: this.whitelistNewIp
      }).then(function (response) {
        window.console.log(response.data)
        if (response.status === 200) {
          if (response.data.error === "") {
            this.conf.ip_whitelist.push(this.whitelistNewIp.ip_address)
            this.whitelistNewIp.ip_address = ""
          } else {
            this.alert = "Error adding ip: " + response.data.error
            this.showAlert("danger")
          }
        } else {
          this.alert = "Response code: " + response.status
          this.showAlert("danger")
        }
      }.bind(this))
    },
    deleteIp: function (ip, index) {
      window.console.log("Deleting ip: ", ip)
      if (!confirm("Do you really want to remove " + ip + " from the whitelist?")) {
        return
      }
      axios({
        url: this.apiRoot + "/config/whitelist",
        method: "PATCH",
        data: {delete: true, ip_address: ip}
      }).then(function (response) {
        if (response.status === 200) {
          if (response.data.error === "") {
            this.conf.ip_whitelist = this.removeIndex(this.conf.ip_whitelist, index)
            this.whitelistNewIp.ip_address = ""
          } else {
            this.alert = "Error adding ip: " + response.data.error
            this.showAlert("danger")
          }
        } else {
          this.alert = "Response code: " + response.status
          this.showAlert("danger")
        }
      }.bind(this))
    }
  }
}
</script>

<style scoped>
form {
  margin: 15px 0 15px 0;
}

.description {
  color: #7e7e7e;
  font-size: 14px;
}
</style>