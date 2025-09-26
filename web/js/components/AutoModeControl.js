const { useState, useEffect, useCallback } = React;

// Auto Mode Control Component
function AutoModeControl({ autoMode, onToggle, disabled }) {
    return (
        <div className="control-panel auto-control">
            <h1 className="control-title">AUTO MODE</h1>
            <button
                className={`auto-toggle ${autoMode ? 'active' : ''}`}
                onClick={onToggle}
                disabled={disabled}
            >
                <div className="auto-icon">{autoMode ? '🔄' : '⏸️'}</div>
                <div className="auto-label">{autoMode ? 'ON' : 'OFF'}</div>
            </button>
        </div>
    );
}

// Export for use in other files
window.AutoModeControl = AutoModeControl;
