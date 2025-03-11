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
  
  computed: {
    hasUnreadMessages() {
      if (!this.conversation.latestMessage) return false;
      
      const myUsername = sessionStorage.getItem('username');
      return this.conversation.latestMessage.status === 'sent' && 
             this.conversation.latestMessage.sender.username !== myUsername;
    }
  },
  
  methods: {
    formatTimestamp(timestamp) {
      if (!timestamp) return '';
      const date = new Date(timestamp * 1000);
      
      const today = new Date();
      const isToday = date.getDate() === today.getDate() && 
                      date.getMonth() === today.getMonth() && 
                      date.getFullYear() === today.getFullYear();
      
      if (isToday) {
        return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
      } else {
        return date.toLocaleDateString([], { day: '2-digit', month: '2-digit' });
      }
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
      <div class="avatar-wrapper">
        <div class="avatar-container" :style="{ backgroundColor: getAvatarColor(conversation.name) }">
          <img v-if="conversation.conversationPhoto" 
               :src="'data:image/jpeg;base64,' + conversation.conversationPhoto"
               class="avatar-image"
               alt="Profile photo">
          <div v-else class="avatar-text">
            {{ getInitials(conversation.name) }}
          </div>
        </div>
        
        <div v-if="hasUnreadMessages && !isActive" class="unread-indicator">
          <span class="pulse"></span>
        </div>
      </div>
      
      <div class="ms-3 overflow-hidden flex-grow-1">
        <div class="d-flex justify-content-between align-items-center">
          <h6 :class="['mb-1 text-truncate', hasUnreadMessages && !isActive ? 'fw-bold text-white' : 'text-white']">
            {{ conversation.name }}
          </h6>
          <div class="d-flex align-items-center">
            <small :class="['timestamp', hasUnreadMessages && !isActive ? 'text-primary fw-bold' : 'text-light']" v-if="conversation.latestMessage">
              {{ formatTimestamp(conversation.latestMessage.timestamp) }}
            </small>
          </div>
        </div>
        <div class="d-flex justify-content-between align-items-center">
          <p :class="['mb-0 small text-truncate flex-grow-1', hasUnreadMessages && !isActive ? 'text-white' : 'text-light-grey']" v-if="conversation.latestMessage">
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
          
          <div v-if="hasUnreadMessages && !isActive" class="unread-badge">
            <span>New</span>
          </div>
        </div>
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

.avatar-wrapper {
  position: relative;
  margin-right: 4px;
}

.avatar-container {
  width: 50px;
  height: 50px;
  min-width: 50px;
  border-radius: 50%;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.2);
}

.avatar-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.avatar-text {
  color: white;
  font-size: 1.4rem;
  font-weight: 600;
  text-align: center;
}

.unread-indicator {
  position: absolute;
  top: 0;
  right: 0;
  width: 14px;
  height: 14px;
  background-color: #00a884;
  border-radius: 50%;
  border: 2px solid #202c33;
  display: flex;
  align-items: center;
  justify-content: center;
}

.pulse {
  display: block;
  width: 100%;
  height: 100%;
  border-radius: 50%;
  background-color: #00a884;
  opacity: 0.8;
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0% {
    transform: scale(0.9);
    box-shadow: 0 0 0 0 rgba(0, 168, 132, 0.7);
  }
  
  70% {
    transform: scale(1);
    box-shadow: 0 0 0 6px rgba(0, 168, 132, 0);
  }
  
  100% {
    transform: scale(0.9);
  }
}

.unread-badge {
  min-width: 22px;
  height: 22px;
  border-radius: 11px;
  background-color: #00a884;
  color: white;
  font-size: 0.75rem;
  font-weight: bold;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0 6px;
  margin-left: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
  animation: fadeIn 0.3s ease-out;
}

.timestamp {
  font-size: 0.75rem;
  margin-left: 4px;
}

.text-primary {
  color: #00a884 !important;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: scale(0.8);
  }
  to {
    opacity: 1;
    transform: scale(1);
  }
}
</style>