<template>
  <div>
    <Navbar/>
    <div class="custom-container">
      <b-row class="custom-row">
        <b-col class="sidebar">
          <h1 id="sidebar-header">Header</h1>
          <div class="text-center">
            <a class="nav-link" aria-controls="Dashboard" aria-selected="true">
              <b-icon icon="layout-text-window-reverse"></b-icon>
              Dashboard</a>
            <a class="nav-link" aria-controls="Logs" aria-selected="false">
              <b-icon icon="journals"></b-icon>
              Logs</a>
          </div>
        </b-col>
        <b-col md="10" class="dashboard-container">
          <line-chart v-if="loaded" :chartdata="chartdata" :options="options"></line-chart>
          <button @click="filldata()">Randomize</button>
        </b-col>
      </b-row>
    </div>
    <Footer/>
  </div>
</template>

<script>
import Navbar from "../components/Navbar";
import Footer from "../components/Footer";
import { Line } from "vue-chartjs";

export default {
  name: "Dashboard",
  components: {Navbar, Footer},
  extends: Line,
  props: {
    chartdata: {
      type: Object,
      default: null,
    },
    options: {
      type: Object,
      default: null
    }
  },
  mounted() {
    this.renderChart(this.chartdata, this.options)
  },
  methods: {
    fillData() {
      this.datacollection = {
        labels: [this.getRandomInt(), this.getRandomInt()],
        datasets: [
          {
            label: 'Data One',
            backgroundColor: '#f87979',
            data: [this.getRandomInt(), this.getRandomInt()]
          }, {
            label: 'Data One',
            backgroundColor: '#f87979',
            data: [this.getRandomInt(), this.getRandomInt()]
          }
        ]
      }
    },
    getRandomInt() {
      return Math.floor(Math.random() * (50 - 5 + 1)) + 5
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

#sidebar-header {
  text-align: center;
  color: black;
  padding: 5px;
  font-family: Poppins, sans-serif;
  font-weight: normal;
  font-size: 25px;
}

.sidebar {
  border: 2px solid #eeefff;
  border-radius: 5px;
  width: 250px !important;
  max-width: 250px !important;
  height: 100%;
  background-color: #eeeeee;
}

.dashboard-container {
  border-radius: 5px;
  background-color: #eeeeee;
  margin-left: 1%;
}

.dashboard-card {

}

.dashboard-card-body {

}

p.large {
  font-family: Poppins, sans-serif;
  font-weight: lighter;
  font-size: x-large;
  padding: 5px;
  margin-bottom: 5px;
}

p.header {
  font-size: large;
  font-weight: normal;
  padding: 5px;
  margin-bottom: 5px;
}


.custom-row {
  height: 100% !important;
}

.custom-container {
  height: calc(100vh - 116px);
  padding: 10px;
}

</style>