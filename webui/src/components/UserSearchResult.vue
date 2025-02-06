// components/UserSearchResult.vue
<script>
export default {
  props: {
    user: {
      type: Object,
      required: true
    }
  }
}
</script>

<template>
  <div class="user-item p-3">
    <div class="d-flex align-items-center">
      <div class="rounded-circle bg-secondary flex-shrink-0" style="width: 40px; height: 40px;">
        <img v-if="user.profilePhoto" 
             :src="'data:image/jpeg;base64,' + user.profilePhoto"
             class="rounded-circle w-100 h-100"
             alt="Profile photo">
      </div>
      <div class="ms-3">
        <h6 class="mb-0 text-white">{{ user.username }}</h6>
      </div>
    </div>
  </div>
</template>

<style>
.user-item {
  cursor: pointer;
  transition: background-color 0.2s ease;
  border-bottom: 1px solid #36434a;
}

.user-item:hover {
  background-color: #2a3942;
}
</style>