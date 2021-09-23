<template>
  <div>
    <Navbar/>
      <div class="container">
        <div class="text-center">
          <h1>Settings</h1>
        </div>
        <b-row>
          <b-col class="settings-list">
            <div class="nav flex-column nav-pills sticky-top" id="settings-list" role="tablist" aria-orientation="vertical">
              <a class="nav-link active show" id="profile-tab" data-toggle="pill" href="#profile" aria-controls="profile" aria-selected="true">Profile</a>
              <a class="nav-link" id="users-tab" hred="#users" data-toggle="pill" href="#users" aria-controls="users" aria-selected="false">Users</a>
              <a class="nav-link" id="images-tab" hred="#images" data-toggle="pill" href="#images" aria-controls="images" aria-selected="false">Images</a>
              <a class="nav-link" id="honeypots-tab" hred="#honeypots" data-toggle="pill" href="#honeypots" aria-controls="honeypots" aria-selected="false">Honeypots</a>
              <a class="nav-link" id="service-tab" hred="#service" data-toggle="pill" href="#service" aria-controls="service" aria-selected="false">Service</a>
            </div>
          </b-col>
          <b-col md="8" class="settings-content">
            <div class="tab-content" id="v-pills-tabContent">
              <div id="profile" aria-labelledby="profile-tab" class="tab-pane fade active show" role="tabpanel">
                <p>Profile</p>
              </div>
              <div id="images" aria-labelledby="images-tab" class="tab-pane fade" role="tabpanel">
                <p>Images</p>
                <b-row>
                  <b-col v-for="image in images" :key="image.image_id">
                    {{image.image_id}}
                    {{image.name}}
                  </b-col>
                </b-row>
              </div>
            </div>
          </b-col>
        </b-row>

      </div>
    <Footer/>
  </div>
</template>

<script>
  import Navbar from "../components/Navbar";
  import Footer from "../components/Footer";
  import axios from "axios";
  export default{
    name: "Settings",
    components: {Footer, Navbar},
    data: function(){
      return{
        images: []

      }
    },
    created() {
      axios.defaults.headers.common['Authorization'] = 'Bearer '+ this.$cookies.get("token")
      this.getImages()
    },
    methods: {
      getImages: function() {
        let that = this
        axios.get(process.env.VUE_APP_API_ROOT+"/images/getImages").then(function(response){
          if (response.status === 200) {
            that.images = response.data
          }
          window.console.log("response data" + response.status)
        })
      }
    }
  }
</script>

<style>
.nav-pills .nav-link {
  color: black;
  text-overflow: ellipsis;
  overflow: hidden;
  white-space: nowrap;
  width: 165px;
  border-radius: 2rem!important;
}
.nav-link.active{
  color: black!important;
  background-color: rgba(204, 200, 200, 0.46) !important;
}
.settings-list{
  width: 200px!important;
  max-width: 200px!important;
  padding: 10px 0 10px 0;
}
.settings-content {
  border-radius: 10px;
  box-shadow: 1px 6px 16px -5px #888888;
  padding: 10px 0 10px 0;
}
</style>