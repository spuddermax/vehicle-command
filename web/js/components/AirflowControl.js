const { useState, useEffect, useCallback } = React;

// Airflow Control Component
function AirflowControl({ airflowPattern, onAirflowChange, disabled }) {
    const airflowPatterns = [
        { value: 'face', label: 'FACE', icon: '👤' },
        { value: 'feet', label: 'FEET', icon: '🦶' },
        { value: 'defrost', label: 'DEFROST', icon: '❄️' },
        { value: 'face-feet', label: 'FACE+FEET', icon: '👤🦶' },
        { value: 'auto', label: 'AUTO', icon: '🔄' }
    ];

    return (
        <div className="control-panel airflow-control">
            <h1 className="control-title">AIRFLOW</h1>
            <div className="airflow-grid">
                {airflowPatterns.map(pattern => (
                    <button
                        key={pattern.value}
                        className={`airflow-button ${airflowPattern === pattern.value ? 'active' : ''}`}
                        onClick={() => onAirflowChange(pattern.value)}
                        disabled={disabled}
                    >
                        <div className="airflow-icon">{pattern.icon}</div>
                        <div className="airflow-label">{pattern.label}</div>
                    </button>
                ))}
            </div>
        </div>
    );
}

// Export for use in other files
window.AirflowControl = AirflowControl;
