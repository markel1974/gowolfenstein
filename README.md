# gowolfenstein 🐺

A rich, retro-style Wolfenstein 3D clone and raycasting engine written entirely in [Go](https://golang.org/). 
Originally a personal test project from over 12 years ago, it brings back the nostalgic feel of classic 90s first-person shooters.

**The most important aspect of this project is that it does not rely on any external game engine.** The engine is completely custom-built from scratch, with its only real dependency being raw OpenGL (`go-gl`).

![Screenshot](https://github.com/markel1974/gowolfenstein/blob/main/bloob/screenshot1.png?raw=true)

## 🎮 Features

*   **Zero Game Engine Dependencies:** No Unity, no Godot, no external frameworks. The only real dependency is raw OpenGL.
*   **Custom Raycasting Engine:** Built from scratch in Go for that authentic 2.5D feel.
*   **Physics Engine:** Includes a custom physics system handling collision detection and dynamic movement.
*   **OpenGL Rendering:** Hardware-accelerated graphics using `go-gl` and GLFW.
*   **Retro Aesthetics:** Classic pixel-art graphics, textures, and lighting effects.
*   **Interactive Environments:** Functional minimap, doors, weapons, and dynamic entities.
*   **Embedded Graphics Library:** Includes `pixels`, a custom 2D wrapper handling sprites, batches, and windowing over OpenGL.

## 🕹️ Controls

| Key / Input | Action |
| :--- | :--- |
| `W`, `A`, `S`, `D` / `Arrows` | Move and turn |
| `Mouse` | Look around (Camera rotation) |
| `Space` | Fire |
| `Left Click` | Fire (Weapon 2) |
| `J`, `K`, `L` | Switch Weapons (1, 2, 3) |
| `Shift` | Run (Double speed) |
| `P` | Kick |
| `Tab` | Open doors / Interact |
| `1`, `2` | Toggle lighting effects |
| `M` / `N` | Toggle mini-map / Disable mouse |
| `Esc` | Quit the game |

## 🚀 Getting Started

### Prerequisites

You need Go (1.16+) installed on your system. Because the project uses `go-gl` for OpenGL bindings, make sure you have a working C compiler (e.g., GCC or Clang) and the necessary OpenGL development headers for your operating system.

### Building from Source

Clone the repository and build the game from the `game` directory:

```bash
git clone https://github.com/markel1974/gowolfenstein.git
cd gowolfenstein/game
go build -o gowolfenstein .
```

*(Note: There is also a pre-built binary and resources folder located in the `build/` directory.)*

### Running the Game

Run the compiled executable. You can use optional command-line flags to adjust the resolution, scale, and display mode:

```bash
./gowolfenstein [flags]
```

**Available Flags:**
*   `-f`: Enable fullscreen mode (default: false)
*   `-w`: Screen width (default: 495)
*   `-h`: Screen height (default: 270)
*   `-s`: Screen scale multiplier (default: 2.0)

Example (Running in a larger window):
```bash
./gowolfenstein -w 800 -h 600 -s 1.5
```

## 🛠️ Architecture Overview

*   `game/`: Contains the core game loop and engine logic. This includes the raycasting mechanics, physics, world state, minimap, entities, and weapons.
*   `pixels/`: A tailored 2D graphics framework wrapper around OpenGL, dealing with shaders, matrices, sprites, and window creation.
*   `build/`: Build artifacts, binaries, and game resources.
*   `bloob/`: Screenshots and media assets.

## 📜 License

This project is open-source. See the [LICENSE](LICENSE) file for more details.
