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
                <b-row>
                  <b-col md="12">
                    <table style="margin: auto;">
                      <tr>
                        <th>Id</th>
                        <th>Name</th>
                        <th>Created</th>
                        <th>Download</th>
                        <th>Delete</th>
                      </tr>
                      <tr v-for="image in images" :key="image.id">
                        <td>{{image.id}}</td>
                        <td>{{image.name}}</td>
                        <td>{{image.date_created}}</td>
                        <td class="text-center">
                          <b-icon-download v-on:click="downloadImage(image)" class="click-icon"></b-icon-download>
                          <progressbar size="medium" :val="image.downloadPercentage" :text="image.downloadPercentage+'%'"></progressbar>
                        </td>
                        <td class="text-center">
                          <b-icon-trash class="click-icon" v-on:click="deleteImage(image.id)" variant="danger"></b-icon-trash>
                        </td>
                      </tr>
                    </table>
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
      downloadImage: function(image) {
        window.console.log(this.images)
        axios({
          url: process.env.VUE_APP_API_ROOT+"/images?download="+image.id,
          method: 'GET',
          responseType: 'blob',
          onDownloadProgress: function (event) {
            image.downloadPercentage = parseInt( Math.round((event.loaded / event.total)*100))
          }.bind(this)
        }).then(function (response){
          window.console.log(response.data)
          const url = window.URL.createObjectURL(new Blob([response.data]));
          const link = document.createElement('a');
          link.href = url;
          link.setAttribute('download', 'raspberrypi.img'); //or any other extension
          document.body.appendChild(link);
          link.click();
        })
      },
      deleteImage: function(imageId) {
        window.console.log("deleting image"+imageId)
      },
      getImages: function() {
        let that = this
        this.images = []
        axios.get(process.env.VUE_APP_API_ROOT+"/images/getImages").then(function(response){
          if (response.status === 200) {
            response.data.forEach(function(image){
              let img = {id: image.image_id, name: image.name, date_created: image.date_created, downloadPercentage: 0}
              that.images.push(img)
            })
            window.console.log(that.images)
          }
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
  height: 500px;
}
.container{
  height: calc(100vh - 116px);
  width: 70%;
}

.click-icon {
  cursor: pointer;
}

table, th, td {
  padding: 5px 15px 5px 15px;
}

</style>