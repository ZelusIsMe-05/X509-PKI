// ─────────────────────────────────────────────────────────────────
// CONFIG
// ─────────────────────────────────────────────────────────────────
const API_URL = "http://localhost:8080/api";

// Keys used for localStorage persistence
const STORAGE_KEYS = {
  ACCESS_TOKEN: "auth_access_token",
  REFRESH_TOKEN: "auth_refresh_token",
  USERNAME: "auth_username",
} as const;

// ─────────────────────────────────────────────────────────────────
// TOKEN STORAGE (localStorage)
// ─────────────────────────────────────────────────────────────────

export const saveTokens = (accessToken: string, refreshToken: string, username: string) => {
  localStorage.setItem(STORAGE_KEYS.ACCESS_TOKEN, accessToken);
  localStorage.setItem(STORAGE_KEYS.REFRESH_TOKEN, refreshToken);
  localStorage.setItem(STORAGE_KEYS.USERNAME, username);
};

export const getStoredAccessToken = (): string | null =>
  localStorage.getItem(STORAGE_KEYS.ACCESS_TOKEN);

export const getStoredRefreshToken = (): string | null =>
  localStorage.getItem(STORAGE_KEYS.REFRESH_TOKEN);

export const getStoredUsername = (): string | null =>
  localStorage.getItem(STORAGE_KEYS.USERNAME);

export const clearTokens = () => {
  localStorage.removeItem(STORAGE_KEYS.ACCESS_TOKEN);
  localStorage.removeItem(STORAGE_KEYS.REFRESH_TOKEN);
  localStorage.removeItem(STORAGE_KEYS.USERNAME);
};

// ─────────────────────────────────────────────────────────────────
// REGISTER
// ─────────────────────────────────────────────────────────────────

export const registerUser = async (username: string, password: string) => {
  const response = await fetch(`${API_URL}/auth/register`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ username, password }),
  });

  if (!response.ok) {
    const errorText = await response.text();
    throw new Error(errorText.trim());
  }
  return response.json();
};

// ─────────────────────────────────────────────────────────────────
// LOGIN
// ─────────────────────────────────────────────────────────────────

export const loginUser = async (username: string, password: string) => {
  const response = await fetch(`${API_URL}/auth/login`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ username, password }),
  });

  if (!response.ok) {
    const errorText = await response.text();
    throw new Error(errorText.trim());
  }

  const data = await response.json();
  // Persist tokens and username in localStorage
  saveTokens(data.access_token, data.refresh_token, data.username);
  return data as { access_token: string; refresh_token: string; username: string };
};

// ─────────────────────────────────────────────────────────────────
// REFRESH TOKEN
// ─────────────────────────────────────────────────────────────────

// refreshAccessToken exchanges the stored refresh token for a new token pair.
// Returns true on success, false if the refresh token is invalid or expired.
export const refreshAccessToken = async (): Promise<boolean> => {
  const refreshToken = getStoredRefreshToken();
  if (!refreshToken) return false;

  try {
    const response = await fetch(`${API_URL}/auth/refresh`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ refresh_token: refreshToken }),
    });

    if (!response.ok) {
      clearTokens(); // Refresh token expired or invalid → force logout
      return false;
    }

    const data = await response.json();
    saveTokens(data.access_token, data.refresh_token, data.username);
    return true;
  } catch {
    return false;
  }
};

// ─────────────────────────────────────────────────────────────────
// VERIFY SESSION (/me)
// ─────────────────────────────────────────────────────────────────

// verifySession checks whether the stored access token is still valid by calling /api/auth/me.
// If the access token is expired, it attempts a silent refresh before retrying.
// Returns the username if the session is valid, or null otherwise.
export const verifySession = async (): Promise<string | null> => {
  const accessToken = getStoredAccessToken();
  if (!accessToken) return null;

  // Try verifying with the current access token
  const result = await callMe(accessToken);
  if (result) return result;

  // Access token expired — attempt a silent refresh
  const refreshed = await refreshAccessToken();
  if (!refreshed) return null;

  // Retry with the new access token
  const newAccessToken = getStoredAccessToken();
  if (!newAccessToken) return null;
  return await callMe(newAccessToken);
};

// callMe is a helper that calls GET /api/auth/me with the given token.
const callMe = async (accessToken: string): Promise<string | null> => {
  try {
    const response = await fetch(`${API_URL}/auth/me`, {
      method: "GET",
      headers: { Authorization: `Bearer ${accessToken}` },
    });

    if (!response.ok) return null;
    const data = await response.json();
    return data.username as string;
  } catch {
    return null;
  }
};

// ─────────────────────────────────────────────────────────────────
// API CALL WRAPPER — automatic token injection and 401 retry
// ─────────────────────────────────────────────────────────────────

// apiCall wraps fetch with automatic Authorization header injection.
// On a 401 response it silently refreshes the token and retries once.
export const apiCall = async (url: string, options: RequestInit = {}): Promise<Response> => {
  const makeRequest = (token: string) =>
    fetch(url, {
      ...options,
      headers: {
        "Content-Type": "application/json",
        ...options.headers,
        Authorization: `Bearer ${token}`,
      },
    });

  const accessToken = getStoredAccessToken();
  if (!accessToken) throw new Error("Not authenticated");

  let response = await makeRequest(accessToken);

  // On 401, attempt token refresh and retry once
  if (response.status === 401) {
    const refreshed = await refreshAccessToken();
    if (!refreshed) throw new Error("Session expired. Please log in again.");

    const newToken = getStoredAccessToken()!;
    response = await makeRequest(newToken);
  }

  return response;
};