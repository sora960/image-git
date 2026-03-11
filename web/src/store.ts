import { create } from 'zustand'

// The TypeScript "Contract" - must match your Go 'Layer' struct
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
}

export const useStore = create<StudioState>((set, get) => ({
    layers: [],
    repoName: 'art-project', // Default repo from your git diff
    isLoading: false,

    fetchLayers: async () => {
        set({ isLoading: true })
        try {
            const res = await fetch(`http://localhost:3000/api/v1/repo/${get().repoName}`)
            if (!res.ok) throw new Error("Backend unreachable")
            const data = await res.json()

            // We ensure layers is always an array to prevent .map() crashes
            set({ layers: data.layers || [], isLoading: false })
        } catch (err) {
            console.error("API Error:", err)
            set({ isLoading: false })
        }
    },
}))