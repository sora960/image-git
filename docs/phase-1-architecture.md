# Phase 1: Core CLI Engine Architecture

## 🏛️ System Design
The engine utilizes a **Content-Addressable Storage (CAS)** pattern. It decouples visual metadata from physical image data to ensure storage efficiency.

### Data Flow
1. **Ingestion:** `File -> SHA-256 Hash -> objects/[hash].png`
2. **Metadata:** `Manifest.json` tracks Layer properties (Z-Index, Opacity, Frames).
3. **Rendering:** `Sort(Layers) -> Filter(CurrentFrame) -> DrawMask(Canvas)`



## 🧬 The Layer Primitive
Every layer in the system is defined by the following mathematical properties:
$$Layer = \{Name, Hash, Opacity, StartFrame, EndFrame, ZIndex\}$$

## 🧪 Verified Milestones
### The "Purple Square" Success
Confirmed the engine's ability to handle complex alpha-blending and depth:
- **Background:** Blue ($Z=0, \alpha=1.0$)
- **Foreground:** Red ($Z=10, \alpha=0.5$)
- **Result:** Correctly rendered Purple frame.

## 💻 CLI Command Reference
- `image-git --status`: View current manifest stack.
- `image-git --delete [name]`: Remove layer from manifest.
- `image-git --prune`: Clean orphaned object files.
- `image-git --composite --frame [n]`: Render specific frame.