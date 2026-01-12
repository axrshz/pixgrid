import { useCallback, useState } from 'react';

interface ImageUploaderProps {
  onUpload: (file: File) => void;
  isUploading: boolean;
}

export function ImageUploader({ onUpload, isUploading }: ImageUploaderProps) {
  const [isDragging, setIsDragging] = useState(false);

  const handleDrag = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
  }, []);

  const handleDragIn = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    if (e.dataTransfer.items && e.dataTransfer.items.length > 0) {
      setIsDragging(true);
    }
  }, []);

  const handleDragOut = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    setIsDragging(false);
  }, []);

  const handleDrop = useCallback(
    (e: React.DragEvent) => {
      e.preventDefault();
      e.stopPropagation();
      setIsDragging(false);

      if (e.dataTransfer.files && e.dataTransfer.files.length > 0) {
        const file = e.dataTransfer.files[0];
        if (file.type.startsWith('image/')) {
          onUpload(file);
        }
      }
    },
    [onUpload]
  );

  const handleFileSelect = useCallback(
    (e: React.ChangeEvent<HTMLInputElement>) => {
      if (e.target.files && e.target.files.length > 0) {
        onUpload(e.target.files[0]);
      }
    },
    [onUpload]
  );

  return (
    <div
      className={`
        relative border-2 border-dashed rounded-2xl p-12 text-center
        transition-all duration-300 cursor-pointer
        ${
          isDragging
            ? 'border-pixel-cyan bg-pixel-cyan/10 scale-[1.02]'
            : 'border-ash hover:border-smoke hover:bg-graphite/50'
        }
        ${isUploading ? 'opacity-50 pointer-events-none' : ''}
      `}
      onDragEnter={handleDragIn}
      onDragLeave={handleDragOut}
      onDragOver={handleDrag}
      onDrop={handleDrop}
      onClick={() => document.getElementById('file-input')?.click()}
    >
      <input
        id="file-input"
        type="file"
        accept="image/png,image/jpeg,image/jpg"
        className="hidden"
        onChange={handleFileSelect}
      />

      <div className="flex flex-col items-center gap-4">
        <div
          className={`
            w-20 h-20 rounded-xl flex items-center justify-center
            transition-all duration-300
            ${isDragging ? 'bg-pixel-cyan/20' : 'bg-slate'}
          `}
        >
          {isUploading ? (
            <svg
              className="w-10 h-10 text-pixel-cyan animate-spin"
              fill="none"
              viewBox="0 0 24 24"
            >
              <circle
                className="opacity-25"
                cx="12"
                cy="12"
                r="10"
                stroke="currentColor"
                strokeWidth="4"
              />
              <path
                className="opacity-75"
                fill="currentColor"
                d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
              />
            </svg>
          ) : (
            <svg
              className={`w-10 h-10 transition-colors ${
                isDragging ? 'text-pixel-cyan' : 'text-smoke'
              }`}
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
          )}
        </div>

        <div>
          <p className="text-pearl font-medium text-lg">
            {isUploading
              ? 'Uploading...'
              : isDragging
              ? 'Drop your image here'
              : 'Drop an image or click to upload'}
          </p>
          <p className="text-smoke text-sm mt-1">PNG, JPG up to 10MB</p>
        </div>
      </div>

      {isDragging && (
        <div className="absolute inset-0 bg-pixel-cyan/5 rounded-2xl pointer-events-none" />
      )}
    </div>
  );
}
