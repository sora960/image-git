import { useEffect } from 'react'
import { useStore } from './store'
import { Layers, RefreshCw } from 'lucide-react'

function App() {
  const { layers, fetchLayers, isLoading, updateLayerOpacity } = useStore()

  // Fetch layers from Go backend on mount
  useEffect(() => {
    fetchLayers()
  }, [])

  return (
    <div className="h-screen bg-zinc-950 text-zinc-100 font-sans flex overflow-hidden">
      {/* Sidebar */}
      <aside className="w-80 border-r border-zinc-800 bg-zinc-900/50 flex flex-col">
        <header className="p-6 border-b border-zinc-800 flex justify-between items-center">
          <h1 className="text-xs font-bold uppercase tracking-tighter italic text-zinc-400">
            Image-Git Studio
          </h1>
          <button
            onClick={fetchLayers}
            className="p-2 hover:bg-zinc-800 rounded-md transition-colors"
          >
            <RefreshCw size={14} className={isLoading ? "animate-spin" : ""} />
          </button>
        </header>

        <div className="flex-1 overflow-y-auto p-4 space-y-2">
          <div className="flex items-center gap-2 text-[10px] font-bold text-zinc-500 uppercase px-2 mb-4">
            <Layers size={12} /> Layers
          </div>

          {layers.length === 0 && !isLoading && (
            <p className="text-xs text-zinc-600 px-2 italic">No layers found in manifest.</p>
          )}

          {layers.map((layer) => (
            <div
              key={layer.hash}
              className="p-3 rounded-lg bg-zinc-800/40 border border-zinc-800 hover:border-zinc-700 transition-all space-y-3 group"
            >
              <div className="flex justify-between items-start">
                <div className="overflow-hidden">
                  <span className="text-sm font-medium truncate block">{layer.name}</span>
                  <p className="text-[10px] font-mono text-zinc-600 truncate">{layer.hash}</p>
                </div>
                <span className="text-[10px] font-mono text-zinc-500 bg-zinc-900 px-1.5 py-0.5 rounded">
                  Z: {layer.z_index}
                </span>
              </div>

              {/* Opacity Slider - Issue #12 */}
              <div className="space-y-1.5">
                <div className="flex justify-between text-[10px] uppercase tracking-wider text-zinc-500 font-bold">
                  <span>Opacity</span>
                  <span className="text-zinc-300">{Math.round(layer.opacity * 100)}%</span>
                </div>
                <input
                  type="range"
                  min="0"
                  max="1"
                  step="0.01"
                  value={layer.opacity}
                  onChange={(e) => updateLayerOpacity(layer.name, parseFloat(e.target.value))}
                  className="w-full h-1 bg-zinc-700 rounded-lg appearance-none cursor-pointer accent-zinc-100 hover:accent-white transition-all"
                />
              </div>
            </div>
          ))}
        </div>
      </aside>

      {/* Main Canvas Area */}
      <main className="flex-1 bg-[radial-gradient(#27272a_1px,transparent_1px)] [background-size:24px_24px] flex items-center justify-center">
        <div className="p-8 border-2 border-dashed border-zinc-800 rounded-3xl text-zinc-700 font-mono text-sm italic">
          Canvas Preview Pipeline Offline
        </div>
      </main>
    </div>
  )
}

export default App