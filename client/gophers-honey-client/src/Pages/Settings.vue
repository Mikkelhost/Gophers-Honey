<!-- TODO: Add/remove image functionality-->
<template>
  <div>
    <Navbar/>
      <div class="container-xl">
        <div class="text-center">
          <h1>Settings</h1>
        </div>
        <b-row>
          <div class="settings-list">
            <div class="nav flex-nowrap nav-pills sticky-top" id="settings-list" role="tablist">
              <a class="nav-link active show" id="profile-tab" data-toggle="pill" href="#profile" aria-controls="profile" aria-selected="true">Profile</a>
              <a class="nav-link" id="users-tab" hred="#users" data-toggle="pill" href="#users" aria-controls="users" aria-selected="false">Users</a>
              <a class="nav-link" id="images-tab" hred="#images" data-toggle="pill" href="#images" aria-controls="images" aria-selected="false">Images</a>
              <a class="nav-link" id="service-tab" hred="#service" data-toggle="pill" href="#service" aria-controls="service" aria-selected="false">Service</a>
            </div>
          </div>
          <b-col md="10" class="settings-content">
            <div class="tab-content" id="v-pills-tabContent" style="height: 100%">
              <div id="profile" aria-labelledby="profile-tab" class="tab-pane fade active show" role="tabpanel">
                <profile></profile>
              </div>
              <div id="users" aria-labelledby="users-tab" class="tab-pane fade" role="tabpanel">
                <users></users>
              </div>
              <div id="images" aria-labelledby="images-tab" class="tab-pane fade" role="tabpanel">
                <images></images>
              </div>
              <div id="service" aria-labelledby="service-tab" class="tab-pane fade" role="tabpanel" style="height: 100%">
                <service></service>
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
  import Images from "../components/Images";
  import Profile from "../components/Profile";
  import Users from "../components/Users";
  import Service from "../components/Service";
  import axios from "axios";

  export default{
    name: "Settings",
    components: {Footer, Navbar, Images, Profile, Users, Service},
    created() {
      axios.defaults.headers.common['Authorization'] = 'Bearer '+ this.$cookies.get("token")
    },
  }
</script>

<style>
.nav-pills .nav-link {
  color: black;
  text-overflow: ellipsis;
  margin: 0 15px 0 15px;
  white-space: nowrap;
  width: 165px;
  border-radius: 2rem!important;
}
.nav-link.active{
  color: black!important;
  background-color: rgba(204, 200, 200, 0.46) !important;
}
.settings-list{
  padding: 10px 0 10px 0;
  margin: auto;
}

.settings-content {
  border-radius: 10px;
  box-shadow: 1px 6px 16px -5px #888888;
  padding: 10px 0 10px 0;
  height: 500px;
  margin: auto;
}
.container-xl{
  height: calc(100vh - 116px);
  width: 70%;
}

.click-icon {
  cursor: pointer;
}

table, th, td {
  padding: 5px 15px 5px 15px;
}

#add-image {
  margin: auto;
  border-radius: 2px;
  color: white;
  background-color: olivedrab;
}

 #add-user {
   margin: auto;
   border-radius: 2px;
   color: white;
   background-color: olivedrab;
 }

</style>