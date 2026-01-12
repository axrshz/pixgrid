import { useState } from 'react';

interface PreviewPanelProps {
  original: string | null;
  converted: string | null;
  isConverting: boolean;
  onDownload: () => void;
  outputSize?: { width: number; height: number };
}

type ViewMode = 'side-by-side' | 'original' | 'result';

export function PreviewPanel({
  original,
  converted,
  isConverting,
  onDownload,
  outputSize,
}: PreviewPanelProps) {
  const [viewMode, setViewMode] = useState<ViewMode>('side-by-side');

  if (!original) {
    return (
      <div className="flex items-center justify-center h-96 bg-obsidian rounded-2xl border border-ash">
        <div className="text-center">
          <div className="w-16 h-16 mx-auto mb-4 rounded-xl bg-slate flex items-center justify-center">
            <svg
              className="w-8 h-8 text-smoke"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"
              />
            </svg>
          </div>
          <p className="text-smoke">Upload an image to see the preview</p>
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-4">
      {/* View mode toggle */}
      <div className="flex items-center justify-between">
        <div className="flex bg-graphite rounded-xl p-1">
          {(['side-by-side', 'original', 'result'] as ViewMode[]).map((mode) => (
            <button
              key={mode}
              onClick={() => setViewMode(mode)}
              className={`
                px-4 py-2 rounded-lg text-sm font-medium transition-all
                ${
                  viewMode === mode
                    ? 'bg-slate text-pearl'
                    : 'text-smoke hover:text-cloud'
                }
              `}
            >
              {mode === 'side-by-side'
                ? 'Compare'
                : mode === 'original'
                ? 'Original'
                : 'Result'}
            </button>
          ))}
        </div>

        <button
          onClick={onDownload}
          disabled={!converted || isConverting}
          className={`
            flex items-center gap-2 px-5 py-2.5 rounded-xl font-medium transition-all
            ${
              converted && !isConverting
                ? 'bg-gradient-to-r from-pixel-cyan to-pixel-purple text-void hover:opacity-90 hover:scale-105'
                : 'bg-slate text-smoke cursor-not-allowed'
            }
          `}
        >
          <svg
            className="w-5 h-5"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4"
            />
          </svg>
          Download PNG
        </button>
      </div>

      {/* Preview area */}
      <div className="relative bg-obsidian rounded-2xl border border-ash overflow-hidden">
        {/* Checkerboard background pattern */}
        <div
          className="absolute inset-0 opacity-10"
          style={{
            backgroundImage: `
              linear-gradient(45deg, #333 25%, transparent 25%),
              linear-gradient(-45deg, #333 25%, transparent 25%),
              linear-gradient(45deg, transparent 75%, #333 75%),
              linear-gradient(-45deg, transparent 75%, #333 75%)
            `,
            backgroundSize: '20px 20px',
            backgroundPosition: '0 0, 0 10px, 10px -10px, -10px 0px',
          }}
        />

        <div
          className={`
            relative p-6 min-h-[400px] flex items-center justify-center
            ${viewMode === 'side-by-side' ? 'gap-6' : ''}
          `}
        >
          {/* Original image */}
          {(viewMode === 'side-by-side' || viewMode === 'original') && (
            <div className="flex-1 flex flex-col items-center">
              <div className="relative group">
                <img
                  src={original}
                  alt="Original"
                  className="max-h-[350px] max-w-full object-contain rounded-lg shadow-2xl"
                />
                <div className="absolute bottom-2 left-2 bg-void/80 backdrop-blur-sm px-2 py-1 rounded text-xs text-cloud">
                  Original
                </div>
              </div>
            </div>
          )}

          {/* Divider for side-by-side */}
          {viewMode === 'side-by-side' && (
            <div className="w-px h-64 bg-gradient-to-b from-transparent via-ash to-transparent" />
          )}

          {/* Converted image */}
          {(viewMode === 'side-by-side' || viewMode === 'result') && (
            <div className="flex-1 flex flex-col items-center">
              <div className="relative">
                {isConverting ? (
                  <div className="w-64 h-64 flex items-center justify-center">
                    <div className="relative">
                      <div className="w-16 h-16 border-4 border-pixel-cyan/30 border-t-pixel-cyan rounded-full animate-spin" />
                      <div className="absolute inset-0 flex items-center justify-center">
                        <div className="w-8 h-8 border-4 border-pixel-magenta/30 border-b-pixel-magenta rounded-full animate-spin animate-reverse" />
                      </div>
                    </div>
                  </div>
                ) : converted ? (
                  <>
                    <img
                      src={converted}
                      alt="Pixel Art"
                      className="max-h-[350px] max-w-full object-contain rounded-lg shadow-2xl"
                      style={{ imageRendering: 'pixelated' }}
                    />
                    <div className="absolute bottom-2 left-2 bg-void/80 backdrop-blur-sm px-2 py-1 rounded text-xs text-pixel-cyan">
                      Pixel Art
                      {outputSize && (
                        <span className="text-smoke ml-2">
                          {outputSize.width}×{outputSize.height}
                        </span>
                      )}
                    </div>
                  </>
                ) : (
                  <div className="w-64 h-64 flex items-center justify-center text-smoke">
                    Processing...
                  </div>
                )}
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
