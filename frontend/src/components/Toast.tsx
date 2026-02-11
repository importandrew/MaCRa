import { addToast } from '@heroui/react';

type ToastColor =
  | 'default'
  | 'primary'
  | 'secondary'
  | 'success'
  | 'warning'
  | 'danger';

interface ShowToastParams {
  title: string;
  description?: string;
  color?: ToastColor;
}

export function showToast({ title, description, color }: ShowToastParams) {
  addToast({
    title,
    description,
    color,
    classNames: {
      title: 'whitespace-normal break-words',
      description: 'whitespace-normal break-words',
    },
  });
}
