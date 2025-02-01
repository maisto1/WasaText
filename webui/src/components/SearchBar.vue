<script>
export default {
  name: 'SearchBar',
  data() {
    return {
      searchQuery: ''
    }
  },
  methods: {
    handleInput() {
      this.$emit('search', this.searchQuery)
    }
  }
}
</script>

<template>
  <div class="search-bar">
    <div class="search-wrapper">
      <svg
        class="search-icon"
        xmlns="http://www.w3.org/2000/svg"
        width="20"
        height="20"
        viewBox="0 0 24 24"
        fill="none"
        stroke="#808080"
        stroke-width="2"
        stroke-linecap="round"
        stroke-linejoin="round"
      >
        <circle cx="11" cy="11" r="8"></circle>
        <line x1="21" y1="21" x2="16.65" y2="16.65"></line>
      </svg>
      <input
        v-model="searchQuery"
        @input="handleInput"
        type="text"
        placeholder="Search or start a new conversation"
        class="search-input"
      >
    </div>
  </div>
</template>

<style scoped>
.search-bar {
  padding: 10px;
  background-color: #202020;
}

.search-wrapper {
  position: relative;
  display: flex;
  align-items: center;
  background-color: #323232;
  border-radius: 8px;
  padding: 0 12px;
}

.search-icon {
  min-width: 20px;
  margin-right: 8px;
}

.search-input {
  width: 100%;
  padding: 8px 0;
  border: none;
  background-color: transparent;
  color: white;
  outline: none;
}

.search-input::placeholder {
  color: #808080;
}
</style>