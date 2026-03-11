import { useEffect } from 'react'
import { useStore } from './store'
import { Layers, RefreshCw, ChevronUp, ChevronDown, Trash2, Plus, X, Download } from 'lucide-react'

function App() {
  const { layers, addLayer, fetchLayers, isLoading, updateLayerOpacity, updateZIndex, removeLayer, renderPreview, previewUrl, setPreviewUrl } = useStore()

  // Fetch layers from Go backend on mount
  useEffect(() => {
    fetchLayers()
  }, [])

  return (
    <div className="h-screen bg-zinc-950 text-zinc-100 font-sans flex overflow-hidden">
      {/* Sidebar */}
      <aside className="w-80 border-r border-zinc-800 bg-zinc-900/50 flex flex-col">
        {/* Replace the header in App.tsx */}
        <header className="p-6 border-b border-zinc-800 flex justify-between items-center">
          <h1 className="text-xs font-bold uppercase tracking-tighter italic text-zinc-400">
            Image-Git Studio
          </h1>
          <div className="flex gap-1">
            <label className="p-2 hover:bg-zinc-800 rounded-md cursor-pointer transition-colors text-zinc-400 hover:text-white">
              <Plus size={14} />
              <input
                type="file"
                className="hidden"
                accept="image/png"
                onChange={(e) => e.target.files?.[0] && addLayer(e.target.files[0])}
              />
            </label>
            <button
              onClick={fetchLayers}
              className="p-2 hover:bg-zinc-800 rounded-md transition-colors text-zinc-400 hover:text-white"
            >
              <RefreshCw size={14} className={isLoading ? "animate-spin" : ""} />
            </button>
          </div>
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
              className="p-3 rounded-lg bg-zinc-800/40 border border-zinc-800 hover:border-zinc-700 transition-all space-y-4 group"
            >
              {/* Header Row: Info + Z-Controls */}
              <div className="flex justify-between items-center gap-3">
                <div className="flex-1 overflow-hidden">
                  <div className="flex items-center gap-2 mb-0.5">
                    <span className="text-sm font-medium truncate text-zinc-200">
                      {layer.name}
                    </span>
                    {/* Issue #15: Trash icon appears on card hover */}
                    <button
                      onClick={() => {
                        if (confirm(`Delete layer "${layer.name}"?`)) removeLayer(layer.name)
                      }}
                      className="opacity-0 group-hover:opacity-100 p-1 text-zinc-500 hover:text-red-400 transition-all"
                    >
                      <Trash2 size={12} />
                    </button>
                  </div>
                  <p className="text-[10px] font-mono text-zinc-600 truncate uppercase">
                    {layer.hash.substring(0, 12)}...
                  </p>
                </div>

                {/* Issue #14: Z-Index Stack */}
                <div className="flex flex-col items-center bg-zinc-950 rounded border border-zinc-800 shrink-0">
                  <button
                    onClick={() => updateZIndex(layer.name, 1)}
                    className="p-1 hover:text-white hover:bg-zinc-800 rounded-t transition-colors"
                  >
                    <ChevronUp size={12} />
                  </button>
                  <span className="text-[10px] font-mono font-bold border-y border-zinc-800 px-2 bg-zinc-900 leading-tight">
                    {layer.z_index}
                  </span>
                  <button
                    onClick={() => updateZIndex(layer.name, -1)}
                    className="p-1 hover:text-white hover:bg-zinc-800 rounded-b transition-colors"
                  >
                    <ChevronDown size={12} />
                  </button>
                </div>
              </div>

              {/* Issue #12: Opacity Slider */}
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


          <div className="p-4 border-t border-zinc-800 bg-zinc-900/80">
            <button
              onClick={renderPreview}
              disabled={isLoading}
              className="w-full py-2.5 bg-zinc-100 hover:bg-white text-zinc-950 rounded-md font-bold text-xs uppercase tracking-widest transition-all disabled:opacity-50 flex items-center justify-center gap-2"
            >
              {isLoading ? <RefreshCw size={14} className="animate-spin" /> : "Render Composite"}
            </button>
          </div>


        </div>
      </aside>

      {/* Main Canvas Area - Issue #13 */}
      <main className="flex-1 bg-[radial-gradient(#27272a_1px,transparent_1px)] [background-size:24px_24px] flex items-center justify-center relative overflow-hidden">
        <div className="relative w-[800px] h-[600px] bg-zinc-900 shadow-2xl rounded-lg border border-zinc-800 overflow-hidden">
          {/* Stack the layers based on their Z-Index */}
          {[...layers]
            .sort((a, b) => a.z_index - b.z_index)
            .map((layer) => (
              <img
                key={layer.hash}
                src={`http://localhost:3000/api/v1/objects/${layer.hash}.png`}
                alt={layer.name}
                className="absolute inset-0 w-full h-full object-contain pointer-events-none transition-opacity duration-200"
                style={{ opacity: layer.opacity }}
              />
            ))}
        </div>
      </main>
      {/* Issue #18: Render Preview Modal */}
      {previewUrl && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/80 backdrop-blur-sm p-10">
          <div className="relative bg-zinc-900 border border-zinc-800 rounded-xl shadow-2xl max-w-5xl w-full overflow-hidden flex flex-col">
            <div className="p-4 border-b border-zinc-800 flex justify-between items-center bg-zinc-900/50">
              <h3 className="text-sm font-bold uppercase tracking-widest text-zinc-400">Final Export Preview</h3>
              <button
                onClick={() => setPreviewUrl(null)}
                className="p-2 hover:bg-zinc-800 rounded-full text-zinc-500 hover:text-white transition-colors"
              >
                <X size={20} />
              </button>
            </div>

            <div className="flex-1 bg-[radial-gradient(#27272a_1px,transparent_1px)] [background-size:20px_20px] p-8 flex items-center justify-center overflow-auto">
              <img src={previewUrl} alt="Render Preview" className="max-w-full max-h-full shadow-2xl rounded border border-zinc-700" />
            </div>

            <div className="p-4 bg-zinc-950/50 border-t border-zinc-800 flex justify-end gap-3">
              <button
                onClick={() => setPreviewUrl(null)}
                className="px-6 py-2 text-xs font-bold uppercase tracking-widest text-zinc-500 hover:text-zinc-300 transition-colors"
              >
                Close
              </button>
              <a
                href={previewUrl}
                download="export.png"
                className="px-6 py-2 bg-zinc-100 hover:bg-white text-zinc-950 rounded font-bold text-xs uppercase tracking-widest flex items-center gap-2 transition-all"
              >
                <Download size={14} /> Download PNG
              </a>
            </div>
          </div>
        </div>
      )}
    </div>

  )
}

export default App