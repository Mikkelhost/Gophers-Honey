<template>
  <h3>Welcome to login</h3>
</template>

<script>
export default{
  name: "Login",
  mounted() {
    console.log(process.env.VUE_APP_API_ROOT)
  }
}
</script>