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
      availableEmojis: [
        { emoji: '👍', color: '#FFD767' },
        { emoji: '❤️', color: '#F7574D' },
        { emoji: '😂', color: '#FECD3D' },
        { emoji: '😮', color: '#6EB1FF' },
        { emoji: '😢', color: '#6EB1FF' },
        { emoji: '🙏', color: '#8A65CF' },
        { emoji: '🎉', color: '#FF8F5F' },
        { emoji: '🔥', color: '#FF5F5F' }
      ],
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
    
    async selectEmoji(emojiObj) {
      this.closeEmojiPicker();
      
      if (this.isAddingReaction) return;
      
      try {
        this.isAddingReaction = true;
        const emoji = emojiObj.emoji;
        
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
    
    getReactionStyle(emoji) {
      const emojiObj = this.availableEmojis.find(e => e.emoji === emoji);
      
      if (emojiObj) {
        if (this.reactions.some(r => r.username === this.myUsername && r.emoji === emoji)) {
          return {
            backgroundColor: emojiObj.color + '40',
            borderLeft: `3px solid ${emojiObj.color}`
          };
        } else {
          return {
            backgroundColor: emojiObj.color + '20'
          };
        }
      }
      
      return {};
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
        class="reaction-badge-wrapper"
      >
        <div 
          class="reaction-badge"
          :class="{ 'my-reaction': reactors.some(r => r.username === myUsername) }"
          :title="showReactors(emoji)"
          @click="handleReactionBadgeClick(emoji)"
          :style="getReactionStyle(emoji)"
        >
          <span class="emoji">{{ emoji }}</span>
          <span class="count">{{ reactors.length }}</span>
        </div>
        
        <div class="reactors-dropdown">
          <div class="reactors-list">
            <div v-for="reactor in reactors" :key="reactor.id" class="reactor-item">
              {{ reactor.username }}
            </div>
          </div>
        </div>
      </div>
    </div>
    
    <Teleport to="body">
      <div v-if="showEmojiPicker" class="emoji-modal-overlay" @click="closeEmojiPicker">
        <div class="emoji-modal" @click.stop>
          <div class="emoji-modal-header">
            <span>Scegli una reazione</span>
            <button class="close-button" @click="closeEmojiPicker">
              <i class="fas fa-times"></i>
            </button>
          </div>
          <div class="emoji-container">
            <button 
              v-for="emojiObj in availableEmojis" 
              :key="emojiObj.emoji"
              class="emoji-button"
              :style="{ backgroundColor: emojiObj.color + '30' }"
              @click="selectEmoji(emojiObj)"
              :title="emojiObj.label"
            >
              <span class="emoji-circle" :style="{ backgroundColor: emojiObj.color + '50' }">
                {{ emojiObj.emoji }}
              </span>
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

.reaction-badge-wrapper {
  position: relative;
  display: inline-block;
  margin-right: 4px;
  margin-bottom: 4px;
}

.reaction-badge-wrapper:hover .reactors-dropdown {
  display: block;
}

.reaction-badge {
  background-color: #2a3942;
  border-radius: 10px;
  padding: 3px 8px;
  display: flex;
  align-items: center;
  font-size: 0.85rem;
  cursor: pointer;
  transition: all 0.2s ease;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.reaction-badge:hover {
  transform: scale(1.05);
}

.my-reaction {
  border-left: 3px solid #00a884;
}

.emoji {
  font-size: 1.1rem;
  margin-right: 5px;
}

.count {
  font-size: 0.8rem;
  font-weight: 600;
  color: rgba(255, 255, 255, 0.9);
}

.reactors-dropdown {
  display: none;
  position: absolute;
  bottom: 100%;
  left: 50%;
  transform: translateX(-50%);
  background-color: #202c33;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
  padding: 8px;
  margin-bottom: 8px;
  z-index: 10;
  min-width: 120px;
  max-width: 200px;
}

.reactors-dropdown:after {
  content: '';
  position: absolute;
  bottom: -5px;
  left: 50%;
  transform: translateX(-50%);
  width: 0;
  height: 0;
  border-left: 6px solid transparent;
  border-right: 6px solid transparent;
  border-top: 6px solid #202c33;
}

.reactor-item {
  padding: 4px 8px;
  font-size: 0.85rem;
  color: #e9edef;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
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
  width: 350px;
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
  gap: 12px;
}

.emoji-button {
  display: flex;
  align-items: center;
  justify-content: center;
  background: transparent;
  border: none;
  border-radius: 12px;
  cursor: pointer;
  padding: 8px;
  transition: all 0.2s ease;
  width: 60px;
  height: 60px;
}

.emoji-button:hover {
  transform: scale(1.1);
}

.emoji-circle {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 44px;
  height: 44px;
  border-radius: 50%;
  font-size: 1.8rem;
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