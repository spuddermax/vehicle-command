const { useState, useEffect, useCallback } = React;

// Tabbed Menu Component
function TabbedMenu({ activeTab, onTabChange }) {
    const tabs = [
        { id: 'hvac', label: 'HVAC', icon: '🌡️' },
        { id: 'audio', label: 'AUDIO', icon: '🔊' }
    ];

    return (
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
    );
}

// Export for use in other files
window.TabbedMenu = TabbedMenu;
