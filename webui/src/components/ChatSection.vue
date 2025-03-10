<script>
import MessageBubble from './MessageBubble.vue';
import ChatInput from './ChatInput.vue';

export default {
  components: {
    MessageBubble,
    ChatInput
  },

  props: {
    conversation: {
      type: Object,
      required: true
    },
    isTempChat: {
      type: Boolean,
      default: false
    }
  },

  data() {
    return {
      messages: [],
      loading: false,
      error: null,
      showToast: false,
      toastMessage: '',
      toastType: 'success',
      sendingMessage: false,
      isGroupChat: false
    }
  },

  created() {
    this.checkIfGroupChat();
  },

  watch: {
    'conversation': {
      immediate: true,
      deep: true,
      handler() {
        this.checkIfGroupChat();
        if (this.conversation?.id && !this.isTempChat) {
          this.fetchMessages();
        } else if (this.isTempChat) {
          this.messages = [];
        }
      }
    }
  },

  methods: {
    checkIfGroupChat() {
      if (this.conversation) {
        if (this.conversation.conversationType) {
          this.isGroupChat = this.conversation.conversationType === 'group';
        }
        
        else if (this.conversation.name) {
          this.isGroupChat = this.conversation.name.toLowerCase().includes('group');
        }
        
        else if (this.conversation.partecipants && this.conversation.partecipants.length > 2) {
          this.isGroupChat = true;
        }
      }
    },

    async fetchMessages() {
      if (!this.conversation?.id || this.isTempChat) return;
      
      this.loading = true;
      try {
        console.log('Fetching messages for conversation:', this.conversation.id);
        const response = await this.$axios.get(`/conversations/${this.conversation.id}`);
        this.messages = response.data;
        this.scrollToBottom();
      } catch (error) {
        console.error('Error fetching messages:', error);
        this.showNotification('Failed to load messages', 'error');
      } finally {
        this.loading = false;
      }
    },

    async sendMessage(messageContent) {
      if (!messageContent.trim() || this.sendingMessage) return;
      
      this.sendingMessage = true;
      
      try {

        if (this.isTempChat) {
          console.log('Sending first message in temporary chat:', messageContent);
          this.$emit('send-first-message', messageContent);
          
        } 

        else {
          console.log('Sending message to existing chat:', this.conversation.id);
          const response = await this.$axios.post(`/conversations/${this.conversation.id}/messages/`, {
            type: "text",
            content: messageContent
          });
          
          this.messages.push(response.data);
          this.scrollToBottom();
        }
      } catch (error) {
        console.error('Error sending message:', error);
        this.showNotification('Failed to send message', 'error');
      } finally {
        this.sendingMessage = false;
      }
    },

    async deleteMessage(messageId) {
      try {
        await this.$axios.delete(`/conversations/${this.conversation.id}/messages/${messageId}`);
        this.messages = this.messages.filter(msg => msg.id !== messageId);
        this.showNotification('Message deleted successfully', 'success');
      } catch (error) {
        console.error('Error deleting message:', error);
        this.showNotification('Failed to delete message', 'error');
      }
    },

    async forwardMessage(messageId, targetConversationId) {
      try {
        await this.$axios.post(
          `/conversations/${this.conversation.id}/messages/${messageId}`,
          { conversationId: targetConversationId }
        );
        this.showNotification('Message forwarded successfully', 'success');
      } catch (error) {
        console.error('Error forwarding message:', error);
        this.showNotification('Failed to forward message', 'error');
      }
    },

    scrollToBottom() {
      this.$nextTick(() => {
        const container = this.$refs.messagesContainer;
        if (container) {
          container.scrollTop = container.scrollHeight;
        }
      });
    },

    showNotification(message, type) {
      this.toastMessage = message;
      this.toastType = type;
      this.showToast = true;
      setTimeout(() => {
        this.showToast = false;
      }, 3000);
    },

    isMyMessage(message) {
      const myUsername = sessionStorage.getItem('username');
      return message.sender.username === myUsername;
    }
  }
}
</script>

<template>
  <div class="chat-section d-flex flex-column h-100">
    <div class="chat-header p-3 d-flex align-items-center border-bottom">
      <div class="d-flex align-items-center">
        <div class="profile-photo rounded-circle bg-secondary" style="width: 40px; height: 40px;">
          <img 
            v-if="conversation.conversationPhoto" 
            :src="'data:image/jpeg;base64,' + conversation.conversationPhoto"
            class="rounded-circle w-100 h-100"
            alt="Profile"
          />
        </div>
        <div class="ms-3">
          <h6 class="mb-0 text-white">{{ conversation.name }}</h6>
          <small v-if="isTempChat" class="text-light-grey">New chat</small>
        </div>
      </div>
    </div>

    <div ref="messagesContainer" class="messages-container flex-grow-1 p-3 overflow-auto">
      <div v-if="loading" class="d-flex justify-content-center align-items-center h-100">
        <div class="spinner-border text-light" role="status">
          <span class="visually-hidden">Loading...</span>
        </div>
      </div>

      <template v-else-if="isTempChat">
        <div class="text-center text-light-grey py-5">
          <p>Type a message to start the conversation</p>
        </div>
      </template>

      <template v-else>
        <MessageBubble
          v-for="message in messages"
          :key="message.id"
          :message="message"
          :isMyMessage="isMyMessage(message)"
          :showUsername="isGroupChat"
          @delete="deleteMessage"
          @forward="forwardMessage"
        />
      </template>
    </div>

    <ChatInput @send="sendMessage" :disabled="sendingMessage" />

    <transition name="toast">
      <div v-if="showToast" class="toast-container position-fixed p-3">
        <div class="custom-toast" :class="toastType">
          <div class="toast-content">
            <i v-if="toastType === 'success'" class="fas fa-check-circle toast-icon"></i>
            <i v-else class="fas fa-exclamation-circle toast-icon"></i>
            <div class="toast-message">{{ toastMessage }}</div>
          </div>
          <div class="toast-progress" :class="toastType"></div>
        </div>
      </div>
    </transition>
  </div>
</template>

<style>
.chat-section {
  background-color: #111b21;
}

.chat-header {
  background-color: #202c33;
}

.messages-container {
  background-color: #0b141a;
}

.messages-container::-webkit-scrollbar {
  width: 6px;
}

.messages-container::-webkit-scrollbar-track {
  background: #0b141a;
}

.messages-container::-webkit-scrollbar-thumb {
  background-color: #374045;
  border-radius: 6px;
}

.text-light-grey {
  color: #8696a0;
}

.toast-enter-active {
  animation: slideIn 0.3s ease-out;
}

.toast-leave-active {
  animation: slideOut 0.3s ease-in;
}

@keyframes slideIn {
  from {
    transform: translateX(100%);
    opacity: 0;
  }
  to {
    transform: translateX(0);
    opacity: 1;
  }
}

@keyframes slideOut {
  from {
    transform: translateX(0);
    opacity: 1;
  }
  to {
    transform: translateX(100%);
    opacity: 0;
  }
}
</style>