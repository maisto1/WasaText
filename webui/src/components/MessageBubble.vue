<script>
import ForwardModal from './ForwardModal.vue';
import ReplyModal from './ReplyModal.vue';
import EmojiReaction from './EmojiReaction.vue';

export default {
  components: {
    ForwardModal,
    ReplyModal,
    EmojiReaction
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
    },
    conversationId: {
      type: Number,
      required: true
    }
  },

  data() {
    return {
      showContextMenu: false,
      contextMenuPosition: { x: 0, y: 0 },
      imageLoaded: false,
      imageError: false,
      showForwardModal: false,
      showReplyModal: false,
      reactions: [],
      comments: [],
      reactionPollingTimer: null
    }
  },

  created() {
    window.addEventListener('click', this.closeContextMenu);
    window.addEventListener('scroll', this.closeContextMenu);
    
    this.startCommentsPolling();
  },
  
  beforeUnmount() {
    window.removeEventListener('click', this.closeContextMenu);
    window.removeEventListener('scroll', this.closeContextMenu);
    
    this.stopCommentsPolling();
  },
  
  computed: {
    myUsername() {
      return sessionStorage.getItem('username');
    }
  },

  methods: {
    formatTimestamp(timestamp) {
      if (!timestamp) return '';
      const date = new Date(timestamp * 1000);
      return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
    },

    handleDelete() {
      if (confirm('Sei sicuro di voler eliminare questo messaggio?')) {
        this.$emit('delete', this.message.id);
      }
      this.closeContextMenu();
    },

    handleForwardClick() {
      this.showForwardModal = true;
      this.closeContextMenu();
    },

    handleForward(messageId, targetConversationId) {
      this.$emit('forward', messageId, targetConversationId);
      this.showForwardModal = false;
    },
    
    handleReplyClick() {
      this.showReplyModal = true;
      this.closeContextMenu();
    },
    
    handleSendReply(originalMessage, replyData) {
      this.$emit('send-reply', originalMessage, replyData);
      this.showReplyModal = false;
    },

    closeForwardModal() {
      this.showForwardModal = false;
    },
    
    closeReplyModal() {
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
    },
    
    showContextMenuHandler(event) {
      event.preventDefault();
      event.stopPropagation();
      
      this.closeContextMenu();
      
      const x = event.clientX;
      const y = event.clientY;

      const menuWidth = 180;
      const menuHeight = 180;
      
      const adjustedX = Math.min(x, window.innerWidth - menuWidth - 10);
      
      const adjustedY = Math.min(y, window.innerHeight - menuHeight - 10);
      
      this.contextMenuPosition = { x: adjustedX, y: adjustedY };
      this.showContextMenu = true;
      
      setTimeout(() => {
        this.addClickOutsideListener();
      }, 100);
    },
    
    addClickOutsideListener() {
      document.addEventListener('click', this.handleClickOutside);
      document.addEventListener('contextmenu', this.handleClickOutside);
    },
    
    handleClickOutside(event) {
      const contextMenu = this.$refs.contextMenu;
      if (contextMenu && !contextMenu.contains(event.target)) {
        this.closeContextMenu();
      }
    },
    
    closeContextMenu() {
      document.removeEventListener('click', this.handleClickOutside);
      document.removeEventListener('contextmenu', this.handleClickOutside);
      
      this.showContextMenu = false;
    },
    
    async fetchComments() {
      try {
        const response = await this.$axios.get(`/conversations/${this.conversationId}/messages/${this.message.id}/comments/`);
        if (!response || !response.data) {
          console.error('Invalid response from server:', response);
          return;
        }
        
        this.comments = response.data;
        
        this.processReactionsFromComments();
      } catch (error) {
        console.error('Error fetching comments/reactions:', error);
      }
    },
    
    processReactionsFromComments() {
      this.reactions = [];
      
      this.comments.forEach(comment => {
        if (comment.content && comment.content.startsWith('reaction:')) {
          const emoji = comment.content.replace('reaction:', '');
          
          this.reactions.push({
            id: comment.id,
            messageId: this.message.id,
            username: comment.sender.username,
            emoji: emoji,
            timestamp: comment.timestamp
          });
        }
      });
      
      this.reactions.sort((a, b) => {
        return new Date(b.timestamp) - new Date(a.timestamp);
      });
    },
    
    startCommentsPolling() {
      this.fetchComments();
      
      this.reactionPollingTimer = setInterval(() => {
        this.fetchComments();
      }, 3000);
    },
    
    stopCommentsPolling() {
      if (this.reactionPollingTimer) {
        clearInterval(this.reactionPollingTimer);
        this.reactionPollingTimer = null;
      }
    },
    
    handleReactionAdded(reaction) {
      this.reactions.push(reaction);
    },
    
    handleReactionUpdated(updatedReaction) {
      const index = this.reactions.findIndex(r => r.username === updatedReaction.username);
      if (index !== -1) {
        this.reactions[index].emoji = updatedReaction.emoji;
      }
    },
    
    handleReactionRemoved(removedReaction) {
      this.reactions = this.reactions.filter(r => r.username !== removedReaction.username);
    },
    
    handleReactClick() {
      this.closeContextMenu();
      this.$refs.emojiReaction.openEmojiPicker();
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
      @contextmenu="showContextMenuHandler"
    >
      <div v-if="message.isForwarded" class="forwarded-label">
        <i class="fas fa-share"></i> Forwarded
      </div>
          
    <div v-if="message.replyTo" class="reply-info mb-2">
      <div class="reply-preview">
        <div class="reply-sender">
          {{ message.replyTo.senderName}}
        </div>
        <div class="reply-content">
          <template v-if="message.replyTo.content">
            {{ message.replyTo.content}}
          </template>
          <template v-else>
            <i class="fas fa-ban me-1"></i> Messaggio eliminato
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
              @contextmenu="showContextMenuHandler"
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
      
      <EmojiReaction 
        ref="emojiReaction"
        :reactions="reactions"
        :message-id="message.id"
        :conversation-id="conversationId"
        :my-username="myUsername"
        @reaction-added="handleReactionAdded"
        @reaction-updated="handleReactionUpdated"
        @reaction-removed="handleReactionRemoved"
      />
    </div>
    

    <Teleport to="body">
      <div 
        v-if="showContextMenu" 
        ref="contextMenu"
        class="message-context-menu"
        :style="{
          top: `${contextMenuPosition.y}px`,
          left: `${contextMenuPosition.x}px`
        }"
        @click.stop
        @contextmenu.prevent
      >
        <div class="menu-item" @click="handleReactClick">
          <i class="fas fa-smile me-2 text-warning"></i>
          <span>Reagisci</span>
        </div>
        
        <div class="menu-item" @click="handleReplyClick">
          <i class="fas fa-reply me-2 text-success"></i>
          <span>Rispondi</span>
        </div>
        
        <div class="menu-item" @click="handleForwardClick">
          <i class="fas fa-share me-2 text-info"></i>
          <span>Inoltra</span>
        </div>
        
        <div v-if="isMyMessage" class="menu-item delete" @click="handleDelete">
          <i class="fas fa-trash me-2 text-danger"></i>
          <span>Elimina</span>
        </div>
      </div>
    </Teleport>
    
    
    <ForwardModal
      v-if="showForwardModal"
      :show="showForwardModal"
      :conversations="availableConversations"
      :message-id="message.id"
      @close="closeForwardModal"
      @forward="handleForward"
    />
    
    
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

.message-context-menu {
  position: fixed;
  z-index: 9999;
  background-color: #202c33;
  border-radius: 6px;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.3);
  min-width: 180px;
  overflow: hidden;
  user-select: none;
  border: 1px solid #36434a;
  animation: fadeInMenu 0.2s ease-out;
}

@keyframes fadeInMenu {
  from {
    opacity: 0;
    transform: scale(0.9);
  }
  to {
    opacity: 1;
    transform: scale(1);
  }
}

.menu-item {
  padding: 12px 16px;
  cursor: pointer;
  color: #e9edef;
  font-size: 0.9rem;
  display: flex;
  align-items: center;
  transition: background-color 0.15s;
}

.menu-item:hover {
  background-color: #2a3942;
}

.menu-item.delete {
  border-top: 1px solid #36434a;
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