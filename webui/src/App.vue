<script setup>
import { RouterLink, RouterView } from 'vue-router'
</script>

<template>
 <div class="bg-dark">
   <RouterLink to="/login" class="navbar-brand fw-bold ps-2">WasaText</RouterLink>
    <main>
      <RouterView />
    </main>
 </div>
</template>

<style>
.navbar-brand {
 font-size: 1.5rem;
 letter-spacing: 1px;
 color: white;
 text-decoration: none;
 padding: 1rem;
 display: block;
}

.navbar-brand,
.navbar-brand:hover,
.navbar-brand:active,
.navbar-brand:visited {
  color: white;
  text-decoration: none;
}
</style>