<script>
export default {
  props: {
    conversation: {
      type: Object,
      required: true
    },
    isActive: {
      type: Boolean,
      default: false
    }
  },
  methods: {
    formatTimestamp(timestamp) {
      if (!timestamp) return '';
      const date = new Date(timestamp * 1000);
      return date.toLocaleString();
    },
    getMessagePreview(message) {
      if (!message) return '';
      if (message.type === 'text') {
        return message.content;
      } else if (message.type === 'media') {
        return 'Photo';
      }
      return '';
    },
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
  <div :class="['conversation-item p-3', {'active': isActive}]">
    <div class="d-flex align-items-center">
      <div class="avatar-container" :style="{ backgroundColor: getAvatarColor(conversation.name) }">
        <img v-if="conversation.conversationPhoto" 
             :src="'data:image/jpeg;base64,' + conversation.conversationPhoto"
             class="avatar-image"
             alt="Profile photo">
        <div v-else class="avatar-text">
          {{ getInitials(conversation.name) }}
        </div>
      </div>
      
      <div class="ms-3 overflow-hidden">
        <div class="d-flex justify-content-between align-items-center">
          <h6 class="mb-1 text-truncate text-white">{{ conversation.name }}</h6>
          <small class="text-light ms-2" v-if="conversation.latestMessage">
            {{ formatTimestamp(conversation.latestMessage.timestamp) }}
          </small>
        </div>
        <p class="mb-0 text-light-grey small text-truncate" v-if="conversation.latestMessage">
          <template v-if="conversation.latestMessage.type === 'media'">
            <i class="fas fa-camera me-1"></i> Photo
            <template v-if="conversation.latestMessage.content">
              - {{ conversation.latestMessage.content }}
            </template>
          </template>
          <template v-else>
            {{ getMessagePreview(conversation.latestMessage) }}
          </template>
        </p>
      </div>
    </div>
  </div>
</template>

<style scoped>
.conversation-item {
  cursor: pointer;
  transition: background-color 0.2s ease;
  border-bottom: 1px solid #36434a;
}

.conversation-item:hover {
  background-color: #2a3942;
}

.conversation-item.active {
  background-color: #2a3942;
}

.text-light-grey {
  color: #8696a0;
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