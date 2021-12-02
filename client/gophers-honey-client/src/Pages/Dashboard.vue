<template>
  <div>
    <Navbar/>
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
    <div class="custom-container">
      <b-row>
        <div class="nav-list">
          <div class="nav flex-nowrap nav-pills sticky-top" id="settings-list" role="tablist">
            <a class="nav-link active show" id="dashboard-tab" data-toggle="pill" href="#dashboard"
               aria-controls="dashboard" aria-selected="true">Dashboard</a>
            <a class="nav-link" id="logs-tab" data-toggle="pill" href="#logs" aria-controls="logs"
               aria-selected="false">Logs</a>
          </div>
        </div>
      </b-row>
      <div class="tab-content settings-content" id="v-pills-tabContent">
        <div id="dashboard" aria-labelledby="dashboard-tab" class="tab-pane fade active show" role="tabpanel">
          <b-row>
            <div style="margin: auto;">
              <b-button
                  style="font-size: 20px; border-width: 0; background-color: white!important; color: grey!important;"
                  v-on:click="refreshData">
                <b-icon-arrow-counterclockwise></b-icon-arrow-counterclockwise>
              </b-button>
            </div>
          </b-row>
          <b-row>
            <div class="chart col-md-4">
              <div class="chart-header">
                <h5 class="text-center">Service Distribution</h5>
                <h6 class="text-center">Total Decoys Deployed: {{ servicesData.totalDevices }}</h6>
              </div>
              <pie-chart v-if="loaded" :chart-data="servicesData" :options="options"></pie-chart>
            </div>
            <div class="chart col-md-4">
              <div class="chart-header">
                <h5 class="text-center">Most Touched Protocols</h5>
              </div>
              <doughnut-chart v-if="loaded" :chart-data="protocolsData" :options="options"></doughnut-chart>
            </div>
            <div class="chart col-md-4">
              <div class="chart-header">
                <h5 class="text-center">Logtype Distribution</h5>
                <h6 class="text-center">Total Log Entries: {{ logtypeData.totalLogs }}</h6>
              </div>
              <pie-chart v-if="loaded" :chart-data="logtypeData" :options="options"></pie-chart>
            </div>
          </b-row>
          <b-row style="height: 300px;">
            <b-col md="12" style="margin-top: 15px;">
              <h5 class="text-center">Incidents the last month</h5>
              <bar-chart v-if="loaded" :chart-data="incidentData" :options="lineOptions"></bar-chart>
            </b-col>
          </b-row>

        </div>
        <div id="logs" aria-labelledby="logs-tab" class="tab-pane fade" role="tabpanel">
        </div>
      </div>
    </div>
    <Footer/>
  </div>

</template>

<script>
import Navbar from "../components/Navbar";
import Footer from "../components/Footer";
import PieChart from "../components/pieChart";
import DoughnutChart from "../components/doughnutChart";
import BarChart from "../components/barChart";
import axios from "axios";

export default {
  name: "Dashboard",
  components: {Navbar, Footer, PieChart, DoughnutChart, BarChart},
  data() {
    return {
      loaded: false,
      dismissCountDown: 0,
      dismissSecs: 10,
      alert: "",
      variant: "",
      randomData: 0,
      incidentData: {
        labels: [],
        datasets: [{
          label: 'Incidents over time',
          data: [],
          fill: false,
          borderColor: '#2554FF',
          backgroundColor: 'rgba(0,0,255,1)',
          borderWidth: 1
        }]
      },
      protocolsData: {
        labels: [],
        datasets: [{
          data: [],
          fill: false,
          borderColor: '#ffffff',
          backgroundColor: [],
          borderWidth: 5
        }]
      },
      servicesData: {
        labels: [],
        datasets: [{
          data: [],
          fill: false,
          borderColor: '#ffffff',
          backgroundColor: [],
          borderWidth: 5
        }],
        totalDevices: 0
      },
      logtypeData: {
        labels: [],
        datasets: [{
          data: [],
          fill: false,
          borderColor: '#ffffff',
          backgroundColor: [],
          borderWidth: 5
        }],
        totalLogs: 0,
      },
      options: {
        legend: {
          display: true
        },
        responsive: true,
        maintainAspectRatio: false
      },
      lineOptions: {
        scales: {
          yAxes: [{
            ticks: {
              beginAtZero: true
            },
            gridLines: {
              display: true
            }
          }],
          xAxes: [{
            gridLines: {
              display: false
            }
          }]
        },
        legend: {
          display: true
        },
        responsive: true,
        maintainAspectRatio: false
      },
      logs: [],
      devices: [],
      protocolMap: {
        22: "SSH",
        23: "TELNET",
        21: "FTP",
        80: "HTTP",
        139: "SMB",
        445: "SMB"
      }
    }
  },
  async mounted() {
    axios.defaults.headers.common['Authorization'] = 'Bearer ' + this.$cookies.get("token")
    this.loaded = false
    await this.getLogs()
    this.generateProtocolsData()
    this.generateLogtypeData()
    this.generateIncidentData()
    await this.getDevices()
    this.generateDeviceData()

    this.loaded = true
  },
  methods: {
    refreshData: async function () {
      this.loaded = false
      await this.getLogs()
      this.generateProtocolsData()
      this.generateLogtypeData()
      this.generateIncidentData()
      await this.getDevices()
      this.generateDeviceData()
      this.loaded = true
    },
    showAlert: function (variant) {
      this.variant = variant
      this.dismissCountDown = this.dismissSecs
    },
    countDownChanged: function (dismissCountDown) {
      this.dismissCountDown = dismissCountDown
    },
    getLogs: function () {
      return axios({
        url: process.env.VUE_APP_API_ROOT + "/logs",
        method: "GET"
      }).then(function (response) {
        if (response.status === 200) {
          if (typeof response.data.error === "undefined") {
            this.logs = response.data
            //window.console.log(this.logs)
          } else {
            this.alert = "Error getting logs: " + response.data.error
            this.showAlert("danger")
          }
        } else {
          this.alert = "Api returned status code: " + response.status
          this.showAlert("danger")
        }
      }.bind(this))
    },
    getDevices: function () {
      return axios({
        url: process.env.VUE_APP_API_ROOT + "/devices",
        method: "GET"
      }).then(function (response) {
        if (response.status === 200) {
          if (typeof response.data.error === "undefined") {
            this.devices = response.data
            //window.console.log(this.devices)
          } else {
            this.alert = "Error getting devices: " + response.data.error
            this.showAlert("danger")
          }
        } else {
          this.alert = "Api returned status code: " + response.status
          this.showAlert("danger")
        }
      }.bind(this))
    },
    generateColor: function (length) {
      let colorArray = []
      let spacing = Math.floor(255 / length)
      for (let i = 1; i <= length; i++) {
        colorArray.push('rgba(0,120,' + spacing * i + ',1)')
      }
      //console.log(colorArray)
      return colorArray
    },
    generateProtocolsData: function () {
      // ProtocolData
      this.protocolsData.labels = []
      this.protocolsData.datasets = []
      let protocolData = [
        {label: "SSH", count: 0},
        {label: "TELNET", count: 0},
        {label: "FTP", count: 0},
        {label: "HTTP", count: 0},
        {label: "SMB", count: 0}
      ]
      this.logs.forEach(function (log) {
        if (log.level === 0) {
          let index = protocolData.findIndex(obj => obj["label"] === this.protocolMap[log.dst_port])
          protocolData[index].count = protocolData[index].count + 1
        }
      }.bind(this))
      let labels = []
      let datasets = [{
        data: [],
        fill: false,
        borderColor: '#ffffff',
        backgroundColor: [],
        borderWidth: 5
      }]
      protocolData.forEach(function (protocol) {
        labels.push(protocol.label)
        datasets[0].data.push(protocol.count)
      }.bind(this))
      datasets[0].backgroundColor = this.generateColor(protocolData.length)
      this.$set(this.protocolsData, 'labels', labels)
      this.$set(this.protocolsData, 'datasets', datasets)
      //console.log(this.protocolsData)
    },
    generateLogtypeData: function () {
      this.logtypeData.totalLogs = 0
      this.logtypeData.labels = []
      this.logtypeData.datasets = []
      let logtypeData = [
        {label: "Portscan", count: 0},
        {label: "Critical", count: 0}
      ]
      this.logs.forEach(function (log) {
        if (log.level === 0) {
          let index = logtypeData.findIndex(obj => obj["label"] === "Critical")
          logtypeData[index].count = logtypeData[index].count + 1
          this.logtypeData.totalLogs += 1
        } else if (log.level === 1) {
          let index = logtypeData.findIndex(obj => obj["label"] === "Portscan")
          logtypeData[index].count = logtypeData[index].count + 1
          this.logtypeData.totalLogs += 1
        }
      }.bind(this))
      //console.log("logtypeData: ", logtypeData)

      let labels = []
      let datasets = [{
        data: [],
        fill: false,
        borderColor: '#ffffff',
        backgroundColor: [],
        borderWidth: 5
      }]
      logtypeData.forEach(function (protocol) {
        labels.push(protocol.label)
        datasets[0].data.push(protocol.count)
      }.bind(this))
      datasets[0].backgroundColor = this.generateColor(logtypeData.length)
      this.$set(this.logtypeData, 'labels', labels)
      this.$set(this.logtypeData, 'datasets', datasets)
      //console.log(this.logtypeData)
    },
    generateDeviceData: function () {
      this.servicesData.totalDevices = this.devices.length
      this.servicesData.labels = []
      this.servicesData.datasets = []
      let serviceData = [
        {label: "SSH", count: 0},
        {label: "TELNET", count: 0},
        {label: "FTP", count: 0},
        {label: "HTTP", count: 0},
        {label: "SMB", count: 0}
      ]
      this.devices.forEach(function (device) {
        //console.log(typeof device.services)
        for (const [key, value] of Object.entries(device.services)) {
          if (value) {
            let index = serviceData.findIndex(obj => obj["label"] === key.toUpperCase())
            serviceData[index].count = serviceData[index].count + 1
          }
        }
      }.bind(this))
      //console.log(serviceData)

      let labels = []
      let datasets = [{
        data: [],
        fill: false,
        borderColor: '#ffffff',
        backgroundColor: [],
        borderWidth: 5
      }]
      serviceData.forEach(function (service) {
        labels.push(service.label)
        datasets[0].data.push(service.count)
      }.bind(this))
      datasets[0].backgroundColor = this.generateColor(serviceData.length)
      this.$set(this.servicesData, 'labels', labels)
      this.$set(this.servicesData, 'datasets', datasets)
      //console.log(this.servicesData)
    },
    getDaysArray: function (start, end) {
      for (var arr = [], dt = new Date(start); dt <= end; dt.setDate(dt.getDate() + 1)) {
        dt.setHours(0, 0, 0, 0)
        arr.push(new Date(dt));
      }
      return arr;
    },
    parseDMY: function (value) {
      var date = value.split("/");
      var d = parseInt(date[0], 10),
          m = parseInt(date[1], 10),
          y = parseInt(date[2], 10);
      return new Date(y, m - 1, d);
    },
    generateIncidentData: function () {
      var today = new Date()
      today.setHours(0, 0, 0, 0)
      var priorDate = new Date().setDate(today.getDate() - 30)

      let getDaysArray = this.getDaysArray(priorDate, today)
      let incidentData = []
      getDaysArray.forEach(function (day) {
        let options = {year: 'numeric', month: 'numeric', day: 'numeric'}
        let dayData = {
          label: day.toLocaleString('en-GB', options),
          count: 0
        }
        incidentData.push(dayData)
      })
      //window.console.log(incidentData)

      this.logs.forEach(function (log) {
        if (log.level === 0 || log.level === 1) {
          incidentData.forEach(function (data) {
            let date = this.parseDMY(data.label)
            let nextDay = new Date(date)
            nextDay.setDate(nextDay.getDate() + 1)
            let logDate = new Date(log.log_time_stamp)
            //console.log(new Date(log.log_time_stamp))
            if (nextDay > logDate && logDate > date) {
              data.count += 1
            }
          }.bind(this))
        }
      }.bind(this))
      //console.log(incidentData)

      let labels = []
      let datasets = [{
        label: "Incidents over time",
        data: [],
        fill: false,
        borderColor: '#ffffff',
        backgroundColor: 'blue',
        borderWidth: 5
      }]
      incidentData.forEach(function (incident) {
        labels.push(incident.label)
        datasets[0].data.push(incident.count)
      }.bind(this))
      this.$set(this.incidentData, 'labels', labels)
      this.$set(this.incidentData, 'datasets', datasets)


    }
  }
}
</script>

<style>

.nav-link {
  color: black;
  text-align: center;
  text-overflow: ellipsis;
  border-radius: 2rem !important;
}

.nav-link:hover {
  color: black !important;
  background-color: rgba(204, 200, 200, 0.46) !important;
}

.nav-list {
  padding: 10px 0 10px 0;
  margin: auto;
}

.tab-content {
  height: calc(100vh - 176px);
  overflow-y: auto;
  overflow-x: hidden;
}

body {
  overflow-y: hidden;
}

.chart {
}

.chart-header {
  height: 52px;
}

.custom-container {
  width: 95vw;
  margin: auto;
}

#bar-chart {
  height: 300px!important;
  width: 1100px!important;
  margin: auto;
}
</style>