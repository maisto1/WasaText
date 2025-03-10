<script>
export default {
  props: {
    user: {
      type: Object,
      required: true
    }
  },
  methods: {
    getInitials(name) {
      if (!name) return '';
      return name.charAt(0).toUpperCase();
    },
    getAvatarColor(name) {
      if (!name) return '#202c33';
      
      let hash = 0;
      for (let i = 0; i < name.length; i++) {
        hash = name.charCodeAt(i) + ((hash << 5) - hash);
      }
      
      const h = hash % 360;
      const s = 65 + (hash % 20);
      const l = 45 + (hash % 10);
      
      return `hsl(${h}, ${s}%, ${l}%)`;
    }
  }
}
</script>

<template>
  <div class="user-item p-3">
    <div class="d-flex align-items-center">
      <div class="avatar-container" :style="{ backgroundColor: getAvatarColor(user.username) }">
        <img v-if="user.profilePhoto" 
             :src="'data:image/jpeg;base64,' + user.profilePhoto"
             class="avatar-image"
             alt="Profile photo">
        <div v-else class="avatar-text">
          {{ getInitials(user.username) }}
        </div>
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

.avatar-container {
  width: 40px;
  height: 40px;
  min-width: 40px;
  border-radius: 50%;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
}

.avatar-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.avatar-text {
  color: white;
  font-size: 1.2rem;
  font-weight: 600;
  text-align: center;
}
</style>