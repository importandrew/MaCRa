export interface MainButtonProps {
  text: string;
  link?: string;
  onPress?: () => void;
  color?:
    | 'primary'
    | 'default'
    | 'secondary'
    | 'success'
    | 'warning'
    | 'danger';
}
