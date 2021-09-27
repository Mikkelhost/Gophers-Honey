<template>
  <div class="b-sidebar">
    <b-row class="test">
      <b-col md="2">
        <b-nav-item href="#" left>
          <b-icon-gear></b-icon-gear>
          Settings
        </b-nav-item>
        <b-nav-item href="#" left>
          <b-icon-gear></b-icon-gear>
          Settings
        </b-nav-item>
      </b-col>
    </b-row>
  </div>

</template>

<script>
export default {
  name: "Dashboard"
}
</script>

<style>
.page-height {
  height: 100vh;
  width: 100%;
}

.test {
  border: teal;
  border-bottom-color: #2c3e50;
  border-bottom-width: 2px;
}

</style>