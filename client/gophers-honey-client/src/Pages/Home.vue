<template>
  <div>
    <Navbar/>
    <p>Hello</p>
    <Footer/>
  </div>
</template>

<script>
  import Navbar from "../components/Navbar";
  import Footer from "../components/Footer";
  export default{
    name: "Home",
    components: {Footer, Navbar},
  }
</script>