<script>
export default {
  props: {
    modelValue: {
      type: String,
      required: true
    }
  },
  emits: ['update:modelValue', 'search'],
  methods: {
    handleInput(event) {
      this.$emit('update:modelValue', event.target.value);
      this.$emit('search');
    }
  }
}
</script>

<template>
  <div class="search-container p-2">
    <div class="search-wrapper">
      <i class="fas fa-search search-icon"></i>
      <input
        type="text"
        :value="modelValue"
        @input="handleInput"
        class="form-control bg-dark text-white search-input"
        placeholder="Search users..."
      >
    </div>
  </div>
</template>

<style>
.search-wrapper {
  position: relative;
  display: flex;
  align-items: center;
}

.search-icon {
  position: absolute;
  left: 12px;
  color: #8696a0;
  z-index: 1;
}

.search-input {
  padding-left: 35px;
  background-color: #202c33 !important;
  border: 1px solid #36434a;
  border-radius: 8px;
}

.search-input:focus {
  border-color: #00a884;
  box-shadow: none;
}

.search-input::placeholder {
  color: #8696a0;
}
</style>