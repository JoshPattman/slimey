# Slimey - Particle Simulation

A mesmerizing particle simulation written in Go using the Pixel game library. Watch as thousands of particles create beautiful patterns through pheromone trails.

![Particle Simulation](particle.png)

## Features

- Real-time particle simulation with configurable parameters
- Pheromone trail system with decay
- Smooth particle movement and rotation
- High-performance rendering using the Pixel game library
- Configurable simulation parameters via command-line flags

## Usage

### Download Pre-built Binary
The easiest way to get started is to download a pre-built binary from the [Releases page](https://github.com/JoshPattman/slimey/releases). Choose the appropriate version for your operating system and architecture.

### Run The Binary
`$ ./slimey-<os>-<arch> <flags>`

### Command-line Flags

- `-h`: See a help page
- `-screen-x`: Screen width (default: 800)
- `-screen-y`: Screen height (default: 800)
- `-num`: Number of particles (default: 20000)
- `-particle-speed`: Speed of particles (default: 100.0)
- `-particle-rotation`: Rotational speed of particles (default: 1.5)
- `-chunk-size`: Size of pheromone chunks in pixels (default: 4)
- `-pher-rate`: Rate of dropping pheromones (default: 1.0)
- `-pher-decay`: Decay rate of pheromones (default: 0.99)
- `-prof`: Enable debug statistics output (default: false)
