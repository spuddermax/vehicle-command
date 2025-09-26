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
                </>
            );
        } else if (activeTab === 'audio') {
            return (
                <div className="controls-grid">
                    <AudioControl
                        volume={audioState.volume}
                        onVolumeChange={setVolume}
                        disabled={false}
                    />
                </div>
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
        </div>
    );
}

// Export for use in other files
window.TeslaHVACApp = TeslaHVACApp;
