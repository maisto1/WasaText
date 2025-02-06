// components/MessageBubble.vue
<script>
export default {
  props: {
    message: {
      type: Object,
      required: true
    },
    isMyMessage: {
      type: Boolean,
      required: true
    }
  },

  data() {
    return {
      showActions: false
    }
  },

  methods: {
    formatTimestamp(timestamp) {
      if (!timestamp) return '';
      const date = new Date(timestamp * 1000);
      return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
    },

    handleDelete() {
      if (confirm('Are you sure you want to delete this message?')) {
        this.$emit('delete', this.message.id);
      }
    },

    handleForward() {
      // In a real app, you'd show a modal for selecting the target conversation
      const targetConversationId = prompt('Enter conversation ID to forward to:');
      if (targetConversationId) {
        this.$emit('forward', this.message.id, parseInt(targetConversationId));
      }
    }
  }
}
</script>

<template>
  <div :class="['message mb-3 d-flex', isMyMessage ? 'justify-content-end' : 'justify-content-start']">
    <div 
      :class="['message-bubble p-2 rounded position-relative', 
               isMyMessage ? 'bg-primary text-white' : 'bg-secondary text-white']"
      style="max-width: 70%;"
      @mouseenter="showActions = true"
      @mouseleave="showActions = false"
    >
      <!-- Sender name for received messages -->
      <div v-if="!isMyMessage" class="message-sender mb-1">
        <small class="text-light">{{ message.sender.username }}</small>
      </div>

      <!-- Message content -->
      <div class="message-content">
        <template v-if="message.type === 'text'">
          <p class="mb-1">{{ message.content }}</p>
        </template>
        <template v-else-if="message.type === 'media'">
          <img 
            :src="'data:image/jpeg;base64,' + message.media"
            class="img-fluid rounded mb-1"
            alt="Media content"
          />
          <p v-if="message.content" class="mb-1">{{ message.content }}</p>
        </template>
      </div>

      <!-- Message meta -->
      <div class="message-meta d-flex justify-content-end align-items-center">
        <small class="text-light me-1">{{ formatTimestamp(message.timestamp) }}</small>
        <span v-if="isMyMessage" class="message-status">
          <i class="fas" :class="message.status === 'read' ? 'fa-check-double' : 'fa-check'"></i>
        </span>
      </div>

      <!-- Message actions -->
      <div v-if="showActions" class="message-actions position-absolute top-0 end-0 mt-1 me-1">
        <div class="btn-group">
          <button 
            v-if="isMyMessage"
            @click="handleDelete" 
            class="btn btn-sm btn-danger"
          >
            <i class="fas fa-trash"></i>
          </button>
          <button 
            @click="handleForward"
            class="btn btn-sm btn-info"
          >
            <i class="fas fa-share"></i>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.message-bubble {
  border-radius: 12px;
  padding: 8px 12px;
}

.message-actions {
  transform: translateY(-100%);
  opacity: 0.9;
}

.message-status {
  font-size: 0.8em;
  color: rgba(255, 255, 255, 0.8);
}

.btn-group .btn {
  padding: 0.25rem 0.5rem;
  font-size: 0.75rem;
}
</style>