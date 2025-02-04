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
    <input
      type="text"
      :value="modelValue"
      @input="handleInput"
      class="form-control bg-dark text-white search-input"
      placeholder="Search users..."
    >
  </div>
</template>