class BrowserController {
    constructor() {
        this.canvas = document.getElementById('vnc');
        this.ctx = this.canvas.getContext('2d');
        this.ws = null;
        this.isRunning = false;

        // UI Elements
        this.startBtn = document.getElementById('startBtn');
        this.stopBtn = document.getElementById('stopBtn');
        this.urlInput = document.getElementById('urlInput');
        this.navigateBtn = document.getElementById('navigateBtn');
        this.reloadBtn = document.getElementById('reloadBtn');

        // Bind event listeners
        this.startBtn.addEventListener('click', () => this.startBrowser());
        this.stopBtn.addEventListener('click', () => this.stopBrowser());
        this.navigateBtn.addEventListener('click', () => this.navigate());
        this.reloadBtn.addEventListener('click', () => this.reload());
        
        // Initial UI state
        this.updateUIState(false);
        
        // Handle registration/auth
        this.handleAuth();
    }

    async handleAuth() {
        const token = localStorage.getItem('reaperToken');
        if (!token) {
            const username = document.querySelector('input').value;
            if (!username) return;

            try {
                const response = await fetch('/register', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ username })
                });
                
                if (!response.ok) throw new Error('Registration failed');
                
                const data = await response.json();
                if (data.user && data.user.token) {
                    localStorage.setItem('reaperToken', data.user.token);
                }
            } catch (error) {
                console.error('Auth error:', error);
            }
        }
    }

    async fetchWithAuth(url, options = {}) {
        const token = localStorage.getItem('reaperToken');
        if (!token) {
            throw new Error('No auth token available');
        }

        const headers = {
            ...options.headers,
            'X-Reaper-Token': token
        };

        return fetch(url, { ...options, headers });
    }

    async startBrowser() {
        try {
            const response = await this.fetchWithAuth('/api/browser/start', { method: 'POST' });
            if (!response.ok) throw new Error('Failed to start browser');
            
            this.isRunning = true;
            this.updateUIState(true);
            this.connectVNC();
        } catch (error) {
            console.error('Start browser error:', error);
            alert('Failed to start browser');
        }
    }

    async stopBrowser() {
        try {
            const response = await this.fetchWithAuth('/api/browser/stop', { method: 'POST' });
            if (!response.ok) throw new Error('Failed to stop browser');
            
            this.isRunning = false;
            this.updateUIState(false);
            if (this.ws) {
                this.ws.close();
                this.ws = null;
            }
        } catch (error) {
            console.error('Stop browser error:', error);
            alert('Failed to stop browser');
        }
    }

    async navigate() {
        const url = this.urlInput.value.trim();
        if (!url) return;

        try {
            const response = await this.fetchWithAuth('/api/browser/navigate', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ url })
            });
            if (!response.ok) throw new Error('Failed to navigate');
        } catch (error) {
            console.error('Navigation error:', error);
            alert('Failed to navigate to URL');
        }
    }

    async reload() {
        try {
            const response = await this.fetchWithAuth('/api/browser/reload', { method: 'POST' });
            if (!response.ok) throw new Error('Failed to reload');
        } catch (error) {
            console.error('Reload error:', error);
            alert('Failed to reload page');
        }
    }

    connectVNC() {
        const token = localStorage.getItem('reaperToken');
        // Use the current window's location to get the host and port
        const wsProtocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const wsUrl = `${wsProtocol}//${window.location.host}/browser/vnc`;
        this.ws = new WebSocket(`${wsUrl}?token=${token}`);
        
        this.ws.onopen = () => {
            console.log('VNC WebSocket connected');
        };

        this.ws.onmessage = (event) => {
            // Handle VNC frame data
            if (event.data instanceof Blob) {
                const reader = new FileReader();
                reader.onload = () => {
                    const img = new Image();
                    img.onload = () => {
                        this.ctx.drawImage(img, 0, 0, this.canvas.width, this.canvas.height);
                    };
                    img.src = reader.result;
                };
                reader.readAsDataURL(event.data);
            }
        };

        this.ws.onclose = () => {
            console.log('VNC WebSocket disconnected');
            this.ws = null;
        };

        this.ws.onerror = (error) => {
            console.error('VNC WebSocket error:', error);
        };
    }

    updateUIState(isRunning) {
        this.startBtn.disabled = isRunning;
        this.stopBtn.disabled = !isRunning;
        this.urlInput.disabled = !isRunning;
        this.navigateBtn.disabled = !isRunning;
        this.reloadBtn.disabled = !isRunning;
    }
}

// Initialize the browser controller when the page loads
window.addEventListener('DOMContentLoaded', () => {
    const controller = new BrowserController();
    
    // Handle sign in button click
    document.querySelector('button').addEventListener('click', () => {
        controller.handleAuth();
    });
});
