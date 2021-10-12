<template>
  <div>
    <b-navbar toggleable="md" type="light" variant="light" class="navbar">
      <b-navbar-brand href="/">Gophers-Honey</b-navbar-brand>
      <b-navbar-toggle target="nav-collapse"></b-navbar-toggle>
      <b-collapse id="nav-collapse" is-nav>
        <b-navbar-nav>
          <b-nav-item href="/" left>
            <b-icon icon="graph-up"></b-icon>
            Dashboard
          </b-nav-item>
          <b-nav-item href="/honeypots" left>
            <b-icon icon="cpu"></b-icon>
            Honeypots
          </b-nav-item>
          <b-nav-item href="/settings" left>
            <b-icon-gear></b-icon-gear>
            Settings
          </b-nav-item>
        </b-navbar-nav>
        <b-navbar-nav class="ml-auto">
          <b-nav-item-dropdown right>
            <template #button-content>
              <em>User</em>
            </template>
            <b-dropdown-item href="#">Profile</b-dropdown-item>
            <b-dropdown-item href="#" @click="signOut">Sign Out</b-dropdown-item>
          </b-nav-item-dropdown>
          <b-nav-item href="#">
            <b-icon-bell></b-icon-bell>
          </b-nav-item>
        </b-navbar-nav>
      </b-collapse>
    </b-navbar>
  </div>
</template>

<script>
import { router } from '../router';

export default {
  name: "Navbar",
  data: function(){
    return{
      userinfo: null
    }
  },
  mounted() {

  },
  methods: {
    signOut: function(){
      this.$cookies.remove("token")
      router.push("/login").catch(()=>{})
    }
  }
}
</script>

<style>
.navbar{
  border-style: solid;
  border-width: 0 0 2px 0;
  border-color: #ffd500;
}
</style>