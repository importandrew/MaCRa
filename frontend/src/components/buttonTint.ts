// Soft color fills for buttons, like HeroUI v2's variant="flat" with a color.
// HeroUI v3 has soft color tokens but no button variants for them (only
// danger-soft), and it dropped the v2 "secondary" (violet) color.
// Use with variant="secondary": these classes override its gray fill.
export type ButtonTint = 'accent' | 'secondary' | 'success' | 'warning';

const TINT_CLASSES: Record<ButtonTint, string> = {
  accent:
    'bg-accent-soft text-accent-soft-foreground hover:bg-accent-soft-hover',
  secondary: 'bg-violet-100 text-violet-700 hover:bg-violet-200',
  success:
    'bg-success-soft text-success-soft-foreground hover:bg-success-soft-hover',
  warning:
    'bg-warning-soft text-warning-soft-foreground hover:bg-warning-soft-hover',
};

// buttonTint returns the classes for a soft-filled button of that color.
export function buttonTint(tint: ButtonTint): string {
  return TINT_CLASSES[tint];
}
