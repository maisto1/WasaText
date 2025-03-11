<script>
import ForwardModal from './ForwardModal.vue';
import ReplyModal from './ReplyModal.vue';

export default {
  components: {
    ForwardModal,
    ReplyModal
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
      showForwardModal: false,
      showReplyModal: false
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
    
    handleReplyClick() {
      console.log('Opening reply modal...');
      this.showReplyModal = true;
    },
    
    handleSendReply(originalMessage, replyData) {
      console.log('Sending reply to message:', originalMessage.id, 'with data:', replyData);
      this.$emit('send-reply', originalMessage, replyData);
      this.showReplyModal = false;
    },

    closeForwardModal() {
      console.log('Closing forward modal');
      this.showForwardModal = false;
    },
    
    closeReplyModal() {
      console.log('Closing reply modal');
      this.showReplyModal = false;
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
      
      <!-- Reply Info -->
      <div v-if="message.replyTo" class="reply-info mb-2">
        <div class="reply-preview">
          <div class="reply-sender">
            {{ message.replyTo.sender }}
          </div>
          <div class="reply-content">
            <template v-if="message.replyTo.content">
              {{ message.replyTo.content }}
            </template>
            <template v-else>
              <i class="fas fa-camera me-1"></i> Photo
            </template>
          </div>
        </div>
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
          <i v-if="message.status === 'read'" class="fas fa-check-double" title="Read"></i>
          <i v-else class="fas fa-check" title="Sent"></i>
        </span>
      </div>

      <!-- Menu delle azioni sotto al messaggio -->
      <div v-if="showActions" class="message-actions-below">
        <div class="action-button-container">
          <button 
            @click.stop="handleReplyClick" 
            class="action-button reply-action"
            title="Reply"
          >
            <i class="fas fa-reply"></i>
          </button>
          
          <button 
            @click.stop="handleForwardClick"
            class="action-button forward-action"
            title="Forward"
          >
            <i class="fas fa-share"></i>
          </button>
          
          <button 
            v-if="isMyMessage"
            @click.stop="handleDelete" 
            class="action-button delete-action"
            title="Delete"
          >
            <i class="fas fa-trash"></i>
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
    
    <!-- Reply Modal -->
    <ReplyModal
      v-if="showReplyModal"
      :show="showReplyModal"
      :message="message"
      @close="closeReplyModal"
      @send-reply="handleSendReply"
    />
  </div>
</template>

<style scoped>
.message-bubble {
  border-radius: 12px;
  padding: 8px 12px;
  word-break: break-word;
}

/* Stile per il menu delle azioni sotto al messaggio */
.message-actions-below {
  position: absolute;
  bottom: -45px;
  left: 50%;
  transform: translateX(-50%);
  z-index: 10;
}

.action-button-container {
  display: flex;
  background-color: #202c33;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.4);
  padding: 2px;
}

.action-button {
  background: transparent;
  border: none;
  color: #fff;
  width: 36px;
  height: 36px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 2px;
  cursor: pointer;
  transition: background-color 0.2s;
}

.action-button:hover {
  background-color: rgba(255, 255, 255, 0.1);
}

.reply-action {
  color: #00a884;
}

.forward-action {
  color: #34B7F1;
}

.delete-action {
  color: #f15c6d;
}

.message-status {
  font-size: 0.8em;
  color: rgba(255, 255, 255, 0.8);
}

.message-status .fa-check {
  color: rgba(255, 255, 255, 0.6);
}

.message-status .fa-check-double {
  color: #34B7F1;
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

.reply-info {
  padding: 8px;
  border-radius: 6px;
  background-color: rgba(0, 0, 0, 0.2);
  position: relative;
  border-left: 3px solid #34B7F1;
  margin-top: 4px;
}

.reply-preview {
  font-size: 0.85rem;
  overflow: hidden;
}

.reply-sender {
  font-weight: 600;
  color: #34B7F1;
  margin-bottom: 2px;
}

.reply-content {
  color: rgba(255, 255, 255, 0.8);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
</style>