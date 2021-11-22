<!-- TODO: Add/remove image functionality-->
<template>
  <div>
    <Navbar/>
      <div class="container-xl">
        <div class="text-center">
          <h1>Settings</h1>
        </div>
        <b-row style="margin: auto;">
          <b-col class="settings-list">
            <div class="nav flex-column nav-pills sticky-top" id="settings-list" role="tablist" aria-orientation="vertical">
              <a class="nav-link active show" id="profile-tab" data-toggle="pill" href="#profile" aria-controls="profile" aria-selected="true">Profile</a>
              <a class="nav-link" id="users-tab" hred="#users" data-toggle="pill" href="#users" aria-controls="users" aria-selected="false">Users</a>
              <a class="nav-link" id="images-tab" hred="#images" data-toggle="pill" href="#images" aria-controls="images" aria-selected="false">Images</a>
              <a class="nav-link" id="service-tab" hred="#service" data-toggle="pill" href="#service" aria-controls="service" aria-selected="false">Service</a>
            </div>
          </b-col>
          <b-col md="8" class="settings-content">
            <div class="tab-content" id="v-pills-tabContent">
              <div id="profile" aria-labelledby="profile-tab" class="tab-pane fade active show" role="tabpanel">
                <profile></profile>
              </div>
              <div id="users" aria-labelledby="users-tab" class="tab-pane fade" role="tabpanel">
                <users></users>
              </div>
              <div id="images" aria-labelledby="images-tab" class="tab-pane fade" role="tabpanel">
                <images></images>
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
  import axios from "axios";

  export default{
    name: "Settings",
    components: {Footer, Navbar, Images, Profile, Users},
    created() {
      axios.defaults.headers.common['Authorization'] = 'Bearer '+ this.$cookies.get("token")
    },
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
  height: 500px;
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