<!-- TODO Add configure device modal -->
<template>
  <div>
    <Navbar></Navbar>
    <b-modal id="configure-honeypot" size="lg" hide-footer>
      <template #modal-title>
        Configure {{ deviceToConfigure }}
      </template>
      <b-form @submit.prevent="submitConfiguration" class="container" style="height: fit-content">
        <b-form-row>
          <div style="margin: auto;">
            <b-button type="submit" class="carousel-button">Submit</b-button>
          </div>
        </b-form-row>
        <template v-if="loading">
          <b-form-row>
            <b-col class="text-center">
              <div class="lds-ellipsis"><div></div><div></div><div></div><div></div></div>
            </b-col>
          </b-form-row>
        </template>
        <b-form-row>
          <b-col class="input">
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
        </b-form-row>
      </b-form>
    </b-modal>
    <div class="custom-container">
      <b-col md="12" class="content">
        <table style="margin: auto;">
          <tr>
            <th>Id</th>
            <th>Ip address</th>
            <th>Last seen</th>
            <th>Configured</th>
            <th class="text-center">Services</th>
            <th>Configure Device</th>
            <th>Delete</th>
          </tr>
          <template v-if="devices.length !== 0">
            <tr v-for="device in devices" :key="device.device_id">
              <td>{{device.device_id}}</td>
              <td>{{device.ip_str}}</td>
              <td>{{device.last_seen}}</td>
              <td>{{device.configured}}</td>
              <template v-if="device.configured">
                <td>
                  <div v-for="(service, index) in device.services" :key="index">
                    <div v-if="service" class="text-center" style="margin: auto">
                      {{index}}
                    </div>
                  </div>
                </td>
              </template>
              <template v-else>
                <td>Device not yet configured</td>
              </template>
              <td>
                <div style="margin: auto; width: fit-content">
                  <b-button v-on:click="setDeviceToConfigure(device.device_id)">Configure</b-button>
                </div>
              </td>
              <td class="text-center">
                <b-icon-trash class="click-icon" v-on:click="removeDevice(device.device_id)" variant="danger"></b-icon-trash>
              </td>
            </tr>
          </template>
        </table>
        <b-col md="12" v-if="devices.length === 0">
          <p class="text-center">No Devices has contacted the C2</p>
        </b-col>
      </b-col>
    </div>
    <Footer></Footer>
  </div>
</template>

<script>
import Navbar from "../components/Navbar";
import Footer from "../components/Footer";
import axios from "axios";
import moment from "moment"
export default {
  name: "Honeypots",
  components: {Navbar, Footer},
  data: function() {
    return {
      devices: [],
      deviceToConfigure: null,
      loading: false,
      dismissCountDown: 0,
      dismissSecs: 3,
      alert: "",
      variant: "",
    }
  },
  created() {
    axios.defaults.headers.common['Authorization'] = 'Bearer '+ this.$cookies.get("token")
    this.getDevices()
  },
  mounted() {
    console.log("Starting connection to WebSocket Server")
    let loc = window.location, new_uri;
    if(loc.protocol === "https:"){
      new_uri = "wss:"
    } else {
      new_uri = "ws:"
    }
    new_uri += "//" + loc.hostname + ":8000/ws";
    window.console.log("Trying to connect to ws on: " + new_uri)
    this.connection = new WebSocket(new_uri)

    this.connection.onmessage = function(event) {
      let data = JSON.parse(event.data)
      window.console.log(data)
      if (data.type == 2) {
        window.console.log("Recieved heartbeat event")
        this.updateDevice(data.device_id)
      } else if (data.type == 3) {
        window.console.log("New device registered, updating device list")
        this.getDevices()
      }
    }.bind(this)

    this.connection.onopen = function() {
      console.log("Successfully connected to the echo websocket server...")
    }
  },
  methods: {
    showAlert: function (variant) {
      this.variant = variant
      this.dismissCountDown = this.dismissSecs
    },
    countDownChanged: function (dismissCountDown) {
      this.dismissCountDown = dismissCountDown
    },
    setDeviceToConfigure: function(device_id){
      this.deviceToConfigure = device_id
      this.$bvModal.show('configure-honeypot')
    },
    submitConfiguration: function(){
      window.console.log("Submitting config")
    },
    removeItem: function(array, key, value) {
      const index = array.findIndex(obj => obj[key] === value)
      return index >= 0 ? [
        ...array.slice(0, index),
        ...array.slice(index+1)
      ] : array;
    },
    updateDevice: function (device_id) {
      let devices = []
      axios({
        url: process.env.VUE_APP_API_ROOT+"/devices",
        method: "GET",
      }).then(function (response){
        if (response.data.error == null) {
          devices = response.data
          devices.forEach(function (device){
            let date = new moment.utc(device.last_seen).format('dddd YYYY-MM-DD, HH:mm:ss[Z]')
            device.last_seen = date
          })
          const index = devices.findIndex(obj => obj["device_id"] === device_id)
          window.console.log("Index of updated device: "+index)
          this.$set(this.devices, index, devices[index])
          window.console.log(this.devices)
        }
      }.bind(this))
    },
    getDevices: function(){
      window.console.log("Getting devices")
      //this.devices = []
      axios({
        url: process.env.VUE_APP_API_ROOT+"/devices",
        method: "GET",
      }).then(function (response){
        if (response.data.error == null) {
          this.devices = response.data
          this.devices.forEach(function (device){
            let date = new moment.utc(device.last_seen).format('dddd YYYY-MM-DD, HH:mm:ss[Z]')
            device.last_seen = date
          })
          window.console.log("Succesfully got devices")
          window.console.log(this.devices)
        }
      }.bind(this))
    },
    removeDevice: function(deviceID) {
      if (confirm("Are you sure you want to delete device with id?: " + deviceID)) {
        let device_id = {device_id: deviceID}
        axios({
          url: process.env.VUE_APP_API_ROOT+"/devices",
          method: "DELETE",
          data: device_id
        }).then(function (response){
          if (response.data.error == "") {
            this.devices = this.removeItem(this.devices, "device_id", deviceID)
          } else {
            window.console.log("Error deleting device: "+response.data.error)
          }
        }.bind(this))
      } else {
        window.console.log("Cancelled")
      }
    }
  }
}
</script>

<style scoped>
  .content {
    box-shadow: 1px 6px 16px -5px #888888;
    border-radius: 10px;
    height: calc(100vh - 150px);
    overflow-y: auto;
  }
  .custom-container {
    height: calc(100vh - 116px);
    padding: 10px;
  }
  table, th, td {
    padding: 5px 15px 5px 15px;
  }
  .enabled {
    color: green;
  }
  .disabled{
    color: red;
  }
</style>