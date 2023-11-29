<template>
    <div>
        <div class="ml-96 mr-12 mt-16 h-full">
            <AdminQuickStats></AdminQuickStats>
            <div class="flex justify-evenly justify-items-center mx-24">
                <div>
                    <AdminPresidenialResults></AdminPresidenialResults>
                </div>

                <div>
                    <AdminResultsPerCounty></AdminResultsPerCounty>
                </div>
            </div>

        </div>
    </div>
</template>

<style scoped></style>

<script>

import AdminQuickStats from '../components/AdminQuickStats.vue'
import AdminPresidenialResults from '../components/AdminPresidentialResults.vue'
import AdminResultsPerCounty from '../components/AdminResultsPerCounty.vue'

export default {
  components: {
    AdminQuickStats,
    AdminPresidenialResults,
    AdminResultsPerCounty
  }
}
</script>
