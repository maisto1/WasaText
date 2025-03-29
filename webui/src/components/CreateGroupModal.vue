<script>
export default {
  props: {
    show: {
      type: Boolean,
      required: true
    }
  },
  
  data() {
    return {
      groupName: '',
      searchQuery: '',
      searchResults: [],
      selectedUsers: [],
      isCreating: false,
      error: null
    }
  },
  
  watch: {
    show(newVal) {
      if (newVal) {
        this.resetForm();
      }
    }
  },
  
  methods: {
    async searchUsers() {
      if (!this.searchQuery.trim()) {
        this.searchResults = [];
        return;
      }
      
      try {
        const response = await this.$axios.get(`/users/?username=${this.searchQuery}`);
        this.searchResults = response.data.filter(user => 
          user.username !== sessionStorage.getItem('username') && 
          !this.selectedUsers.some(selectedUser => selectedUser.id === user.id)
        );
      } catch (error) {
        console.error('Error searching users for group:', error);
      }
    },

    addUser(user) {
      this.selectedUsers.push(user);
      this.searchResults = this.searchResults.filter(u => u.id !== user.id);
      this.searchQuery = '';
    },

    removeUser(user) {
      this.selectedUsers = this.selectedUsers.filter(u => u.id !== user.id);
    },
    
    resetForm() {
      this.groupName = '';
      this.searchQuery = '';
      this.searchResults = [];
      this.selectedUsers = [];
      this.error = null;
    },
    
    closeModal() {
      this.$emit('close');
    },
    
    async createGroup() {
      if (!this.groupName.trim()) {
        this.error = 'Please enter a group name';
        return;
      }
      
      if (this.groupName.length < 3 || this.groupName.length > 16) {
        this.error = 'Group name must be between 3 and 16 characters';
        return;
      }
      
      this.isCreating = true;
      this.error = null;
      
      try {
        const payload = {
          groupName: this.groupName,
          conversationType: "group",
          partecipant: ""
        };
        
        const response = await this.$axios.post('/conversations/', payload);
        
        const groupId = response.data.id;
        
        if (this.selectedUsers.length > 0) {
          const addPromises = [];
          
          for (const user of this.selectedUsers) {
            addPromises.push(
              this.$axios.post(`/conversations/${groupId}/members/`, {
                username: user.username
              }).catch(error => {
                console.error(`Error adding user ${user.username} to group:`, error);
                return null;
              })
            );
          }
          
          await Promise.all(addPromises);
        }
        
        this.$emit('created');
        this.closeModal();
        this.resetForm();
      } catch (error) {
        console.error('Error creating group:', error);
        if (error.response) {
          console.error('Response status:', error.response.status);
          console.error('Response data:', error.response.data);
          this.error = `Error (${error.response.status}): Failed to create group`;
        } else if (error.request) {
          console.error('No response received:', error.request);
          this.error = 'No response received from server. Please try again.';
        } else {
          console.error('Error message:', error.message);
          this.error = 'Failed to create group. Please try again.';
        }
      } finally {
        this.isCreating = false;
      }
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
  <transition name="modal-fade">
    <div v-if="show" class="modal-overlay" @click="closeModal">
      <div class="modal-dialog create-group-modal" @click.stop>
        <div class="modal-content">
        
   
          <div class="modal-header">
            <div class="d-flex align-items-center">
              <div class="header-icon-container">
                <i class="fas fa-users"></i>
              </div>
              <h5 class="modal-title">Create New Group</h5>
            </div>
            <button @click="closeModal" class="close-button">
              <i class="fas fa-times"></i>
            </button>
          </div>
          
          <div class="modal-body">

            <div v-if="error" class="error-message mb-3">
              <i class="fas fa-exclamation-circle me-2"></i>
              {{ error }}
            </div>
            

            <div class="group-name-container mb-4">
              <label class="form-label mb-2">Group Name</label>
              <div class="input-group">
                <span class="input-group-text">
                  <i class="fas fa-user-friends"></i>
                </span>
                <input 
                  v-model="groupName" 
                  type="text" 
                  class="form-control"
                  placeholder="Enter group name"
                  :maxlength="16"
                  :minlength="3"
                  :disabled="isCreating"
                >
              </div>
              <small class="hint-text">Group name must be between 3-16 characters</small>
            </div>
            
  
            <div class="group-participants-container mb-4">
              <label class="form-label d-flex justify-content-between align-items-center mb-2">
                <span>Add Participants (Optional)</span>
                <span class="participants-counter">
                  {{ selectedUsers.length }} selected
                </span>
              </label>
              

              <div class="input-group mb-3">
                <span class="input-group-text">
                  <i class="fas fa-search"></i>
                </span>
                <input 
                  v-model="searchQuery"
                  type="text"
                  class="form-control"
                  placeholder="Search users to add..."
                  @input="searchUsers"
                  :disabled="isCreating"
                >
              </div>
        
              <div v-if="selectedUsers.length > 0" class="selected-users-container mb-3">
                <div 
                  v-for="user in selectedUsers" 
                  :key="user.id"
                  class="selected-user-item"
                >
                  <div class="d-flex align-items-center">
                    <div class="avatar-container" :style="{ backgroundColor: getAvatarColor(user.username) }">
                      <img 
                        v-if="user.profilePhoto" 
                        :src="'data:image/jpeg;base64,' + user.profilePhoto"
                        class="avatar-image"
                        alt="Profile photo"
                      >
                      <div v-else class="avatar-text">
                        {{ getInitials(user.username) }}
                      </div>
                    </div>
                    <div class="user-info ms-2">
                      <div class="username">{{ user.username }}</div>
                    </div>
                    <button 
                      @click="removeUser(user)" 
                      class="btn-remove-user ms-auto"
                      title="Remove from group"
                      :disabled="isCreating"
                    >
                      <i class="fas fa-times-circle"></i>
                    </button>
                  </div>
                </div>
              </div>
              
     
              <div v-if="searchResults.length > 0" class="search-results-container">
                <div 
                  v-for="user in searchResults" 
                  :key="user.id"
                  class="search-result-item"
                  @click="addUser(user)"
                >
                  <div class="d-flex align-items-center">
                    <div class="avatar-container" :style="{ backgroundColor: getAvatarColor(user.username) }">
                      <img 
                        v-if="user.profilePhoto" 
                        :src="'data:image/jpeg;base64,' + user.profilePhoto"
                        class="avatar-image"
                        alt="Profile photo"
                      >
                      <div v-else class="avatar-text">
                        {{ getInitials(user.username) }}
                      </div>
                    </div>
                    <div class="user-info ms-2">
                      <div class="username">{{ user.username }}</div>
                    </div>
                    <i class="fas fa-plus-circle add-icon ms-auto"></i>
                  </div>
                </div>
              </div>
              
              <div v-else-if="searchQuery && !searchResults.length" class="empty-results">
                <i class="fas fa-search mb-2"></i>
                <div>No users found</div>
              </div>
            </div>
          </div>
          
          <div class="modal-footer">
            <button @click="closeModal" class="btn-cancel" :disabled="isCreating">
              Cancel
            </button>
            <button 
              @click="createGroup" 
              class="btn-create"
              :disabled="!groupName.trim() || groupName.length < 3 || isCreating"
            >
              <span v-if="!isCreating">Create Group</span>
              <span v-else class="creating-spinner"></span>
            </button>
          </div>
          
        </div>
      </div>
    </div>
  </transition>
</template>

<style scoped>

.modal-fade-enter-active, .modal-fade-leave-active {
  transition: all 0.3s ease;
}

.modal-fade-enter-from, .modal-fade-leave-to {
  opacity: 0;
  transform: translateY(20px);
}

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.7);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  backdrop-filter: blur(4px);
}

.create-group-modal {
  width: 450px;
  max-width: 92%;
  animation: modalFadeIn 0.3s ease;
}

@keyframes modalFadeIn {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.modal-content {
  border-radius: 16px;
  background-color: #111b21;
  border: 1px solid #2a3942;
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.5);
  overflow: hidden;
}

.modal-header {
  padding: 16px 20px;
  background-color: #202c33;
  border-bottom: 1px solid #2a3942;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-icon-container {
  width: 34px;
  height: 34px;
  min-width: 34px;
  border-radius: 50%;
  background-color: #00a884;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 12px;
  box-shadow: 0 3px 8px rgba(0, 168, 132, 0.3);
}

.header-icon-container i {
  color: white;
  font-size: 1rem;
}

.modal-title {
  font-size: 1.2rem;
  font-weight: 600;
  margin: 0;
  color: #e9edef;
}

.close-button {
  background: transparent;
  border: none;
  color: #8696a0;
  font-size: 1.1rem;
  cursor: pointer;
  transition: color 0.2s;
  padding: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border-radius: 50%;
}

.close-button:hover {
  color: #e9edef;
  background-color: #2a3942;
}


.modal-body {
  padding: 20px;
}

.error-message {
  background-color: rgba(231, 76, 60, 0.15);
  color: #e74c3c;
  padding: 12px;
  border-radius: 8px;
  font-size: 0.9rem;
  border-left: 3px solid #e74c3c;
}

.form-label {
  color: #e9edef;
  font-size: 0.95rem;
  font-weight: 500;
}

.participants-counter {
  background-color: #00a884;
  color: white;
  font-size: 0.8rem;
  padding: 4px 10px;
  border-radius: 12px;
  font-weight: 500;
}

.input-group-text {
  background-color: #2a3942;
  border: 1px solid #2a3942;
  color: #8696a0;
}

.form-control {
  background-color: #2a3942;
  border: 1px solid #2a3942;
  color: white;
}

.form-control:focus {
  background-color: #36434a;
  border-color: #00a884;
  box-shadow: 0 0 0 2px rgba(0, 168, 132, 0.25);
  color: white;
}

.form-control::placeholder {
  color: #8696a0;
}

.hint-text {
  color: #8696a0;
  font-size: 0.8rem;
  margin-top: 6px;
  display: block;
}


.selected-users-container {
  max-height: 160px;
  overflow-y: auto;
  background-color: #202c33;
  border-radius: 8px;
  padding: 6px;
}

.selected-user-item {
  padding: 8px 10px;
  margin-bottom: 4px;
  border-radius: 8px;
  background-color: #2a3942;
  transition: all 0.2s ease;
}

.selected-user-item:last-child {
  margin-bottom: 0;
}

.avatar-container {
  width: 36px;
  height: 36px;
  min-width: 36px;
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
  font-size: 1rem;
  font-weight: 600;
  text-align: center;
}

.user-info {
  overflow: hidden;
}

.username {
  color: #e9edef;
  font-size: 0.9rem;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.btn-remove-user {
  background: none;
  border: none;
  color: #e74c3c;
  font-size: 1rem;
  opacity: 0.7;
  transition: all 0.2s ease;
  padding: 0;
}

.btn-remove-user:hover {
  opacity: 1;
}

.btn-remove-user:disabled {
  opacity: 0.3;
  cursor: not-allowed;
}


.search-results-container {
  max-height: 200px;
  overflow-y: auto;
  background-color: #202c33;
  border-radius: 8px;
  padding: 6px;
}

.search-result-item {
  padding: 8px 10px;
  margin-bottom: 4px;
  border-radius: 8px;
  background-color: #2a3942;
  transition: all 0.2s ease;
  cursor: pointer;
}

.search-result-item:last-child {
  margin-bottom: 0;
}

.search-result-item:hover {
  background-color: #36434a;
}

.add-icon {
  color: #00a884;
  font-size: 1.1rem;
  transition: transform 0.2s ease;
}

.search-result-item:hover .add-icon {
  transform: scale(1.2);
}

.empty-results {
  text-align: center;
  color: #8696a0;
  padding: 20px;
  font-size: 0.9rem;
  background-color: #202c33;
  border-radius: 8px;
}

.empty-results i {
  font-size: 1.5rem;
  opacity: 0.7;
  display: block;
}

.modal-footer {
  padding: 16px 20px;
  background-color: #202c33;
  border-top: 1px solid #2a3942;
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

.btn-cancel {
  padding: 10px 16px;
  background-color: transparent;
  border: 1px solid #2a3942;
  border-radius: 8px;
  color: #e9edef;
  font-size: 0.95rem;
  font-weight: 500;
  cursor: pointer;
  transition: background-color 0.2s;
}

.btn-cancel:hover:not(:disabled) {
  background-color: #2a3942;
}

.btn-cancel:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-create {
  padding: 10px 16px;
  background-color: #00a884;
  border: none;
  border-radius: 8px;
  color: #fff;
  font-size: 0.95rem;
  font-weight: 500;
  cursor: pointer;
  transition: background-color 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;
  min-width: 130px;
}

.btn-create:hover:not(:disabled) {
  background-color: #008f6f;
}

.btn-create:disabled {
  background-color: rgba(0, 168, 132, 0.5);
  cursor: not-allowed;
}

.creating-spinner {
  display: inline-block;
  width: 18px;
  height: 18px;
  border-radius: 50%;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top-color: white;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}


.selected-users-container::-webkit-scrollbar,
.search-results-container::-webkit-scrollbar {
  width: 4px;
}

.selected-users-container::-webkit-scrollbar-track,
.search-results-container::-webkit-scrollbar-track {
  background: transparent;
}

.selected-users-container::-webkit-scrollbar-thumb,
.search-results-container::-webkit-scrollbar-thumb {
  background-color: #374045;
  border-radius: 4px;
}
</style>