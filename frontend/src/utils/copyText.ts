// copyText copies text to the clipboard and reports whether it worked.
//
// The Clipboard API exists only on secure pages (HTTPS or localhost); a
// site opened over plain HTTP (a test stand by its address) doesn't have
// it. There the text is selected in a hidden field and copied with the
// older execCommand('copy'), which browsers still support for a user's tap.
export async function copyText(text: string): Promise<boolean> {
  if (window.isSecureContext && navigator.clipboard) {
    try {
      await navigator.clipboard.writeText(text);
      return true;
    } catch {
      // Denied: try the older way below.
    }
  }

  const field = document.createElement('textarea');
  field.value = text;
  field.readOnly = true; // no keyboard on phones
  field.style.position = 'fixed';
  field.style.opacity = '0';
  document.body.appendChild(field);
  field.select();
  field.setSelectionRange(0, text.length); // iOS selects this way
  try {
    return document.execCommand('copy');
  } catch {
    return false;
  } finally {
    field.remove();
  }
}
