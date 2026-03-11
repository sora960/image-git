# Image-Git Studio (Frontend)

The web-based control center for the Image-Git engine. This studio allows for real-time layer management, opacity control, and repository manifest visualization.

## 🚀 Tech Stack

- **Framework:** React 19 (TypeScript + SWC)
- **Styling:** Tailwind CSS v4 (Engine-integrated)
- **State:** Zustand (Minimalist Store)
- **UI Architecture:** shadcn/ui (Manual Implementation)
- **Build Tool:** Vite

## 🛠️ Development

### Prerequisites
- Node.js (v20+)
- Image-Git Backend (Running on :3000)

### Getting Started
1. Install dependencies: `npm install`
2. Start dev server: `npm run dev`

## 🌉 API Bridge
This frontend communicates with the Go Fiber backend via REST. 
- **Endpoint:** `http://localhost:3000/api/v1`
- **Primary Contract:** `interface Layer { name, hash, opacity, z_index }`