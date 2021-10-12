<template>
  <div>
    <Navbar></Navbar>
    <div class="custom-container">
      <b-col md="12" class="content">
        <table style="margin: auto;">
          <tr>
            <th>Id</th>
            <th>Ip address</th>
            <th>Last seen</th>
            <th>Configured</th>
            <th class="text-center">Services</th>
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
export default {
  name: "Honeypots",
  components: {Navbar, Footer},
  data: function() {
    return {
      devices: [],
    }
  },
  created() {
    axios.defaults.headers.common['Authorization'] = 'Bearer '+ this.$cookies.get("token")
    this.getDevices()
  },
  methods: {
    removeItem: function(array, key, value) {
      const index = array.findIndex(obj => obj[key] === value)
      return index >= 0 ? [
        ...array.slice(0, index),
        ...array.slice(index+1)
      ] : array;
    },
    getDevices: function(){
      window.console.log("Getting devices")
      this.devices = []
      axios({
        url: process.env.VUE_APP_API_ROOT+"/devices",
        method: "GET",
      }).then(function (response){
        if (response.data.error == null) {
          this.devices = response.data
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