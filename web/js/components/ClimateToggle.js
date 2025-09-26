const { useState, useEffect, useCallback } = React;

// Climate Toggle Component
function ClimateToggle({ isOn, onToggle, disabled }) {
    return (
        <div className="control-panel climate-toggle">
            <h1 className="control-title">CLIMATE</h1>
            <button
                className={`climate-button ${isOn ? 'active' : ''}`}
                onClick={onToggle}
                disabled={disabled}
            >
                <div className="climate-icon">{isOn ? '🌡️' : '⭕'}</div>
                <div className="climate-label">{isOn ? 'ON' : 'OFF'}</div>
            </button>
        </div>
    );
}

// Export for use in other files
window.ClimateToggle = ClimateToggle;
