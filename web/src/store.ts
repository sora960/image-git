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
    removeLayer: (name: string) => Promise<void>
    addLayer: (file: File) => Promise<void>

}

export const useStore = create<StudioState>((set, get) => ({
    layers: [],
    repoName: 'art-project',
    isLoading: false,


    // Inside useStore in store.ts
    addLayer: async (file: File) => {
        const { repoName, fetchLayers } = get()
        const formData = new FormData()
        formData.append('image', file)
        formData.append('name', file.name.split('.')[0])
        formData.append('z', '0')

        try {
            const res = await fetch(`http://localhost:3000/api/v1/repo/${repoName}/layers`, {
                method: 'POST',
                body: formData,
            })
            if (res.ok) await fetchLayers()
        } catch (err) {
            console.error("Upload failed", err)
        }
    },

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
    },

    // Add inside useStore:
    removeLayer: async (name) => {
        const { repoName, layers } = get()

        // 1. Optimistic UI update
        const previousLayers = layers
        set({ layers: layers.filter(l => l.name !== name) })

        try {
            // 2. Sync to Go Backend (using your existing DELETE route)
            const res = await fetch(`http://localhost:3000/api/v1/repo/${repoName}/layers/${name}`, {
                method: 'DELETE',
            })

            if (!res.ok) throw new Error("Delete failed on server")
        } catch (err) {
            console.error("Delete sync failed", err)
            // Rollback if server fails
            set({ layers: previousLayers })
        }
    }



}))