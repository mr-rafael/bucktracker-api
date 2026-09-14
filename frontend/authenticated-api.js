(function () {
  const API_BASE_URL = "http://localhost:8080";
  let refreshPromise = null;

  class AuthenticationRedirectError extends Error {
    constructor() {
      super("Authentication failed; redirecting to the homepage.");
      this.name = "AuthenticationRedirectError";
    }
  }

  function redirectHome() {
    sessionStorage.removeItem("access_token");
    window.location.replace("/");
    throw new AuthenticationRedirectError();
  }

  function isAccessTokenExpired(token) {
    if (!token) {
      return true;
    }

    try {
      const encodedPayload = token.split(".")[1];
      const normalizedPayload = encodedPayload.replace(/-/g, "+").replace(/_/g, "/");
      const paddedPayload = normalizedPayload.padEnd(
        normalizedPayload.length + ((4 - (normalizedPayload.length % 4)) % 4),
        "=",
      );
      const payload = JSON.parse(atob(paddedPayload));

      return typeof payload.exp !== "number" || Date.now() >= payload.exp * 1000;
    } catch {
      return true;
    }
  }

  async function requestNewAccessToken() {
    const response = await fetch(`${API_BASE_URL}/app/refresh`, {
      method: "POST",
      credentials: "include",
    });

    if (response.status === 401) {
      redirectHome();
    }

    if (!response.ok) {
      const error = new Error(await apiErrors.formatResponse(response));
      error.isBackendResponse = true;
      throw error;
    }

    let responseBody;
    try {
      responseBody = await response.json();
    } catch (error) {
      error.message = `INVALID_RESPONSE: ${error.message}`;
      error.isInvalidResponse = true;
      throw error;
    }

    if (!responseBody.access_token) {
      const error = new Error(
        "INVALID_RESPONSE: The refresh endpoint returned no access token.",
      );
      error.isInvalidResponse = true;
      throw error;
    }

    sessionStorage.setItem("access_token", responseBody.access_token);
    return responseBody.access_token;
  }

  function refreshAccessToken() {
    if (!refreshPromise) {
      refreshPromise = requestNewAccessToken().finally(() => {
        refreshPromise = null;
      });
    }

    return refreshPromise;
  }

  function sendRequest(url, options, accessToken) {
    const headers = new Headers(options.headers || {});
    headers.set("Authorization", `Bearer ${accessToken}`);

    return fetch(url, {
      ...options,
      credentials: options.credentials || "include",
      headers,
    });
  }

  async function authenticatedFetch(url, options = {}) {
    let accessToken = sessionStorage.getItem("access_token");
    let refreshedBeforeRequest = false;

    if (isAccessTokenExpired(accessToken)) {
      accessToken = await refreshAccessToken();
      refreshedBeforeRequest = true;
    }

    let response = await sendRequest(url, options, accessToken);
    if (response.status !== 401) {
      return response;
    }

    if (refreshedBeforeRequest) {
      redirectHome();
    }

    accessToken = await refreshAccessToken();
    response = await sendRequest(url, options, accessToken);

    if (response.status === 401) {
      redirectHome();
    }

    return response;
  }

  window.authenticatedApi = Object.freeze({
    fetch: authenticatedFetch,
    isRedirectError(error) {
      return error instanceof AuthenticationRedirectError;
    },
  });
})();
