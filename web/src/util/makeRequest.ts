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
    throw new Error(`API request failed with status ${res.status}`);
  }

  const response: T = (await res.json()) as T;

  return response;
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
