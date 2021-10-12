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
            <th>Services</th>
            <th>Delete</th>
          </tr>
          <template v-if="devices.length !== 0">
            <tr v-for="device in devices" :key="device.id">
              <td>{{device.device_id}}</td>
              <td>{{device.ip_str}}</td>
              <td>{{device.last_seen}}</td>
              <td>{{device.configured}}</td>
              <template v-if="device.configured">
                <td>

                </td>
              </template>
              <template v-else>
                <td>Device not yet configured</td>
              </template>

              <td class="text-center">
                <b-icon-trash class="click-icon" v-on:click="deleteImage(device.id)" variant="danger"></b-icon-trash>
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
        }
      }.bind(this))
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
</style>