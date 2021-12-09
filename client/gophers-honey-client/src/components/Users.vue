<template>
  <div>
    <b-modal id="new-user" size="md" hide-footer>
      <template #modal-title>
        Create a new user
      </template>
      <b-form @submit.prevent="submitUser" class="container" style="height: fit-content">
        <b-row class="form-row">
          <b-col class="input">
            <b-form-group
                id="input-group-1"
                label="First name*"
                label-for="input-1"
            >
              <b-form-input
                  id="input-1"
                  v-model="form.firstName"
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
                  v-model="form.lastName"
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
                  v-model="form.email"
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
                  v-model="form.username"
                  type="text"
                  placeholder="Username"
                  required
                  @input.native="checkUserForm()"
              >
              </b-form-input>
            </b-form-group>
          </b-col>
          <b-col class="input">
            <b-form-group
              id="input-group-7"
              label="Role"
              label-for="input-7"
            >
              <b-form-select v-model="form.role">
                <b-form-select-option value="Admin">Admin</b-form-select-option>
                <b-form-select-option value="User">User</b-form-select-option>
              </b-form-select>

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
                  v-model="form.password"
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
                  v-model="form.confirmPw"
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
          <div style="margin: auto;">
            <b-button type="submit" class="carousel-button" :disabled="!formValid">Submit</b-button>
          </div>
        </b-form-row>
        <b-form-row>
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
        </b-form-row>
      </b-form>
    </b-modal>
    <b-row>
      <b-button id="add-user" @click="$bvModal.show('new-user')">
        <b-icon-plus></b-icon-plus>
        Add User
      </b-button>
    </b-row>
    <b-row>
      <table style="margin: auto;">
        <tr>
          <th>First name</th>
          <th>Last name</th>
          <th>Username</th>
          <th>Email</th>
          <th>Role</th>
          <th>Delete</th>
        </tr>

        <tr v-for="user in users" :key="user.id">
          <td>{{ user.first_name }}</td>
          <td>{{ user.last_name }}</td>
          <td>{{ user.username }}</td>
          <td>{{ user.email }}</td>
          <td>{{ user.role }}</td>
          <td>
            <b-icon-trash class="click-icon" v-on:click="deleteUser(user.username)" variant="danger"></b-icon-trash>
          </td>
        </tr>
      </table>
    </b-row>
  </div>
</template>

<script>
import axios from "axios";

export default {
  name: "Users",
  data: function () {
    return {
      users: [],
      form: {
        firstName: "",
        lastName: "",
        username: "",
        email: "",
        role: "",
        password: "",
        confirmPw: "",
      },
      formValid: false,
      dismissCountDown: 0,
      dismissSecs: 3,
      alert: "",
      variant: "",
    }
  },
  created() {
    axios.defaults.headers.common['Authorization'] = 'Bearer ' + this.$cookies.get("token")
    this.getUsers()
  },
  methods: {
    showAlert: function (variant) {
      this.variant = variant
      this.dismissCountDown = this.dismissSecs
    },
    countDownChanged: function (dismissCountDown) {
      this.dismissCountDown = dismissCountDown
    },
    removeItem: function (array, key, value) {
      const index = array.findIndex(obj => obj[key] === value)
      return index >= 0 ? [
        ...array.slice(0, index),
        ...array.slice(index + 1)
      ] : array;
    },
    checkUserForm: function () {
      if (this.form.firstName.length === 0 || this.form.lastName.length === 0 ||
          this.form.email.length === 0 || this.form.username.length === 0 ||
          this.form.password.length === 0 || this.form.confirmPw.length === 0) {
        this.formValid = false
      } else if (this.form.password !== this.form.confirmPw) {
        this.formValid = false
      } else {
        this.formValid = true
      }
    },
    getUsers: function () {
      axios({
        url: process.env.VUE_APP_API_ROOT + "/users",
        method: "GET",
      }).then(function (response) {
        if (response.status === 200) {
          if (typeof response.data.error === "undefined") {
            window.console.log(response.data)
            this.users = response.data
          } else {
            this.alert = response.data.error
            this.showAlert("danger")
          }
        } else {
          this.alert = "Statuscode" + response.status
          window.console.log(this.alert)
        }

      }.bind(this))
    },
    submitUser: function () {
      window.console.log("Submitting user" + this.form)
      axios({
        url: process.env.VUE_APP_API_ROOT + "/users",
        method: "POST",
        data: this.form,
      }).then(function (response) {
        if (response.status === 200) {
          if (response.data.error === "") {
            this.getUsers()
            this.$bvModal.hide('new-user')
          } else {
            this.alert = response.data.error
            this.showAlert("danger")
          }
        } else {
          this.alert = "Statuscode: " + response.status
          this.showAlert("danger")
        }

      }.bind(this))
    },
    deleteUser: function (username) {
      window.console.log("Deleting user: " + username)
      let data = {username: username}
      if (confirm("Do you really want to delete user: " + username)) {
        axios({
          url: process.env.VUE_APP_API_ROOT + "/users",
          method: "DELETE",
          data: data,
        }).then(function (response) {
          if (response.status === 200) {
            if (response.data.error === "") {
              this.getUsers()
            } else {
              window.console.log("error deleting user: " + response.data.error)
            }
          } else {
            window.console.log("Statuscode: " + response.status)
          }
        }.bind(this))
      }
    }
  }
}
</script>

<style scoped>

</style>