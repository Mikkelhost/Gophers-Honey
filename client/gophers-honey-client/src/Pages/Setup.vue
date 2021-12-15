<template>
  <div class="container">
    <b-form @submit.prevent="submitSetup" class="container">
      <div class="text-center" style="margin-top: 20px">
        <h1>Gophers-Honey setup</h1>
        <p class="text-center">Welcome to the Gophers-Honey setup.<br>
          This is the first time the service has been run,<br>
          you will therefore have to create the first user account<br>
          and thereafter the first image for your honeypots.</p>
      </div>
      <b-carousel
          id="carousel"
          :interval="0"
          ref="setupCarousel"
          class="setup-carousel"
      >
        <!-- Text slides with image -->
        <b-carousel-slide class="carousel-height">
          <template slot="img" class="h-100">
            <div class="text-center" style="margin-top: 5px;">
              <h3>Create an administrative account</h3>
            </div>
            <b-row class="form-row">
              <b-col class="input">
                <b-form-group
                    id="input-group-1"
                    label="First name*"
                    label-for="input-1"
                >
                  <b-form-input
                      id="input-1"
                      v-model="form.userInfo.firstName"
                      type="text"
                      placeholder="First name"
                      required
                      @input.native="checkUserForm()"
                  >
                  </b-form-input>
                </b-form-group>
              </b-col>
              <b-col class="input">
                <b-form-group
                    id="input-group-2"
                    label="Last name*"
                    label-for="input-2"
                >
                  <b-form-input
                      id="input-2"
                      v-model="form.userInfo.lastName"
                      type="text"
                      placeholder="Last name"
                      required
                      @input.native="checkUserForm()"
                  >
                  </b-form-input>
                </b-form-group>
              </b-col>
            </b-row>
            <b-form-row>
              <b-col class="input">
                <b-form-group
                    id="input-group-3"
                    label="Email address*"
                    label-for="input-3"
                >
                  <b-form-input
                      id="input-3"
                      v-model="form.userInfo.email"
                      type="email"
                      placeholder="Email address"
                      required
                      @input.native="checkUserForm()"
                  >
                  </b-form-input>
                </b-form-group>
              </b-col>
            </b-form-row>
            <b-form-row>
              <b-col class="input">
                <b-form-group
                    id="input-group-4"
                    label="Username*"
                    label-for="input-4"
                >
                  <b-form-input
                      id="input-4"
                      v-model="form.userInfo.username"
                      type="text"
                      placeholder="Username"
                      required
                      @input.native="checkUserForm()"
                  >
                  </b-form-input>
                </b-form-group>
              </b-col>
            </b-form-row>
            <b-form-row>
              <b-col class="input">
                <b-form-group
                    id="input-group-5"
                    label="Password*"
                    label-for="input-5"
                >
                  <b-form-input
                      id="input-5"
                      v-model="form.userInfo.password"
                      type="password"
                      placeholder="Password"
                      required
                      @input.native="checkUserForm()"
                  >
                  </b-form-input>
                </b-form-group>
              </b-col>
              <b-col class="input">
                <b-form-group
                    id="input-group-6"
                    label="Confirm password*"
                    label-for="input-6"
                >
                  <b-form-input
                      id="input-6"
                      v-model="form.userInfo.confirmPw"
                      type="password"
                      placeholder="Confirm password"
                      required
                      @input.native="checkUserForm()"
                  >
                  </b-form-input>
                </b-form-group>
              </b-col>
            </b-form-row>
            <b-form-row>
              <b-button @click="$refs.setupCarousel.next()"
                        :disabled="!formValid"
                        style="margin: auto;"
                        class="carousel-button"
              >Next</b-button>
            </b-form-row>
          </template>
        </b-carousel-slide>


        <!-- Slides with custom text -->
        <b-carousel-slide>
          <template slot="img" class="h-100">
            <div class="text-center" style="margin-top: 5px;">
              <h3>Create your first image</h3>
              <p class="text-center">The following parameters will be put into your custom raspian image. <br>
              After image creation, you will be able to download the image<br> and use it for all honeypots
              that is used for this specific host.</p>
            </div>
            <b-form-row>
              <b-col class="input">
                <b-form-group
                    id="input-group-7"
                    label="Imagename*"
                    label-for="input-7"
                >
                  <b-form-input
                      id="input-7"
                      v-model="form.imageInfo.name"
                      type="text"
                      placeholder="Imagename"
                      required
                      @input.native="checkUserForm()"
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
                      v-model="form.imageInfo.c2_protocol"
                      :options="protocolOptions"
                      @input.native="checkUserForm()"
                  >
                  </b-form-radio-group>
                </b-form-group>
              </b-col>
            </b-form-row>
            <b-form-row>
              <b-col style="margin-top: -10px" class="input">
                <b-form-group
                    id="input-group-9"
                    label="C2 Hostname*"
                    label-for="input-9"
                    description="The C2 hostname is what comes after https:// or http:// for your api"
                >
                  <b-form-input
                      id="input-9"
                      v-model="form.imageInfo.c2"
                      type="text"
                      placeholder="C2 Hostname"
                      required
                      @input.native="checkUserForm()"
                  >
                  </b-form-input>
                </b-form-group>
              </b-col>
              <b-col style="margin-top: -10px" class="input">
                <b-form-group
                    id="input-group-10"
                    label="Port*"
                    label-for="input-10"
                    description="The api port"
                >
                  <b-form-input
                      id="input-10"
                      v-model="form.imageInfo.port"
                      type="number"
                      min="0"
                      placeholder="Port"
                      required
                      @input.native="checkUserForm()"
                  >
                  </b-form-input>
                </b-form-group>
              </b-col>
            </b-form-row>
            <b-form-row>
              <div style="margin: auto;">
                <b-button @click="$refs.setupCarousel.prev()" class="carousel-button" style="margin-right: 30px">Prev</b-button>
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
          </template>
        </b-carousel-slide>

        <!-- Slides with img slot -->
        <!-- Note the classes .d-block and .img-fluid to prevent browser default image alignment -->
      </b-carousel>
    </b-form>
  </div>
</template>

<script>
import axios from "axios";
import {router} from "../router";

export default {
  name: "Setup",
  data() {
    return {
      form: {
        userInfo: {
          firstName: "",
          lastName: "",
          email: "",
          username: "",
          password: "",
          confirmPw: ""
        },
        imageInfo: {
          name: "",
          c2: "",
          port: null,
          c2_protocol: "https",
        }
      },
      protocolOptions: [
        { text: 'HTTPS', value: 'https' },
        { text: 'HTTP', value: 'http' },
      ],
      formValid: false,
      dismissCountDown: 0,
      dismissSecs: 3,
      alert: "",
      variant: "",
      loading: false,
    }
  },
  async beforeCreate() {
    const resp = await axios.get(process.env.VUE_APP_API_ROOT+"/config/configured")
    if (resp.status === 200) {
      window.console.log(resp.data.configured)
      if (resp.data.configured) {
        router.push("/")
      }
    }
  },
  methods: {
    countDownChanged: function (dismissCountDown) {
      this.dismissCountDown = dismissCountDown
    },
    showAlert: function (variant) {
      this.variant = variant
      this.dismissCountDown = this.dismissSecs
    },
    checkUserForm: function() {
      if(this.form.userInfo.firstName.length === 0 || this.form.userInfo.lastName.length === 0 ||
          this.form.userInfo.email.length === 0 || this.form.userInfo.username.length === 0 ||
          this.form.userInfo.password.length === 0 || this.form.userInfo.confirmPw.length === 0){
        this.formValid = false
      } else if (this.form.userInfo.password !== this.form.userInfo.confirmPw) {
        this.formValid = false
      } else {
        this.formValid = true
      }
    },
    submitSetup: function() {
      window.console.log(this.form)
      let that = this
      let setupInfoJson = JSON.stringify(this.form)
      this.loading = true
      this.dismissCountDown = 0
      axios.post(
          process.env.VUE_APP_API_ROOT+"/config/configured", setupInfoJson
      ).then(response => {
        if (response.status === 200) {
          that.loading = false
          window.console.log(response.data.error)
          if (response.data.token != null) {
            that.$cookies.set("token",response.data.token,"24h","/")
            router.push('/').catch(()=>{})
          } else {
            that.alert = response.data.error
            that.showAlert("danger")
          }
        }
      })
    },
  }
}
</script>

<style>
.body {
  width: 100vh;
  height: 100vh;
}

.form-row {
  margin: 0;
}

.input {
  margin: 10px 10px 5px 10px;
}
.carousel-height{
  height: 500px;
}
.setup-carousel {
  position: relative;
  width: 500px;
  height: 500px;
  top: 30%;
  left: 50%;
  transform: translate(-50%, -50%);
  border-radius: 5px;
  box-shadow: 2px 3px 15px 5px #888888;
}

@media screen and (max-height: 920px){
  .setup-carousel{
    top: 35%;
  }
}

@media screen and (max-height: 800px){
  .setup-carousel{
    top: 37%;
  }
}

@media screen and (max-height: 760px){
  .setup-carousel{
    top: 40%;
  }
}

@media screen and (max-height: 720px){
  .setup-carousel{
    top: 43%;
  }
}

@media screen and (max-height: 685px){
  .setup-carousel{
    top: 45%;
  }
}
@media screen and (max-height: 650px){
  .setup-carousel{
    top: 48%;
  }
}

.carousel-button{
  width: 75px;
}

.lds-ellipsis {
  display: inline-block;
  position: relative;
  width: 80px;
  height: 80px;
  margin: auto;
}
.lds-ellipsis div {
  position: absolute;
  top: 33px;
  width: 13px;
  height: 13px;
  border-radius: 50%;
  background: #7e7e7e;
  animation-timing-function: cubic-bezier(0, 1, 1, 0);
}
.lds-ellipsis div:nth-child(1) {
  left: 8px;
  animation: lds-ellipsis1 0.6s infinite;
}
.lds-ellipsis div:nth-child(2) {
  left: 8px;
  animation: lds-ellipsis2 0.6s infinite;
}
.lds-ellipsis div:nth-child(3) {
  left: 32px;
  animation: lds-ellipsis2 0.6s infinite;
}
.lds-ellipsis div:nth-child(4) {
  left: 56px;
  animation: lds-ellipsis3 0.6s infinite;
}
@keyframes lds-ellipsis1 {
  0% {
    transform: scale(0);
  }
  100% {
    transform: scale(1);
  }
}
@keyframes lds-ellipsis3 {
  0% {
    transform: scale(1);
  }
  100% {
    transform: scale(0);
  }
}
@keyframes lds-ellipsis2 {
  0% {
    transform: translate(0, 0);
  }
  100% {
    transform: translate(24px, 0);
  }
}
</style>