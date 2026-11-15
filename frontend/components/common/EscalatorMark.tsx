type Props = {
  className?: string;
  railColor?: string;
  figureColor?: string;
};

export default function EscalatorMark({
  className,
  railColor = "rgba(255,255,255,0.34)",
  figureColor = "#EF5B6E",
}: Props) {
  return (
    <svg
      viewBox="0 0 100 100"
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
      className={className}
      aria-hidden="true"
    >
      <line x1="10" y1="92" x2="82" y2="28" stroke={railColor} strokeWidth="5" strokeLinecap="round" />
      <line x1="10" y1="78" x2="82" y2="14" stroke={railColor} strokeWidth="4" strokeLinecap="round" opacity="0.7" />
      <line x1="24.4" y1="78.8" x2="24.4" y2="64.8" stroke={railColor} strokeWidth="3" strokeLinecap="round" opacity="0.55" />
      <line x1="46" y1="60" x2="46" y2="46" stroke={railColor} strokeWidth="3" strokeLinecap="round" opacity="0.55" />
      <line x1="67.6" y1="41.2" x2="67.6" y2="27.2" stroke={railColor} strokeWidth="3" strokeLinecap="round" opacity="0.55" />
      <circle cx="58.6" cy="13.1" r="5.5" fill={figureColor} />
      <line x1="58.6" y1="21.1" x2="58.6" y2="35.1" stroke={figureColor} strokeWidth="4.2" strokeLinecap="round" />
      <line x1="58.6" y1="35.1" x2="55.6" y2="47.1" stroke={figureColor} strokeWidth="4.2" strokeLinecap="round" />
      <line x1="58.6" y1="35.1" x2="61.6" y2="47.1" stroke={figureColor} strokeWidth="4.2" strokeLinecap="round" />
      <line x1="58.6" y1="21.1" x2="55" y2="32" stroke={figureColor} strokeWidth="4" strokeLinecap="round" />
      <line x1="58.6" y1="21.1" x2="62.2" y2="32" stroke={figureColor} strokeWidth="4" strokeLinecap="round" />
    </svg>
  );
}
