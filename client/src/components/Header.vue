<script setup lang="ts">
import { Github, Code, Sandwich } from "lucide-vue-next";
</script>

<template>
  <header class="flex justify-between py-2">
    <div class="flex gap-4">
      <Sandwich class="text-3xl self-center" />
      <h1 class="text-3xl">Prueba</h1>
    </div>
    <div class="flex justify-center gap-4 items-center">
      <a
        href="https://github.com/andresprogramacion123"
        target="_blank"
        rel="noopener noreferrer"
        class="flex gap-2 items-center hover:scale-105 duration-200"
      >
        <Github :size="32" />
        <span>/Julian Montoya</span>
      </a>
      <a
        href="https://github.com/"
        target="_blank"
        rel="noopener noreferrer"
        class="flex gap-2 items-center hover:scale-105 duration-200"
      >
        <Code :size="32" />
        <span>Codigo</span>
      </a>
    </div>
  </header>
</template>
