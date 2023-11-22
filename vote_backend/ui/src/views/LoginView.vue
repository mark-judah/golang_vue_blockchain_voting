<!-- <template>
  <div class="home">
    <img alt="Vue logo" src="../assets/logo.png">
    <HelloWorld msg="Welcome to Your Vue.js App"/>
  </div>
</template>

<script>
// @ is an alias to /src
import HelloWorld from '@/components/HelloWorld.vue'

export default {
  name: 'HomeView',
  components: {
    HelloWorld
  }
}
</script> -->
<template>

<div class="h-screen flex items-center justify-center flex-col">
  <div class="flex items-center justify-center pb-5">
      <h1 class="text-[#565656] font-extrabold text-3xl">Chainvote Admin Dashboard</h1>
    </div>
    <div class="w-full bg-[#565656] rounded-lg shadow dark:border md:mt-0 sm:max-w-md xl:p-0 ">
        <div class="p-6 space-y-4 md:space-y-6 sm:p-8">
            <form ref="form" @submit.prevent="submitLogin" class="space-y-4 md:space-y-6">
                <div>
                    <label for="email" class="block mb-2 text-sm font-medium text-gray-900 dark:text-white">Username</label>
                    <input type="email" name="email" id="email" v-model="loginCreds.email"
                        class="bg-white border border-gray-300 text-gray-900 sm:text-sm  focus:ring-primary-600 focus:border-primary-600 block w-full p-2.5 i"
                        placeholder="username" required="">
                </div>
                <div>
                    <label for="password"
                        class="block mb-2 text-sm font-medium text-gray-900 dark:text-white">Password</label>
                    <input type="password" name="password" id="password"  v-model="loginCreds.password"
                        placeholder="••••••••"
                        class="bg-white border border-gray-300 text-gray-900 sm:text-sm focus:ring-primary-600 focus:border-primary-600 block w-full p-2.5 i"
                        required="">
                </div>
                    <div class="flex items-center justify-center">
                        <button type="submit"
                            class="bg-white text-sm px-5 py-2.5 text-center mr-3 md:mr-0">Login</button>

                    </div>

            </form>
        </div>
    </div>
</div>
</template>

<script>
import SecureLS from 'secure-ls'
import axios from 'axios'

export default {

  data () {
    return {
      page: 'login',
      ls: new SecureLS(),

      loginCreds: {
        email: '',
        password: ''
      }
    }
  },

  methods: {
    submitLogin () {
      // console.log(this.$refs.form.formEmail.value)
      console.log(this.loginCreds)
      axios.post('http://127.0.0.1:3500/api/login', this.loginCreds).then(result => {
        console.log(result.data)
        this.ls.set('user', {
          username: result.data.username,
          email: result.data.email,
          token: result.data.token
        })
        console.log(this.ls.get('user'))
        if (this.ls.get('user') != null) {
          window.location.href = 'admin'
        }
      })
    }
  }
}
</script>
