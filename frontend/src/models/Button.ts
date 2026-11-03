export interface MainButtonProps {
  text: string;
  link?: string;
  onPress?: () => void;
  variant?: 'primary' | 'secondary' | 'tertiary';
}
