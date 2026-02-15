const API_BASE_URL = import.meta.env.VITE_BASE_URL || "";

export const makeRequest = async <T>(
  endpoint: string,
  version: string = "v1",
  options: RequestInit = {},
  baseUrl: string = API_BASE_URL
): Promise<T> => {
  const res = await fetch(`${baseUrl}/api/${version}/${endpoint}`, {
    ...options,
    headers: {
      "Content-Type": "application/json",
      ...(options.headers || {}),
    },
  });

  if (!res.ok) {
    const message = await res.text();
    if (message) {
      throw new Error(message);
    } else {
      throw new Error(`API request ${res.status}`);
    }
  }

  // Some endpoints intentionally return no body (e.g. 204 No Content).
  const raw = await res.text();
  if (!raw) {
    return undefined as T;
  }

  return JSON.parse(raw) as T;
};

export const makePostRequest = async <T>(
  endpoint: string,
  body: any,
  version: string = "v1",
  options: RequestInit = {}
): Promise<T> => {
  return makeRequest<T>(endpoint, version, {
    method: "POST",
    body: JSON.stringify(body),
    ...options,
  });
};

export const makeGetRequest = async <T>(
  endpoint: string,
  version: string = "v1",
  options: RequestInit = {}
): Promise<T> => {
  return makeRequest<T>(endpoint, version, {
    method: "GET",
    ...options,
  });
};
