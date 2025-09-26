const { useState, useEffect, useCallback } = React;

// Fan Speed Control Component
function FanSpeedControl({ fanSpeed, onFanSpeedChange, disabled }) {
    const fanSpeeds = [
        { value: 0, label: 'OFF', icon: '⭕' },
        { value: 1, label: '1', icon: '💨' },
        { value: 2, label: '2', icon: '💨💨' },
        { value: 3, label: '3', icon: '💨💨💨' },
        { value: 4, label: '4', icon: '💨💨💨💨' },
        { value: 5, label: '5', icon: '💨💨💨💨💨' },
        { value: -1, label: 'AUTO', icon: '🔄' }
    ];

    return (
        <div className="control-panel fan-control">
            <h1 className="control-title">FAN SPEED</h1>
            <div className="fan-speed-grid">
                {fanSpeeds.map(speed => (
                    <button
                        key={speed.value}
                        className={`fan-speed-button ${fanSpeed === speed.value ? 'active' : ''}`}
                        onClick={() => onFanSpeedChange(speed.value)}
                        disabled={disabled}
                    >
                        <div className="fan-icon">{speed.icon}</div>
                        <div className="fan-label">{speed.label}</div>
                    </button>
                ))}
            </div>
        </div>
    );
}

// Export for use in other files
window.FanSpeedControl = FanSpeedControl;
