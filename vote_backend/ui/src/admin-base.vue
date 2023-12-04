<template>
    <AdminNavAndSideBar v-if="shouldDisplayAdminNav"></AdminNavAndSideBar>
    <router-view></router-view>
</template>

<style></style>
<script>
import AdminNavAndSideBar from './components/AdminNavAndSideBar.vue'

export default {
  components: {
    AdminNavAndSideBar
  },

  computed: {
    shouldDisplayAdminNav () {
      const allowedRoutes = [
        'admin',
        'counties',
        'constituencies',
        'wards',
        'polling-stations',
        'candidates',
        'voters',
        'desktop-clients',
        'transaction-pool',
        'connected-nodes',
        'users',
        'new-user',
        'new-county',
        'new-constituency',
        'new-ward',
        'new-polling-station',
        'new-candidate',
        'new-voter'
      ]

      return allowedRoutes.includes(this.$route.name)
    }
  }
}
</script>
