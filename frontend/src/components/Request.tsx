function PostRequest(url: string, body: unknown = {}): Request {
  return new Request(import.meta.env.VITE_API_ADDR + url, {
    method: 'POST',
    credentials: 'include',
    body: JSON.stringify(body),
  });
}

function GetRequest(url: string, params?: Record<string, string | number>): Request {
  let fullUrl = import.meta.env.VITE_API_ADDR + url;
  if (params) {
    const searchParams = new URLSearchParams();
    Object.entries(params).forEach(([key, value]) => {
      if (value !== undefined && value !== null && value !== '') {
        searchParams.append(key, String(value));
      }
    });
    const queryString = searchParams.toString();
    if (queryString) {
      fullUrl += '?' + queryString;
    }
  }
  return new Request(fullUrl, {
    credentials: 'include',
  });
}

export { GetRequest, PostRequest };
