<template>
  <div>
    <Navbar/>
    <apexchart width="500" type="bar" :options="options" :series="series"></apexchart>
    test test test
    <Footer/>
  </div>
</template>

<script>
  import Navbar from "../components/Navbar";
  import Footer from "../components/Footer";
  export default{
    name: "Home",
    components: {Footer, Navbar},
    data: function(){
      return{
        options: {
          chart: {
            id: 'vuechart-example'
          },
          xaxis: {
            categories: [1991, 1992, 1994, 1995, 1996, 1997, 1998]
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