async function ParseResponseError(result: Response): Promise<string> {
  try {
    const data = JSON.parse(await result.text());
    if (data && data.error) {
      return data.error;
    }
    return 'error:' + data;
  } catch (err) {
    return 'error parsing response:' + err;
  }
}

export { ParseResponseError };
