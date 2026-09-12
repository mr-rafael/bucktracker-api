(function () {
  async function formatResponse(response) {
    const body = (await response.text()).trim();
    if (!body) {
      return String(response.status);
    }

    let message = body;

    try {
      const parsedBody = JSON.parse(body);
      message = parsedBody.error || parsedBody.message || JSON.stringify(parsedBody);
    } catch {
      // The backend returned plain text, so display it as-is.
    }

    return `${response.status}: ${message}`;
  }

  function formatNetworkError(error) {
    return `NETWORK_ERROR: ${error.message}`;
  }

  window.apiErrors = Object.freeze({
    formatResponse,
    formatNetworkError,
  });
})();
