import { useState, useEffect, useCallback } from 'react';
import { ImageUploader } from './components/ImageUploader';
import { ParameterSliders, type Parameters } from './components/ParameterSliders';
import { PreviewPanel } from './components/PreviewPanel';
import { uploadImage, convertImage, downloadImage } from './api';
import { useDebounce } from './hooks/useDebounce';

interface ImageState {
  sessionId: string;
  original: string;
  width: number;
  height: number;
}

function App() {
  const [imageState, setImageState] = useState<ImageState | null>(null);
  const [converted, setConverted] = useState<string | null>(null);
  const [outputSize, setOutputSize] = useState<{ width: number; height: number } | null>(null);
  const [isUploading, setIsUploading] = useState(false);
  const [isConverting, setIsConverting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const [params, setParams] = useState<Parameters>({
    size: 64,
    scale: 8,
    colors: 32,
  });

  const debouncedParams = useDebounce(params, 300);

  const handleUpload = useCallback(async (file: File) => {
    setIsUploading(true);
    setError(null);
    setConverted(null);

    try {
      const response = await uploadImage(file);
      setImageState({
        sessionId: response.sessionId,
        original: response.original,
        width: response.width,
        height: response.height,
      });
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Upload failed');
    } finally {
      setIsUploading(false);
    }
  }, []);

  // Convert image when parameters change
  useEffect(() => {
    if (!imageState) return;

    const doConvert = async () => {
      setIsConverting(true);
      try {
        const response = await convertImage({
          sessionId: imageState.sessionId,
          ...debouncedParams,
        });
        setConverted(response.image);
        setOutputSize({ width: response.width, height: response.height });
        setError(null);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Conversion failed');
      } finally {
        setIsConverting(false);
      }
    };

    doConvert();
  }, [imageState, debouncedParams]);

  const handleDownload = useCallback(async () => {
    if (!imageState) return;

    try {
      await downloadImage({
        sessionId: imageState.sessionId,
        ...params,
      });
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Download failed');
    }
  }, [imageState, params]);

  const handleReset = useCallback(() => {
    setImageState(null);
    setConverted(null);
    setOutputSize(null);
    setError(null);
    setParams({ size: 64, scale: 8, colors: 32 });
  }, []);

  return (
    <div className="min-h-screen">
      {/* Header */}
      <header className="border-b border-ash/50">
        <div className="max-w-7xl mx-auto px-6 py-4 flex items-center justify-between">
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 bg-gradient-to-br from-pixel-cyan to-pixel-purple rounded-xl flex items-center justify-center">
              <svg
                className="w-6 h-6 text-void"
                fill="currentColor"
                viewBox="0 0 24 24"
              >
                <rect x="4" y="4" width="4" height="4" />
                <rect x="10" y="4" width="4" height="4" />
                <rect x="16" y="4" width="4" height="4" />
                <rect x="4" y="10" width="4" height="4" />
                <rect x="10" y="10" width="4" height="4" />
                <rect x="16" y="10" width="4" height="4" />
                <rect x="4" y="16" width="4" height="4" />
                <rect x="10" y="16" width="4" height="4" />
                <rect x="16" y="16" width="4" height="4" />
              </svg>
            </div>
            <div>
              <h1 className="font-display text-xl text-pearl tracking-wide">
                PIXGRID
              </h1>
              <p className="text-smoke text-xs">Pixel Art Converter</p>
            </div>
          </div>

          {imageState && (
            <button
              onClick={handleReset}
              className="flex items-center gap-2 px-4 py-2 rounded-lg text-smoke hover:text-pearl hover:bg-slate transition-all"
            >
              <svg
                className="w-4 h-4"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"
                />
              </svg>
              New Image
            </button>
          )}
        </div>
      </header>

      <main className="max-w-7xl mx-auto px-6 py-8">
        {/* Error banner */}
        {error && (
          <div className="mb-6 p-4 bg-pixel-magenta/10 border border-pixel-magenta/30 rounded-xl flex items-center gap-3">
            <svg
              className="w-5 h-5 text-pixel-magenta shrink-0"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
              />
            </svg>
            <p className="text-pearl text-sm">{error}</p>
            <button
              onClick={() => setError(null)}
              className="ml-auto text-smoke hover:text-pearl"
            >
              <svg
                className="w-4 h-4"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M6 18L18 6M6 6l12 12"
                />
              </svg>
            </button>
          </div>
        )}

        {!imageState ? (
          /* Upload state */
          <div className="max-w-2xl mx-auto">
            <div className="text-center mb-8">
              <h2 className="font-display text-3xl text-pearl tracking-wide mb-3">
                TRANSFORM YOUR IMAGES
              </h2>
              <p className="text-smoke text-lg">
                Convert photos into retro pixel art with real-time preview
              </p>
            </div>
            <ImageUploader onUpload={handleUpload} isUploading={isUploading} />

            {/* Features */}
            <div className="grid grid-cols-3 gap-6 mt-12">
              {[
                {
                  icon: '⚡',
                  title: 'Real-time',
                  desc: 'Instant preview as you adjust',
                },
                {
                  icon: '🎨',
                  title: 'Custom Palette',
                  desc: 'Control color count',
                },
                {
                  icon: '📐',
                  title: 'Scalable',
                  desc: 'Any size output',
                },
              ].map((feature) => (
                <div
                  key={feature.title}
                  className="text-center p-4 rounded-xl bg-graphite/30 border border-ash/30"
                >
                  <div className="text-2xl mb-2">{feature.icon}</div>
                  <h3 className="text-pearl font-medium text-sm">
                    {feature.title}
                  </h3>
                  <p className="text-smoke text-xs mt-1">{feature.desc}</p>
                </div>
              ))}
            </div>
          </div>
        ) : (
          /* Editor state */
          <div className="grid grid-cols-12 gap-8">
            {/* Sidebar with controls */}
            <div className="col-span-3">
              <div className="bg-graphite/50 rounded-2xl border border-ash p-6 sticky top-8">
                <ParameterSliders
                  params={params}
                  onChange={setParams}
                  disabled={isConverting}
                />
              </div>
            </div>

            {/* Preview area */}
            <div className="col-span-9">
              <PreviewPanel
                original={imageState.original}
                converted={converted}
                isConverting={isConverting}
                onDownload={handleDownload}
                outputSize={outputSize ?? undefined}
              />
            </div>
          </div>
        )}
      </main>

      {/* Footer */}
      <footer className="border-t border-ash/30 mt-16">
        <div className="max-w-7xl mx-auto px-6 py-6 text-center text-smoke text-sm">
          Built with Go + React • Open Source
        </div>
      </footer>
    </div>
  );
}

export default App;
