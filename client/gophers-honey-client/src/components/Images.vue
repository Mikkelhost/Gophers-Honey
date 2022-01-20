<template>
  <div>
    <b-modal id="new-image" size="md" hide-footer>
      <template #modal-title>
        Create a new image
      </template>
      <b-form @submit.prevent="submitImage" class="container" style="height: fit-content">
        <b-form-row>
          <b-col class="input">
            <b-form-group
                id="input-group-7"
                label="Imagename*"
                label-for="input-7"
            >
              <b-form-input
                  id="input-7"
                  v-model="imageInfo.name"
                  type="text"
                  placeholder="Imagename"
                  required
              >
              </b-form-input>
            </b-form-group>
          </b-col>
          <b-col class="input">
            <b-form-group
                id="input-group-8"
                label="C2 Protocol*"
                label-for="input-8"
                description="The C2 protocol can be either https or http, for https please ensure that certificates has been set up correctly"
            >
              <b-form-radio-group
                  id="input-8"
                  v-model="imageInfo.c2_protocol"
                  :options="protocolOptions"
              >
              </b-form-radio-group>
            </b-form-group>
          </b-col>
        </b-form-row>
        <b-form-row>
          <b-col style="margin-top: -20px" class="input">
            <b-form-group
                id="input-group-8"
                label="C2 Hostname*"
                label-for="input-8"
                description="The C2 hostname is the url for your the api"
            >
              <b-form-input
                  id="input-8"
                  v-model="imageInfo.c2"
                  type="text"
                  placeholder="C2 Hostname"
                  required
              >
              </b-form-input>
            </b-form-group>
          </b-col>
          <b-col style="margin-top: -20px" class="input">
            <b-form-group
                id="input-group-9"
                label="Port*"
                label-for="input-9"
                description="The api port"
            >
              <b-form-input
                  id="input-9"
                  v-model="imageInfo.port"
                  type="number"
                  min="0"
                  placeholder="Port"
                  required
              >
              </b-form-input>
            </b-form-group>
          </b-col>
        </b-form-row>
        <b-form-row>
          <div style="margin: auto;">
            <b-button type="submit" class="carousel-button">Submit</b-button>
          </div>
        </b-form-row>
        <template v-if="loading">
          <b-form-row>
            <b-col class="text-center">
              <div class="lds-ellipsis"><div></div><div></div><div></div><div></div></div>
            </b-col>
          </b-form-row>
        </template>
        <b-form-row>
          <b-col class="input">
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
          </b-col>
        </b-form-row>
      </b-form>
    </b-modal>
    <b-row>
      <b-button id="add-image" @click="$bvModal.show('new-image')">
        <b-icon-plus></b-icon-plus>
        Add Image
      </b-button>
    </b-row>
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
          <template v-if="images.length !== 0">
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
          </template>
        </table>
        <b-col md="12" v-if="images.length === 0">
          <p class="text-center">No images available</p>
        </b-col>
      </b-col>
    </b-row>
  </div>
</template>

<script>
import getEnv from '../utils/env'
import axios from "axios";
export default {
  name: "Images",
  data: function(){
    return{
      apiRoot: getEnv('VUE_APP_API_ROOT'),
      loading: false,
      dismissCountDown: 0,
      dismissSecs: 3,
      alert: "",
      variant: "",
      images: [],
      protocolOptions: [
        { text: 'HTTPS', value: 'https' },
        { text: 'HTTP', value: 'http' },
      ],
      imageInfo: {
        name: "",
        c2: "",
        c2_protocol: "https",
        port: null,
      }
    }
  },
  created() {
    axios.defaults.headers.common['Authorization'] = 'Bearer '+ this.$cookies.get("token")
    this.getImages()
  },
  methods: {
    showAlert: function (variant) {
      this.variant = variant
      this.dismissCountDown = this.dismissSecs
    },
    countDownChanged: function (dismissCountDown) {
      this.dismissCountDown = dismissCountDown
    },
    removeItem: function(array, key, value) {
      const index = array.findIndex(obj => obj[key] === value)
      return index >= 0 ? [
        ...array.slice(0, index),
        ...array.slice(index+1)
      ] : array;
    },
    submitImage: function() {
      let that = this
      this.loading = true
      axios({
        url: this.apiRoot+"/api/images",
        method: 'POST',
        data: that.imageInfo
      }).then(function (response){
        if (response.data.error === "") {
          this.loading = false
          this.getImages()
          this.$bvModal.hide('new-image')
        } else {
          this.alert = response.data.error
          this.showAlert("danger")
        }

      }.bind(this))
    },
    downloadImage: function(image) {
      //window.console.log(this.images)
      axios({
        url: this.apiRoot+"/api/images?download="+image.id,
        method: 'GET',
        responseType: 'blob',
        onDownloadProgress: function (event) {
          image.downloadPercentage = parseInt(Math.round((event.loaded / event.total)*100))
        }.bind(this)
      }).then(function (response){
        //window.console.log(response.data)
        const url = window.URL.createObjectURL(new Blob([response.data]));
        const link = document.createElement('a');
        link.href = url;
        link.setAttribute('download', 'raspberrypi.img'); //or any other extension
        document.body.appendChild(link);
        link.click();
      })
    },
    deleteImage: function(imageId) {
      if (confirm("Are you sure you want to delete image with id?: " + imageId)) {
        let image_id = {image_id: imageId}
        axios({
          url: this.apiRoot+"/api/images",
          method: "DELETE",
          data: image_id,
        }).then(function (response){
          if (response.data.error === "") {
            this.images = this.removeItem(this.images, "id", imageId)
          } else {
            window.console.log("Error deleting image: "+response.data.error)
          }
        }.bind(this))
      } else {
        window.console.log("Canceling deletion")
      }
    },
    getImages: function() {
      let that = this
      this.images = []
      axios.get(this.apiRoot+"/api/images").then(function(response){
        if (response.status === 200) {
          if (response.data.error !== "No images in DB") {
            response.data.forEach(function(image){
              let img = {id: image.image_id, name: image.name, date_created: image.date_created, downloadPercentage: 0}
              that.images.push(img)
            })
            window.console.log(that.images)
          } else {
            window.console.log(response.data)
          }

        }
      })
    }
  }
}

</script>


<style scoped>

</style>