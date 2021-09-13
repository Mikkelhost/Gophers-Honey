<template>
  <b-navbar type="light" variant="light" class="navbar">
    <b-icon class="sidebar-toggle" v-b-toggle.sidebar icon="list"></b-icon>
    <b-navbar-brand href="/">Gophers-Honey</b-navbar-brand>
    <b-navbar-nav class="ml-auto">
      <b-nav-item-dropdown right>
        <!-- Using 'button-content' slot -->
        <template #button-content>
          <em>User</em>
        </template>
        <b-dropdown-item href="#">Profile</b-dropdown-item>
        <b-dropdown-item href="#">Sign Out</b-dropdown-item>
      </b-nav-item-dropdown>
      <b-nav-item href="#">
        <b-icon-bell></b-icon-bell>
      </b-nav-item>
    </b-navbar-nav>
  </b-navbar>
</template>

<script>
export default {
  name: "Navbar"
}
</script>

<style>
.sidebar-toggle {
  margin-right: 10px;
}
.navbar{
  border-style: solid;
  border-width: 0 0 2px 0;
  border-color: #ff0000;
}
</style>