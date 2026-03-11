# Image-Git (MVP) 🎨🚀

**Image-Git** is a specialized version control system designed for artists and animators. It treats image layers like lines of code, enabling asynchronous collaboration, branching, and non-destructive merging for visual assets.

## 🌟 The Vision
Current collaboration tools (Figma, Canva) are "synchronous"—if one person edits, everyone sees it immediately. **Image-Git** allows artists to:
- **Fork** a 10-second animation clip or a high-res illustration.
- **Branch** to work on specific tasks (e.g., Shading, Line-art, FX) without affecting the Master version.
- **Pull Request** their specific layers for review by an Art Director.
- **Merge** pixels seamlessly using content-addressable storage.

## 🌟 The Vision: Unified Visual Version Control
Image-Git is designed to handle both static and temporal visual assets:

1. **Image Collab (Layers):** Asynchronous branching for illustrations. 
   - *Example:* Artist A branches to try a "Cyberpunk" lighting layer while the Master branch remains "Fantasy" style.
   
2. **Video/Animation Collab (Sequences):** Branching per scene or frame range.
   - *Example:* Animator A takes Frames 1-120 (Scene 1) as a branch. They commit their "Keyframes." Animator B pulls those keyframes to start "In-betweening" on a sub-branch.

## 🏗️ Tech Stack (Optimized for 2026)
- **Backend:** Go (Golang) for high-performance binary processing and concurrency.
- **Frontend:** React + WebGPU (for high-speed layer compositing).
- **Storage:** Content-Addressable Storage (CAS) via SHA-256 hashing.
- **Environment:** Developed on CachyOS (Arch Linux) using the Unix philosophy.

## 📂 Project Structure
- `cmd/server`: The Go entry point for the backend API.
- `internal/gitlogic`: Core "Image-Git" logic (Hashing, Merging, Diffing).
- `web/ui`: The React-based studio dashboard.
- `data/repositories`: Local storage for hashed image objects and manifests.

## 🛠️ Status: Conceptual MVP
Current Milestone: Implementing **Layer-as-a-Blob** storage using Go.
