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
    updateZIndex: (name: string, delta: number) => Promise<void>
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

    updateZIndex: async (name, delta) => {
        const { repoName, layers } = get()
        const updatedLayers = layers.map(l =>
            l.name === name ? { ...l, z_index: l.z_index + delta } : l
        )

        set({ layers: updatedLayers })

        try {
            // URL updated to generic PATCH route (removed /z)
            await fetch(`http://localhost:3000/api/v1/repo/${repoName}/layers/${name}`, {
                method: 'PATCH',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ z_index: updatedLayers.find(l => l.name === name)?.z_index })
            })
        } catch (err) {
            console.error("Z-Index sync failed", err)
        }
    },

    updateLayerOpacity: async (name, opacity) => {
        const previousLayers = get().layers
        set({
            layers: previousLayers.map(l => l.name === name ? { ...l, opacity } : l)
        })

        try {
            await fetch(`http://localhost:3000/api/v1/repo/${get().repoName}/layers/${name}`, {
                method: 'PATCH',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ opacity })
            })
        } catch (err) {
            console.error("Sync failed:", err)
            set({ layers: previousLayers })
        }
    }
}))