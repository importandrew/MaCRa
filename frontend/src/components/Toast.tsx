import { toast } from '@heroui/react';

type ToastVariant = 'default' | 'accent' | 'success' | 'warning' | 'danger';

interface ShowToastParams {
  title: string;
  description?: string;
  variant?: ToastVariant;
}

const TOAST_TIMEOUT_MS = 3000;

export function showToast({ title, description, variant }: ShowToastParams) {
  toast(title, { description, variant, timeout: TOAST_TIMEOUT_MS });
}
