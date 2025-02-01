<script>
import SearchBar from '../components/SearchBar.vue'

export default {
  components: {
    SearchBar
  },
  data: function() {
    return {
      users: [],
      loading: false,
      defaultImage: '../public/default-icon.png'
    }
  },
  methods: {
    async searchUsers(query) {
      if (!query) {
        this.users = []
        return
      }
      this.loading = true
      try {
        const token = sessionStorage.getItem('authToken')
        const response = await this.$axios.get(`/users/`, {
          params: {
            username: query
          },
          headers: {
            'Authorization': token
          }
        })
        console.log("Search results:", response.data)
        this.users = response.data
      } catch (e) {
        console.error("Error while searching users:", e)
        this.users = []
      } finally {
        this.loading = false
      }
    },
    handleImageError(e) {
      e.target.src = this.defaultImage
    }
  }
}
</script>

<template>
  <div class="conversations-container">
    <SearchBar @search="searchUsers" />
    
    <div v-if="loading" class="loading-indicator">
      Searching...
    </div>
    
    <div v-else-if="users.length > 0" class="search-results">
      <div v-for="user in users"
        :key="user.id"
        class="user-item"
      >
        <img
          :src="user.photo ? `data:image/jpeg;base64,${user.photo}` : defaultImage"
          :alt="user.username"
          class="user-avatar"
          @error="handleImageError"
        >
        <div class="user-info">
          <div class="username">{{ user.username }}</div>
        </div>
      </div>
    </div>
  </div>
</template>

<style>
.loading-indicator {
  color: #808080;
  text-align: center;
  padding: 20px;
}

.conversations-container {
  width: 100%;
  max-width: 400px;
  height: 100vh;
  border-right: 1px solid #2c2c2c;
  background-color: #111;
}

.search-results {
  overflow-y: auto;
  max-height: calc(100vh - 60px);
}

.user-item {
  display: flex;
  align-items: center;
  padding: 12px;
  cursor: pointer;
  border-bottom: 1px solid #2c2c2c;
}

.user-item:hover {
  background-color: #202020;
}

.user-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  margin-right: 12px;
  object-fit: cover;
  background-color: #111;
  border: 2px solid #111;
}

.user-info {
  display: flex;
  flex-direction: column;
}

.username {
  color: white;
  font-weight: 500;
}
</style>