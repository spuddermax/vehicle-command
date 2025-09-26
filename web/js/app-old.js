const { useState, useEffect, useCallback } = React;

// Main HVAC Control Application
function TeslaHVACApp() {
    const [connectionStatus, setConnectionStatus] = useState('disconnected');
    const [hvacState, setHvacState] = useState({
        isOn: false,
        driverTemp: 72.0, // 72°F = 22°C
        passengerTemp: 72.0, // 72°F = 22°C
        fanSpeed: 0,
        airflowPattern: 'auto',
        autoMode: false,
        insideTemp: 68.0, // 68°F = 20°C
        outsideTemp: 59.0 // 59°F = 15°C
    });
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState(null);
    const [orientation, setOrientation] = useState('portrait');

    // Handle orientation changes
    useEffect(() => {
        const handleOrientationChange = () => {
            const isLandscape = window.innerWidth > window.innerHeight;
            setOrientation(isLandscape ? 'landscape' : 'portrait');
        };

        // Initial orientation check
        handleOrientationChange();

        // Listen for orientation changes
        window.addEventListener('resize', handleOrientationChange);
        window.addEventListener('orientationchange', handleOrientationChange);

        return () => {
            window.removeEventListener('resize', handleOrientationChange);
            window.removeEventListener('orientationchange', handleOrientationChange);
        };
    }, []);

    // Simulate connection to Tesla vehicle
    const connectToVehicle = useCallback(async () => {
        setIsLoading(true);
        setError(null);
        
        try {
            // Simulate connection delay (reduced for faster startup)
            await new Promise(resolve => setTimeout(resolve, 500));
            
            // Simulate successful connection
            setConnectionStatus('connected');
            setHvacState(prev => ({
                ...prev,
                isOn: true,
                driverTemp: 72.0, // 72°F
                passengerTemp: 72.0, // 72°F
                fanSpeed: 2,
                autoMode: true
            }));
        } catch (err) {
            setError('Failed to connect to vehicle');
            setConnectionStatus('error');
        } finally {
            setIsLoading(false);
        }
    }, []);

    // HVAC Control Functions
    const setTemperature = useCallback(async (driverTemp, passengerTemp) => {
        if (connectionStatus !== 'connected') return;
        
        // Update UI immediately for instant feedback
        setHvacState(prev => ({
            ...prev,
            driverTemp,
            passengerTemp
        }));
        
        // Send to backend in background (no loading state)
        try {
            // Simulate API call without blocking UI
            await new Promise(resolve => setTimeout(resolve, 100));
        } catch (err) {
            setError('Failed to set temperature');
        }
    }, [connectionStatus]);

    const setFanSpeed = useCallback(async (speed) => {
        if (connectionStatus !== 'connected') return;
        
        // Update UI immediately for instant feedback
        setHvacState(prev => ({
            ...prev,
            fanSpeed: speed
        }));
        
        // Send to backend in background (no loading state)
        try {
            // Simulate API call without blocking UI
            await new Promise(resolve => setTimeout(resolve, 50));
        } catch (err) {
            setError('Failed to set fan speed');
        }
    }, [connectionStatus]);

    const setAirflowPattern = useCallback(async (pattern) => {
        if (connectionStatus !== 'connected') return;
        
        // Update UI immediately for instant feedback
        setHvacState(prev => ({
            ...prev,
            airflowPattern: pattern
        }));
        
        // Send to backend in background (no loading state)
        try {
            // Simulate API call without blocking UI
            await new Promise(resolve => setTimeout(resolve, 50));
        } catch (err) {
            setError('Failed to set airflow pattern');
        }
    }, [connectionStatus]);

    const toggleAutoMode = useCallback(async () => {
        if (connectionStatus !== 'connected') return;
        
        // Update UI immediately for instant feedback
        setHvacState(prev => ({
            ...prev,
            autoMode: !prev.autoMode
        }));
        
        // Send to backend in background (no loading state)
        try {
            // Simulate API call without blocking UI
            await new Promise(resolve => setTimeout(resolve, 50));
        } catch (err) {
            setError('Failed to toggle auto mode');
        }
    }, [connectionStatus]);

    const toggleClimate = useCallback(async () => {
        if (connectionStatus !== 'connected') return;
        
        // Update UI immediately for instant feedback
        setHvacState(prev => ({
            ...prev,
            isOn: !prev.isOn
        }));
        
        // Send to backend in background (no loading state)
        try {
            // Simulate API call without blocking UI
            await new Promise(resolve => setTimeout(resolve, 100));
        } catch (err) {
            setError('Failed to toggle climate');
        }
    }, [connectionStatus]);

    // Auto-connect on app start
    useEffect(() => {
        connectToVehicle();
    }, [connectToVehicle]);

    return (
        <div className={`app ${orientation}`}>
            <Header 
                connectionStatus={connectionStatus}
                onConnect={connectToVehicle}
                isLoading={isLoading}
            />
            
            <main className="main-content">
                {error && (
                    <ErrorMessage 
                        message={error} 
                        onDismiss={() => setError(null)} 
                    />
                )}
                
                <div className="controls-grid">
                    <TemperatureControl
                        driverTemp={hvacState.driverTemp}
                        passengerTemp={hvacState.passengerTemp}
                        onTemperatureChange={setTemperature}
                        disabled={!hvacState.isOn}
                    />
                    
                    <FanSpeedControl
                        fanSpeed={hvacState.fanSpeed}
                        onFanSpeedChange={setFanSpeed}
                        disabled={!hvacState.isOn}
                    />
                    
                    <AirflowControl
                        airflowPattern={hvacState.airflowPattern}
                        onAirflowChange={setAirflowPattern}
                        disabled={!hvacState.isOn}
                    />
                    
                    <AutoModeControl
                        autoMode={hvacState.autoMode}
                        onToggle={toggleAutoMode}
                        disabled={!hvacState.isOn}
                    />
                    
                    <ClimateToggle
                        isOn={hvacState.isOn}
                        onToggle={toggleClimate}
                        disabled={false}
                    />
                </div>
                
                <StatusDisplay
                    hvacState={hvacState}
                    connectionStatus={connectionStatus}
                />
            </main>
        </div>
    );
}

// Header Component
function Header({ connectionStatus, onConnect, isLoading }) {
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

    return (
        <header className="header">
            <div className="header-content">
                <div className="logo">
                    <h1>TESLA HVAC</h1>
                    <div className="logo-subtitle">Climate Control</div>
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

// Error Message Component
function ErrorMessage({ message, onDismiss }) {
    return (
        <div className="error-message">
            <div className="error-content">
                <span className="error-icon">⚠</span>
                <span className="error-text">{message}</span>
                <button className="error-dismiss" onClick={onDismiss}>×</button>
            </div>
        </div>
    );
}

// Temperature Control Component
function TemperatureControl({ driverTemp, passengerTemp, onTemperatureChange, disabled }) {
    const handleTempChange = (side, delta) => {
        const newTemp = side === 'driver' 
            ? Math.max(59, Math.min(86, driverTemp + delta)) // 59°F to 86°F range
            : Math.max(59, Math.min(86, passengerTemp + delta)); // 59°F to 86°F range
        
        if (side === 'driver') {
            onTemperatureChange(newTemp, passengerTemp);
        } else {
            onTemperatureChange(driverTemp, newTemp);
        }
    };

    return (
        <div className="control-panel temperature-control">
            <h3 className="control-title">TEMPERATURE</h3>
            <div className="temperature-display">
                <div className="temp-side">
                    <div className="temp-label">DRIVER</div>
                    <div className="temp-value">{Math.round(driverTemp)}°F</div>
                    <div className="temp-buttons">
                        <button 
                            className="temp-button temp-down"
                            onClick={() => handleTempChange('driver', -1)}
                            disabled={disabled}
                        >
                            −
                        </button>
                        <button 
                            className="temp-button temp-up"
                            onClick={() => handleTempChange('driver', 1)}
                            disabled={disabled}
                        >
                            +
                        </button>
                    </div>
                </div>
                
                <div className="temp-divider"></div>
                
                <div className="temp-side">
                    <div className="temp-label">PASSENGER</div>
                    <div className="temp-value">{Math.round(passengerTemp)}°F</div>
                    <div className="temp-buttons">
                        <button 
                            className="temp-button temp-down"
                            onClick={() => handleTempChange('passenger', -1)}
                            disabled={disabled}
                        >
                            −
                        </button>
                        <button 
                            className="temp-button temp-up"
                            onClick={() => handleTempChange('passenger', 1)}
                            disabled={disabled}
                        >
                            +
                        </button>
                    </div>
                </div>
            </div>
        </div>
    );
}

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
            <h3 className="control-title">FAN SPEED</h3>
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
            <h3 className="control-title">AIRFLOW</h3>
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

// Auto Mode Control Component
function AutoModeControl({ autoMode, onToggle, disabled }) {
    return (
        <div className="control-panel auto-control">
            <h3 className="control-title">AUTO MODE</h3>
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

// Climate Toggle Component
function ClimateToggle({ isOn, onToggle, disabled }) {
    return (
        <div className="control-panel climate-toggle">
            <h3 className="control-title">CLIMATE</h3>
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

// Status Display Component
function StatusDisplay({ hvacState, connectionStatus }) {
    return (
        <div className="status-display">
            <div className="status-grid">
                <div className="status-item">
                    <div className="status-label">INSIDE</div>
                    <div className="status-value">{Math.round(hvacState.insideTemp)}°F</div>
                </div>
                <div className="status-item">
                    <div className="status-label">OUTSIDE</div>
                    <div className="status-value">{Math.round(hvacState.outsideTemp)}°F</div>
                </div>
                <div className="status-item">
                    <div className="status-label">FAN</div>
                    <div className="status-value">
                        {hvacState.fanSpeed === -1 ? 'AUTO' : hvacState.fanSpeed}
                    </div>
                </div>
                <div className="status-item">
                    <div className="status-label">MODE</div>
                    <div className="status-value">
                        {hvacState.autoMode ? 'AUTO' : 'MANUAL'}
                    </div>
                </div>
            </div>
        </div>
    );
}

// Render the app
ReactDOM.render(<TeslaHVACApp />, document.getElementById('root'));
