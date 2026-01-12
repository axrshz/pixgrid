const API_BASE = '/api';

export interface UploadResponse {
  sessionId: string;
  width: number;
  height: number;
  original: string;
}

export interface ConvertResponse {
  image: string;
  width: number;
  height: number;
}

export interface ConvertParams {
  sessionId: string;
  size: number;
  scale: number;
  colors: number;
}

export async function uploadImage(file: File): Promise<UploadResponse> {
  const formData = new FormData();
  formData.append('image', file);

  const response = await fetch(`${API_BASE}/upload`, {
    method: 'POST',
    body: formData,
  });

  if (!response.ok) {
    throw new Error(`Upload failed: ${response.statusText}`);
  }

  return response.json();
}

export async function convertImage(params: ConvertParams): Promise<ConvertResponse> {
  const response = await fetch(`${API_BASE}/convert`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(params),
  });

  if (!response.ok) {
    throw new Error(`Convert failed: ${response.statusText}`);
  }

  return response.json();
}

export async function downloadImage(params: ConvertParams): Promise<void> {
  const response = await fetch(`${API_BASE}/download`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(params),
  });

  if (!response.ok) {
    throw new Error(`Download failed: ${response.statusText}`);
  }

  const blob = await response.blob();
  const url = URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  a.download = 'pixelart.png';
  document.body.appendChild(a);
  a.click();
  document.body.removeChild(a);
  URL.revokeObjectURL(url);
}
