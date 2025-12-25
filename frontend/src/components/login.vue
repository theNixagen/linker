<script setup lang="ts">
import { reactive } from 'vue';

const form = reactive({
  email: "",
  password: "",
  name: "",
  username: ""
})

function submit() {

}

</script>

<template>
  <main class="bg-lime-500">
    <form @submit.prevent="submit">
      <label>
        Email
        <input v-model.trim="form.email" type="email" autocomplete="email" />
      </label>
    </form>
  </main>
</template>
