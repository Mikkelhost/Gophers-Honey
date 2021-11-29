<template>
  <div>
    <Navbar/>
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
            <div class="chart col-md-4">
              <div class="chart-header">
                <h5 class="text-center">Service Distribution</h5>
                <h6 class="text-center">Total Decoys Deployed: </h6>
              </div>
              <pie-chart :chart-data="chartData" :options="options"></pie-chart>
            </div>
            <div class="chart col-md-4">
              <div class="chart-header">
                <h5 class="text-center">Most Touched Protocols</h5>
              </div>
              <doughnut-chart :chart-data="chartData" :options="options"></doughnut-chart>
            </div>
            <div class="chart col-md-4">
              <div class="chart-header">
                <h5 class="text-center">Logtype Distribution</h5>
                <h6 class="text-center">Total Log Entries: </h6>
              </div>
              <pie-chart :chart-data="chartData" :options="options"></pie-chart>
            </div>
          </b-row>
          <b-row>
            <b-col md="12">
              <line-chart :chart-data="lineData" :options="lineOptions"></line-chart>
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
import LineChart from "../components/lineChart";

export default {
  name: "Dashboard",
  components: {Navbar, Footer, PieChart, DoughnutChart, LineChart},
  data() {
    return {
      randomData: 0,
      lineData: {
        labels: ["Babol",	"Cabanatuan",	"Daegu",	"Jerusalem",	"Fairfield",	"New York",	"Gangtok", "Buenos Aires", "Hafar Al-Batin", "Idlib"],
        datasets: [{
          label: 'Line Chart',
          data: [600, 1150, 342, 6050, 2522, 3241, 1259, 157, 1545, 9841],
          fill: false,
          borderColor: '#2554FF',
          backgroundColor: '#2554FF',
          borderWidth: 1
        }]

      },
      chartData: {
        labels: ["Italy", "India", "Japan", "USA",],
        datasets: [{
          data: [1000, 500, 1500, 1000],
          fill: false,
          borderColor: '#ffffff',
          backgroundColor: '#2554FF',
          borderWidth: 5
        }]
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
      }
    }
  },
  mounted() {
  },
  methods: {
    generateNumber: function () {
      window.console.log("Generating number")
      this.randomData = Math.floor(Math.random() * 100)
      window.console.log("Number is: ", this.randomData)
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
</style>