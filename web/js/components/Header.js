const { useState, useEffect, useCallback } = React;

// Header Component
function Header({ connectionStatus, onConnect, isLoading, activeTab, onTabChange, brightness, onBrightnessChange }) {
    const getStatusColor = () => {
        switch (connectionStatus) {
            case 'connected': return '#00ff88';
            case 'connecting': return '#ffaa00';
            case 'error': return '#ff4444';
            default: return '#666666';
        }
    };

    const getStatusText = () => {
        switch (connectionStatus) {
            case 'connected': return 'CONNECTED';
            case 'connecting': return 'CONNECTING...';
            case 'error': return 'CONNECTION ERROR';
            default: return 'DISCONNECTED';
        }
    };

    const tabs = [
        { id: 'instruments', label: 'INSTRUMENTS', icon: '📊' },
        { id: 'hvac', label: 'HVAC', icon: '🌡️' },
        { id: 'audio', label: 'AUDIO', icon: '🔊' }
    ];

    return (
        <header className="header">
            <div className="header-content">
                <div className="tabbed-menu">
                    <div className="tab-header">
                        {tabs.map(tab => (
                            <button
                                key={tab.id}
                                className={`tab-button ${activeTab === tab.id ? 'active' : ''}`}
                                onClick={() => onTabChange(tab.id)}
                            >
                                <div className="tab-icon">{tab.icon}</div>
                                <div className="tab-label">{tab.label}</div>
                            </button>
                    ))}
                </div>
            </div>

                <div className="brightness-control">
                    <div className="brightness-slider-container">
                        <input
                            type="range"
                            min="0"
                            max="100"
                            value={brightness}
                            onChange={(e) => onBrightnessChange(parseInt(e.target.value))}
                            className="brightness-slider"
                        />
                        <div className="brightness-label">BRIGHTNESS</div>
                    </div>
                </div>
                
                <div className="connection-status">
                    <div 
                        className="status-indicator"
                        style={{ backgroundColor: getStatusColor() }}
                    ></div>
                    <span className="status-text">{getStatusText()}</span>
                    {connectionStatus === 'disconnected' && (
                        <button 
                            className="connect-button"
                            onClick={onConnect}
                            disabled={isLoading}
                        >
                            {isLoading ? 'CONNECTING...' : 'CONNECT'}
                        </button>
                    )}
                </div>
            </div>
        </header>
    );
}

// Export for use in other files
window.Header = Header;
