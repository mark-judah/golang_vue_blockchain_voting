<template>
  <div class="">
    <div class="ml-96 mr-12 mt-32 flex flex-col justify-start  bg-[#ffffff] h-full">
      <div class="flex flex-row justify-between items-center w-full">
        <div>
          <p class="px-12">Users</p>
        </div>

        <div class="flex flex-row px-12">
          <router-link :to="{ name: 'users' }" @click="$router.go()" class="flex items-center p-2">
            <img src="../assets/images/refresh.svg" class="w-6 h-8">
          </router-link>

          <router-link :to="{ name: 'new-user' }" class="flex items-center p-2">
            <img src="../assets/images/add.svg" class="w-6 h-8">
          </router-link>
        </div>
      </div>

      <div class="px-12 flex justify-center w-fit">
        <table class="table-auto space-y-6 w-fit">
          <thead class="">

            <tr class="rounded-3xl">
              <th class="px-1.5 py-2 border">Index</th>
              <th class="px-1.5 py-2 border">Name</th>
              <th class="px-1.5 py-2 border">Email</th>
              <th class="px-1.5 py-2 border">Contact</th>
              <th class="px-1.5 py-2 border">Role</th>

            </tr>
          </thead>
          <tbody>
            <tr v-for="(user, index) in this.allUsers" :key="index" class="bg-white border-b">
              <td class="px-1.5 py-2   border ">{{ index + 1 }} </td>
              <td class="px-1.5 py-2   border ">{{ user.name }}</td>
              <td class="px-1.5 py-2   border ">{{ user.email }}</td>
              <td class="px-1.5 py-2   border ">{{ user.contact }}</td>
              <td class="px-1.5 py-2   border ">{{ user.role }}</td>

            </tr>
          </tbody>
        </table>
      </div>

    </div>
  </div>
</template>

<script>
import axios from 'axios'
import SecureLS from 'secure-ls'

export default {

  data () {
    this.getUsers()
    return {
      allUsers: {}

    }
  },

  methods: {
    getUsers () {
      const ls = new SecureLS()
      const config = {
        headers: { Authorization: `Bearer ${ls.get('user').token}` }
      }
      axios.get(
        'http://127.0.0.1:3500/api/secured/get-all-users', config
      ).then((response) => {
        console.log(response.data)
        this.allUsers = response.data
      }).catch(function (error) {
        if (error.response.status === 401) {
          ls.removeAll()
          window.location.href = '/'
        } if (error.toJSON().message === 'Network Error') {
          alert('no internet connection')
        }
      })
    }
  }

}
</script>

<style scoped></style>
