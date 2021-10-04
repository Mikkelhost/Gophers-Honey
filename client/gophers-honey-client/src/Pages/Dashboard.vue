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
            <a class="nav-link" aria-controls="Devices" aria-selected="false">
              <b-icon icon="cpu"></b-icon>
              Devices</a>
            <a class="nav-link" aria-controls="Logs" aria-selected="false">
              <b-icon icon="journals"></b-icon>
              Logs</a>
            <a class="nav-link" aria-controls="Users" aria-selected="false">
              <b-icon icon="people"></b-icon>
              Users</a>
          </div>
        </b-col>
          <b-col md="10" class="dashboard-container">
            <b-row class="custom-row">
              <b-col class="dashboard-card">
                <b-col>
                  <p class="header">Dashboard</p>
                  <p class="large">Alerts over time</p>
                </b-col>
                <b-col class="dashboard-card-body">
                  <apexchart type="bar" height="250px" :options="options" :series="series"></apexchart>
                </b-col>
              </b-col>
            </b-row>
          </b-col>
      </b-row>
    </div>
    <Footer/>
  </div>
</template>

<script>
import Navbar from "../components/Navbar";
import Footer from "../components/Footer";

export default {
  name: "Dashboard",
  components: {Navbar, Footer},
  data: function () {
    return {
      options: {
        chart: {
          id: "vuechart-example"
        },
        xaxis: {
          categories: [1991, 1992, 1993, 1994, 1995, 1996, 1997, 1998]
        }
      },
      series: [{
        name: 'series-1',
        data: [30, 40, 45, 50, 49, 60, 70, 91]
      }]
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