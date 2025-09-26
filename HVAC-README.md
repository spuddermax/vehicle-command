# Tesla Model Y HVAC Interface

A web-based Bluetooth interface for Tesla Model Y HVAC controls, inspired by 1990s automotive design.

## Project Structure

```
vehicle-command/
├── cmd/
│   └── tesla-hvac-server/     # Main HVAC server application
├── internal/
│   ├── hvac/                  # HVAC-specific functionality
│   └── web/                   # Web server functionality
├── web/
│   ├── css/                   # CSS styles
│   └── js/                    # JavaScript files
├── config/                    # Configuration files
└── pkg/                       # Tesla vehicle-command library packages
```

## Development Setup

1. **Fork and Clone**: This repository is a fork of Tesla's vehicle-command library
2. **Go Environment**: Requires Go 1.21+ (currently using 1.23.4)
3. **Dependencies**: All dependencies managed via go.mod

## Key Features

- **Web-based Interface**: HTML5/CSS3/JavaScript touchscreen interface
- **Tesla Integration**: Uses official Tesla vehicle-command library
- **HVAC Controls**: Temperature, fan speed, airflow patterns
- **Responsive Design**: Portrait and landscape orientations
- **Customizable Layout**: Drag-and-drop layout customization
- **1990s Design**: Chevy Suburban-inspired color scheme and typography

## Next Steps

1. Implement Tesla vehicle communication backend
2. Create web-based touchscreen interface
3. Add HVAC control functionality
4. Implement layout customization features
5. Testing and deployment

## Original Tesla Library

This project is built on top of the official Tesla vehicle-command library:
- Repository: https://github.com/teslamotors/vehicle-command
- License: Apache-2.0
- Documentation: See original README.md
