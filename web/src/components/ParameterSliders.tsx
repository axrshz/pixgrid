interface SliderProps {
  label: string;
  value: number;
  min: number;
  max: number;
  step?: number;
  onChange: (value: number) => void;
  unit?: string;
  description?: string;
}

function Slider({
  label,
  value,
  min,
  max,
  step = 1,
  onChange,
  unit = '',
  description,
}: SliderProps) {
  const percentage = ((value - min) / (max - min)) * 100;

  return (
    <div className="space-y-3">
      <div className="flex items-center justify-between">
        <div>
          <span className="text-pearl font-medium">{label}</span>
          {description && (
            <p className="text-smoke text-xs mt-0.5">{description}</p>
          )}
        </div>
        <div className="bg-slate px-3 py-1 rounded-lg">
          <span className="text-pixel-cyan font-mono font-bold">
            {value}
            {unit}
          </span>
        </div>
      </div>
      <div className="relative">
        <div className="h-2 bg-graphite rounded-full overflow-hidden">
          <div
            className="h-full bg-gradient-to-r from-pixel-purple to-pixel-cyan rounded-full transition-all duration-150"
            style={{ width: `${percentage}%` }}
          />
        </div>
        <input
          type="range"
          min={min}
          max={max}
          step={step}
          value={value}
          onChange={(e) => onChange(Number(e.target.value))}
          className="absolute inset-0 w-full h-full opacity-0 cursor-pointer"
        />
      </div>
      <div className="flex justify-between text-xs text-smoke">
        <span>{min}{unit}</span>
        <span>{max}{unit}</span>
      </div>
    </div>
  );
}

export interface Parameters {
  size: number;
  scale: number;
  colors: number;
}

interface ParameterSlidersProps {
  params: Parameters;
  onChange: (params: Parameters) => void;
  disabled?: boolean;
}

export function ParameterSliders({
  params,
  onChange,
  disabled = false,
}: ParameterSlidersProps) {
  return (
    <div
      className={`space-y-6 ${disabled ? 'opacity-50 pointer-events-none' : ''}`}
    >
      <h3 className="text-pearl font-display text-lg tracking-wide">
        PARAMETERS
      </h3>

      <Slider
        label="Pixel Size"
        description="Target width in pixels"
        value={params.size}
        min={8}
        max={256}
        step={4}
        unit="px"
        onChange={(size) => onChange({ ...params, size })}
      />

      <Slider
        label="Scale Factor"
        description="How much to enlarge the result"
        value={params.scale}
        min={1}
        max={16}
        unit="x"
        onChange={(scale) => onChange({ ...params, scale })}
      />

      <Slider
        label="Color Palette"
        description="Number of colors (0 = original)"
        value={params.colors}
        min={0}
        max={128}
        step={4}
        onChange={(colors) => onChange({ ...params, colors })}
      />

      <div className="pt-4 border-t border-ash">
        <div className="flex items-center gap-2 text-smoke text-sm">
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
              d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
            />
          </svg>
          <span>Preview updates automatically</span>
        </div>
      </div>
    </div>
  );
}
