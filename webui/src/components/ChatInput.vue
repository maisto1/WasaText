// components/ChatInput.vue
<script>
export default {
  data() {
    return {
      message: '',
      isSubmitting: false
    }
  },

  methods: {
    async handleSubmit() {
      if (!this.message.trim() || this.isSubmitting) return;

      this.isSubmitting = true;
      try {
        await this.$emit('send', this.message);
        this.message = '';
      } finally {
        this.isSubmitting = false;
      }
    }
  }
}
</script>

<template>
  <div class="chat-input p-3 bg-dark border-top">
    <form @submit.prevent="handleSubmit" class="d-flex">
      <input
        v-model="message"
        type="text"
        class="form-control form-control-lg bg-secondary text-white me-2"
        placeholder="Type a message..."
        :disabled="isSubmitting"
      />
      <button 
        type="submit" 
        class="btn btn-primary btn-lg"
        :disabled="!message.trim() || isSubmitting"
      >
        <i class="fas" :class="isSubmitting ? 'fa-spinner fa-spin' : 'fa-paper-plane'"></i>
      </button>
    </form>
  </div>
</template>

<style scoped>
.chat-input {
  background-color: #202c33;
}

.form-control {
  background-color: #2a3942 !important;
  border: none;
}

.form-control:focus {
  box-shadow: none;
  border-color: #00a884;
}

.btn-primary {
  background-color: #00a884;
  border-color: #00a884;
}

.btn-primary:hover:not(:disabled) {
  background-color: #008f6f;
  border-color: #008f6f;
}

.btn-primary:disabled {
  background-color: #2a3942;
  border-color: #2a3942;
}
</style>