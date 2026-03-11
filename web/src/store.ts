import { create } from 'zustand'

export interface Layer {
    name: string
    hash: string
    opacity: number
    z_index: number
}

interface StudioState {
    layers: Layer[]
    repoName: string
    isLoading: boolean
    fetchLayers: () => Promise<void>
    // NEW: Update function for Issue #12
    updateLayerOpacity: (name: string, opacity: number) => Promise<void>
}

export const useStore = create<StudioState>((set, get) => ({
    layers: [],
    repoName: 'art-project',
    isLoading: false,

    fetchLayers: async () => {
        set({ isLoading: true })
        try {
            const res = await fetch(`http://localhost:3000/api/v1/repo/${get().repoName}`)
            const data = await res.json()
            set({ layers: data.layers || [], isLoading: false })
        } catch (err) {
            set({ isLoading: false })
        }
    },

    // Logic for Issue #12
    updateLayerOpacity: async (name, opacity) => {
        // 1. Optimistic Update (Immediate UI response)
        const previousLayers = get().layers
        set({
            layers: previousLayers.map(l => l.name === name ? { ...l, opacity } : l)
        })

        try {
            // 2. Sync to Go Backend
            await fetch(`http://localhost:3000/api/v1/repo/${get().repoName}/layers/${name}`, {
                method: 'PATCH',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ opacity })
            })
        } catch (err) {
            console.error("Sync failed:", err)
            // Rollback if server fails
            set({ layers: previousLayers })
        }
    }
}))