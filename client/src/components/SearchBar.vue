<script setup lang="ts">
import { ref } from "vue";
const searchText = ref("");
const emit = defineEmits(["search"]);

const onSearch = () => {
  emit("search", searchText.value);
};
</script>

<template>
  <section class="flex justify-end gap-4 p-2">
    <input
      type="text"
      class="text-black rounded-md p-2 w-10/12"
      v-model="searchText"
      placeholder="Enter your search term..."
      @keyup.enter="onSearch"
    />
    <button
      @click="onSearch"
      class="bg-blue-500 text-white rounded-md p-2 w-2/12 opacity-90 hover:opacity-100" :style="{ backgroundColor: '#80a1f7' }"
    >
      Buscar
    </button>
  </section>
</template>
