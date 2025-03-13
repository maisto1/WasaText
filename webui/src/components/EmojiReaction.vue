<script>
export default {
  props: {
    reactions: {
      type: Array,
      default: () => []
    },
    messageId: {
      type: Number,
      required: true
    },
    conversationId: {
      type: Number,
      required: true
    },
    myUsername: {
      type: String,
      required: true
    }
  },
  
  data() {
    return {
      showEmojiPicker: false,
      availableEmojis: ['👍', '❤️', '😂', '😮', '😢', '🙏'],
      isAddingReaction: false
    }
  },
  
  computed: {
    groupedReactions() {
      const grouped = {};
      
      this.reactions.forEach(reaction => {
        if (!grouped[reaction.emoji]) {
          grouped[reaction.emoji] = [];
        }
        grouped[reaction.emoji].push(reaction);
      });
      
      return grouped;
    },
    
    hasReacted() {
      return this.reactions.some(reaction => reaction.username === this.myUsername);
    },
    
    myReaction() {
      return this.reactions.find(reaction => reaction.username === this.myUsername);
    }
  },
  
  methods: {
    openEmojiPicker() {
      this.showEmojiPicker = true;
    },
    
    closeEmojiPicker() {
      this.showEmojiPicker = false;
    },
    
    async selectEmoji(emoji) {
      this.closeEmojiPicker();
      
      if (this.isAddingReaction) return;
      
      try {
        this.isAddingReaction = true;
        
        if (this.hasReacted) {
          if (this.myReaction.emoji === emoji) {
            await this.removeReaction();
          } else {
            await this.updateReaction(emoji);
          }
        } else {
    
          await this.addReaction(emoji);
        }
      } catch (error) {
        console.error('Error handling reaction:', error);
      } finally {
        this.isAddingReaction = false;
      }
    },
    

    async handleReactionBadgeClick(emoji) {
      if (this.isAddingReaction) return;
      
      try {
        this.isAddingReaction = true;
        
        const userReactionWithEmoji = this.reactions.find(r => 
          r.username === this.myUsername && r.emoji === emoji
        );
        
        if (userReactionWithEmoji) {
 
          await this.removeReaction();
        } else if (this.hasReacted) {
         
          await this.updateReaction(emoji);
        } else {
         
          await this.addReaction(emoji);
        }
      } catch (error) {
        console.error('Error handling reaction badge click:', error);
      } finally {
        this.isAddingReaction = false;
      }
    },
    

    async addReaction(emoji) {
      try {

        const response = await this.$axios.post(`/conversations/${this.conversationId}/messages/${this.messageId}/comments/`, {
          content: `reaction:${emoji}`
        });
        
        this.$emit('reaction-added', {
          id: response.data.id,
          messageId: this.messageId,
          emoji: emoji,
          username: this.myUsername,
          timestamp: response.data.timestamp
        });
      } catch (error) {
        console.error('Error adding reaction:', error);
      }
    },
    

    async updateReaction(emoji) {
      try {

        await this.removeReaction();
        

        const response = await this.$axios.post(`/conversations/${this.conversationId}/messages/${this.messageId}/comments/`, {
          content: `reaction:${emoji}`
        });
        
      
        this.$emit('reaction-updated', {
          id: response.data.id,
          messageId: this.messageId,
          emoji: emoji,
          username: this.myUsername,
          previousEmoji: this.myReaction.emoji,
          timestamp: response.data.timestamp
        });
      } catch (error) {
        console.error('Error updating reaction:', error);
      }
    },
    
   
    async removeReaction() {
      try {
        if (!this.myReaction || !this.myReaction.id) {
          console.error('Cannot remove reaction: reaction ID not found');
          return;
        }
        
        
        await this.$axios.delete(`/conversations/${this.conversationId}/messages/${this.messageId}/comments/${this.myReaction.id}`);
        
       
        this.$emit('reaction-removed', {
          messageId: this.messageId,
          username: this.myUsername,
          emoji: this.myReaction.emoji
        });
      } catch (error) {
        console.error('Error removing reaction:', error);
      }
    },
    
    
    showReactors(emoji) {
      if (!this.groupedReactions[emoji] || this.groupedReactions[emoji].length === 0) return '';
      
      const reactors = this.groupedReactions[emoji].map(r => r.username);
      
      if (reactors.length === 1) {
        return reactors[0];
      } else if (reactors.length === 2) {
        return `${reactors[0]} e ${reactors[1]}`;
      } else {
        return `${reactors[0]}, ${reactors[1]} e altri ${reactors.length - 2}`;
      }
    }
  }
}
</script>

<template>
  <div class="message-reactions">

    <div v-if="Object.keys(groupedReactions).length > 0" class="reaction-summary">
      <div 
        v-for="(reactors, emoji) in groupedReactions" 
        :key="emoji"
        class="reaction-badge"
        :class="{ 'my-reaction': reactors.some(r => r.username === myUsername) }"
        :title="showReactors(emoji)"
        @click="handleReactionBadgeClick(emoji)"
      >
        <span class="emoji">{{ emoji }}</span>
        <span class="count">{{ reactors.length }}</span>
      </div>
    </div>
    
    
    <Teleport to="body">
      <div v-if="showEmojiPicker" class="emoji-modal-overlay" @click="closeEmojiPicker">
        <div class="emoji-modal" @click.stop>
          <div class="emoji-modal-header">
            <span>Choose a reaction</span>
            <button class="close-button" @click="closeEmojiPicker">
              <i class="fas fa-times"></i>
            </button>
          </div>
          <div class="emoji-container">
            <button 
              v-for="emoji in availableEmojis" 
              :key="emoji"
              class="emoji-button"
              @click="selectEmoji(emoji)"
            >
              {{ emoji }}
            </button>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<style scoped>
.message-reactions {
  display: flex;
  align-items: center;
  margin-top: 4px;
}

.reaction-summary {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
}

.reaction-badge {
  background-color: #2a3942;
  border-radius: 10px;
  padding: 2px 6px;
  display: flex;
  align-items: center;
  font-size: 0.85rem;
  cursor: pointer;
  transition: background-color 0.2s;
}

.reaction-badge:hover {
  background-color: #36434a;
}

.my-reaction {
  background-color: #045d4c;
}

.my-reaction:hover {
  background-color: #06735e;
}

.emoji {
  font-size: 1rem;
  margin-right: 3px;
}

.count {
  font-size: 0.75rem;
  color: rgba(255, 255, 255, 0.7);
}


.emoji-modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 2000;
  backdrop-filter: blur(3px);
}

.emoji-modal {
  background-color: #202c33;
  border-radius: 12px;
  width: 300px;
  max-width: 90vw;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.5);
  animation: modalFadeIn 0.2s ease-out;
}

.emoji-modal-header {
  padding: 12px 16px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-bottom: 1px solid #36434a;
}

.emoji-modal-header span {
  font-size: 0.95rem;
  font-weight: 500;
  color: #e9edef;
}

.close-button {
  background: transparent;
  border: none;
  color: #8696a0;
  width: 24px;
  height: 24px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: background-color 0.2s;
}

.close-button:hover {
  background-color: #36434a;
  color: #e9edef;
}

.emoji-container {
  display: flex;
  flex-wrap: wrap;
  padding: 16px;
  justify-content: center;
  gap: 8px;
}

.emoji-button {
  background: transparent;
  border: none;
  font-size: 1.8rem;
  padding: 10px;
  border-radius: 50%;
  cursor: pointer;
  transition: background-color 0.2s;
  height: 50px;
  width: 50px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.emoji-button:hover {
  background-color: #36434a;
}

@keyframes modalFadeIn {
  from {
    opacity: 0;
    transform: scale(0.9);
  }
  to {
    opacity: 1;
    transform: scale(1);
  }
}
</style>