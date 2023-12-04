<template>
    <nav class="fixed top-0 z-50 w-full bg-white">
        <div class="flex items-center justify-start ml-60">

            <BlockExplorerNav></BlockExplorerNav>

            <router-link :to="{ name: 'results' }" class="flex items-center p-2 ml-32">
                <div class="px-1">
                    <button class="bg-transparent  font-semibold text-black py-2 px-4 ">
                        Results</button>

                </div>
            </router-link>

            <router-link :to="{ name: 'block-explorer' }" class="flex items-center p-2">
                <div class="px-1">
                    <button class="bg-transparent  font-semibold text-black py-2 px-4 ">
                        Public Ledger</button>

                </div>
            </router-link>
        </div>
    </nav>
    <div>
        <div class="mt-32  mr-12 mt-16 h-full">
            <div class="flex justify-evenly justify-items-center mx-24">
                <div>
                    <PublicPresidenialResults></PublicPresidenialResults>
                </div>

                <div>
                    <PublicResultsPerCounty></PublicResultsPerCounty>
                </div>
            </div>

        </div>
    </div>
</template>

<style scoped></style>

<script>

import PublicPresidenialResults from '../components/PublicPresidentialResults.vue'
import PublicResultsPerCounty from '../components/PublicResultsPerCounty.vue'
import BlockExplorerNav from '../components/BlockExplorerNav.vue'

export default {
  components: {
    BlockExplorerNav,
    PublicPresidenialResults,
    PublicResultsPerCounty
  }
}
</script>
