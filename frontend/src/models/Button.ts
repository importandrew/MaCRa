export interface MainButtonProps {
  text: string;
  link?: string;
  onPress?: () => void;
  // tint is the soft fill color; 'default' is a plain gray button.
  tint?: 'accent' | 'secondary' | 'default';
}
