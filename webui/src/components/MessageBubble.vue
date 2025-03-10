<script>
import ForwardModal from './ForwardModal.vue';

export default {
  components: {
    ForwardModal
  },
  
  props: {
    message: {
      type: Object,
      required: true
    },
    isMyMessage: {
      type: Boolean,
      required: true
    },
    showUsername: {
      type: Boolean,
      default: false 
    },
    availableConversations: {
      type: Array,
      default: () => []
    }
  },

  data() {
    return {
      showActions: false,
      imageLoaded: false,
      imageError: false,
      showForwardModal: false
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

    handleForwardClick() {
      console.log('Opening forward modal...');
      if (this.availableConversations.length === 0) {
        alert('No conversations available to forward to');
      } else {
        this.showForwardModal = true;
        console.log('Forward modal should be open now:', this.showForwardModal);
      }
    },

    handleForward(messageId, targetConversationId) {
      console.log('Forwarding message', messageId, 'to conversation', targetConversationId);
      this.$emit('forward', messageId, targetConversationId);
      this.showForwardModal = false;
    },

    closeForwardModal() {
      console.log('Closing forward modal');
      this.showForwardModal = false;
    },

    handleImageLoad() {
      this.imageLoaded = true;
    },

    handleImageError() {
      this.imageError = true;
    },

    openMediaInNewTab() {
      if (this.message.media) {
        const imgWindow = window.open();
        imgWindow.document.write(`
          <html>
            <head>
              <title>WasaText Media</title>
              <style>
                body { 
                  margin: 0; 
                  padding: 0; 
                  display: flex; 
                  justify-content: center; 
                  align-items: center; 
                  height: 100vh; 
                  background-color: #121212; 
                }
                img { 
                  max-width: 90%; 
                  max-height: 90vh; 
                  object-fit: contain; 
                }
              </style>
            </head>
            <body>
              <img src="data:image/jpeg;base64,${this.message.media}" alt="Full size media" />
            </body>
          </html>
        `);
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
      :style="[message.type === 'media' ? {'max-width': '300px'} : {'max-width': '70%'}]"
      @mouseenter="showActions = true"
      @mouseleave="showActions = false"
    >
      <!-- Forwarded Label -->
      <div v-if="message.isForwarded" class="forwarded-label">
        <i class="fas fa-share"></i> Forwarded
      </div>
      
      <div v-if="!isMyMessage && showUsername" class="message-sender mb-1">
        <small class="text-light">{{ message.sender.username }}</small>
      </div>

      <div class="message-content">
        <template v-if="message.type === 'media' && message.media">
          <div class="media-container mb-2">
            <div v-if="!imageLoaded && !imageError" class="text-center p-2">
              <div class="spinner-border spinner-border-sm text-light" role="status">
                <span class="visually-hidden">Loading...</span>
              </div>
            </div>
            
            <div v-if="imageError" class="text-center p-2">
              <i class="fas fa-exclamation-circle"></i>
              <small class="d-block">Failed to load image</small>
            </div>
            
            <img 
              v-show="!imageError"
              :src="'data:image/jpeg;base64,' + message.media"
              class="img-fluid rounded cursor-pointer"
              alt="Media content"
              @load="handleImageLoad"
              @error="handleImageError"
              @click="openMediaInNewTab"
            />
          </div>
          
          <p v-if="message.content" class="mb-1">{{ message.content }}</p>
        </template>
        
        <template v-else-if="message.type === 'text'">
          <p class="mb-1">{{ message.content }}</p>
        </template>
      </div>

      <div class="message-meta d-flex justify-content-end align-items-center">
        <small class="text-light me-1">{{ formatTimestamp(message.timestamp) }}</small>
        <span v-if="isMyMessage" class="message-status">
          <i class="fas" :class="message.status === 'read' ? 'fa-check-double' : 'fa-check'"></i>
        </span>
      </div>

      <div v-if="showActions" class="message-actions position-absolute top-0 end-0 mt-1 me-1">
        <div class="btn-group">
          <button 
            v-if="isMyMessage"
            @click.stop="handleDelete" 
            class="btn btn-sm btn-danger"
          >
            <i class="fas fa-trash"></i>
          </button>
          <button 
            @click.stop="handleForwardClick"
            class="btn btn-sm btn-info"
          >
            <i class="fas fa-share"></i>
          </button>
        </div>
      </div>
    </div>
    
    <!-- Forward Modal -->
    <ForwardModal
      v-if="showForwardModal"
      :show="showForwardModal"
      :conversations="availableConversations"
      :message-id="message.id"
      @close="closeForwardModal"
      @forward="handleForward"
    />
  </div>
</template>

<style scoped>
.message-bubble {
  border-radius: 12px;
  padding: 8px 12px;
  word-break: break-word;
}

.message-actions {
  transform: translateY(-100%);
  opacity: 0.9;
  z-index: 10;
}

.message-status {
  font-size: 0.8em;
  color: rgba(255, 255, 255, 0.8);
}

.btn-group .btn {
  padding: 0.25rem 0.5rem;
  font-size: 0.75rem;
}

.media-container {
  position: relative;
  min-height: 50px;
}

.cursor-pointer {
  cursor: pointer;
}

.media-container img {
  max-height: 300px;
  object-fit: contain;
  width: 100%;
}

.forwarded-label {
  color: #8696a0;
  font-size: 0.75rem;
  margin-bottom: 4px;
  font-style: italic;
}
</style>