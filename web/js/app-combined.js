const { useState, useEffect, useCallback } = React;

// Global state for card expansion
let expandedCardId = null;
const cardExpansionListeners = new Set();

function setExpandedCard(cardId) {
    expandedCardId = cardId;
    cardExpansionListeners.forEach(listener => listener(cardId));
}

function getExpandedCard() {
    return expandedCardId;
}

function addExpansionListener(listener) {
    cardExpansionListeners.add(listener);
    return () => cardExpansionListeners.delete(listener);
}

// Standardized Card Component
function Card({ 
    title, 
    children, 
    borderColor = 'var(--primary-blue)', 
    backgroundColor = 'var(--bg-panel)',
    className = '',
    ...props 
}) {
    const [isExpanded, setIsExpanded] = useState(false);
    const cardId = React.useMemo(() => Math.random().toString(36), []);

    const handleTitleClick = () => {
        const newExpanded = !isExpanded;
        setIsExpanded(newExpanded);
        
        if (newExpanded) {
            setExpandedCard(cardId);
        } else {
            setExpandedCard(null);
        }
    };

    // Listen for other card expansions
    React.useEffect(() => {
        const removeListener = addExpansionListener((expandedId) => {
            if (expandedId !== cardId) {
                setIsExpanded(false);
            }
        });
        return removeListener;
    }, [cardId]);

    const cardStyle = {
        borderColor: borderColor,
        backgroundColor: backgroundColor
    };

    return (
        <div 
            className={`card ${className} ${isExpanded ? 'expanded' : ''}`}
            style={cardStyle}
            {...props}
        >
            {title && (
                <h2 className="card-title" onClick={handleTitleClick}>
                    <span className="card-title-text">{title}</span>
                    <span className="card-expand-icon">
                        {isExpanded ? '🔍➖' : '🔍➕'}
                    </span>
                </h2>
            )}
            <div className="card-content">
                {children}
            </div>
        </div>
    );
}

// Speedometer Component
function Speedometer({ speed, maxSpeed = 120 }) {
    const percentage = (speed / maxSpeed) * 100;
    const rotation = (percentage / 100) * 180 - 90; // -90 to 90 degrees
    
    return (
        <Card title="SPEED" className="speedometer">
            <div className="speedometer-container">
                <div className="speedometer-gauge">
                    <div className="speedometer-needle" style={{ transform: `rotate(${rotation}deg)` }}></div>
                    <div className="speedometer-center"></div>
                </div>
                <div className="speed-display">
                    <span className="speed-value">{speed}</span>
                    <span className="speed-unit">MPH</span>
                </div>
            </div>
        </Card>
    );
}

// Battery Gauge Component
function BatteryGauge({ batteryLevel, range = 250 }) {
    const percentage = Math.max(0, Math.min(100, batteryLevel));
    const batteryColor = percentage > 50 ? '#00ff88' : percentage > 20 ? '#ffaa00' : '#ff4444';
    
    return (
        <Card title="BATTERY" className="battery-gauge">
            <div className="battery-container">
                <div className="battery-level">
                    <div 
                        className="battery-fill" 
                        style={{ 
                            width: `${percentage}%`,
                            backgroundColor: batteryColor
                        }}
                    ></div>
                </div>
                <div className="battery-info">
                    <span className="battery-percentage">{Math.round(percentage)}%</span>
                    <span className="battery-range">{range} mi</span>
                </div>
            </div>
        </Card>
    );
}

// Power Meter Component
function PowerMeter({ power, maxPower = 100 }) {
    const percentage = Math.abs(power / maxPower) * 100;
    const isRegen = power < 0;
    
    return (
        <Card title="POWER" className="power-meter">
            <div className="power-container">
                <div className="power-bar">
                    <div 
                        className={`power-fill ${isRegen ? 'regen' : 'accel'}`}
                        style={{ 
                            width: `${percentage}%`,
                            right: isRegen ? '0' : 'auto',
                            left: isRegen ? 'auto' : '0'
                        }}
                    ></div>
                </div>
                <div className="power-info">
                    <span className="power-value">{power > 0 ? '+' : ''}{power}%</span>
                    <span className="power-label">{isRegen ? 'REGEN' : 'ACCEL'}</span>
                </div>
            </div>
        </Card>
    );
}

// Temperature Gauge Component
function TemperatureGauge({ temperature, label = "MOTOR" }) {
    const tempColor = temperature > 80 ? '#ff4444' : temperature > 60 ? '#ffaa00' : '#00ff88';
    
    return (
        <Card title={label} className="temperature-gauge">
            <div className="temp-container">
                <div className="temp-display">
                    <span className="temp-value" style={{ color: tempColor }}>{temperature}°</span>
                    <span className="temp-unit">F</span>
                </div>
                <div className="temp-bar">
                    <div 
                        className="temp-fill" 
                        style={{ 
                            width: `${Math.min(100, (temperature / 100) * 100)}%`,
                            backgroundColor: tempColor
                        }}
                    ></div>
                </div>
            </div>
        </Card>
    );
}

// Grid with Overlay Component
function GridWithOverlay({ children }) {
    const [hasExpandedCard, setHasExpandedCard] = useState(false);

    React.useEffect(() => {
        const removeListener = addExpansionListener((expandedId) => {
            setHasExpandedCard(expandedId !== null);
        });
        return removeListener;
    }, []);

    return (
        <div className="controls-grid">
            <div className={`grid-overlay ${hasExpandedCard ? 'active' : ''}`}></div>
            {children}
        </div>
    );
}

// Main Dashboard Component
function Dashboard({ dashboardState }) {
    return (
        <GridWithOverlay>
            <Speedometer speed={dashboardState.speed} />
            <BatteryGauge batteryLevel={dashboardState.batteryLevel} range={dashboardState.range} />
            <PowerMeter power={dashboardState.power} />
            <TemperatureGauge temperature={dashboardState.motorTemp} label="MOTOR" />
            <TemperatureGauge temperature={dashboardState.batteryTemp} label="BATTERY" />
        </GridWithOverlay>
    );
}

// Tabbed Menu Component
function TabbedMenu({ activeTab, onTabChange }) {
    const tabs = [
        { id: 'instruments', label: 'INSTRUMENTS', icon: '📊' },
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

// Audio Control Component
function AudioControl({ volume, onVolumeChange, disabled }) {
    const handleVolumeChange = (delta) => {
        const newVolume = Math.max(0, Math.min(100, volume + delta));
        onVolumeChange(newVolume);
    };

    return (
        <Card title="VOLUME" className="audio-control">
            <div className="volume-display">
                <div className="volume-value">{volume}%</div>
                <div className="volume-buttons">
                    <button 
                        className="volume-button volume-down"
                        onClick={() => handleVolumeChange(-5)}
                        disabled={disabled}
                    >
                        −
                    </button>
                    <button 
                        className="volume-button volume-up"
                        onClick={() => handleVolumeChange(5)}
                        disabled={disabled}
                    >
                        +
                    </button>
                </div>
            </div>
        </Card>
    );
}

// Header Component
function Header({ connectionStatus, onConnect, isLoading, activeTab, onTabChange }) {
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
            ? Math.max(60, Math.min(85, driverTemp + delta)) // 60°F to 85°F range
            : Math.max(60, Math.min(85, passengerTemp + delta)); // 60°F to 85°F range
        
        if (side === 'driver') {
            onTemperatureChange(newTemp, passengerTemp);
        } else {
            onTemperatureChange(driverTemp, newTemp);
        }
    };

    return (
        <Card title="TEMPERATURE" className="temperature-control">
            <div className="temperature-display">
                <div className="temp-side">
                    <div className="temp-label">DRIVER</div>
                    <div className="temp-value">{Math.round(driverTemp)}°F</div>
                    <div className="temp-slider-container">
                        <input
                            type="range"
                            min="60"
                            max="85"
                            value={driverTemp}
                            onChange={(e) => handleTempChange('driver', parseFloat(e.target.value) - driverTemp)}
                            disabled={disabled}
                            className="temp-slider"
                        />
                        <div className="temp-labels">
                            <span onClick={() => handleTempChange('driver', 60 - driverTemp)}>60</span>
                            <span onClick={() => handleTempChange('driver', 65 - driverTemp)}>65</span>
                            <span onClick={() => handleTempChange('driver', 70 - driverTemp)}>70</span>
                            <span onClick={() => handleTempChange('driver', 75 - driverTemp)}>75</span>
                            <span onClick={() => handleTempChange('driver', 80 - driverTemp)}>80</span>
                            <span onClick={() => handleTempChange('driver', 85 - driverTemp)}>85</span>
                        </div>
                    </div>
                </div>
                
                <div className="temp-divider"></div>
                
                <div className="temp-side">
                    <div className="temp-label">PASSENGER</div>
                    <div className="temp-value">{Math.round(passengerTemp)}°F</div>
                    <div className="temp-slider-container">
                        <input
                            type="range"
                            min="60"
                            max="85"
                            value={passengerTemp}
                            onChange={(e) => handleTempChange('passenger', parseFloat(e.target.value) - passengerTemp)}
                            disabled={disabled}
                            className="temp-slider"
                        />
                        <div className="temp-labels">
                            <span onClick={() => handleTempChange('passenger', 60 - passengerTemp)}>60</span>
                            <span onClick={() => handleTempChange('passenger', 65 - passengerTemp)}>65</span>
                            <span onClick={() => handleTempChange('passenger', 70 - passengerTemp)}>70</span>
                            <span onClick={() => handleTempChange('passenger', 75 - passengerTemp)}>75</span>
                            <span onClick={() => handleTempChange('passenger', 80 - passengerTemp)}>80</span>
                            <span onClick={() => handleTempChange('passenger', 85 - passengerTemp)}>85</span>
                        </div>
                    </div>
                </div>
            </div>
        </Card>
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
        <Card title="FAN SPEED" className="fan-control">
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
        </Card>
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
        <Card title="AIRFLOW" className="airflow-control">
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
        </Card>
    );
}

// Auto Mode Control Component
function AutoModeControl({ autoMode, onToggle, disabled }) {
    return (
        <Card title="AUTO MODE" className="auto-control">
            <div className="auto-buttons">
                <button
                    className={`auto-button auto-on ${autoMode ? 'active' : ''}`}
                    onClick={() => onToggle(true)}
                    disabled={disabled}
                >
                    <div className="auto-icon">🔄</div>
                    <div className="auto-label">ON</div>
                </button>
                <button
                    className={`auto-button auto-off ${!autoMode ? 'active' : ''}`}
                    onClick={() => onToggle(false)}
                    disabled={disabled}
                >
                    <div className="auto-icon">⏸️</div>
                    <div className="auto-label">OFF</div>
                </button>
            </div>
        </Card>
    );
}

// Climate Toggle Component
function ClimateToggle({ isOn, onToggle, disabled }) {
    return (
        <Card title="CLIMATE SYSTEM" className="climate-toggle">
            <div className="climate-buttons">
                <button
                    className={`climate-button climate-on ${isOn ? 'active' : ''}`}
                    onClick={() => onToggle(true)}
                    disabled={disabled}
                >
                    <div className="climate-icon">🌡️</div>
                    <div className="climate-label">ON</div>
                </button>
                <button
                    className={`climate-button climate-off ${!isOn ? 'active' : ''}`}
                    onClick={() => onToggle(false)}
                    disabled={disabled}
                >
                    <div className="climate-icon">⭕</div>
                    <div className="climate-label">OFF</div>
                </button>
            </div>
        </Card>
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

// Main HVAC Control Application
function TeslaHVACApp() {
    const [connectionStatus, setConnectionStatus] = useState('disconnected');
    const [brightness, setBrightness] = useState(100);
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
    const [activeTab, setActiveTab] = useState('instruments');
    const [audioState, setAudioState] = useState({
        volume: 50
    });
    const [dashboardState, setDashboardState] = useState({
        speed: 45,
        batteryLevel: 78,
        range: 245,
        power: 25,
        motorTemp: 72,
        batteryTemp: 68
    });

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

    const toggleAutoMode = useCallback(async (autoMode) => {
        if (connectionStatus !== 'connected') return;
        
        // Update UI immediately for instant feedback
        setHvacState(prev => ({
            ...prev,
            autoMode: autoMode
        }));
        
        // Send to backend in background (no loading state)
        try {
            // Simulate API call without blocking UI
            await new Promise(resolve => setTimeout(resolve, 50));
        } catch (err) {
            setError('Failed to set auto mode');
        }
    }, [connectionStatus]);

    const toggleClimate = useCallback(async (isOn) => {
        if (connectionStatus !== 'connected') return;
        
        // Update UI immediately for instant feedback
        setHvacState(prev => ({
            ...prev,
            isOn: isOn
        }));
        
        // Send to backend in background (no loading state)
        try {
            // Simulate API call without blocking UI
            await new Promise(resolve => setTimeout(resolve, 100));
        } catch (err) {
            setError('Failed to set climate system');
        }
    }, [connectionStatus]);

    const handleBrightnessChange = useCallback((newBrightness) => {
        setBrightness(newBrightness);
    }, []);

    // Audio Control Functions
    const setVolume = useCallback(async (volume) => {
        // Update UI immediately for instant feedback
        setAudioState(prev => ({
            ...prev,
            volume
        }));
        
        // Send to backend in background (no loading state)
        try {
            // Simulate API call without blocking UI
            await new Promise(resolve => setTimeout(resolve, 50));
        } catch (err) {
            setError('Failed to set volume');
        }
    }, []);

    // Auto-connect on app start
    useEffect(() => {
        connectToVehicle();
    }, [connectToVehicle]);

    const renderTabContent = () => {
        if (activeTab === 'instruments') {
            return (
                <Dashboard dashboardState={dashboardState} />
            );
        } else if (activeTab === 'hvac') {
            return (
                <>
                    <GridWithOverlay>
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
                    </GridWithOverlay>
                    
                    <StatusDisplay
                        hvacState={hvacState}
                        connectionStatus={connectionStatus}
                    />
                </>
            );
        } else if (activeTab === 'audio') {
            return (
                <GridWithOverlay>
                    <AudioControl
                        volume={audioState.volume}
                        onVolumeChange={setVolume}
                        disabled={false}
                    />
                </GridWithOverlay>
            );
        }
        return null;
    };

    return (
        <div className={`app ${orientation}`}>
            <Header 
                connectionStatus={connectionStatus}
                onConnect={connectToVehicle}
                isLoading={isLoading}
                activeTab={activeTab}
                onTabChange={setActiveTab}
            />
            
            <main className="main-content">
                {error && (
                    <ErrorMessage 
                        message={error} 
                        onDismiss={() => setError(null)} 
                    />
                )}
                
                {renderTabContent()}
            </main>
            
            {/* Brightness Overlay - covers entire page */}
            <div 
                className={`brightness-overlay ${brightness < 100 ? 'active' : ''}`}
                style={{
                    background: `rgba(0, 0, 0, ${(100 - brightness) / 100})`
                }}
            />
            
            {/* Brightness Control - positioned absolutely above overlay */}
            <div 
                className="brightness-control-absolute"
                style={{
                    opacity: Math.max(0.2, brightness / 100)
                }}
            >
                <div className="brightness-slider-container">
                    <input
                        type="range"
                        min="0"
                        max="100"
                        value={brightness}
                        onChange={(e) => handleBrightnessChange(parseInt(e.target.value))}
                        className="brightness-slider"
                    />
                    <div className="brightness-label">BRIGHTNESS</div>
                </div>
            </div>
        </div>
    );
}

// Render the app
ReactDOM.render(<TeslaHVACApp />, document.getElementById('root'));
