// Turning on theme toggler as soon as the page loads
initTheme();

let token = '';
let ws = null;

// DOM elements
const authSection = document.getElementById('authSection');
const chatSection = document.getElementById('chatSection');
const messages = document.getElementById('messages');
const messageInput = document.getElementById('messageInput');
const sendButton = document.getElementById('sendButton');

// Api func
async function register() {
    const username = document.getElementById('regUsername').value;
    const password = document.getElementById('regPassword').value;
    
    const response = await fetch('/register', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ username, password }),
    });
    
    if (response.ok) {
        alert('Registration successful! Please login.');
    } else {
        alert('Registration failed: ' + (await response.text()));
    }
}

async function login() {
    const username = document.getElementById('loginUsername').value;
    const password = document.getElementById('loginPassword').value;
    
    const response = await fetch('/login', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ username, password }),
    });
    
    if (response.ok) {
        const data = await response.json();
        token = data.token;
        
        authSection.classList.add('hidden');
        chatSection.classList.remove('hidden');
        
        connectWebSocket();
        loadMessageHistory();
        setupReply();
    } else {
        alert('Login failed: ' + (await response.text()));
    }
}

function connectWebSocket() {
    const protocol = window.location.protocol === 'https:' ? 'wss://' : 'ws://';
    const wsUrl = protocol + window.location.host + '/ws';
    
    console.log('Connecting to WebSocket:', wsUrl);
    ws = new WebSocket(wsUrl);

    // Reconnection timeout
    const reconnectTimeout = 3000;

    ws.onopen = function() {
        console.log('WebSocket connected');
        const token = localStorage.getItem('token');
        if (token) {
            console.log(token);
        }
    };

    ws.onmessage = function(event) {
        try {
            const msg = JSON.parse(event.data);
            const isCurrentUser = msg.username === document.getElementById('loginUsername').value;
            
            const messageElement = document.createElement('div');
            messageElement.className = isCurrentUser 
                ? 'message my-message' 
                : 'message other-message';
            
            messageElement.innerHTML = `
                <strong>${isCurrentUser ? 'Вы' : msg.username}</strong>
                <div class="text">${msg.text}</div>
                <span class="time">${new Date().toLocaleTimeString()}</span>
            `;
            
            document.getElementById('messages').appendChild(messageElement);
            scrollToBottom('auto');
        } catch (e) {
            console.error('Error parsing message:', e);
        }
    }

    ws.onerror = function(error) {
        console.error('WebSocket error:', error);
        setTimeout(connectWebSocket, reconnectTimeout);
    };

    ws.onclose = function(event) {
        console.log(`WebSocket closed: ${event.code} ${event.reason}`);
        if (event.code === 1006) {
            setTimeout(connectWebSocket, reconnectTimeout);
        }
    };
}

function sendMessage() {
    const text = messageInput.value.trim();
    if (!text || !ws || ws.readyState !== WebSocket.OPEN) {
        console.log('Cannot send message:', {text, wsReady: ws?.readyState});
        return;
    }
    
    try {
        ws.send(text); 
        messageInput.value = '';
    } catch (e) {
        console.error('Error sending message:', e);
    }
}

async function loadMessageHistory() {
    try {
        const currentUsername = document.getElementById('loginUsername').value;
        const response = await fetch('/messages', {
            headers: {
                'Authorization': 'Bearer ' + localStorage.getItem('token')
            }
        });
        
        const messages = await response.json();
        const messagesContainer = document.getElementById('messages');
        messagesContainer.innerHTML = '';
        
        messages.forEach(msg => {
            const isCurrentUser = msg.username === currentUsername;
            
            const messageElement = document.createElement('div');
            messageElement.className = isCurrentUser 
                ? 'message my-message' 
                : 'message other-message';
            
            messageElement.innerHTML = `
                <strong>${isCurrentUser ? 'Вы' : msg.username}</strong>
                <div class="text">${msg.text}</div>
                <span class="time">${formatTime(msg.time)}</span>
            `;
            
            messagesContainer.appendChild(messageElement);
        });
        scrollToBottom('auto');
    } catch (error) {
        console.error('Error loading messages:', error);
    }
}

// Additional func for formatting timestamps
function formatTime(isoString) {
    const date = new Date(isoString);
    return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
}

function addMessageToChat(msg) {
    const messageElement = document.createElement('div');
    const isCurrentUser = msg.username === document.getElementById('loginUsername').value;
    
    messageElement.className = isCurrentUser ? 'message my-message' : 'message other-message';
    messageElement.innerHTML = `
        <strong>${isCurrentUser ? 'Вы' : msg.username}</strong>
        <div class="text">${msg.text}</div>
        <span class="time">${new Date(msg.time).toLocaleTimeString()}</span>
    `;
    
    messages.appendChild(messageElement);
    scrollToBottom('auto');
}

function scrollToBottom(behavior = 'smooth') {
    const messagesContainer = document.getElementById('messages');
    if (!messagesContainer) return;
    
    requestAnimationFrame(() => {
        messagesContainer.scrollTo({
            top: messagesContainer.scrollHeight,
            behavior: behavior
        });
    });
    
    console.log('Scrolling to bottom:', messagesContainer.scrollHeight); // logs
}

// Reply functions
function setupReply() {
    let replyingTo = null;
    
    document.addEventListener('click', function(e) {
        if (e.target.closest('.reply-btn')) {
            const msgElement = e.target.closest('.message');
            replyingTo = {
                id: msgElement.dataset.id,
                username: msgElement.dataset.username,
                text: msgElement.querySelector('.text').textContent
            };
            showReplyPreview();
        }
    });
    
    function showReplyPreview() {
        const preview = document.createElement('div');
        preview.className = 'reply-preview';
        preview.innerHTML = `
            <div>Replying to ${replyingTo.username}: ${replyingTo.text.substring(0, 30)}...</div>
            <button class="cancel-reply">×</button>
        `;
        document.getElementById('messageInput').before(preview);
    }
}

// Theme changer
function toggleTheme() {
    const currentTheme = document.documentElement.getAttribute('data-theme') || 'dark';
    const newTheme = currentTheme === 'dark' ? 'light' : 'dark';
    
    document.documentElement.setAttribute('data-theme', newTheme);
    localStorage.setItem('theme', newTheme);
    
    // Icon switch
    document.querySelector('.theme-icon').textContent = newTheme === 'dark' ? '🌙' : '☀️';
}

// Theme initialization
function initTheme() {
    const savedTheme = localStorage.getItem('theme') || 'dark';
    document.documentElement.setAttribute('data-theme', savedTheme);
    document.querySelector('.theme-icon').textContent = savedTheme === 'dark' ? '🌙' : '☀️';
}

// Authentication handlers
document.getElementById('registerBtn').addEventListener('click', register);
document.getElementById('loginBtn').addEventListener('click', login);

// Auto-scrolling
window.addEventListener('load', () => scrollToBottom('auto'));

// Sending messages by pressing ENTER
sendButton.addEventListener('click', sendMessage);
messageInput.addEventListener('keypress', function(e) {
    if (e.key === 'Enter') {
        sendMessage();
    }
});

// Theme Toggler
document.getElementById('themeToggle').addEventListener('click', toggleTheme);